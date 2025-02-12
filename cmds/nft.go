package cmds

type NFTCommand struct {
	RegisterModel     RegisterModelCommand     `cmd:"" name:"register-model" help:"register new nft service"`
	UpdateModelConfig UpdateModelConfigCommand `cmd:"" name:"update-model-config" help:"update model config"`
	Mint              MintCommand              `cmd:"" name:"mint" help:"mint new nft to collection"`
	Transfer          TransferCommand          `cmd:"" name:"transfer" help:"transfer nfts to receiver"`
	Delegate          DelegateCommand          `cmd:"" name:"delegate" help:"delegate operator or cancel operator delegation"`
	Approve           ApproveCommand           `cmd:"" name:"approve" help:"approve account for nft"`
	Sign              SignCommand              `cmd:"" name:"sign" help:"sign nft as creator | copyrighter"`
}
