package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/volmedo/almendruco.git/internal/raices"
)

func TestNotify(t *testing.T) {
	chatID := ChatID(123456789)

	msg := raices.Message{
		ID:                  123456,
		SentDate:            time.Date(2021, time.Month(11), 11, 0, 0, 0, 0, time.UTC),
		Sender:              "Test Sender",
		Subject:             "Test Subject",
		Body:                "Hi you, this is a test message",
		ContainsAttachments: true,
		Attachments: []raices.Attachment{
			{
				ID:       98765,
				FileName: "attachment.file",
			},
		},
		ReadDate: time.Now(),
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		require.NoError(t, err)

		reqChatID := r.Form.Get(chatIDParam)
		assert.Equal(t, "123456789", reqChatID)

		reqText := r.Form.Get(textParam)
		expectedText := "New message arrived to Raíces!\nFrom: Test Sender\nSubject: Test Subject\nHi you, this is a test message\nAttachments: true"
		assert.Equal(t, expectedText, reqText)

		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()

	tn, err := NewTelegramNotifier(svr.URL, "test_token")
	require.NoError(t, err)

	err = tn.Notify(chatID, []raices.Message{msg})
	assert.NoError(t, err)
}
