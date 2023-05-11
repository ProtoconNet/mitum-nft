package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	nftcollection "github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type CollectionPolicyUpdaterCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender     cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract   cmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Collection string              `arg:"" name:"collection" help:"collection id" required:"true"`
	Name       string              `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty    uint                `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Currency   cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	URI        string              `name:"uri" help:"collection uri" optional:""`
	White      cmds.AddressFlag    `name:"white" help:"whitelisted address" optional:""`
	sender     base.Address
	contract   base.Address
	collection extensioncurrency.ContractID
	name       nftcollection.CollectionName
	royalty    nft.PaymentParameter
	uri        nft.URI
	white      []base.Address
}

func NewCollectionPolicyUpdaterCommand() CollectionPolicyUpdaterCommand {
	cmd := NewbaseCommand()
	return CollectionPolicyUpdaterCommand{baseCommand: *cmd}
}

func (cmd *CollectionPolicyUpdaterCommand) Run(pctx context.Context) error {
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

func (cmd *CollectionPolicyUpdaterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid sender address format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	if cmd.White.String() != "" {
		if a, err := cmd.White.Encode(enc); err != nil {
			return errors.Wrapf(err, "invalid whitelist address format, %q", cmd.White)
		} else {
			cmd.white = []base.Address{a}
		}
	}

	if a, err := cmd.Contract.Encode(enc); err != nil {
		return errors.Wrapf(err, "invalid contract address format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	collection := extensioncurrency.ContractID(cmd.Collection)
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

	return nil
}

func (cmd *CollectionPolicyUpdaterCommand) createOperation() (base.Operation, error) {
	e := util.StringErrorFunc("failed to create collection-policy-updater operation")

	fact := nftcollection.NewCollectionPolicyUpdaterFact(
		[]byte(cmd.Token),
		cmd.sender,
		cmd.contract,
		cmd.collection,
		cmd.name,
		cmd.royalty,
		cmd.uri,
		cmd.white,
		cmd.Currency.CID,
	)

	op, err := nftcollection.NewCollectionPolicyUpdater(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
