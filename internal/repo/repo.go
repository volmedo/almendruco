package repo

//go:generate mockery --case underscore --inpkg --name Repo
type Repo interface {
	GetChats() ([]Chat, error)
	UpdateLastNotifiedMessage(chatID ChatID, lastNotifiedMessage uint64) error
}

type ChatID uint64

type Chat struct {
	ID                  ChatID
	Credentials         Credentials
	LastNotifiedMessage uint64
}

type Credentials struct {
	UserName string
	Password string
}
