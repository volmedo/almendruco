package raices

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	c.login(creds)

	msgs := []Message{}

	u, _ := url.Parse(c.baseURL.String())
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

		// For some reason, the server is using ISO 8859-1 to encode its responses instead of UTF-8
		utf8Reader := transform.NewReader(bytes.NewReader(data), charmap.ISO8859_1.NewDecoder())
		utf8Data, _ := ioutil.ReadAll(utf8Reader)

		var msgResp messagesResponse
		if err := json.Unmarshal(utf8Data, &msgResp); err != nil {
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

func (c *client) login(creds repo.Credentials) error {
	params := url.Values{}
	params.Set(userParam, creds.UserName)
	params.Set(passParam, creds.Password)
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

	return c.checkSessionCookie()
}

func (c *client) checkSessionCookie() error {
	cks := c.http.Jar.Cookies(c.baseURL)
	if len(cks) == 0 {
		return errors.New("no cookies received")
	}

	for _, ck := range cks {
		if ck.Name == loginCookieName {
			return nil
		}
	}

	return errors.New("no login cookies found")
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
