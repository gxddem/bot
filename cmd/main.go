package main

import (
	"foxhole-bot/internal/handlers"
	"foxhole-bot/internal/services"
	"foxhole-bot/internal/storage"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("TOKEN")
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("failed to init discord session; err: %v", err)
	}
	if err := client.Open(); err != nil {
		log.Fatalf("failed to open ws connection; err: %v", err)
	}

	storage := storage.NewStorage()
	service := services.NewPanelService(storage)
	handlers.NewPanelController(service, client)

	select {}
}
