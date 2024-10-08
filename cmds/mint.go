package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type MintCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	Receiver currencycmds.AddressFlag    `arg:"" name:"receiver" help:"receiver address" required:"true"`
	Hash     string                      `arg:"" name:"hash" help:"nft hash" required:"true"`
	Uri      string                      `arg:"" name:"uri" help:"nft uri" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Creator  SignerFlag                  `name:"creator" help:"nft contents creator \"<address>,<share>\"" optional:""`
	sender   base.Address
	contract base.Address
	receiver base.Address
	hash     types.NFTHash
	uri      types.URI
	creators types.Signers
}

func (cmd *MintCommand) Run(pctx context.Context) error { // nolint:dupl
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

func (cmd *MintCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender address format, %v", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract address format, %v", cmd.Contract)
	} else {
		cmd.contract = a
	}

	a, err = cmd.Receiver.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid receiver address format, %v", cmd.Receiver)
	} else {
		cmd.receiver = a
	}

	hash := types.NFTHash(cmd.Hash)
	if err := hash.IsValid(nil); err != nil {
		return err
	} else {
		cmd.hash = hash
	}

	uri := types.URI(cmd.Uri)
	if err := uri.IsValid(nil); err != nil {
		return err
	} else {
		cmd.uri = uri
	}

	var crts []types.Signer
	if len(cmd.Creator.address) > 0 {
		a, err := cmd.Creator.Encode(cmd.Encoders.JSON())
		if err != nil {
			return errors.Wrapf(err, "invalid creator address format, %v", cmd.Creator)
		}

		signer := types.NewSigner(a, cmd.Creator.share, false)
		if err = signer.IsValid(nil); err != nil {
			return err
		}

		crts = append(crts, signer)
	}

	creators := types.NewSigners(crts)
	if err := creators.IsValid(nil); err != nil {
		return err
	} else {
		cmd.creators = creators
	}

	return nil

}

func (cmd *MintCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create mint operation")

	item := nft.NewMintItem(cmd.contract, cmd.receiver, cmd.hash, cmd.uri, cmd.creators, cmd.Currency.CID)
	fact := nft.NewMintFact([]byte(cmd.Token), cmd.sender, []nft.MintItem{item})

	op, err := nft.NewMint(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
