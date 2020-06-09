package message

import (
	"github.com/CosmWasm/wasmd/x/message/internal/keeper"
	"github.com/CosmWasm/wasmd/x/message/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgTextSend:
			return handleMsgTextSend(ctx, k, msg)

		// case types.MsgMultiSend:
		// 	return handleMsgMultiSend(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized bank message type: %T", msg)
		}
	}
}

// Handle MsgSend.
func handleMsgTextSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgTextSend) (*sdk.Result, error) {
	if !k.GetSendEnabled(ctx) {
		return nil, types.ErrSendDisabled
	}

	if k.BlacklistedAddr(msg.ToAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", msg.ToAddress)
	}

	err := k.SendTextMessage(ctx, msg.FromAddress, msg.ToAddress, msg.Text)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// // Handle MsgMultiSend.
// func handleMsgMultiSend(ctx sdk.Context, k keeper.Keeper, msg types.MsgMultiSend) (*sdk.Result, error) {
// 	// NOTE: totalIn == totalOut should already have been checked
// 	if !k.GetSendEnabled(ctx) {
// 		return nil, types.ErrSendDisabled
// 	}

// 	for _, out := range msg.Outputs {
// 		if k.BlacklistedAddr(out.Address) {
// 			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", out.Address)
// 		}
// 	}

// 	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(
// 			sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
// 		),
// 	)

// 	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
// }
