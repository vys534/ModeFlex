package api

// User represents an osu! user.
type User struct {
	Cover struct {
		URL string `json:"url"`
	} `json:"cover"`
	CountryCode string         `json:"country_code"`
	ID          int64          `json:"id"`
	Username    string         `json:"username"`
	Statistics  UserStatistics `json:"statistics"`
}

// UserStatistics is a severely cut down version of what actually is contained here.
type UserStatistics struct {
	GlobalRank  int `json:"global_rank"`
	CountryRank int `json:"country_rank"`
}
