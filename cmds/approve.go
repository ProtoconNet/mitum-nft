package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type ApproveCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Approved currencycmds.AddressFlag    `arg:"" name:"approved" help:"approved account address" required:"true"`
	NFTidx   uint64                      `arg:"" name:"nft" help:"target nft idx to approve"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   mitumbase.Address
	contract mitumbase.Address
	approved mitumbase.Address
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

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *ApproveCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract format, %q", cmd.Sender)
	} else {
		cmd.contract = a
	}

	if a, err := cmd.Approved.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid approved format, %q", cmd.Approved)
	} else {
		cmd.approved = a
	}

	return nil

}

func (cmd *ApproveCommand) createOperation() (mitumbase.Operation, error) {
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
