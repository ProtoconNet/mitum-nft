package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	"github.com/ProtoconNet/mitum-nft/v2/operation/nft"
	"github.com/ProtoconNet/mitum-nft/v2/operation/processor"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-nft-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *currencyprocessor.OperationProcessor
	var set *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		currencycmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &set,
	); err != nil {
		return pctx, err
	}

	//err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	//if err != nil {
	//	return pctx, err
	//}
	err := opr.SetGetNewProcessorFunc(processor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}
	if err := opr.SetProcessor(
		nft.CreateCollectionHint,
		nft.NewCreateCollectionProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.UpdateCollectionPolicyHint,
		nft.NewUpdateCollectionPolicyProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.MintHint,
		nft.NewMintProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.TransferHint,
		nft.NewTransferProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.DelegateHint,
		nft.NewDelegateProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.ApproveHint,
		nft.NewApproveProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		nft.SignHint,
		nft.NewSignProcessor(),
	); err != nil {
		return pctx, err
	}

	_ = set.Add(nft.CreateCollectionHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.UpdateCollectionPolicyHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.MintHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.TransferHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.DelegateHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.ApproveHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	_ = set.Add(nft.SignHint,
		func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
			return opr.New(
				height,
				getStatef,
				nil,
				nil,
			)
		})

	pctx = context.WithValue(pctx, currencycmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter

	return pctx, nil
}
