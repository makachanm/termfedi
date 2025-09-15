package layer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/makachanm/flogger-lib"
)

type MastodonFetch struct {
	token        string
	instance_url *url.URL
}

func NewMastodonFetch(token string, insurl string) *MastodonFetch {
	var err error
	flogger.Printf("Initializing MastodonFetch for instance %s", insurl)

	mastodon := new(MastodonFetch)
	mastodon.instance_url, err = url.Parse(insurl)
	mastodon.token = token

	if err != nil {
		flogger.Errorf("Failed to parse instance url: %v", err)
		panic(err)
	}

	return mastodon
}

func (m *MastodonFetch) getData(method string, path string, data interface{}, authneed bool, hasdata bool) ([]byte, error) {
	spath := strings.Split(path, "/")
	xurl := m.instance_url.JoinPath(spath...)
	flogger.Printf("Mastodon API call: %s %s", method, xurl.String())

	var bodyreader *bytes.Reader
	if hasdata {
		bodyreader = bytes.NewReader(marshallJSON(&data))
	} else {
		bodyreader = bytes.NewReader([]byte{})
	}

	httpreq, err := http.NewRequest(method, xurl.String(), bodyreader)
	if err != nil {
		return []byte{}, flogger.Errorf("Failed to create http request: %v", err)
	}

	if authneed {
		btoken := fmt.Sprintf("Bearer %s", m.token)
		httpreq.Header.Set("Authorization", btoken)
	}

	resp, herr := http.DefaultClient.Do(httpreq)
	if herr != nil {
		return []byte{}, flogger.Errorf("Failed to execute http request: %v", herr)
	}
	defer resp.Body.Close()

	bytes, rerr := io.ReadAll(resp.Body)

	if rerr != nil {
		return []byte{}, flogger.Errorf("Failed to read response body: %v", rerr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		flogger.Errorf("Mastodon API error: %s, %s", resp.Status, string(bytes))
		return bytes, fmt.Errorf("mastodon API error: %s", resp.Status)
	}

	return bytes, nil
}

func (m *MastodonFetch) getQueryData(method string, path string, authneed bool, querys map[string]string) ([]byte, error) {
	spath := strings.Split(path, "/")
	xurl := m.instance_url.JoinPath(spath...)

	query := url.Values{}
	for val, key := range querys {
		query.Add(val, key)
	}
	fullUrl := fmt.Sprintf("%s?%s", xurl.String(), query.Encode())
	flogger.Printf("Mastodon API call: %s %s", method, fullUrl)

	httpreq, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return []byte{}, flogger.Errorf("Failed to create http request: %v", err)
	}

	if authneed {
		btoken := fmt.Sprintf("Bearer %s", m.token)
		httpreq.Header.Set("Authorization", btoken)
	}

	resp, herr := http.DefaultClient.Do(httpreq)
	if herr != nil {
		return []byte{}, flogger.Errorf("Failed to execute http request: %v", herr)
	}
	defer resp.Body.Close()

	bytes, rerr := io.ReadAll(resp.Body)

	if rerr != nil {
		return []byte{}, flogger.Errorf("Failed to read response body: %v", rerr)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		flogger.Errorf("Mastodon API error: %s, %s", resp.Status, string(bytes))
		return bytes, fmt.Errorf("mastodon API error: %s", resp.Status)
	}

	return bytes, nil
}

