package notifier

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
				Contents: []byte{1, 2, 3, 4, 5, 6},
			},
		},
		ReadDate: time.Now(),
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.String(), sendMessagePath) {
			err := r.ParseForm()
			require.NoError(t, err)

			reqChatID := r.Form.Get("chat_id")
			assert.Equal(t, "123456789", reqChatID)

			reqText := r.Form.Get("text")
			expectedText := "Nuevo mensaje en Raíces!\n\n<b>Fecha:</b> 11/11/2021 00:00\n<b>De:</b> Test Sender\n<b>Asunto:</b> Test Subject\n\nHi you, this is a test message\n\n<b>Adjuntos:</b>\n\t\t\tattachment.file\n"
			assert.Equal(t, expectedText, reqText)

			w.WriteHeader(http.StatusOK)
		} else if strings.HasSuffix(r.URL.String(), sendDocumentPath) {
			err := r.ParseMultipartForm(10)
			require.NoError(t, err)

			reqChatID := r.MultipartForm.Value["chat_id"][0]
			assert.Equal(t, "123456789", reqChatID)

			reqDocument := r.MultipartForm.File["document"][0]
			assert.Equal(t, "attachment.file", reqDocument.Filename)
			assert.Equal(t, int64(6), reqDocument.Size)

			f, err := reqDocument.Open()
			require.NoError(t, err)
			defer f.Close()
			contents, err := io.ReadAll(f)
			require.NoError(t, err)
			assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, contents)

			w.WriteHeader(http.StatusOK)
		}

	}))
	defer svr.Close()

	tn, err := NewTelegramNotifier(svr.URL, "test_token")
	require.NoError(t, err)

	lastNotifiedMessage, err := tn.Notify(chatID, []raices.Message{msg})
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), lastNotifiedMessage)
}
