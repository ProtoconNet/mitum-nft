package cmds

import (
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
	NFT      NFTIDFlag                   `arg:"" name:"nft" help:"target nft; \"<symbol>,<idx>\""`
	sender   base.Address
	nft      nft.NFTID
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
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	cmd.nft = n

	return nil

}

func (cmd *BurnCommand) createOperation() (operation.Operation, error) {
	item := collection.NewBurnItem(cmd.nft, cmd.Currency.CID)

	fact := collection.NewBurnFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.BurnItem{item},
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
