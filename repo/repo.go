package repo

//go:generate mockery --case underscore --inpkg --name Repo
type Repo interface {
	GetUserData() (UserData, error)
	GetLastNotifiedMessage() (uint64, error)
	SetLastNotifiedMessage(uint64) error
}

type UserData struct {
	User     string
	Password string
}
