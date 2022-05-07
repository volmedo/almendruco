package raices

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/volmedo/almendruco.git/internal/repo"
)

const msgsPage = 10

func TestFetchMessages(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(loginPath, http.HandlerFunc(happyLoginHandler))
	mux.Handle(msgPath, http.HandlerFunc(happyMessagesHandler))
	mux.Handle(attachmentPath, http.HandlerFunc(happyAttachmentHandler))

	svr := httptest.NewServer(mux)
	defer svr.Close()

	c, err := NewClient(svr.URL)
	require.NoError(t, err, "Unable to create client")

	testCreds := repo.Credentials{
		User: "Some User",
		Pass: "s0m3p4ss",
	}

	msgs, err := c.FetchMessages(testCreds, 0)

	require.NoError(t, err, "Unexpected error fetching messages")
	require.Equal(t, 1, len(msgs), "Expected 1 message")

	// Time strings reported by Raices are always CET/CEST
	cet, err := time.LoadLocation("CET")
	require.NoError(t, err, "Failed to load CET/CEST timezone data")

	expected := Message{
		ID:                  12345678,
		SentDate:            time.Date(2021, time.October, 1, 18, 27, 0, 0, cet),
		Sender:              "Jon Doe (Director)",
		Subject:             "SOME SUBJECT",
		Body:                "A message with some HTML entities&nbsp; and <div>markup</div>",
		ContainsAttachments: true,
		Attachments: []Attachment{
			{
				ID:       123456,
				FileName: "Some File.ext",
				Contents: []byte{1, 2, 3, 4, 5, 6},
			},
		},
		ReadDate: time.Date(2021, time.October, 2, 19, 3, 00, 00, cet),
	}

	if diff := cmp.Diff(expected, msgs[0]); diff != "" {
		t.Fatalf("Message not equal to expected:\n%s", diff)
	}
}

func happyLoginHandler(w http.ResponseWriter, r *http.Request) {
	testResp := `
		{
			"ESTADO": {
				"CODIGO": "C"
			}
		}
		`

	loginCk := &http.Cookie{
		Name:     loginCookieName,
		Value:    "abcd",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, loginCk)
	w.WriteHeader(200)
	_, _ = w.Write([]byte(testResp))
}

func happyMessagesHandler(w http.ResponseWriter, r *http.Request) {
	testResp := `
		{
			"ESTADO": {
				"CODIGO": "C"
			},
			"RESULTADO": [
				{
					"X_NOTMENSAL": 12345678,
					"F_ENVIO": "01/10/2021 18:27",
					"L_ADJUNTO": "S",
					"T_ASUNTO": "SOME SUBJECT",
					"F_LECTURA": "02/10/2021 19:03",
					"CENTRO": "12345678 - SOME SCHOOL",
					"REMITIDO": "Jon Doe (Director)",
					"T_MENSAJE": "A message with some HTML entities&nbsp; and <div>markup</div>",
					"ADJUNTOS": [
						{
							"X_ADJMENSAL": 123456,
							"T_NOMFIC": "Some File.ext"
						}
					]
				}
			]
		}
		`
	fmt.Fprint(w, testResp)
}

func happyAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte{1, 2, 3, 4, 5, 6})
}

func TestMultiPageMessages(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(loginPath, http.HandlerFunc(happyLoginHandler))
	mux.Handle(msgPath, http.HandlerFunc(multiPageHandler))

	svr := httptest.NewServer(mux)
	defer svr.Close()

	c, err := NewClient(svr.URL)
	require.NoError(t, err, "Unable to create client")

	testCreds := repo.Credentials{
		User: "Some User",
		Pass: "s0m3p4ss",
	}

	msgs, err := c.FetchMessages(testCreds, 0)

	assert.NoError(t, err, "Unexpected error fetching messages")
	assert.Equal(t, 15, len(msgs), "Expected 15 messages")
}

func TestOnlyReturnsNonNotifiedMsgs(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(loginPath, http.HandlerFunc(happyLoginHandler))
	mux.Handle(msgPath, http.HandlerFunc(multiPageHandler))

	svr := httptest.NewServer(mux)
	defer svr.Close()

	c, err := NewClient(svr.URL)
	require.NoError(t, err, "Unable to create client")

	testCreds := repo.Credentials{
		User: "Some User",
		Pass: "s0m3p4ss",
	}

	msgs, err := c.FetchMessages(testCreds, 4)

	assert.NoError(t, err, "Unexpected error fetching messages")
	assert.Equal(t, 11, len(msgs), "Expected 11 messages")
}

func multiPageHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	page, ok := params[pageParam]
	if !ok || page[0] == "" {
		http.Error(w, "page param missing", http.StatusBadRequest)
	}

	pageNum, err := strconv.Atoi(page[0])
	if err != nil {
		http.Error(w, "page param is not a number", http.StatusBadRequest)
	}

	// We'll return <msgsPage> messages with IDs in descending order starting from <numMsgs>
	numMsgs := 15
	firstMsg := numMsgs - msgsPage*(pageNum-1)
	lastMsg := firstMsg - msgsPage + 1
	lastMsg = int(math.Max(float64(lastMsg), 1))
	msgs := make([]rawMessage, 0)
	for i := firstMsg; i >= lastMsg; i-- {
		msg := rawMessage{
			ID:                  uint64(i),
			SentDate:            "01/10/2021 18:27",
			Sender:              "Jon Doe (Director)",
			Subject:             "SOME SUBJECT",
			Body:                "A message with some HTML entities&nbsp; and <div>markup</div>",
			ContainsAttachments: "S",
			Attachments: []rawAttachment{
				{
					ID:       123456,
					FileName: "Some File.ext",
				},
			},
			ReadDate: "02/10/2021 19:03",
		}

		msgs = append(msgs, msg)
	}

	messagesResp := messagesResponse{
		Status:   status{Code: statusCodeOK},
		Messages: msgs,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(messagesResp)
}
