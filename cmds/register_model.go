package cmds

import (
	"context"

	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type RegisterModelCommand struct {
	BaseCommand
	ccmds.OperationFlags
	Sender    ccmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract  ccmds.AddressFlag    `arg:"" name:"contract" help:"contract account to register policy" required:"true"`
	Name      string               `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty   uint                 `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Currency  ccmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	URI       string               `name:"uri" help:"collection uri" optional:""`
	White     ccmds.AddressFlag    `name:"white" help:"whitelisted address" optional:""`
	sender    base.Address
	contract  base.Address
	name      types.CollectionName
	royalty   types.PaymentParameter
	uri       types.URI
	whitelist []base.Address
}

func (cmd *RegisterModelCommand) Run(pctx context.Context) error {
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

func (cmd *RegisterModelCommand) parseFlags() error {
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

	var white base.Address = nil
	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(cmd.Encoders.JSON()); err != nil {
			return errors.Wrapf(err, "invalid whitelist address format, %v", cmd.White)
		} else {
			white = a
		}
	}

	name := types.CollectionName(cmd.Name)
	if err := name.IsValid(nil); err != nil {
		return err
	} else {
		cmd.name = name
	}

	royalty := types.PaymentParameter(cmd.Royalty)
	if err := royalty.IsValid(nil); err != nil {
		return err
	} else {
		cmd.royalty = royalty
	}

	uri := types.URI(cmd.URI)
	if err := uri.IsValid(nil); err != nil {
		return err
	} else {
		cmd.uri = uri
	}

	whitelist := []base.Address{}
	if white != nil {
		whitelist = append(whitelist, white)
	} else {
		cmd.whitelist = whitelist
	}

	return nil
}

func (cmd *RegisterModelCommand) createOperation() (base.Operation, error) {
	e := util.StringError("failed to create register-model operation")

	fact := nft.NewRegisterModelFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.name,
		cmd.royalty,
		cmd.uri,
		cmd.whitelist,
		cmd.Currency.CID,
	)

	op, err := nft.NewRegisterModel(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
