package notifier

import "github.com/volmedo/almendruco.git/internal/raices"

type ChatID uint64

type Notifier interface {
	Notify(chatID ChatID, msgs []raices.Message) (uint64, error)
}
