package layer

import "time"

type MastodonUser struct {
	Name   string `json:"display_name"`
	Finger string `json:"acct"`
	Desc   string `json:"note"`
}

type MastodonMedia struct {
	URL string `json:"url"`
}

type MastodonNote struct {
	Id              string          `json:"id"`
	Content         string          `json:"content"`
	Spoilerwarning  string          `json:"spoiler_text"`
	User            MastodonUser    `json:"account"`
	ReblogsCount    int             `json:"reblogs_count"`
	FavouritesCount int             `json:"favourites_count"`
	Visibility      string          `json:"visibility"`
	Medias          []MastodonMedia `json:"media_attachments"`
	RenotedField    *MastodonNote   `json:"reblog,omitempty"`
}

type MastodonNotification struct {
	Id      string           `json:"id"`
	Type    string           `json:"type"`
	Account MastodonNotiUser `json:"account"`
	Status  MastodonStatus   `json:"status,omitempty"`
}

type MastodonNotiUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
}

type MastodonMentions struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Acct     string `json:"acct"`
}

type MastodonStatus struct {
	ID                 string             `json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	InReplyToID        *string            `json:"in_reply_to_id,omitempty"`
	InReplyToAccountID *string            `json:"in_reply_to_account_id,omitempty"`
	Account            MastodonNotiUser   `json:"account"`
	Contents           string             `json:"content,omitempty"`
	Mentions           []MastodonMentions `json:"mentions,omitempty"`
}
