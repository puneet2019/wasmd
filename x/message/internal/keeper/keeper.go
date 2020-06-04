package keeper

import (
	"github.com/CosmWasm/wasmd/x/message/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	MessageKeeper
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseMessageKeeper
	ak         types.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, blacklistedAddrs map[string]bool,
) BaseKeeper {

	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return BaseKeeper{
		BaseMessageKeeper: NewBaseMessageKeeper(cdc, key, ps, ak, blacklistedAddrs),
		ak:                ak,
		paramSpace:        ps,
	}
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type MessageKeeper interface {
	SendTextMessage(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, text string) error

	SetMessage(ctx sdk.Context, addr sdk.AccAddress, toAddr sdk.AccAddress, text string) error
	GetMessages(ctx sdk.Context, addr sdk.AccAddress) types.MsgTextsSend
	GetSendEnabled(ctx sdk.Context) bool
	SetSendEnabled(ctx sdk.Context, enabled bool)
	BlacklistedAddr(addr sdk.AccAddress) bool
}

var _ MessageKeeper = (*BaseMessageKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseMessageKeeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	ak         types.AccountKeeper
	paramSpace params.Subspace

	// list of addresses that are restricted from receiving transactions
	blacklistedAddrs map[string]bool
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseMessageKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	ak types.AccountKeeper, blacklistedAddrs map[string]bool,
) BaseMessageKeeper {

	return BaseMessageKeeper{
		storeKey:         key,
		cdc:              cdc,
		ak:               ak,
		paramSpace:       paramSpace,
		blacklistedAddrs: blacklistedAddrs,
	}
}

// SendCoins moves coins from one account to another
func (keeper BaseMessageKeeper) SendTextMessage(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, text string) error {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeText,
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute("Text", text),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
		),
	})

	err := keeper.SetMessage(ctx, fromAddr, toAddr, text)
	if err != nil {
		return err
	}

	return nil
}

//GetMessages gets messages sent by user
func (keeper BaseMessageKeeper) GetMessages(ctx sdk.Context, addr sdk.AccAddress) types.MsgTextsSend {
	kvstore := ctx.KVStore(keeper.storeKey)

	existingMsgs := kvstore.Get([]byte(addr))
	var messages types.MsgTextsSend
	err := keeper.cdc.UnmarshalJSON(existingMsgs, messages)
	if err != nil {
		messages = types.MsgTextsSend{}
	}
	return messages
}

// SetCoins sets the coins at the addr.
func (keeper BaseMessageKeeper) SetMessage(ctx sdk.Context, addr sdk.AccAddress, toAddr sdk.AccAddress, text string) error {
	if text == "" {
		sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, text)
	}
	//for from addr
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	keeper.ak.SetAccount(ctx, acc)

	kvstore := ctx.KVStore(keeper.storeKey)

	existingMsgs := kvstore.Get([]byte(addr))
	var messages types.MsgTextsSend
	err := keeper.cdc.UnmarshalJSON(existingMsgs, messages)
	if err != nil {
		messages = types.MsgTextsSend{}
	}

	addmsg, err := keeper.cdc.MarshalJSON(types.MsgTextsSend{append(messages.MsgTexts, types.NewMsgTextSend(addr, toAddr, text))})
	if err != nil {
		return err
	}

	kvstore.Set([]byte(addr), addmsg)

	//for to addr
	acc = keeper.ak.GetAccount(ctx, toAddr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, toAddr)
	}
	keeper.ak.SetAccount(ctx, acc)

	existingMsgs = kvstore.Get([]byte(toAddr))
	var messagesto types.MsgTextsSend
	err = keeper.cdc.UnmarshalJSON(existingMsgs, messagesto)
	if err != nil {
		messages = types.MsgTextsSend{}
	}

	addmsg, err = keeper.cdc.MarshalJSON(types.MsgTextsSend{append(messagesto.MsgTexts, types.NewMsgTextSend(addr, toAddr, text))})
	if err != nil {
		return err
	}

	kvstore.Set([]byte(toAddr), addmsg)

	return nil
}

// GetSendEnabled returns the current SendEnabled
func (keeper BaseMessageKeeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseMessageKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

// BlacklistedAddr checks if a given address is blacklisted (i.e restricted from
// receiving funds)
func (keeper BaseMessageKeeper) BlacklistedAddr(addr sdk.AccAddress) bool {
	return keeper.blacklistedAddrs[addr.String()]
}
