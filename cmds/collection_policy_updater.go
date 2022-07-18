package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type CollectionPolicyUpdaterCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	CSymbol  string                      `arg:"" name:"symbol" help:"collection symbol" required:"true"`
	Name     string                      `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty  uint                        `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Uri      string                      `name:"uri" help:"collection uri" optional:""`
	White    AddressFlag                 `name:"white" help:"whitelisted address" optional:""`
	sender   base.Address
	policy   collection.CollectionPolicy
}

func NewCollectionPolicyUpdaterCommand() CollectionPolicyUpdaterCommand {
	return CollectionPolicyUpdaterCommand{
		BaseCommand: NewBaseCommand("collection-policy-updater-operation"),
	}
}

func (cmd *CollectionPolicyUpdaterCommand) Run(version util.Version) error {
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

func (cmd *CollectionPolicyUpdaterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	var white base.Address = nil
	if cmd.White.s != "" {
		if a, err := cmd.White.Encode(jenc); err != nil {
			return errors.Wrapf(err, "invalid white format; %q", cmd.White)
		} else {
			white = a
		}
	}

	symbol := extensioncurrency.ContractID(cmd.CSymbol)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}

	name := collection.CollectionName(cmd.Name)
	if err := name.IsValid(nil); err != nil {
		return err
	}

	royalty := nft.PaymentParameter(cmd.Royalty)
	if err := royalty.IsValid(nil); err != nil {
		return err
	}

	uri := nft.URI(cmd.Uri)
	if err := uri.IsValid(nil); err != nil {
		return err
	}

	whites := []base.Address{}
	if white != nil {
		whites = append(whites, white)
	}

	policy := collection.NewCollectionPolicy(name, royalty, uri, whites)
	if err := policy.IsValid(nil); err != nil {
		return err
	}
	cmd.policy = policy

	return nil
}

func (cmd *CollectionPolicyUpdaterCommand) createOperation() (operation.Operation, error) {
	fact := collection.NewCollectionPolicyUpdaterFact([]byte(cmd.Token), cmd.sender, extensioncurrency.ContractID(cmd.CSymbol), cmd.policy, cmd.Currency.CID)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewCollectionPolicyUpdater(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collection-policy-updater operation")
	}
	return op, nil
}
