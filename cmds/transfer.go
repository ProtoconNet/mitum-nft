package cmds

import (
	"context"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type TransferCommand struct {
	BaseCommand
	ccmds.OperationFlags
	Sender     ccmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Receiver   ccmds.AddressFlag    `arg:"" name:"receiver" help:"nft owner" required:"true"`
	Contract   ccmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	NFT        uint64               `arg:"" name:"nft" help:"target nft"`
	Currency   ccmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender     base.Address
	receiver   base.Address
	contract   base.Address
	collection types.ContractID
}

func (cmd *TransferCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *TransferCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender address format, %v", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Receiver.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid receiver format, %v", cmd.Receiver.String())
	} else {
		cmd.receiver = a
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract address format, %v", cmd.Contract.String())
	} else {
		cmd.contract = a
	}

	return nil

}

func (cmd *TransferCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create transfer operation")

	item := nft.NewTransferItem(cmd.contract, cmd.receiver, cmd.NFT, cmd.Currency.CID)
	fact := nft.NewTransferFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]nft.TransferItem{item},
	)

	op, err := nft.NewTransfer(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
