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

type SignCommand struct {
	*BaseCommand
	OperationFlags
	Sender        AddressFlag                 `arg:"" name:"sender" help:"sender address; nft owner or agent" required:"true"`
	Currency      currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	NFT           NFTIDFlag                   `arg:"" name:"nft" help:"target nft; \"<symbol>,<idx>\""`
	Qualification string                      `name:"qualification" help:"target qualification; creator | copyrighter" optional:""`
	sender        base.Address
	nft           nft.NFTID
	qualification collection.Qualification
}

func NewSignCommand() SignCommand {
	return SignCommand{
		BaseCommand: NewBaseCommand("sign-operation"),
	}
}

func (cmd *SignCommand) Run(version util.Version) error {
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

func (cmd *SignCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Wrapf(err, "invalid sender format; %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	n := nft.NewNFTID(cmd.NFT.collection, cmd.NFT.idx)
	if err := n.IsValid(nil); err != nil {
		return err
	}
	cmd.nft = n

	if cmd.Qualification == "" {
		cmd.qualification = collection.CreatorQualification
	} else {
		q := collection.Qualification(cmd.Qualification)
		if err := q.IsValid(nil); err != nil {
			return err
		}
		cmd.qualification = q
	}

	return nil

}

func (cmd *SignCommand) createOperation() (operation.Operation, error) {
	item := collection.NewSignItem(cmd.qualification, cmd.nft, cmd.Currency.CID)

	fact := collection.NewSignFact(
		[]byte(cmd.Token),
		cmd.sender,
		[]collection.SignItem{item},
	)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := collection.NewSign(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sign operation")
	}

	return op, nil
}
