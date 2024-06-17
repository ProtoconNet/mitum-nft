package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type DelegateCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Operator currencycmds.AddressFlag    `arg:"" name:"operator" help:"operator account address"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Mode     string                      `name:"mode" help:"delegate mode" optional:""`
	sender   base.Address
	contract base.Address
	operator base.Address
	mode     nft.ApproveAllMode
}

func (cmd *DelegateCommand) Run(pctx context.Context) error {
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

func (cmd *DelegateCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender address format; %v", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract address format; %v", cmd.Contract)
	} else {
		cmd.contract = a
	}

	if a, err := cmd.Operator.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid operator address format; %v", cmd.Operator)
	} else {
		cmd.operator = a
	}

	if len(cmd.Mode) < 1 {
		cmd.mode = nft.ApproveAllAllow
	} else {
		mode := nft.ApproveAllMode(cmd.Mode)
		if err := mode.IsValid(nil); err != nil {
			return err
		}
		cmd.mode = mode
	}

	return nil

}

func (cmd *DelegateCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create delegate operation")

	items := []nft.ApproveAllItem{nft.NewApproveAllItem(cmd.contract, cmd.operator, cmd.mode, cmd.Currency.CID)}

	fact := nft.NewApproveAllFact([]byte(cmd.Token), cmd.sender, items)

	op, err := nft.NewDelegate(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
