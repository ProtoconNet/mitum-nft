package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type DelegateCommand struct {
	*BaseCommand
	OperationFlags
	Sender     AddressFlag                 `arg:"" name:"sender" help:"sender address" required:"true"`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	Collection string                      `arg:"" name:"collection" help:"collection symbol" required:"true"`
	Agent      AddressFlag                 `arg:"" name:"agent" help:"agent account address"`
	Mode       string                      `name:"mode" help:"delegate mode" optional:""`
	sender     base.Address
	symbol     extensioncurrency.ContractID
	agent      base.Address
	mode       collection.DelegateMode
}

func NewDelegateCommand() DelegateCommand {
	return DelegateCommand{
		BaseCommand: NewBaseCommand("delegate-operation"),
	}
}

func (cmd *DelegateCommand) Run(version util.Version) error {
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

func (cmd *DelegateCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	symbol := extensioncurrency.ContractID(cmd.Collection)
	if err := symbol.IsValid(nil); err != nil {
		return err
	}
	cmd.symbol = symbol

	if a, err := cmd.Agent.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid agent format; %q", cmd.Agent)
	} else {
		cmd.agent = a
	}

	if len(cmd.Mode) < 1 {
		cmd.mode = collection.DelegateAllow
	} else {
		mode := collection.DelegateMode(cmd.Mode)
		if err := mode.IsValid(nil); err != nil {
			return err
		}
		cmd.mode = mode
	}

	return nil

}

func (cmd *DelegateCommand) createOperation() (operation.Operation, error) {
	items := []collection.DelegateItem{collection.NewDelegateItem(cmd.symbol, cmd.agent, cmd.mode, cmd.Currency.CID)}

	fact := collection.NewDelegateFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewDelegate(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create delegate operation")
	}
	return op, nil
}
