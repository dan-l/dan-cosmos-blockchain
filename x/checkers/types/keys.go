package types

const (
	// ModuleName defines the module name
	ModuleName = "checkers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_checkers"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	SystemInfoKey = "SystemInfo-value-"
)

// GameCreated events
const (
	GameCreatedEventType      = "new-game-created"
	GameCreatedEventCreator   = "creator"
	GameCreatedEventGameIndex = "game-index"
	GameCreatedEventBlack     = "black"
	GameCreatedEventRed       = "red"
)

// PlayMove events
const (
	PlayMovedEventType      = "play-moved"
	PlayMovedEventCreator   = "creator"
	PlayMovedEventGameIndex = "game-index"
	PlayMovedEventCapturedX = "captured-x"
	PlayMovedEventCapturedY = "captured-y"
	PlayMovedEventWinner    = "winner"
)