func (m *MastodonFetch) GetGlobalTimeline() []Note {
	d, err := m.getQueryData(http.MethodGet, "api/v1/timelines/public", true, map[string]string{"limit": "40"})
	if err != nil {
		flogger.Errorf("Failed to get global timeline: %v", err)
		return make([]Note, 0)
	}

	var mNotes []MastodonNote
	errs := unmarshallJSON[[]MastodonNote](&mNotes, d)
	if !errs {
		flogger.Errorf("Failed to unmarshal Mastodon notes for global timeline")
		return make([]Note, 0)
	}
	flogger.Printf("Fetched %d notes from global timeline", len(mNotes))

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)

		if mNotes[i].RenotedField != nil {
			rnotes[i].IsRenote = true
			rnotes[i].Renote = "Renote of " + mNotes[i].RenotedField.User.Name + "'s Note: <br>" + mNotes[i].RenotedField.Content
		} else {
			rnotes[i].IsRenote = false
		}

		if len(mNotes[i].Medias) >= 1 {
			rnotes[i].HasMedia = true
			for _, data := range mNotes[i].Medias {
				rnotes[i].Medias = append(rnotes[i].Medias, data.URL)
			}
		} else {
			rnotes[i].HasMedia = false
		}

		rnotes[i].Spoiler = mNotes[i].Spoilerwarning

	}

	return rnotes
}
func (m *MastodonFetch) GetLocalTimeline() []Note {
	d, err := m.getQueryData(http.MethodGet, "api/v1/timelines/public", true, map[string]string{"local": "true", "limit": "40"})
	if err != nil {
		flogger.Errorf("Failed to get local timeline: %v", err)
		return make([]Note, 0)
	}

	var mNotes []MastodonNote
	errs := unmarshallJSON[[]MastodonNote](&mNotes, d)
	if !errs {
		flogger.Errorf("Failed to unmarshal Mastodon notes for local timeline")
		return make([]Note, 0)
	}
	flogger.Printf("Fetched %d notes from local timeline", len(mNotes))

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)
		if mNotes[i].RenotedField != nil {
			rnotes[i].IsRenote = true
			rnotes[i].Renote = "Renote of " + mNotes[i].RenotedField.User.Name + "'s Note: <br>" + mNotes[i].RenotedField.Content
		} else {
			rnotes[i].IsRenote = false
		}

		if len(mNotes[i].Medias) >= 1 {
			rnotes[i].HasMedia = true
			for _, data := range mNotes[i].Medias {
				rnotes[i].Medias = append(rnotes[i].Medias, data.URL)
			}
		} else {
			rnotes[i].HasMedia = false
		}

		rnotes[i].Spoiler = mNotes[i].Spoilerwarning

	}

	return rnotes
}
func (m *MastodonFetch) GetHomeTimeline() []Note {
	d, err := m.getQueryData(http.MethodGet, "api/v1/timelines/home", true, map[string]string{"limit": "40"})
	if err != nil {
		flogger.Errorf("Failed to get home timeline: %v", err)
		return make([]Note, 0)
	}

	var mNotes []MastodonNote
	errs := unmarshallJSON[[]MastodonNote](&mNotes, d)
	if !errs {
		flogger.Errorf("Failed to unmarshal Mastodon notes for home timeline")
		return make([]Note, 0)
	}
	flogger.Printf("Fetched %d notes from home timeline", len(mNotes))

	var rnotes []Note = make([]Note, len(mNotes))
	for i := 0; i < len(mNotes); i++ {
		rnotes[i].Id = mNotes[i].Id
		rnotes[i].Author_finger = mNotes[i].User.Finger
		rnotes[i].Author_name = mNotes[i].User.Name
		rnotes[i].Content = mNotes[i].Content
		rnotes[i].RenoteCount = mNotes[i].ReblogsCount
		rnotes[i].ReactionCount = mNotes[i].FavouritesCount
		rnotes[i].Visiblity = platformVisblityToValue(mNotes[i].Visibility)

		if mNotes[i].RenotedField != nil {
			rnotes[i].IsRenote = true
			rnotes[i].Renote = "Renote of " + mNotes[i].RenotedField.User.Name + "'s Note: <br>" + mNotes[i].RenotedField.Content
		} else {
			rnotes[i].IsRenote = false
		}

		if len(mNotes[i].Medias) >= 1 {
			rnotes[i].HasMedia = true
			for _, data := range mNotes[i].Medias {
				rnotes[i].Medias = append(rnotes[i].Medias, data.URL)
			}
		} else {
			rnotes[i].HasMedia = false
		}

		rnotes[i].Spoiler = mNotes[i].Spoilerwarning
	}

	return rnotes
}

func (m *MastodonFetch) GetNotifications() []Notification {
	d, err := m.getData(http.MethodGet, "api/v1/notifications", nil, true, false)
	if err != nil {
		flogger.Errorf("Failed to get notifications: %v", err)
		return make([]Notification, 0)
	}

	var mnotis []MastodonNotification
	errs := unmarshallJSON[[]MastodonNotification](&mnotis, d)
	if !errs {
		flogger.Errorf("Failed to unmarshal Mastodon notifications")
		return make([]Notification, 0)
	}
	flogger.Printf("Fetched %d notifications", len(mnotis))

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
			if len(mnotis[i].Status.Mentions) > 0 {
				rnotis[i].ReactedUser = User{
					Id:          mnotis[i].Status.Mentions[0].Id,
					User_name:   mnotis[i].Status.Mentions[0].Username,
					User_finger: mnotis[i].Status.Mentions[0].Acct,
				}
			}

		case NOTI_FOLLOW, NOTI_FAVOURITE, NOTI_RENOTE:
			rnotis[i].Content = mnotis[i].Status.Contents
			rnotis[i].ReactedUser = User{
				Id:          mnotis[i].Account.Id,
				User_name:   mnotis[i].Account.Username,
				User_finger: mnotis[i].Account.Acct,
			}
		}
	}

	return rnotis
}

func (m *MastodonFetch) PostRenote(note_id string) bool {
	flogger.Printf("Posting renote for note %s", note_id)
	_, err := m.getData(http.MethodPost, fmt.Sprintf("api/v1/statuses/%s/reblog", note_id), nil, true, false)
	if err != nil {
		flogger.Errorf("Failed to post renote: %v", err)
	}
	return err == nil
}

func (m *MastodonFetch) PostReaction(note_id string) bool {
	flogger.Printf("Posting reaction for note %s", note_id)
	// Mastodon only has favourite
	_, err := m.getData(http.MethodPost, fmt.Sprintf("api/v1/statuses/%s/favourite", note_id), nil, true, false)
	if err != nil {
		flogger.Errorf("Failed to post reaction: %v", err)
	}
	return err == nil
}
