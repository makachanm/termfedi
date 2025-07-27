package layer

import "time"

type MisskeyNoteTransaction struct {
	Limit int `json:"limit"`
}

type MisskeyUser struct {
	Name   string `json:"name"`
	Finger string `json:"username"`
	Host   string `json:"host"`
}

type MisskeyNote struct {
	Id              string      `json:"id"`
	Content         string      `json:"text"`
	Spoilerwarning  *string     `json:"cw"`
	User            MisskeyUser `json:"user"`
	RenoteCount     int         `json:"renoteCount"`
	FavouritesCount int         `json:"reactionCount"`
	Visibility      string      `json:"visibility"`
}

type MisskeyNotification struct {
	Id      string           `json:"id"`
	Type    string           `json:"type"`
	Account MastodonNotiUser `json:"account"`
	Status  MastodonStatus   `json:"status,omitempty"`
}

type MisskeyNotiUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
}

type MisskeyStatus struct {
	ID                 string             `json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	InReplyToID        *string            `json:"in_reply_to_id,omitempty"`
	InReplyToAccountID *string            `json:"in_reply_to_account_id,omitempty"`
	Account            MastodonNotiUser   `json:"account"`
	Contents           string             `json:"content,omitempty"`
	Mentions           []MastodonMentions `json:"mentions,omitempty"`
}
