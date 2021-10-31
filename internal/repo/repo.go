package repo

//go:generate mockery --case underscore --inpkg --name Repo
type Repo interface {
	GetChats() ([]Chat, error)
	UpdateLastNotifiedMessage(chatID string, lastNotifiedMessage uint64) error
}

type Chat struct {
	ID                  string
	Credentials         Credentials
	LastNotifiedMessage uint64
}

type Credentials struct {
	User string
	Pass string
}
