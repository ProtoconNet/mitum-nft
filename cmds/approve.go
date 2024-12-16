package cmds

import (
	"context"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type ApproveCommand struct {
	BaseCommand
	ccmds.OperationFlags
	Sender   ccmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract ccmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Approved ccmds.AddressFlag    `arg:"" name:"approved" help:"approved account address" required:"true"`
	NFTidx   uint64               `arg:"" name:"nft" help:"target nft idx to approve"`
	Currency ccmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	contract base.Address
	approved base.Address
}

func (cmd *ApproveCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender format, %v", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract format, %v", cmd.Contract)
	} else {
		cmd.contract = a
	}

	if a, err := cmd.Approved.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid approved format, %v", cmd.Approved)
	} else {
		cmd.approved = a
	}

	return nil

}

func (cmd *ApproveCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create approve operation")

	item := nft.NewApproveItem(cmd.contract, cmd.approved, cmd.NFTidx, cmd.Currency.CID)

	fact := nft.NewApproveFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]nft.ApproveItem{item},
	)

	op, err := nft.NewApprove(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
