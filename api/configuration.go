package api

type Configuration struct {
	Port              int
	ClientID          string
	ClientSecret      string
	Whitelist         []int64
	WhitelistDisabled bool
	UpdateInterval    int64
}
