package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type BurnCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag                 `arg:"" name:"sender" help:"sender address; nft owner or agent" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	CSymbol  string                      `arg:"" name:"collection" help:"collection symbol" required:"true"`
	NFTs     []NFTIDFlag                 `arg:"" name:"nft" help:"target nft; \"<symbol>,<idx>\""`
	sender   base.Address
	nfts     []nft.NFTID
}

func NewBurnCommand() BurnCommand {
	return BurnCommand{
		BaseCommand: NewBaseCommand("burn-nfts-operation"),
	}
}

func (cmd *BurnCommand) Run(version util.Version) error {
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Wrap(err, "failed to initialize command")
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	bs, err := operation.NewBaseSeal(
		cmd.Privatekey,
		[]operation.Operation{op},
		cmd.NetworkID.NetworkID(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create operation.Seal")
	}
	PrettyPrint(cmd.Out, cmd.Pretty, bs)

	return nil
}

func (cmd *BurnCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	if len(cmd.NFTs) < 1 {
		return errors.Errorf("empty nfts; at least one nft is necessary")
	}

	nfts := make([]nft.NFTID, len(cmd.NFTs))
	for i := range cmd.NFTs {
		nft := nft.NewNFTID(cmd.NFTs[i].collection, cmd.NFTs[i].idx)
		if err := nft.IsValid(nil); err != nil {
			return err
		}
		nfts[i] = nft
	}
	cmd.nfts = nfts

	return nil

}

func (cmd *BurnCommand) createOperation() (operation.Operation, error) {
	items := make([]collection.BurnItem, 1)

	if len(cmd.nfts) > 1 {
		items[0] = collection.NewBurnItemMultiNFTs(extensioncurrency.ContractID(cmd.CSymbol), cmd.nfts, cmd.Currency.CID)
	} else {
		items[0] = collection.NewBurnItemSingleNFT(extensioncurrency.ContractID(cmd.CSymbol), cmd.nfts[0], cmd.Currency.CID)
	}

	fact := collection.NewBurnFact(
		[]byte(cmd.Token),
		cmd.sender,
		items,
	)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewBurn(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create burn nfts operation")
	}
	return op, nil
}
