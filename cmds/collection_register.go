package cmds

import (
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type CollectionRegisterCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Target   AddressFlag                 `arg:"" name:"target" help:"target account to register policy" required:"true"`
	CSymbol  string                      `arg:"" name:"symbol" help:"collection symbol" required:"true"`
	Name     string                      `arg:"" name:"name" help:"collection name" required:"true"`
	Royalty  uint                        `arg:"" name:"royalty" help:"royalty parameter; 0 <= royalty param < 100" required:"true"`
	Uri      string                      `name:"uri" help:"collection uri" optional:""`
	sender   base.Address
	target   base.Address
	policy   collection.CollectionPolicy
}

func NewCollectionRegisterCommand() CollectionRegisterCommand {
	return CollectionRegisterCommand{
		BaseCommand: NewBaseCommand("collection-register-operation"),
	}
}

func (cmd *CollectionRegisterCommand) Run(version util.Version) error {
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

func (cmd *CollectionRegisterCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender.String())
	} else {
		cmd.sender = a
	}

	if a, err := cmd.Target.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid approved format; %q", cmd.Target.String())
	} else {
		cmd.target = a
	}

	symbol := nft.Symbol(cmd.CSymbol)
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

	var uri = collection.CollectionUri(cmd.Uri)
	if len(cmd.Uri) > 0 {
		if err := uri.IsValid(nil); err != nil {
			return err
		}
	}

	policy := collection.NewCollectionPolicy(symbol, name, cmd.sender, royalty, uri)
	if err := policy.IsValid(nil); err != nil {
		return err
	}
	cmd.policy = policy

	return nil

}

func (cmd *CollectionRegisterCommand) createOperation() (operation.Operation, error) {
	fact := collection.NewCollectionRegisterFact([]byte(cmd.Token), cmd.sender, cmd.target, cmd.policy, cmd.Currency.CID)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewCollectionRegister(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collection-register operation")
	}
	return op, nil
}
