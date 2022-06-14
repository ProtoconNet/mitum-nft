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
	NFT      NFTIDFlag                   `arg:"" name:"nft" help:"target nft to approve; \"<symbol>,<idx>\""`
	sender   base.Address
	approved base.Address
	nft      nft.NFTID
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
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Approved.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid approved format; %q", cmd.Approved)
	} else {
		cmd.approved = a
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	cmd.nft = n

	return nil

}

func (cmd *ApproveCommand) createOperation() (operation.Operation, error) {
	item := collection.NewApproveItemSingleNFT(cmd.approved, cmd.nft, cmd.Currency.CID)

	fact := collection.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.ApproveItem{item},
	)

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
