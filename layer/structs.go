package layer

import (
	"time"
)

type Visiblity int

const (
	VISIBLITY_PUBLIC Visiblity = iota + 1
	VISIBLITY_QUIET
	VISIBLITY_FOLLOWER
	VISIBLITY_DIRECT
)

func VisiblityToText(v Visiblity) string {
	switch v {
	case VISIBLITY_PUBLIC:
		return "Global"
	case VISIBLITY_QUIET:
		return "Quiet"
	case VISIBLITY_FOLLOWER:
		return "Follower-only"
	case VISIBLITY_DIRECT:
		return "Direct-only"

	default:
		return "Unknown"
	}
}

type Note struct {
	Id string `json:"id"`

	Author_name   string
	Author_finger string

	Visiblity Visiblity
	Time      time.Time

	Spoiler *string `json:"spoiler_text"`
	Content string  `json:"content"`
}

type User struct {
	Id string

	User_name        string
	User_finger      string
	User_description string

	User_desc_field map[string]string
}

type Notification struct {
	Id      string
	Content string
}
