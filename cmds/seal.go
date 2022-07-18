package cmds

import (
	extensioncmds "github.com/ProtoconNet/mitum-currency-extension/cmds"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
)

type SealCommand struct {
	Send                    SendCommand                                `cmd:"" name:"send" help:"send seal to remote mitum node"`
	CreateAccount           currencycmds.CreateAccountCommand          `cmd:"" name:"create-account" help:"create new account"`
	CreateContractAccount   extensioncmds.CreateContractAccountCommand `cmd:"" name:"create-contract-account" help:"create new contract account"`
	Withdraw                extensioncmds.WithdrawCommand              `cmd:"" name:"withdraw" help:"withdraw contract account"`
	Delegate                DelegateCommand                            `cmd:"" name:"delegate" help:"delegate agent or cancel agent delegation"`
	Approve                 ApproveCommand                             `cmd:"" name:"approve" help:"approve account for nft"`
	CollectionRegister      CollectionRegisterCommand                  `cmd:"" name:"collection-register" help:"register collection to contract account"`
	CollectionPolicyUpdater CollectionPolicyUpdaterCommand             `cmd:"" name:"collection-policy-updater" help:"update collection policy"`
	Mint                    MintCommand                                `cmd:"" name:"mint" help:"mint nft to collection"`
	TransferNFTs            TransferCommand                            `cmd:"" name:"transfer-nfts" help:"transfer nfts"`
	Burn                    BurnCommand                                `cmd:"" name:"burn" help:"burn nfts"`
	SignNFTs                SignCommand                                `cmd:"" name:"sign-nfts" help:"sign nfts; creator | copyrighter"`
	Transfer                currencycmds.TransferCommand               `cmd:"" name:"transfer" help:"transfer big"`
	KeyUpdater              currencycmds.KeyUpdaterCommand             `cmd:"" name:"key-updater" help:"update keys"`
	CurrencyRegister        extensioncmds.CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater   extensioncmds.CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"`  // revive:disable-line:line-length-limit
	SuffrageInflation       currencycmds.SuffrageInflationCommand      `cmd:"" name:"suffrage-inflation" help:"suffrage inflation operation"` // revive:disable-line:line-length-limit
	Sign                    currencycmds.SignSealCommand               `cmd:"" name:"sign" help:"sign seal"`
	SignFact                currencycmds.SignFactCommand               `cmd:"" name:"sign-fact" help:"sign facts of operation seal"`
}

func NewSealCommand() SealCommand {
	return SealCommand{
		Send:                    NewSendCommand(),
		CreateAccount:           currencycmds.NewCreateAccountCommand(),
		CreateContractAccount:   extensioncmds.NewCreateContractAccountCommand(),
		Withdraw:                extensioncmds.NewWithdrawCommand(),
		Delegate:                NewDelegateCommand(),
		Approve:                 NewApproveCommand(),
		CollectionRegister:      NewCollectionRegisterCommand(),
		CollectionPolicyUpdater: NewCollectionPolicyUpdaterCommand(),
		Mint:                    NewMintCommand(),
		TransferNFTs:            NewTransferCommand(),
		Burn:                    NewBurnCommand(),
		SignNFTs:                NewSignCommand(),
		Transfer:                currencycmds.NewTransferCommand(),
		KeyUpdater:              currencycmds.NewKeyUpdaterCommand(),
		CurrencyRegister:        extensioncmds.NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater:   extensioncmds.NewCurrencyPolicyUpdaterCommand(),
		SuffrageInflation:       currencycmds.NewSuffrageInflationCommand(),
		Sign:                    currencycmds.NewSignSealCommand(),
		SignFact:                currencycmds.NewSignFactCommand(),
	}
}
