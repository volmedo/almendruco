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

func (tn *telegramNotifier) Notify(chatID ChatID, msgs []raices.Message) (uint64, error) {
	u, _ := url.Parse(tn.baseURL.String())
	u.Path = path.Join(u.Path, sendMessagePath)

	params := url.Values{}
	params.Set(chatIDParam, strconv.FormatUint(uint64(chatID), 10))
	params.Set(parseModeParam, parseModeHTML)

	var lastNotifiedMessage uint64
	for _, m := range msgs {
		text := formatText(m)

		params.Set(textParam, text)

		resp, err := tn.http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
		if err != nil {
			return lastNotifiedMessage, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return lastNotifiedMessage, fmt.Errorf("received status code %d", resp.StatusCode)
		}

		lastNotifiedMessage = m.ID
	}

	return lastNotifiedMessage, nil
}

func formatText(m raices.Message) string {
	var sb strings.Builder
	sb.WriteString("Nuevo mensaje en Raíces!")
	sb.WriteString(fmt.Sprintf("\n\n<b>Fecha:</b> %s", m.SentDate.Format(dateFormat)))
	sb.WriteString(fmt.Sprintf("\n<b>De:</b> %s", m.Sender))
	sb.WriteString(fmt.Sprintf("\n<b>Asunto:</b> %s", m.Subject))
	sb.WriteString(fmt.Sprintf("\n\n%s", formatBody(m.Body)))

	if m.ContainsAttachments {
		sb.WriteString(fmt.Sprintf("\n\n<b>Adjuntos:</b>\n%s", formatAttachments(m.Attachments)))
	}

	return sb.String()
}

func formatBody(body string) string {
	processed := strings.Replace(body, "<div>", "\n", -1)

	p := bluemonday.StrictPolicy()
	return p.Sanitize(processed)
}

func formatAttachments(attachments []raices.Attachment) string {
	var sb strings.Builder
	for _, a := range attachments {
		sb.WriteString(fmt.Sprintf("\t\t\t%s\n", a.FileName))
	}

	return sb.String()
}
