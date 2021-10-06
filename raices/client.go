package raices

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const (
	msgPath     = "/raiz_app/jsp/pasendroid/mensajeria"
	pageParam   = "PAGINA"
	msgsPerPage = 10
)

type Client interface {
	FetchMessages() ([]Message, error)
}

type client struct {
	http    http.Client
	baseURL *url.URL
}

func NewClient(baseURL string) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return &client{}, err
	}

	return &client{
		http:    http.Client{},
		baseURL: u,
	}, nil
}

func (c *client) FetchMessages() ([]Message, error) {
	msgs := []Message{}

	u := c.baseURL
	u.Path = path.Join(u.Path, msgPath)
	numMsgs := msgsPerPage
	for i := 1; numMsgs == msgsPerPage; i++ {
		q := url.Values{}
		q.Set(pageParam, fmt.Sprint(i))
		u.RawQuery = q.Encode()
		resp, err := c.http.Get(u.String())
		if err != nil {
			return []Message{}, err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []Message{}, err
		}

		var msgResp messagesResponse
		if err := json.Unmarshal(data, &msgResp); err != nil {
			return []Message{}, err
		}

		parsed, err := parseMessages(msgResp.Messages)
		if err != nil {
			return []Message{}, err
		}

		msgs = append(msgs, parsed...)

		numMsgs = len(msgs)
	}

	return msgs, nil
}

func parseMessages(raw []rawMessage) ([]Message, error) {
	parsed := make([]Message, 0, len(raw))
	for _, r := range raw {
		m, err := parseMessage(r)
		if err != nil {
			return []Message{}, err
		}

		parsed = append(parsed, m)
	}

	return parsed, nil
}
