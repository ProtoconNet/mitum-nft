package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	timestampservice "github.com/ProtoconNet/mitum-nft/v2/timestamp/service"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type ServiceRegisterCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender   cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Target   cmds.AddressFlag    `arg:"" name:"target" help:"target account to register policy" required:"true"`
	Service  ContractIDFlag      `arg:"" name:"service" help:"STO ID" required:"true"`
	Currency cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender   base.Address
	target   base.Address
}

func NewServiceRegisterCommand() ServiceRegisterCommand {
	cmd := NewbaseCommand()
	return ServiceRegisterCommand{baseCommand: *cmd}
}

func (cmd *ServiceRegisterCommand) Run(pctx context.Context) error {
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

func (cmd *ServiceRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Target.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid target format; %q", cmd.Target)
	} else {
		cmd.target = a
	}

	return nil
}

func (cmd *ServiceRegisterCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create service-register operation")

	fact := timestampservice.NewServiceRegisterFact([]byte(cmd.Token), cmd.sender, cmd.target, cmd.Service.CID, cmd.Currency.CID)

	op, err := timestampservice.NewServiceRegister(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
