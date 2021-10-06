package raices

import "time"

type messagesResponse struct {
	Status   status       `json:"ESTADO"`
	Messages []rawMessage `json:"RESULTADO"`
}

type status struct {
	Code string `json:"CODIGO"`
}

type rawMessage struct {
	ID                  uint64       `json:"X_NOTMENSAL"`
	SentDate            string       `json:"F_ENVIO"`
	Sender              string       `json:"REMITIDO"`
	Subject             string       `json:"T_ASUNTO"`
	Body                string       `json:"T_MENSAJE"`
	ContainsAttachments string       `json:"L_ADJUNTO"`
	Attachments         []Attachment `json:"ADJUNTOS"`
	ReadDate            string       `json:"F_LECTURA"`
}

type Attachment struct {
	ID       uint64 `json:"X_ADJMENSAL"`
	FileName string `json:"T_NOMFIC"`
}

const dateFormat = "02/01/2006 15:04"

func parseMessage(rm rawMessage) (Message, error) {
	// Time strings reported by Raices are always CET/CEST
	cet, err := time.LoadLocation("CET")
	if err != nil {
		return Message{}, err
	}

	sentDate, err := time.ParseInLocation(dateFormat, rm.SentDate, cet)
	if err != nil {
		return Message{}, err
	}

	readDate, err := time.ParseInLocation(dateFormat, rm.ReadDate, cet)
	if err != nil {
		return Message{}, err
	}

	return Message{
		ID:                  rm.ID,
		SentDate:            sentDate,
		Sender:              rm.Sender,
		Subject:             rm.Subject,
		Body:                rm.Body,
		ContainsAttachments: rm.ContainsAttachments == "S",
		Attachments:         rm.Attachments,
		ReadDate:            readDate,
	}, nil
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
