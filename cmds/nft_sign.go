package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/pkg/errors"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type NFTSignCommand struct {
	baseCommand
	OperationFlags
	Sender     AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Collection string         `arg:"" name:"collection" help:"collection id" required:"true"`
	NFT        uint64         `arg:"" name:"nft" help:"target nft; \"<collection>,<idx>\""`
	Currency   CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender     base.Address
	contract   base.Address
	collection currencybase.ContractID
}

func NewNFTSignCommand() NFTSignCommand {
	cmd := NewbaseCommand()
	return NFTSignCommand{baseCommand: *cmd}
}

func (cmd *NFTSignCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *NFTSignCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender address format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid contract address format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	col := currencybase.ContractID(cmd.Collection)
	if err := col.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = col
	}

	return nil

}

func (cmd *NFTSignCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create nft-sign operation")

	item := collection.NewNFTSignItem(cmd.contract, cmd.collection, cmd.NFT, cmd.Currency.CID)

	fact := collection.NewNFTSignFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.NFTSignItem{item},
	)

	op, err := collection.NewNFTSign(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
