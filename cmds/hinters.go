package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: types.SignerHint, Instance: types.Signer{}},
	{Hint: types.SignersHint, Instance: types.Signers{}},
	{Hint: types.NFTHint, Instance: types.NFT{}},
	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.OperatorsBookHint, Instance: types.OperatorsBook{}},
	{Hint: types.CollectionPolicyHint, Instance: types.CollectionPolicy{}},
	{Hint: types.CollectionDesignHint, Instance: types.CollectionDesign{}},
	{Hint: types.NFTBoxHint, Instance: types.NFTBox{}},

	{Hint: nft.CreateCollectionHint, Instance: nft.CreateCollection{}},
	{Hint: nft.UpdateCollectionPolicyHint, Instance: nft.UpdateCollectionPolicy{}},
	{Hint: nft.MintItemHint, Instance: nft.MintItem{}},
	{Hint: nft.MintHint, Instance: nft.Mint{}},
	{Hint: nft.TransferItemHint, Instance: nft.TransferItem{}},
	{Hint: nft.TransferHint, Instance: nft.Transfer{}},
	{Hint: nft.DelegateItemHint, Instance: nft.DelegateItem{}},
	{Hint: nft.DelegateHint, Instance: nft.Delegate{}},
	{Hint: nft.ApproveItemHint, Instance: nft.ApproveItem{}},
	{Hint: nft.ApproveHint, Instance: nft.Approve{}},
	{Hint: nft.SignItemHint, Instance: nft.SignItem{}},
	{Hint: nft.SignHint, Instance: nft.Sign{}},

	{Hint: state.LastNFTIndexStateValueHint, Instance: state.LastNFTIndexStateValue{}},
	{Hint: state.NFTStateValueHint, Instance: state.NFTStateValue{}},
	{Hint: state.NFTBoxStateValueHint, Instance: state.NFTBoxStateValue{}},
	{Hint: state.OperatorsBookStateValueHint, Instance: state.OperatorsBookStateValue{}},
	{Hint: state.CollectionStateValueHint, Instance: state.CollectionStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: nft.CreateCollectionFactHint, Instance: nft.CreateCollectionFact{}},
	{Hint: nft.UpdateCollectionPolicyFactHint, Instance: nft.UpdateCollectionPolicyFact{}},
	{Hint: nft.MintFactHint, Instance: nft.MintFact{}},
	{Hint: nft.TransferFactHint, Instance: nft.TransferFact{}},
	{Hint: nft.DelegateFactHint, Instance: nft.DelegateFact{}},
	{Hint: nft.ApproveFactHint, Instance: nft.ApproveFact{}},
	{Hint: nft.SignFactHint, Instance: nft.SignFact{}},
}

func init() {
	defaultLen := len(launch.Hinters)
	currencyExtendedLen := defaultLen + len(currencycmds.AddedHinters)
	allExtendedLen := currencyExtendedLen + len(AddedHinters)

	Hinters = make([]encoder.DecodeDetail, allExtendedLen)
	copy(Hinters, launch.Hinters)
	copy(Hinters[defaultLen:currencyExtendedLen], currencycmds.AddedHinters)
	copy(Hinters[currencyExtendedLen:], AddedHinters)

	defaultSupportedLen := len(launch.SupportedProposalOperationFactHinters)
	currencySupportedExtendedLen := defaultSupportedLen + len(currencycmds.AddedSupportedHinters)
	allSupportedExtendedLen := currencySupportedExtendedLen + len(AddedSupportedHinters)

	SupportedProposalOperationFactHinters = make(
		[]encoder.DecodeDetail,
		allSupportedExtendedLen)
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[defaultSupportedLen:currencySupportedExtendedLen], currencycmds.AddedSupportedHinters)
	copy(SupportedProposalOperationFactHinters[currencySupportedExtendedLen:], AddedSupportedHinters)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
