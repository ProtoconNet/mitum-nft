package cmds

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	operationisaac "github.com/ProtoconNet/mitum-currency/v3/operation/isaac"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-nft/v2/digest"
	digestisaac "github.com/ProtoconNet/mitum-nft/v2/digest/isaac"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/state"
	"github.com/ProtoconNet/mitum-nft/v2/types"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: common.BaseStateHint, Instance: common.BaseState{}},
	{Hint: common.NodeHint, Instance: common.BaseNode{}},
	{Hint: currencytypes.AddressHint, Instance: currencytypes.Address{}},
	{Hint: currencytypes.AmountHint, Instance: currencytypes.Amount{}},
	{Hint: currencytypes.AccountHint, Instance: currencytypes.Account{}},
	{Hint: statecurrency.AccountStateValueHint, Instance: statecurrency.AccountStateValue{}},
	{Hint: statecurrency.BalanceStateValueHint, Instance: statecurrency.BalanceStateValue{}},
	{Hint: currencytypes.AccountKeysHint, Instance: currencytypes.BaseAccountKeys{}},
	{Hint: currencytypes.AccountKeyHint, Instance: currencytypes.BaseAccountKey{}},
	{Hint: currency.CreateAccountsItemMultiAmountsHint, Instance: currency.CreateAccountsItemMultiAmounts{}},
	{Hint: currency.CreateAccountsItemSingleAmountHint, Instance: currency.CreateAccountsItemSingleAmount{}},
	{Hint: currency.CreateAccountsHint, Instance: currency.CreateAccounts{}},
	{Hint: currency.KeyUpdaterHint, Instance: currency.KeyUpdater{}},
	{Hint: currency.TransfersItemMultiAmountsHint, Instance: currency.TransfersItemMultiAmounts{}},
	{Hint: currency.TransfersItemSingleAmountHint, Instance: currency.TransfersItemSingleAmount{}},
	{Hint: currency.TransfersHint, Instance: currency.Transfers{}},
	{Hint: currency.SuffrageInflationHint, Instance: currency.SuffrageInflation{}},

	{Hint: currencytypes.NilFeeerHint, Instance: currencytypes.NilFeeer{}},
	{Hint: currencytypes.FixedFeeerHint, Instance: currencytypes.FixedFeeer{}},
	{Hint: currencytypes.RatioFeeerHint, Instance: currencytypes.RatioFeeer{}},
	{Hint: currencytypes.CurrencyPolicyHint, Instance: currencytypes.CurrencyPolicy{}},
	{Hint: currencytypes.CurrencyDesignHint, Instance: currencytypes.CurrencyDesign{}},
	{Hint: statecurrency.CurrencyDesignStateValueHint, Instance: statecurrency.CurrencyDesignStateValue{}},
	{Hint: currency.CurrencyRegisterHint, Instance: currency.CurrencyRegister{}},
	{Hint: currency.CurrencyPolicyUpdaterHint, Instance: currency.CurrencyPolicyUpdater{}},
	{Hint: currencytypes.ContractAccountKeysHint, Instance: currencytypes.ContractAccountKeys{}},
	{Hint: stateextension.ContractAccountStateValueHint, Instance: stateextension.ContractAccountStateValue{}},
	{Hint: extension.CreateContractAccountsItemMultiAmountsHint, Instance: extension.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extension.CreateContractAccountsItemSingleAmountHint, Instance: extension.CreateContractAccountsItemSingleAmount{}},
	{Hint: extension.CreateContractAccountsHint, Instance: extension.CreateContractAccounts{}},
	{Hint: extension.WithdrawsItemMultiAmountsHint, Instance: extension.WithdrawsItemMultiAmounts{}},
	{Hint: extension.WithdrawsItemSingleAmountHint, Instance: extension.WithdrawsItemSingleAmount{}},
	{Hint: extension.WithdrawsHint, Instance: extension.Withdraws{}},
	{Hint: currency.GenesisCurrenciesHint, Instance: currency.GenesisCurrencies{}},
	{Hint: currency.GenesisCurrenciesFactHint, Instance: currency.GenesisCurrenciesFact{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},

	{Hint: operationisaac.NetworkPolicyHint, Instance: operationisaac.NetworkPolicy{}},
	{Hint: operationisaac.NetworkPolicyStateValueHint, Instance: operationisaac.NetworkPolicyStateValue{}},
	{Hint: operationisaac.GenesisNetworkPolicyHint, Instance: operationisaac.GenesisNetworkPolicy{}},
	{Hint: operationisaac.SuffrageGenesisJoinHint, Instance: operationisaac.SuffrageGenesisJoin{}},
	{Hint: operationisaac.SuffrageCandidateHint, Instance: operationisaac.SuffrageCandidate{}},
	{Hint: operationisaac.SuffrageJoinHint, Instance: operationisaac.SuffrageJoin{}},
	{Hint: operationisaac.SuffrageDisjoinHint, Instance: operationisaac.SuffrageDisjoin{}},
	{Hint: operationisaac.FixedSuffrageCandidateLimiterRuleHint, Instance: operationisaac.FixedSuffrageCandidateLimiterRule{}},
	{Hint: operationisaac.MajoritySuffrageCandidateLimiterRuleHint, Instance: operationisaac.MajoritySuffrageCandidateLimiterRule{}},

	{Hint: types.SignerHint, Instance: types.Signer{}},
	{Hint: types.SignersHint, Instance: types.Signers{}},
	{Hint: types.NFTHint, Instance: types.NFT{}},
	{Hint: types.DesignHint, Instance: types.Design{}},

	{Hint: state.LastNFTIndexStateValueHint, Instance: state.LastNFTIndexStateValue{}},
	{Hint: state.NFTStateValueHint, Instance: state.NFTStateValue{}},
	{Hint: state.NFTBoxStateValueHint, Instance: state.NFTBoxStateValue{}},
	{Hint: types.NFTBoxHint, Instance: types.NFTBox{}},
	{Hint: state.OperatorsBookStateValueHint, Instance: state.OperatorsBookStateValue{}},
	{Hint: types.OperatorsBookHint, Instance: types.OperatorsBook{}},
	{Hint: types.CollectionPolicyHint, Instance: types.CollectionPolicy{}},
	{Hint: types.CollectionDesignHint, Instance: types.CollectionDesign{}},
	{Hint: state.CollectionStateValueHint, Instance: state.CollectionStateValue{}},
	{Hint: nft.CollectionRegisterHint, Instance: nft.CollectionRegister{}},
	{Hint: nft.CollectionPolicyUpdaterHint, Instance: nft.CollectionPolicyUpdater{}},
	{Hint: nft.MintItemHint, Instance: nft.MintItem{}},
	{Hint: nft.MintHint, Instance: nft.Mint{}},
	{Hint: nft.NFTTransferItemHint, Instance: nft.NFTTransferItem{}},
	{Hint: nft.NFTTransferHint, Instance: nft.NFTTransfer{}},
	{Hint: nft.DelegateItemHint, Instance: nft.DelegateItem{}},
	{Hint: nft.DelegateHint, Instance: nft.Delegate{}},
	{Hint: nft.ApproveItemHint, Instance: nft.ApproveItem{}},
	{Hint: nft.ApproveHint, Instance: nft.Approve{}},
	{Hint: nft.NFTSignItemHint, Instance: nft.NFTSignItem{}},
	{Hint: nft.NFTSignHint, Instance: nft.NFTSign{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: currency.CreateAccountsFactHint, Instance: currency.CreateAccountsFact{}},
	{Hint: currency.KeyUpdaterFactHint, Instance: currency.KeyUpdaterFact{}},
	{Hint: currency.TransfersFactHint, Instance: currency.TransfersFact{}},
	{Hint: currency.SuffrageInflationFactHint, Instance: currency.SuffrageInflationFact{}},

	{Hint: currency.CurrencyRegisterFactHint, Instance: currency.CurrencyRegisterFact{}},
	{Hint: currency.CurrencyPolicyUpdaterFactHint, Instance: currency.CurrencyPolicyUpdaterFact{}},
	{Hint: extension.CreateContractAccountsFactHint, Instance: extension.CreateContractAccountsFact{}},
	{Hint: extension.WithdrawsFactHint, Instance: extension.WithdrawsFact{}},

	{Hint: operationisaac.GenesisNetworkPolicyFactHint, Instance: operationisaac.GenesisNetworkPolicyFact{}},
	{Hint: operationisaac.SuffrageGenesisJoinFactHint, Instance: operationisaac.SuffrageGenesisJoinFact{}},
	{Hint: operationisaac.SuffrageCandidateFactHint, Instance: operationisaac.SuffrageCandidateFact{}},
	{Hint: operationisaac.SuffrageJoinFactHint, Instance: operationisaac.SuffrageJoinFact{}},
	{Hint: operationisaac.SuffrageDisjoinFactHint, Instance: operationisaac.SuffrageDisjoinFact{}},

	{Hint: nft.CollectionRegisterFactHint, Instance: nft.CollectionRegisterFact{}},
	{Hint: nft.CollectionPolicyUpdaterFactHint, Instance: nft.CollectionPolicyUpdaterFact{}},
	{Hint: nft.MintFactHint, Instance: nft.MintFact{}},
	{Hint: nft.NFTTransferFactHint, Instance: nft.NFTTransferFact{}},
	{Hint: nft.DelegateFactHint, Instance: nft.DelegateFact{}},
	{Hint: nft.ApproveFactHint, Instance: nft.ApproveFact{}},
	{Hint: nft.NFTSignFactHint, Instance: nft.NFTSignFact{}},
}

func init() {
	Hinters = make([]encoder.DecodeDetail, len(launch.Hinters)+len(hinters))
	copy(Hinters, launch.Hinters)
	copy(Hinters[len(launch.Hinters):], hinters)

	SupportedProposalOperationFactHinters = make([]encoder.DecodeDetail, len(launch.SupportedProposalOperationFactHinters)+len(supportedProposalOperationFactHinters))
	copy(SupportedProposalOperationFactHinters, launch.SupportedProposalOperationFactHinters)
	copy(SupportedProposalOperationFactHinters[len(launch.SupportedProposalOperationFactHinters):], supportedProposalOperationFactHinters)
}

func LoadHinters(enc encoder.Encoder) error {
	for _, hinter := range Hinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	for _, hinter := range SupportedProposalOperationFactHinters {
		if err := enc.Add(hinter); err != nil {
			return errors.Wrap(err, "failed to add to encoder")
		}
	}

	return nil
}
