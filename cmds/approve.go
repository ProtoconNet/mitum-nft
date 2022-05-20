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

type ApproveCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Approved AddressFlag                 `arg:"" name:"approved" help:"approved account address" required:"true"`
	NFTs     []NFTIDFlag                 `arg:"" name:"nft" help:"target nft to approve; \"<symbol>,<idx>\""`
	sender   base.Address
	approved base.Address
	nfts     []nft.NFTID
}

func NewApproveCommand() ApproveCommand {
	return ApproveCommand{
		BaseCommand: NewBaseCommand("approve-operation"),
	}
}

func (cmd *ApproveCommand) Run(version util.Version) error {
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

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Approved.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid approved format; %q", cmd.Approved.String())
	} else {
		cmd.approved = a
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

func (cmd *ApproveCommand) createOperation() (operation.Operation, error) {
	fact := collection.NewApproveFact([]byte(cmd.Token), cmd.sender, cmd.approved, cmd.nfts, cmd.Currency.CID)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewApprove(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create approve operation")
	}
	return op, nil
}
