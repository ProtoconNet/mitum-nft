package cmds

import (
	"context"

	"github.com/ProtoconNet/mitum-nft/nft"
	nftcollection "github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	currencybase "github.com/ProtoconNet/mitum-currency/v3/base"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type CollectionRegisterCommand struct {
	baseCommand
	OperationFlags
	Sender     AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   AddressFlag    `arg:"" name:"contract" help:"contract account to register policy" required:"true"`
	Collection string         `arg:"" name:"collection" help:"collection id" required:"true"`
	Name       string         `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty    uint           `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Currency   CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	URI        string         `name:"uri" help:"collection uri" optional:""`
	White      AddressFlag    `name:"white" help:"whitelisted address" optional:""`
	sender     base.Address
	contract   base.Address
	collection currencybase.ContractID
	name       nftcollection.CollectionName
	royalty    nft.PaymentParameter
	uri        nft.URI
	whitelist  []base.Address
}

func NewCollectionRegisterCommand() CollectionRegisterCommand {
	cmd := NewbaseCommand()
	return CollectionRegisterCommand{baseCommand: *cmd}
}

func (cmd *CollectionRegisterCommand) Run(pctx context.Context) error {
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

func (cmd *CollectionRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender address format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Contract.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid contract address format; %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	var white base.Address = nil
	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(enc); err != nil {
			return errors.Wrapf(err, "invalid whitelist address format, %q", cmd.White)
		} else {
			white = a
		}
	}

	collection := currencybase.ContractID(cmd.Collection)
	if err := collection.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = collection
	}

	name := nftcollection.CollectionName(cmd.Name)
	if err := name.IsValid(nil); err != nil {
		return err
	} else {
		cmd.name = name
	}

	royalty := nft.PaymentParameter(cmd.Royalty)
	if err := royalty.IsValid(nil); err != nil {
		return err
	} else {
		cmd.royalty = royalty
	}

	uri := nft.URI(cmd.URI)
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

func (cmd *CollectionRegisterCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create collection-register operation")

	fact := nftcollection.NewCollectionRegisterFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.collection,
		cmd.name,
		cmd.royalty,
		cmd.uri,
		cmd.whitelist,
		cmd.Currency.CID,
	)

	op, err := nftcollection.NewCollectionRegister(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
