package main

type config struct {
	Raices   RaicesConfig
	Telegram TelegramConfig
}

type RaicesConfig struct {
	BaseURL string `default:"https://raices.madrid.org"`
}

type TelegramConfig struct {
	BaseURL  string `default:"https://api.telegram.org"`
	BotToken string `required:"true"`
}
