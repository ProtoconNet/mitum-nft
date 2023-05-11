package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type ApproveCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender     cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   cmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Collection string              `arg:"" name:"collection" help:"collection id" required:"true"`
	Approved   cmds.AddressFlag    `arg:"" name:"approved" help:"approved account address" required:"true"`
	NFTidx     uint64              `arg:"" name:"nft" help:"target nft idx to approve"`
	Currency   cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	collection extensioncurrency.ContractID
	approved   base.Address
}

func NewApproveCommand() ApproveCommand {
	cmd := NewbaseCommand()
	return ApproveCommand{baseCommand: *cmd}
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	encs = cmd.encs
	enc = cmd.enc

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid contract format, %q", cmd.Sender)
	} else {
		cmd.contract = a
	}

	collection := extensioncurrency.ContractID(cmd.Collection)
	if err := collection.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = collection
	}

	if a, err := cmd.Approved.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved)
	} else {
		cmd.approved = a
	}

	return nil

}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create approve operation")

	item := collection.NewApproveItem(cmd.contract, cmd.collection, cmd.approved, cmd.NFTidx, cmd.Currency.CID)

	fact := collection.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.ApproveItem{item},
	)

	op, err := collection.NewApprove(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
