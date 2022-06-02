package cmds

import (
	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/ProtoconNet/mitum-nft/nft"
	"github.com/ProtoconNet/mitum-nft/nft/collection"

	"github.com/ProtoconNet/mitum-nft/digest"
	"github.com/spikeekips/mitum-currency/currency"
)

var (
	Hinters []hint.Hinter
	Types   []hint.Type
)

var types = []hint.Type{
	currency.AccountType,
	currency.AddressType,
	currency.AmountType,
	currency.CreateAccountsFactType,
	currency.CreateAccountsItemMultiAmountsType,
	currency.CreateAccountsItemSingleAmountType,
	currency.CreateAccountsType,
	currency.AccountKeyType,
	currency.KeyUpdaterFactType,
	currency.KeyUpdaterType,
	currency.AccountKeysType,
	currency.TransfersFactType,
	currency.TransfersItemMultiAmountsType,
	currency.TransfersItemSingleAmountType,
	currency.TransfersType,
	extensioncurrency.CurrencyDesignType,
	extensioncurrency.CurrencyPolicyType,
	extensioncurrency.CurrencyPolicyUpdaterFactType,
	extensioncurrency.CurrencyPolicyUpdaterType,
	extensioncurrency.CurrencyRegisterFactType,
	extensioncurrency.CurrencyRegisterType,
	extensioncurrency.FeeOperationFactType,
	extensioncurrency.FeeOperationType,
	extensioncurrency.FixedFeeerType,
	extensioncurrency.GenesisCurrenciesFactType,
	extensioncurrency.GenesisCurrenciesType,
	extensioncurrency.NilFeeerType,
	extensioncurrency.RatioFeeerType,
	extensioncurrency.SuffrageInflationFactType,
	extensioncurrency.SuffrageInflationType,
	extensioncurrency.ContractAccountKeysType,
	extensioncurrency.ContractAccountType,
	extensioncurrency.CreateContractAccountsFactType,
	extensioncurrency.CreateContractAccountsType,
	extensioncurrency.CreateContractAccountsItemMultiAmountsType,
	extensioncurrency.CreateContractAccountsItemSingleAmountType,
	extensioncurrency.DeactivateFactType,
	extensioncurrency.DeactivateType,
	extensioncurrency.WithdrawsFactType,
	extensioncurrency.WithdrawsType,
	extensioncurrency.WithdrawsItemMultiAmountsType,
	extensioncurrency.WithdrawsItemSingleAmountType,
	nft.NFTIDType,
	nft.NFTType,
	nft.DesignType,
	collection.PolicyType,
	collection.MintFormType,
	collection.DelegateFactType,
	collection.DelegateType,
	collection.DelegateItemType,
	collection.ApproveFactType,
	collection.ApproveType,
	collection.ApproveItemMultiNFTsType,
	collection.ApproveItemSingleNFTType,
	collection.CollectionRegisterFactType,
	collection.CollectionRegisterType,
	collection.MintFactType,
	collection.MintType,
	collection.TransferFactType,
	collection.TransferType,
	collection.TransferItemMultiNFTsType,
	collection.TransferItemSingleNFTType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	digest.OperationValueType,
}

var hinters = []hint.Hinter{
	currency.AccountHinter,
	currency.AddressHinter,
	currency.AmountHinter,
	currency.CreateAccountsFactHinter,
	currency.CreateAccountsItemMultiAmountsHinter,
	currency.CreateAccountsItemSingleAmountHinter,
	currency.CreateAccountsHinter,
	currency.KeyUpdaterFactHinter,
	currency.KeyUpdaterHinter,
	currency.AccountKeysHinter,
	currency.AccountKeyHinter,
	currency.TransfersFactHinter,
	currency.TransfersItemMultiAmountsHinter,
	currency.TransfersItemSingleAmountHinter,
	currency.TransfersHinter,
	extensioncurrency.CurrencyDesignHinter,
	extensioncurrency.CurrencyPolicyUpdaterFactHinter,
	extensioncurrency.CurrencyPolicyUpdaterHinter,
	extensioncurrency.CurrencyPolicyHinter,
	extensioncurrency.CurrencyRegisterFactHinter,
	extensioncurrency.CurrencyRegisterHinter,
	extensioncurrency.FeeOperationFactHinter,
	extensioncurrency.FeeOperationHinter,
	extensioncurrency.FixedFeeerHinter,
	extensioncurrency.GenesisCurrenciesFactHinter,
	extensioncurrency.GenesisCurrenciesHinter,
	extensioncurrency.NilFeeerHinter,
	extensioncurrency.RatioFeeerHinter,
	extensioncurrency.SuffrageInflationFactHinter,
	extensioncurrency.SuffrageInflationHinter,
	extensioncurrency.ContractAccountKeysHinter,
	extensioncurrency.ContractAccountHinter,
	extensioncurrency.CreateContractAccountsFactHinter,
	extensioncurrency.CreateContractAccountsHinter,
	extensioncurrency.CreateContractAccountsItemMultiAmountsHinter,
	extensioncurrency.CreateContractAccountsItemSingleAmountHinter,
	extensioncurrency.DeactivateFactHinter,
	extensioncurrency.DeactivateHinter,
	extensioncurrency.WithdrawsFactHinter,
	extensioncurrency.WithdrawsHinter,
	extensioncurrency.WithdrawsItemMultiAmountsHinter,
	extensioncurrency.WithdrawsItemSingleAmountHinter,
	nft.NFTIDHinter,
	nft.NFTHinter,
	nft.DesignHinter,
	collection.PolicyHinter,
	collection.MintFormHinter,
	collection.DelegateFactHinter,
	collection.DelegateHinter,
	collection.DelegateItemHinter,
	collection.ApproveFactHinter,
	collection.ApproveHinter,
	collection.ApproveItemMultiNFTsHinter,
	collection.ApproveItemSingleNFTHinter,
	collection.CollectionRegisterFactHinter,
	collection.CollectionRegisterHinter,
	collection.MintFactHinter,
	collection.MintHinter,
	collection.TransferFactHinter,
	collection.TransferHinter,
	collection.TransferItemMultiNFTsHinter,
	collection.TransferItemSingleNFTHinter.BaseTransferItem,
	digest.AccountValue{},
	digest.BaseHal{},
	digest.NodeInfo{},
	digest.OperationValue{},
	digest.Problem{},
}

func init() {
	Hinters = make([]hint.Hinter, len(launch.EncoderHinters)+len(hinters))
	copy(Hinters, launch.EncoderHinters)
	copy(Hinters[len(launch.EncoderHinters):], hinters)

	Types = make([]hint.Type, len(launch.EncoderTypes)+len(types))
	copy(Types, launch.EncoderTypes)
	copy(Types[len(launch.EncoderTypes):], types)
}
