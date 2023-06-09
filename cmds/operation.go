package cmds

import (
	"github.com/ProtoconNet/mitum-currency/v3/cmds"
)

type OperationCommand struct {
	CreateAccount           cmds.CreateAccountCommand         `cmd:"" name:"create-account" help:"create new account"`
	KeyUpdater              cmds.KeyUpdaterCommand            `cmd:"" name:"key-updater" help:"update account keys"`
	Transfer                cmds.TransferCommand              `cmd:"" name:"transfer" help:"transfer amounts to receiver"`
	CreateContractAccount   cmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw                cmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw amounts from target contract account"`
	CurrencyRegister        cmds.CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater   cmds.CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`
	SuffrageInflation       cmds.SuffrageInflationCommand     `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"`
	CollectionRegister      CollectionRegisterCommand         `cmd:"" name:"collection-register" help:"register new collection design"`
	CollectionPolicyUpdater CollectionPolicyUpdaterCommand    `cmd:"" name:"collection-policy-updater" help:"update collection design"`
	Mint                    MintCommand                       `cmd:"" name:"mint" help:"mint new nft to collection"`
	NFTTransfer             NFTTransferCommand                `cmd:"" name:"nft-transfer" help:"transfer nfts to receiver"`
	Delegate                DelegateCommand                   `cmd:"" name:"delegate" help:"delegate operator or cancel operator delegation"`
	Approve                 ApproveCommand                    `cmd:"" name:"approve" help:"approve account for nft"`
	NFTSign                 NFTSignCommand                    `cmd:"" name:"nft-sign" help:"sign nft as creator | copyrighter"`
	SuffrageCandidate       cmds.SuffrageCandidateCommand     `cmd:"" name:"suffrage-candidate" help:"suffrage candidate operation"`
	SuffrageJoin            cmds.SuffrageJoinCommand          `cmd:"" name:"suffrage-join" help:"suffrage join operation"`
	SuffrageDisjoin         cmds.SuffrageDisjoinCommand       `cmd:"" name:"suffrage-disjoin" help:"suffrage disjoin operation"` // revive:disable-line:line-length-limit
}

func NewOperationCommand() OperationCommand {
	return OperationCommand{
		CreateAccount:           cmds.NewCreateAccountCommand(),
		KeyUpdater:              cmds.NewKeyUpdaterCommand(),
		Transfer:                cmds.NewTransferCommand(),
		CreateContractAccount:   cmds.NewCreateContractAccountCommand(),
		Withdraw:                cmds.NewWithdrawCommand(),
		CurrencyRegister:        cmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater:   cmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:       cmds.NewSuffrageInflationCommand(),
		CollectionRegister:      NewCollectionRegisterCommand(),
		CollectionPolicyUpdater: NewCollectionPolicyUpdaterCommand(),
		Mint:                    NewMintCommand(),
		NFTTransfer:             NewNFTTranfserCommand(),
		Delegate:                NewDelegateCommand(),
		Approve:                 NewApproveCommand(),
		NFTSign:                 NewNFTSignCommand(),
		SuffrageCandidate:       cmds.NewSuffrageCandidateCommand(),
		SuffrageJoin:            cmds.NewSuffrageJoinCommand(),
		SuffrageDisjoin:         cmds.NewSuffrageDisjoinCommand(),
	}
}
