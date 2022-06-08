package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type MintCommand struct {
	*BaseCommand
	OperationFlags
	Sender      AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	CSymbol     string                      `arg:"" name:"collection" help:"collection symbol" required:"true"`
	Hash        string                      `arg:"" name:"hash" help:"nft hash" required:"true"`
	Uri         string                      `arg:"" name:"uri" help:"nft uri" required:"true"`
	Copyrighter AddressFlag                 `name:"copyrighter" help:"nft copyrighter" optional:""`
	sender      base.Address
	form        collection.MintForm
}

func NewMintCommand() MintCommand {
	return MintCommand{
		BaseCommand: NewBaseCommand("mint-operation"),
	}
}

func (cmd *MintCommand) Run(version util.Version) error {
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Wrap(err, "failed to initialize command")
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	bs, err := operation.NewBaseSeal(
		cmd.Privatekey,
		[]operation.Operation{op},
		cmd.NetworkID.NetworkID(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create operation.Seal")
	}
	PrettyPrint(cmd.Out, cmd.Pretty, bs)

	return nil
}

func (cmd *MintCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	hash := nft.NFTHash(cmd.Hash)
	if err := hash.IsValid(nil); err != nil {
		return err
	}

	uri := nft.URI(cmd.Uri)
	if err := uri.IsValid(nil); err != nil {
		return err
	}

	var copyrighter base.Address
	if len(cmd.Copyrighter.String()) < 1 {
		copyrighter = currency.Address{}
	} else {
		if a, err := cmd.Copyrighter.Encode(jenc); err != nil {
			return errors.Wrapf(err, "invalid copyrighter format; %q", cmd.Copyrighter.String())
		} else {
			copyrighter = a
		}

		if err := copyrighter.IsValid(nil); err != nil {
			return err
		}
	}

	form := collection.NewMintForm(hash, uri, copyrighter)
	if err := form.IsValid(nil); err != nil {
		return err
	}
	cmd.form = form

	return nil

}

func (cmd *MintCommand) createOperation() (operation.Operation, error) {
	item := collection.NewMintItemSingleNFT(extensioncurrency.ContractID(cmd.CSymbol), cmd.form, cmd.Currency.CID)
	if err := item.IsValid(nil); err != nil {
		return nil, err
	}

	fact := collection.NewMintFact([]byte(cmd.Token), cmd.sender, []collection.MintItem{item})

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewMint(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mint operation")
	}
	return op, nil
}
