package cmds

import (
	"context"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	"github.com/ProtoconNet/mitum-currency/v2/cmds"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

type MintCommand struct {
	baseCommand
	cmds.OperationFlags
	Sender           cmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract         cmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Collection       string              `arg:"" name:"collection" help:"collection symbol" required:"true"`
	Hash             string              `arg:"" name:"hash" help:"nft hash" required:"true"`
	Uri              string              `arg:"" name:"uri" help:"nft uri" required:"true"`
	Currency         cmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Creator          SignerFlag          `name:"creator" help:"nft contents creator \"<address>,<share>\"" optional:""`
	Copyrighter      SignerFlag          `name:"copyrighter" help:"nft contents copyrighter \"<address>,<share>\"" optional:""`
	CreatorTotal     uint                `name:"creator-total" help:"creators total share" optional:""`
	CopyrighterTotal uint                `name:"copyrighter-total" help:"copyrighters total share" optional:""`
	sender           base.Address
	contract         base.Address
	collection       extensioncurrency.ContractID
	form             collection.MintForm
}

func NewMintCommand() MintCommand {
	cmd := NewbaseCommand()
	return MintCommand{baseCommand: *cmd}
}

func (cmd *MintCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *MintCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender address format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Contract.Encode(enc)
	if err != nil {
		return errors.Wrapf(err, "invalid contract address format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	col := extensioncurrency.ContractID(cmd.Collection)
	if err = col.IsValid(nil); err != nil {
		return err
	} else {
		cmd.collection = col
	}

	hash := nft.NFTHash(cmd.Hash)
	if err := hash.IsValid(nil); err != nil {
		return err
	}

	uri := nft.URI(cmd.Uri)
	if err := uri.IsValid(nil); err != nil {
		return err
	}

	var crts = []nft.Signer{}
	if len(cmd.Creator.address) > 0 {
		a, err := cmd.Creator.Encode(enc)
		if err != nil {
			return errors.Wrapf(err, "invalid creator address format, %q", cmd.Creator)
		}

		signer := nft.NewSigner(a, cmd.Creator.share, false)
		if err = signer.IsValid(nil); err != nil {
			return err
		}

		crts = append(crts, signer)
	}

	var cprs = []nft.Signer{}
	if len(cmd.Copyrighter.address) > 0 {
		a, err := cmd.Copyrighter.Encode(enc)
		if err != nil {
			return errors.Wrapf(err, "invalid copyrighter address format, %q", cmd.Copyrighter)
		}

		signer := nft.NewSigner(a, cmd.Copyrighter.share, false)
		if err = signer.IsValid(nil); err != nil {
			return err
		}

		cprs = append(cprs, signer)
	}

	creators := nft.NewSigners(cmd.CreatorTotal, crts)
	if err := creators.IsValid(nil); err != nil {
		return err
	}

	copyrighters := nft.NewSigners(cmd.CopyrighterTotal, cprs)
	if err := copyrighters.IsValid(nil); err != nil {
		return err
	}

	form := collection.NewMintForm(hash, uri, creators, copyrighters)
	if err := form.IsValid(nil); err != nil {
		return err
	}
	cmd.form = form

	return nil

}

func (cmd *MintCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringErrorFunc("failed to create mint operation")

	item := collection.NewMintItem(cmd.contract, cmd.collection, cmd.form, cmd.Currency.CID)
	fact := collection.NewMintFact([]byte(cmd.Token), cmd.sender, []collection.MintItem{item})

	op, err := collection.NewMint(fact)
	if err != nil {
		return nil, e(err, "")
	}
	err = op.HashSign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e(err, "")
	}

	return op, nil
}
