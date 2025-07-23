package layer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type MastodonFetch struct {
	token        string
	instance_url *url.URL
}

func NewMastodonFetch(token string, insurl string) *MastodonFetch {
	var err error

	mastodon := new(MastodonFetch)
	mastodon.instance_url, err = url.Parse(insurl)
	mastodon.token = token

	if err != nil {
		panic(err)
	}

	return mastodon
}

func (m *MastodonFetch) getData(method string, path string, data interface{}, authneed bool, hasdata bool) []byte {
	spath := strings.Split(path, "/")
	xurl := m.instance_url.JoinPath(spath...)

	var bodyreader *bytes.Reader
	if hasdata {
		bodyreader = bytes.NewReader(marshallJSON(&data))
	} else {
		bodyreader = bytes.NewReader([]byte{})
	}

	httpreq, err := http.NewRequest(method, xurl.String(), bodyreader)
	if err != nil {
		panic(err)
	}

	if authneed {
		btoken := fmt.Sprintf("Bearer %s", m.token)
		httpreq.Header.Set("Authorization", btoken)
	}

	resp, herr := http.DefaultClient.Do(httpreq)
	if herr != nil {
		panic(herr)
	}
	defer resp.Body.Close()

	bytes, rerr := io.ReadAll(resp.Body)

	if rerr != nil {
		panic(err)
	}

	if resp.Status != "200 OK" {
		fmt.Println("ERROR:", string(bytes))
		os.Exit(-1)
	}

	return bytes
}

func (m *MastodonFetch) getQueryData(method string, path string, authneed bool, querys map[string]string) []byte {
	spath := strings.Split(path, "/")
	xurl := m.instance_url.JoinPath(spath...)
	uquery := xurl.Query()
	for x, v := range querys {
		uquery.Set(x, v)
	}

	httpreq, err := http.NewRequest(method, xurl.String(), nil)
	if err != nil {
		panic(err)
	}

	if authneed {
		btoken := fmt.Sprintf("Bearer %s", m.token)
		httpreq.Header.Set("Authorization", btoken)
	}

	resp, herr := http.DefaultClient.Do(httpreq)
	if herr != nil {
		panic(herr)
	}
	defer resp.Body.Close()

	bytes, rerr := io.ReadAll(resp.Body)

	if rerr != nil {
		panic(err)
	}

	if resp.Status != "200 OK" {
		fmt.Println("ERROR:", string(bytes))
		os.Exit(-1)
	}

	return bytes
}

func (m *MastodonFetch) GetGlobalTimeline() []Note {
	d := m.getData(http.MethodGet, "api/v1/timelines/public", nil, false, false)

	var mNotes []MastodonNote
	unmarshallJSON[[]MastodonNote](&mNotes, d)

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)

		if mNotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = mNotes[i].Spoilerwarning
		}
	}

	return rnotes
}
func (m *MastodonFetch) GetLocalTimeline() []Note {
	d := m.getQueryData(http.MethodGet, "api/v1/timelines/public", true, map[string]string{"local": "true"})

	var mNotes []MastodonNote
	unmarshallJSON[[]MastodonNote](&mNotes, d)

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)

		if mNotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = mNotes[i].Spoilerwarning
		}
	}

	return rnotes
}
func (m *MastodonFetch) GetHomeTimeline() []Note {
	d := m.getData(http.MethodGet, "api/v1/timelines/home", nil, true, false)

	var mNotes []MastodonNote
	unmarshallJSON[[]MastodonNote](&mNotes, d)

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)

		if mNotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = mNotes[i].Spoilerwarning
		}
	}

	return rnotes
}

func (m *MastodonFetch) GetPost(id string) Note { return Note{} }

func (m *MastodonFetch) GetNotifications() []Notification {
	d := m.getData(http.MethodGet, "api/v1/notifications", nil, true, false)

	var mnotis []MastodonNotification
	unmarshallJSON[[]MastodonNotification](&mnotis, d)

	var rnotis []Notification = make([]Notification, len(mnotis))
	for i := 0; i < len(mnotis); i++ {
		n_type := platformNotiTypeToValue(mnotis[i].Type)
		if n_type == NOTI_UNSUPPORTED {
			continue
		}

		rnotis[i].Id = mnotis[i].Id
		rnotis[i].Type = n_type

		switch n_type {
		case NOTI_MENTION:
			rnotis[i].Content = mnotis[i].Status.Contents
			rnotis[i].ReactedUser = User{
				Id:          mnotis[i].Status.Mentions[0].Id,
				User_name:   mnotis[i].Status.Mentions[0].Username,
				User_finger: mnotis[i].Status.Mentions[0].Acct,
			}

		case NOTI_FOLLOW, NOTI_FAVOURITE, NOTI_RENOTE:
			rnotis[i].ReactedUser = User{
				Id:          mnotis[i].Status.ID,
				User_name:   mnotis[i].Status.Account.Username,
				User_finger: mnotis[i].Status.Account.Acct,
			}
		}
	}

	return rnotis
}

func (m *MastodonFetch) GetNotification(id string) Notification { return Notification{} }

func (m *MastodonFetch) GetUser(id string) User { return User{} }
