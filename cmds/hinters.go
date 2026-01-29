package cmds

import (
	ccmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-nft/operation/nft"
	"github.com/ProtoconNet/mitum-nft/state"
	"github.com/ProtoconNet/mitum-nft/types"
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
	{Hint: types.AllApprovedBookHint, Instance: types.AllApprovedBook{}},
	{Hint: types.CollectionPolicyHint, Instance: types.CollectionPolicy{}},
	{Hint: types.CollectionDesignHint, Instance: types.CollectionDesign{}},

	{Hint: nft.RegisterModelHint, Instance: nft.RegisterModel{}},
	{Hint: nft.UpdateModelConfigHint, Instance: nft.UpdateModelConfig{}},
	{Hint: nft.MintItemHint, Instance: nft.MintItem{}},
	{Hint: nft.MintHint, Instance: nft.Mint{}},
	{Hint: nft.TransferItemHint, Instance: nft.TransferItem{}},
	{Hint: nft.TransferHint, Instance: nft.Transfer{}},
	{Hint: nft.ApproveAllItemHint, Instance: nft.ApproveAllItem{}},
	{Hint: nft.ApproveAllHint, Instance: nft.ApproveAll{}},
	{Hint: nft.ApproveItemHint, Instance: nft.ApproveItem{}},
	{Hint: nft.ApproveHint, Instance: nft.Approve{}},
	{Hint: nft.AddSignatureItemHint, Instance: nft.AddSignatureItem{}},
	{Hint: nft.AddSignatureHint, Instance: nft.AddSignature{}},

	{Hint: state.LastNFTIndexStateValueHint, Instance: state.LastNFTIndexStateValue{}},
	{Hint: state.NFTStateValueHint, Instance: state.NFTStateValue{}},
	{Hint: state.OperatorsBookStateValueHint, Instance: state.OperatorsBookStateValue{}},
	{Hint: state.CollectionStateValueHint, Instance: state.CollectionStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: nft.RegisterModelFactHint, Instance: nft.RegisterModelFact{}},
	{Hint: nft.UpdateModelConfigFactHint, Instance: nft.UpdateModelConfigFact{}},
	{Hint: nft.MintFactHint, Instance: nft.MintFact{}},
	{Hint: nft.TransferFactHint, Instance: nft.TransferFact{}},
	{Hint: nft.ApproveAllFactHint, Instance: nft.ApproveAllFact{}},
	{Hint: nft.ApproveFactHint, Instance: nft.ApproveFact{}},
	{Hint: nft.AddSignatureFactHint, Instance: nft.AddSignatureFact{}},
}

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
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
