package layer

import (
	"time"
)

type Visiblity int
type NotificationType int

const (
	VISIBLITY_PUBLIC Visiblity = iota + 1
	VISIBLITY_QUIET
	VISIBLITY_FOLLOWER
	VISIBLITY_DIRECT
)

const (
	NOTI_MENTION NotificationType = iota + 1
	NOTI_FAVOURITE
	NOTI_RENOTE
	NOTI_FOLLOW
	NOTI_UNSUPPORTED
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

func platformVisblityToValue(s string) Visiblity {
	switch s {
	case "public":
		return VISIBLITY_PUBLIC
	case "unlisted", "home":
		return VISIBLITY_QUIET
	case "private", "followers":
		return VISIBLITY_FOLLOWER
	case "direct", "specified":
		return VISIBLITY_DIRECT

	default:
		return VISIBLITY_QUIET
	}
}

func platformNotiTypeToValue(s string) NotificationType {
	switch s {
	case "mention":
		return NOTI_MENTION
	case "reblog":
		return NOTI_RENOTE
	case "favourite":
		return NOTI_FAVOURITE
	case "follow":
		return NOTI_FOLLOW

	default:
		return NOTI_UNSUPPORTED
	}
}

type Note struct {
	Id string `json:"id"`

	Author_name   string
	Author_finger string

	Visiblity Visiblity
	Time      time.Time

	Spoiler string
	Content string

	RenoteCount   int
	ReactionCount int

	HasMedia bool
	IsRenote bool

	Medias []string
	Renote string
}

type User struct {
	Id string

	User_name        string
	User_finger      string
	User_description string

	User_desc_field map[string]string
}

type Notification struct {
	Id   string
	Type NotificationType

	Content     string
	ReactedUser User
}
