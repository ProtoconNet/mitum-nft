package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-nft/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type UpdateCollectionPolicyCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Name     string                      `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty  uint                        `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	URI      string                      `name:"uri" help:"collection uri" optional:""`
	White    currencycmds.AddressFlag    `name:"white" help:"whitelisted address" optional:""`
	sender   mitumbase.Address
	contract mitumbase.Address
	name     types.CollectionName
	royalty  types.PaymentParameter
	uri      types.URI
	white    []mitumbase.Address
}

func (cmd *UpdateCollectionPolicyCommand) Run(pctx context.Context) error {
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

func (cmd *UpdateCollectionPolicyCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid sender address format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(cmd.Encoders.JSON()); err != nil {
			return errors.Wrapf(err, "invalid whitelist address format, %q", cmd.White)
		} else {
			cmd.white = []mitumbase.Address{a}
		}
	}

	if a, err := cmd.Contract.Encode(cmd.Encoders.JSON()); err != nil {
		return errors.Wrapf(err, "invalid contract address format, %q", cmd.Contract)
	} else {
		cmd.contract = a
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

	return nil
}

func (cmd *UpdateCollectionPolicyCommand) createOperation() (mitumbase.Operation, error) {
	e := util.StringError("failed to create update-collection-policy operation")

	fact := nft.NewUpdateCollectionPolicyFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.name,
		cmd.royalty,
		cmd.uri,
		cmd.white,
		cmd.Currency.CID,
	)

	op, err := nft.NewUpdateCollectionPolicy(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
