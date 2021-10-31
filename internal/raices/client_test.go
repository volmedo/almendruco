package raices

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/volmedo/almendruco.git/internal/repo"
)

func TestFetchMessages(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle(loginPath, http.HandlerFunc(happyLoginHandler))
	mux.Handle(msgPath, http.HandlerFunc(happyMessagesHandler))

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
	assert.Equal(t, 1, len(msgs), "Expected 1 message")

	expected := Message{
		ID:                  12345678,
		SentDate:            time.Date(2021, time.October, 1, 18, 27, 0, 0, time.Local),
		Sender:              "Jon Doe (Director)",
		Subject:             "SOME SUBJECT",
		Body:                "A message with some HTML entities&nbsp; and <div>markup</div>",
		ContainsAttachments: true,
		Attachments:         []Attachment{{ID: 123456, FileName: "Some File.ext"}},
		ReadDate:            time.Date(2021, time.October, 2, 19, 3, 00, 00, time.Local),
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
	w.Write([]byte(testResp))
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
