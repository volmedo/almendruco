package repo

//go:generate mockery --case underscore --inpkg --name Repo
type Repo interface {
	GetPassword(userName string) (string, error)
	GetLastNotifiedMessage(userName string) (uint64, error)
	SetLastNotifiedMessage(userName string, id uint64) error
}
