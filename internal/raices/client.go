package raices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/volmedo/almendruco.git/internal/repo"
)

const (
	loginPath       = "/raiz_app/jsp/pasendroid/login"
	userParam       = "USUARIO"
	passParam       = "CLAVE"
	verParam        = "p"
	verString       = `{"version":"1.0.23"}`
	loginCookieName = "JSESSIONID"

	msgPath     = "/raiz_app/jsp/pasendroid/mensajeria"
	pageParam   = "PAGINA"
	msgsPerPage = 10
)

type Client interface {
	FetchMessages(creds repo.Credentials, lastNotifiedMessage uint64) ([]Message, error)
}

type client struct {
	http    *http.Client
	baseURL *url.URL
}

func NewClient(baseURL string) (Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return &client{}, err
	}

	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return &client{}, err
	}

	hc := &http.Client{Jar: j}

	return &client{
		http:    hc,
		baseURL: u,
	}, nil
}

func (c *client) FetchMessages(creds repo.Credentials, lastNotifiedMessage uint64) ([]Message, error) {
	// Login if needed
	if err := c.login(creds); err != nil {
		return []Message{}, err
	}

	u, _ := url.Parse(c.baseURL.String())
	u.Path = path.Join(u.Path, msgPath)

	msgs := []Message{}
	numMsgs := msgsPerPage
	for i := 1; numMsgs == msgsPerPage; i++ {
		rawMsgs, err := c.fetchPage(u, i)
		if err != nil {
			return []Message{}, err
		}

		rawMsgs = filterNotified(rawMsgs, lastNotifiedMessage)

		parsed, err := parse(rawMsgs)
		if err != nil {
			return []Message{}, err
		}

		numMsgs = len(parsed)

		msgs = append(msgs, parsed...)
	}

	return reverse(msgs), nil
}

func (c *client) login(creds repo.Credentials) error {
	params := url.Values{}
	params.Set(userParam, creds.User)
	params.Set(passParam, creds.Pass)
	params.Set(verParam, verString)

	u, _ := url.Parse(c.baseURL.String())
	u.Path = path.Join(u.Path, loginPath)
	resp, err := c.http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResp loginResponse
	if err := json.Unmarshal(data, &loginResp); err != nil {
		return err
	}

	if loginResp.Status.Code != statusCodeOK {
		return fmt.Errorf("code %s in login response: %s", loginResp.Status.Code, loginResp.Status.Description)
	}

	return nil
}

func (c *client) fetchPage(u *url.URL, pageNum int) ([]rawMessage, error) {
	q := url.Values{}
	q.Set(pageParam, fmt.Sprint(pageNum))
	u.RawQuery = q.Encode()
	resp, err := c.http.Get(u.String())
	if err != nil {
		return []rawMessage{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []rawMessage{}, err
	}

	// For some reason, the server is using ISO 8859-1 to encode its responses instead of UTF-8
	utf8Reader := transform.NewReader(bytes.NewReader(data), charmap.ISO8859_1.NewDecoder())
	utf8Data, _ := io.ReadAll(utf8Reader)

	var msgResp messagesResponse
	if err := json.Unmarshal(utf8Data, &msgResp); err != nil {
		return []rawMessage{}, err
	}

	return msgResp.Messages, nil
}

func filterNotified(rawMsgs []rawMessage, lastNotifiedMessage uint64) []rawMessage {
	lastMessageToNotify := len(rawMsgs)
	for j, r := range rawMsgs {
		if r.ID <= lastNotifiedMessage {
			lastMessageToNotify = j
			break
		}
	}

	return rawMsgs[:lastMessageToNotify]
}

func parse(raw []rawMessage) ([]Message, error) {
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

func reverse(msgs []Message) []Message {
	reversed := make([]Message, 0, len(msgs))
	for i := len(msgs) - 1; i >= 0; i-- {
		reversed = append(reversed, msgs[i])
	}

	return reversed
}
