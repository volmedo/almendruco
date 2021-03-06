package raices

import "time"

type loginResponse struct {
	Status status `json:"ESTADO"`
}

type status struct {
	Code        string `json:"CODIGO"`
	Description string `json:"DESCRIPCION,omitempty"`
}

const statusCodeOK string = "C"

type messagesResponse struct {
	Status   status       `json:"ESTADO"`
	Messages []rawMessage `json:"RESULTADO"`
}

type rawMessage struct {
	ID                  uint64          `json:"X_NOTMENSAL"`
	SentDate            string          `json:"F_ENVIO"`
	Sender              string          `json:"REMITIDO"`
	Subject             string          `json:"T_ASUNTO"`
	Body                string          `json:"T_MENSAJE"`
	ContainsAttachments string          `json:"L_ADJUNTO"`
	Attachments         []rawAttachment `json:"ADJUNTOS"`
	ReadDate            string          `json:"F_LECTURA"`
}

type rawAttachment struct {
	ID       uint64 `json:"X_ADJMENSAL"`
	FileName string `json:"T_NOMFIC"`
	Contents []byte
}

type Message struct {
	ID                  uint64
	SentDate            time.Time
	Sender              string
	Subject             string
	Body                string
	ContainsAttachments bool
	Attachments         []Attachment
	ReadDate            time.Time
}

type Attachment struct {
	ID       uint64
	FileName string
	Contents []byte
}

const dateFormat = "02/01/2006 15:04"

func parseMessage(rm rawMessage) (Message, error) {
	// Time strings reported by Raices are always CET/CEST
	cet, err := time.LoadLocation("CET")
	if err != nil {
		return Message{}, err
	}

	var sentDate time.Time
	if rm.SentDate != "" {
		sentDate, err = time.ParseInLocation(dateFormat, rm.SentDate, cet)
		if err != nil {
			return Message{}, err
		}
	}

	var readDate time.Time
	if rm.ReadDate != "" {
		readDate, err = time.ParseInLocation(dateFormat, rm.ReadDate, cet)
		if err != nil {
			return Message{}, err
		}
	}

	attachments := make([]Attachment, len(rm.Attachments))
	for i, ra := range rm.Attachments {
		contents := make([]byte, len(ra.Contents))
		copy(contents, ra.Contents)
		attachments[i] = Attachment{
			ID:       ra.ID,
			FileName: ra.FileName,
			Contents: contents,
		}
	}

	return Message{
		ID:                  rm.ID,
		SentDate:            sentDate,
		Sender:              rm.Sender,
		Subject:             rm.Subject,
		Body:                rm.Body,
		ContainsAttachments: rm.ContainsAttachments == "S",
		Attachments:         attachments,
		ReadDate:            readDate,
	}, nil
}
