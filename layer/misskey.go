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

type MisskeyFetch struct {
	token        string
	instance_url *url.URL
}

func NewMisskeyFetch(token string, insurl string) *MisskeyFetch {
	var err error

	misskey := new(MisskeyFetch)
	misskey.instance_url, err = url.Parse(insurl)
	misskey.token = token

	if err != nil {
		panic(err)
	}

	return misskey
}

func (m *MisskeyFetch) getData(path string, data interface{}, authneed bool, hasdata bool) []byte {
	spath := strings.Split(path, "/")
	xurl := m.instance_url.JoinPath(spath...)

	fmt.Println(xurl)

	var bodyreader *bytes.Reader
	if hasdata {
		bodyreader = bytes.NewReader(marshallJSON(&data))
	} else {
		bodyreader = bytes.NewReader([]byte{})
	}

	httpreq, err := http.NewRequest(http.MethodPost, xurl.String(), bodyreader)
	if err != nil {
		panic(err)
	}

	if authneed {
		btoken := fmt.Sprintf("Bearer %s", m.token)
		httpreq.Header.Set("Authorization", btoken)
	}

	httpreq.Header.Set("Content-Type", "application/json")

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

func (m *MisskeyFetch) GetGlobalTimeline() []Note {
	d := m.getData("api/notes/global-timeline", MisskeyNoteTransaction{Limit: 100}, false, true)

	var mnotes []MisskeyNote
	unmarshallJSON[[]MisskeyNote](&mnotes, d)

	var rnotes []Note = make([]Note, len(mnotes))
	for i := 0; i < len(mnotes); i++ {
		rnotes[i].Id = mnotes[i].Id
		rnotes[i].Author_finger = mnotes[i].User.Finger
		rnotes[i].Author_name = mnotes[i].User.Name
		rnotes[i].Content = mnotes[i].Content
		rnotes[i].RenoteCount = mnotes[i].RenoteCount
		rnotes[i].ReactionCount = mnotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mnotes[i].Visibility)

		if mnotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = *mnotes[i].Spoilerwarning
		}
	}

	return rnotes
}
func (m *MisskeyFetch) GetLocalTimeline() []Note {
	d := m.getData("api/notes/local-timeline", MisskeyNoteTransaction{Limit: 100}, false, true)

	var mnotes []MisskeyNote
	unmarshallJSON[[]MisskeyNote](&mnotes, d)

	var rnotes []Note = make([]Note, len(mnotes))
	for i := 0; i < len(mnotes); i++ {
		rnotes[i].Id = mnotes[i].Id
		rnotes[i].Author_finger = mnotes[i].User.Finger
		rnotes[i].Author_name = mnotes[i].User.Name
		rnotes[i].Content = mnotes[i].Content
		rnotes[i].RenoteCount = mnotes[i].RenoteCount
		rnotes[i].ReactionCount = mnotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mnotes[i].Visibility)

		if mnotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = *mnotes[i].Spoilerwarning
		}
	}

	return rnotes
}
func (m *MisskeyFetch) GetHomeTimeline() []Note {
	d := m.getData("api/notes/timeline", MisskeyNoteTransaction{Limit: 100}, true, true)

	var mnotes []MisskeyNote
	unmarshallJSON[[]MisskeyNote](&mnotes, d)

	var rnotes []Note = make([]Note, len(mnotes))
	for i := 0; i < len(mnotes); i++ {
		rnotes[i].Id = mnotes[i].Id
		rnotes[i].Author_finger = mnotes[i].User.Finger
		rnotes[i].Author_name = mnotes[i].User.Name
		rnotes[i].Content = mnotes[i].Content
		rnotes[i].RenoteCount = mnotes[i].RenoteCount
		rnotes[i].ReactionCount = mnotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mnotes[i].Visibility)

		if mnotes[i].Spoilerwarning != nil {
			rnotes[i].Spoiler = *mnotes[i].Spoilerwarning
		}
	}

	return rnotes
}

func (m *MisskeyFetch) GetPost(id string) Note { return Note{} }

func (m *MisskeyFetch) GetNotifications() []Notification {
	return []Notification{}
}

func (m *MisskeyFetch) GetNotification(id string) Notification { return Notification{} }

func (m *MisskeyFetch) GetUser(id string) User { return User{} }
