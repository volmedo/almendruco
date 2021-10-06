package raices

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestLogin(t *testing.T) {

}

func TestFetchMessages(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)

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

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testResp)
	}))
	defer svr.Close()

	c, err := NewClient(svr.URL)
	if err != nil {
		t.Fatalf("Unable to create client: %s", err)
	}

	msgs, err := c.FetchMessages()
	if err != nil {
		t.Fatalf("Unexpected error fetching messages: %s", err)
	}

	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(msgs))
	}

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
