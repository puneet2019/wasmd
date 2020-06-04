package message

// nolint

import (
	"github.com/CosmWasm/wasmd/x/message/internal/keeper"
	"github.com/CosmWasm/wasmd/x/message/internal/types"
)

const (
	QueryBalance       = keeper.QueryBalance
	ModuleName         = types.ModuleName
	StoreKey           = types.StoreKey
	QuerierRoute       = types.QuerierRoute
	RouterKey          = types.RouterKey
	DefaultParamspace  = types.DefaultParamspace
	DefaultSendEnabled = types.DefaultSendEnabled

	EventTypeText          = types.EventTypeText
	AttributeKeyRecipient  = types.AttributeKeyRecipient
	AttributeKeySender     = types.AttributeKeySender
	AttributeValueCategory = types.AttributeValueCategory
)

var (
	// RegisterInvariants          = keeper.RegisterInvariants
	// NonnegativeBalanceInvariant = keeper.NonnegativeBalanceInvariant
	NewBaseKeeper        = keeper.NewBaseKeeper
	NewBaseMessageKeeper = keeper.NewBaseMessageKeeper
	// NewBaseViewKeeper           = keeper.NewBaseViewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	// ErrNoInputs                 = types.ErrNoInputs
	// ErrNoOutputs                = types.ErrNoOutputs
	// ErrInputOutputMismatch      = types.ErrInputOutputMismatch
	ErrSendDisabled     = types.ErrSendDisabled
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewMsgTextSend      = types.NewMsgTextSend
	// NewMsgMultiSend             = types.NewMsgMultiSend
	// NewInput                    = types.NewInput
	// NewOutput                   = types.NewOutput
	// ValidateInputsOutputs       = types.ValidateInputsOutputs
	ParamKeyTable            = types.ParamKeyTable
	NewQueryBalanceParams    = types.NewQueryBalanceParams
	ModuleCdc                = types.ModuleCdc
	ParamStoreKeySendEnabled = types.ParamStoreKeySendEnabled
)

type (
	Keeper            = keeper.Keeper
	BaseKeeper        = keeper.BaseKeeper
	MessageKeeper     = keeper.MessageKeeper
	BaseMessageKeeper = keeper.BaseMessageKeeper
	// ViewKeeper     = keeper.ViewKeeper
	// BaseViewKeeper = keeper.BaseViewKeeper
	GenesisState = types.GenesisState
	MsgTextSend  = types.MsgTextSend
	MsgTextsSend = types.MsgTextsSend
	// MsgMultiSend       = types.MsgMultiSend
	// Input              = types.Input
	// Output             = types.Output
	QueryBalanceParams = types.QueryBalanceParams
)
