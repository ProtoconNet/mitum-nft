package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/v2/currency"
	"github.com/ProtoconNet/mitum-currency-extension/v2/digest"
	"github.com/ProtoconNet/mitum-currency/v2/currency"
	digestisaac "github.com/ProtoconNet/mitum-currency/v2/digest/isaac"
	isaacoperation "github.com/ProtoconNet/mitum-currency/v2/isaac"
	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var hinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit
	{Hint: currency.BaseStateHint, Instance: currency.BaseState{}},
	{Hint: currency.NodeHint, Instance: currency.BaseNode{}},
	{Hint: currency.AddressHint, Instance: currency.Address{}},
	{Hint: currency.AmountHint, Instance: currency.Amount{}},
	{Hint: currency.AccountHint, Instance: currency.Account{}},
	{Hint: currency.AccountStateValueHint, Instance: currency.AccountStateValue{}},
	{Hint: currency.BalanceStateValueHint, Instance: currency.BalanceStateValue{}},
	{Hint: currency.AccountKeysHint, Instance: currency.BaseAccountKeys{}},
	{Hint: currency.AccountKeyHint, Instance: currency.BaseAccountKey{}},
	{Hint: currency.CreateAccountsItemMultiAmountsHint, Instance: currency.CreateAccountsItemMultiAmounts{}},
	{Hint: currency.CreateAccountsItemSingleAmountHint, Instance: currency.CreateAccountsItemSingleAmount{}},
	{Hint: currency.CreateAccountsHint, Instance: currency.CreateAccounts{}},
	{Hint: currency.KeyUpdaterHint, Instance: currency.KeyUpdater{}},
	{Hint: currency.TransfersItemMultiAmountsHint, Instance: currency.TransfersItemMultiAmounts{}},
	{Hint: currency.TransfersItemSingleAmountHint, Instance: currency.TransfersItemSingleAmount{}},
	{Hint: currency.TransfersHint, Instance: currency.Transfers{}},
	{Hint: currency.SuffrageInflationHint, Instance: currency.SuffrageInflation{}},

	{Hint: extensioncurrency.NilFeeerHint, Instance: extensioncurrency.NilFeeer{}},
	{Hint: extensioncurrency.FixedFeeerHint, Instance: extensioncurrency.FixedFeeer{}},
	{Hint: extensioncurrency.RatioFeeerHint, Instance: extensioncurrency.RatioFeeer{}},
	{Hint: extensioncurrency.CurrencyPolicyHint, Instance: extensioncurrency.CurrencyPolicy{}},
	{Hint: extensioncurrency.CurrencyDesignHint, Instance: extensioncurrency.CurrencyDesign{}},
	{Hint: extensioncurrency.CurrencyDesignStateValueHint, Instance: extensioncurrency.CurrencyDesignStateValue{}},
	{Hint: extensioncurrency.CurrencyRegisterHint, Instance: extensioncurrency.CurrencyRegister{}},
	{Hint: extensioncurrency.CurrencyPolicyUpdaterHint, Instance: extensioncurrency.CurrencyPolicyUpdater{}},
	{Hint: extensioncurrency.ContractAccountKeysHint, Instance: extensioncurrency.ContractAccountKeys{}},
	{Hint: extensioncurrency.ContractAccountStateValueHint, Instance: extensioncurrency.ContractAccountStateValue{}},
	{Hint: extensioncurrency.CreateContractAccountsItemMultiAmountsHint, Instance: extensioncurrency.CreateContractAccountsItemMultiAmounts{}},
	{Hint: extensioncurrency.CreateContractAccountsItemSingleAmountHint, Instance: extensioncurrency.CreateContractAccountsItemSingleAmount{}},
	{Hint: extensioncurrency.CreateContractAccountsHint, Instance: extensioncurrency.CreateContractAccounts{}},
	{Hint: extensioncurrency.WithdrawsItemMultiAmountsHint, Instance: extensioncurrency.WithdrawsItemMultiAmounts{}},
	{Hint: extensioncurrency.WithdrawsItemSingleAmountHint, Instance: extensioncurrency.WithdrawsItemSingleAmount{}},
	{Hint: extensioncurrency.WithdrawsHint, Instance: extensioncurrency.Withdraws{}},
	{Hint: extensioncurrency.GenesisCurrenciesHint, Instance: extensioncurrency.GenesisCurrencies{}},
	{Hint: extensioncurrency.GenesisCurrenciesFactHint, Instance: extensioncurrency.GenesisCurrenciesFact{}},

	{Hint: digestisaac.ManifestHint, Instance: digestisaac.Manifest{}},
	{Hint: digest.AccountValueHint, Instance: digest.AccountValue{}},
	{Hint: digest.OperationValueHint, Instance: digest.OperationValue{}},

	{Hint: isaacoperation.NetworkPolicyHint, Instance: isaacoperation.NetworkPolicy{}},
	{Hint: isaacoperation.NetworkPolicyStateValueHint, Instance: isaacoperation.NetworkPolicyStateValue{}},
	{Hint: isaacoperation.GenesisNetworkPolicyHint, Instance: isaacoperation.GenesisNetworkPolicy{}},
	{Hint: isaacoperation.SuffrageGenesisJoinHint, Instance: isaacoperation.SuffrageGenesisJoin{}},
	{Hint: isaacoperation.SuffrageCandidateHint, Instance: isaacoperation.SuffrageCandidate{}},
	{Hint: isaacoperation.SuffrageJoinHint, Instance: isaacoperation.SuffrageJoin{}},
	{Hint: isaacoperation.SuffrageDisjoinHint, Instance: isaacoperation.SuffrageDisjoin{}},
	{Hint: isaacoperation.FixedSuffrageCandidateLimiterRuleHint, Instance: isaacoperation.FixedSuffrageCandidateLimiterRule{}},
	{Hint: isaacoperation.MajoritySuffrageCandidateLimiterRuleHint, Instance: isaacoperation.MajoritySuffrageCandidateLimiterRule{}},

	{Hint: nft.SignerHint, Instance: nft.Signer{}},
	{Hint: nft.SignersHint, Instance: nft.Signers{}},
	{Hint: nft.NFTIDHint, Instance: nft.NFTID{}},
	{Hint: nft.NFTHint, Instance: nft.NFT{}},
	{Hint: nft.DesignHint, Instance: nft.Design{}},

	{Hint: collection.LastNFTIndexStateValueHint, Instance: collection.LastNFTIndexStateValue{}},
	{Hint: collection.NFTStateValueHint, Instance: collection.NFTStateValue{}},
	{Hint: collection.NFTBoxStateValueHint, Instance: collection.NFTBoxStateValue{}},
	{Hint: collection.NFTBoxHint, Instance: collection.NFTBox{}},
	{Hint: collection.OperatorsBookStateValueHint, Instance: collection.OperatorsBookStateValue{}},
	{Hint: collection.OperatorsBookHint, Instance: collection.OperatorsBook{}},
	{Hint: collection.CollectionPolicyHint, Instance: collection.CollectionPolicy{}},
	{Hint: collection.CollectionDesignHint, Instance: collection.CollectionDesign{}},
	{Hint: collection.CollectionStateValueHint, Instance: collection.CollectionStateValue{}},
	{Hint: collection.CollectionRegisterFormHint, Instance: collection.CollectionRegisterForm{}},
	{Hint: collection.CollectionRegisterHint, Instance: collection.CollectionRegister{}},
	{Hint: collection.CollectionPolicyUpdaterHint, Instance: collection.CollectionPolicyUpdater{}},
	{Hint: collection.MintFormHint, Instance: collection.MintForm{}},
	{Hint: collection.MintItemHint, Instance: collection.MintItem{}},
	{Hint: collection.MintHint, Instance: collection.Mint{}},
	{Hint: collection.NFTTransferItemHint, Instance: collection.NFTTransferItem{}},
	{Hint: collection.NFTTransferHint, Instance: collection.NFTTransfer{}},
	{Hint: collection.DelegateItemHint, Instance: collection.DelegateItem{}},
	{Hint: collection.DelegateHint, Instance: collection.Delegate{}},
	{Hint: collection.ApproveItemHint, Instance: collection.ApproveItem{}},
	{Hint: collection.ApproveHint, Instance: collection.Approve{}},
	{Hint: collection.NFTSignItemHint, Instance: collection.NFTSignItem{}},
	{Hint: collection.NFTSignHint, Instance: collection.NFTSign{}},
}

