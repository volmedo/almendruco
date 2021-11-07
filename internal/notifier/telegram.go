package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"github.com/volmedo/almendruco.git/internal/raices"
)

const (
	chatIDParam    = "chat_id"
	textParam      = "text"
	parseModeParam = "parse_mode"
	parseModeHTML  = "HTML"

	sendMessagePath = "sendMessage"

	dateFormat = "02/01/2006 15:04"
)

type telegramNotifier struct {
	baseURL *url.URL
	http    *http.Client
}

func NewTelegramNotifier(baseURL, botToken string) (Notifier, error) {
	u, err := url.Parse(fmt.Sprintf("%s/bot%s", baseURL, botToken))
	if err != nil {
		return &telegramNotifier{}, fmt.Errorf("bad baseURL and/or botToken: %s", err)
	}

	return &telegramNotifier{
		baseURL: u,
		http:    &http.Client{},
	}, nil
}

func (tn *telegramNotifier) Notify(chatID ChatID, msgs []raices.Message) error {
	u, _ := url.Parse(tn.baseURL.String())
	u.Path = path.Join(u.Path, sendMessagePath)

	params := url.Values{}
	params.Set(chatIDParam, strconv.FormatUint(uint64(chatID), 10))
	params.Set(parseModeParam, parseModeHTML)

	for _, m := range msgs {
		text := formatText(m)

		params.Set(textParam, text)

		resp, err := tn.http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("received status code %d", resp.StatusCode)
		}
	}

	return nil
}

func formatText(m raices.Message) string {
	attachments := "no"
	if m.ContainsAttachments {
		attachments = "sí"
	}

	sentDate := m.SentDate.Format(dateFormat)
	body := preprocess(m.Body)

	return fmt.Sprintf("Nuevo mensaje en Raíces!\n\n<b>Fecha:</b> %s\n<b>De:</b> %s\n<b>Asunto:</b> %s\n\n%s\n\n<b>Adjuntos:</b> %s",
		sentDate, m.Sender, m.Subject, body, attachments)
}

func preprocess(html string) string {
	processed := strings.Replace(html, "<div>", "\n", -1)

	p := bluemonday.StrictPolicy()
	return p.Sanitize(processed)
}
