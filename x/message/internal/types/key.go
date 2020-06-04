package types

const (
	// module name
	ModuleName = "message"
	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey    = "transient_" + ModuleName
	QuerierRoute = ModuleName
)
