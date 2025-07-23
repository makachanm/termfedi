package config

type Configuration struct {
	Session SessionConfiguration `json:"session"`
}

type SessionConfiguration struct {
	Url   string `json:"instance_url"`
	Token string `json:"access_token"`
}
