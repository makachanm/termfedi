package config

type AppType string

const (
	Mastodon = "mastodon"
	Misskey  = "misskey"
)

type Configuration struct {
	Session SessionConfiguration `json:"session"`
	UI      UIConfiguration      `json:"ui"`
}

type UIConfiguration struct {
	MaxItemHeight int `json:"maxheight"`
}

type SessionConfiguration struct {
	Type  AppType `json:"type"`
	Url   string  `json:"instance_url"`
	Token string  `json:"access_token"`
}
