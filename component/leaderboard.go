package component

// Leaderboard is the response from the API.
type Leaderboard struct {
	NetWorth     []UserRank `yaml:"netWorth"`
	UserNetWorth UserRank   `yaml:"userNetWorth"`
}

// UserRank is detail about a user in the leaderboard response from the API.
type UserRank struct {
	Username string `yaml:"username"`
	NetWorth int64  `yaml:"netWorth"`
	Rank     int64  `yaml:"rank"`
}
