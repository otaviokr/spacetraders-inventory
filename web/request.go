package web

// Proxy represents the interface the Proxy to the game API needs to provide.
type Proxy interface {
	// GetUserDetails collects user details.
	//
	// https://api.spacetraders.io/#api-account-GetAccount
	GetUserDetails() ([]byte, error)

	// GetShipList collects the list of ships the user owns.
	//
	// https://api.spacetraders.io/#api-ships-GetShips
	GetShipList() ([]byte, error)

	// GetLeaderboard collects the leaderboard and the user rank.
	//
	// https://api.spacetraders.io/#api-leaderboard-netWorth
	GetLeaderboard() ([]byte, error)

	// GetGameStatus collects the game status (if it is online and accessible).
	//
	// https://api.spacetraders.io/#api-game-Status
	GetGameStatus() ([]byte, error)
}
