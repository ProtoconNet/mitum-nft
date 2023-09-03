package cmds

type NFTCommand struct {
	CreateCollection       CreateCollectionCommand       `cmd:"" name:"create-collection" help:"register new collection design"`
	UpdateCollectionPolicy UpdateCollectionPolicyCommand `cmd:"" name:"update-collection-policy" help:"update collection design"`
	Mint                   MintCommand                   `cmd:"" name:"mint" help:"mint new nft to collection"`
	Transfer               TransferCommand               `cmd:"" name:"transfer" help:"transfer nfts to receiver"`
	Delegate               DelegateCommand               `cmd:"" name:"delegate" help:"delegate operator or cancel operator delegation"`
	Approve                ApproveCommand                `cmd:"" name:"approve" help:"approve account for nft"`
	Sign                   SignCommand                   `cmd:"" name:"sign" help:"sign nft as creator | copyrighter"`
}
