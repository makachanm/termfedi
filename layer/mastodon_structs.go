package layer

type MastodonUser struct {
	Name   string `json:"display_name"`
	Finger string `json:"acct"`
	Desc   string `json:"note"`
}

type MastodonNote struct {
	Id              string       `json:"id"`
	Content         string       `json:"content"`
	Spoilerwarning  *string      `json:"spoiler_text"`
	User            MastodonUser `json:"account"`
	ReblogsCount    int          `json:"reblogs_count"`
	FavouritesCount int          `json:"favourites_count"`
	Visibility      string       `json:"visibility"`
}