var supportedProposalOperationFactHinters = []encoder.DecodeDetail{
	{Hint: currency.CreateAccountsFactHint, Instance: currency.CreateAccountsFact{}},
	{Hint: currency.KeyUpdaterFactHint, Instance: currency.KeyUpdaterFact{}},
	{Hint: currency.TransfersFactHint, Instance: currency.TransfersFact{}},
	{Hint: currency.SuffrageInflationFactHint, Instance: currency.SuffrageInflationFact{}},

	{Hint: extensioncurrency.CurrencyRegisterFactHint, Instance: extensioncurrency.CurrencyRegisterFact{}},
	{Hint: extensioncurrency.CurrencyPolicyUpdaterFactHint, Instance: extensioncurrency.CurrencyPolicyUpdaterFact{}},
	{Hint: extensioncurrency.CreateContractAccountsFactHint, Instance: extensioncurrency.CreateContractAccountsFact{}},
	{Hint: extensioncurrency.WithdrawsFactHint, Instance: extensioncurrency.WithdrawsFact{}},

	{Hint: isaacoperation.GenesisNetworkPolicyFactHint, Instance: isaacoperation.GenesisNetworkPolicyFact{}},
	{Hint: isaacoperation.SuffrageGenesisJoinFactHint, Instance: isaacoperation.SuffrageGenesisJoinFact{}},
	{Hint: isaacoperation.SuffrageCandidateFactHint, Instance: isaacoperation.SuffrageCandidateFact{}},
	{Hint: isaacoperation.SuffrageJoinFactHint, Instance: isaacoperation.SuffrageJoinFact{}},
	{Hint: isaacoperation.SuffrageDisjoinFactHint, Instance: isaacoperation.SuffrageDisjoinFact{}},

	{Hint: collection.CollectionRegisterFactHint, Instance: collection.CollectionRegisterFact{}},
	{Hint: collection.CollectionPolicyUpdaterFactHint, Instance: collection.CollectionPolicyUpdaterFact{}},
	{Hint: collection.MintFactHint, Instance: collection.MintFact{}},
	{Hint: collection.NFTTransferFactHint, Instance: collection.NFTTransferFact{}},
	{Hint: collection.DelegateFactHint, Instance: collection.DelegateFact{}},
	{Hint: collection.ApproveFactHint, Instance: collection.ApproveFact{}},
	{Hint: collection.NFTSignFactHint, Instance: collection.NFTSignFact{}},
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
