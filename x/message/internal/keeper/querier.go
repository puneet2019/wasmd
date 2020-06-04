package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/CosmWasm/wasmd/x/message/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// query balance path
	QueryMessage = "message"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryMessage:
			return queryMessages(ctx, req, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryMessages(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBalanceParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coins := k.GetMessages(ctx, params.Address)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
