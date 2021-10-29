package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/volmedo/almendruco.git/internal/notifier"
	"github.com/volmedo/almendruco.git/internal/raices"
	"github.com/volmedo/almendruco.git/internal/repo"
	"github.com/volmedo/almendruco.git/internal/repo/dynamodbrepo"
)

const appName = "almendruco"

func main() {
	cfg := config{}
	if err := envconfig.Process(appName, &cfg); err != nil {
		log.Fatalf("Configuration processing failed: %s", err)
	}

	r, err := dynamodbrepo.NewRepo()
	if err != nil {
		log.Fatalf("Unable to initialize repository: %s", err)
	}

	rc, err := raices.NewClient(cfg.Raices.BaseURL)
	if err != nil {
		log.Fatalf("Error creating Raíces client: %s", err)
	}

	n, err := notifier.NewTelegramNotifier(cfg.Telegram.BaseURL, cfg.Telegram.BotToken)
	if err != nil {
		log.Fatalf("Error creating notifier: %s", err)
	}

	if err := notifyMessages(r, rc, n); err != nil {
		log.Fatalf("Error notifying messages: %s", err)
	}

	log.Println("Success!")
}

func notifyMessages(r repo.Repo, rc raices.Client, n notifier.Notifier) error {
	chats, err := r.GetChats()
	if err != nil {
		return fmt.Errorf("unable to fetch chats from repo: %s", err)
	}

	for _, c := range chats {
		msgs, err := rc.FetchMessages(c.Credentials, c.LastNotifiedMessage)
		if err != nil {
			return fmt.Errorf("error fetching messages from Raíces: %s", err)
		}

		if len(msgs) == 0 {
			//lint:ignore ST1005 the word "Raíces" is capitalized as it is the name of the application
			return errors.New("Raíces client returned no messages")
		}

		if err := n.Notify(notifier.ChatID(c.ID), msgs[0:5]); err != nil {
			return fmt.Errorf("error notifying messages: %s", err)
		}
	}

	return nil
}
