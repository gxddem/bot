package main

import (
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

	client.AddHandler(interactionCreate)
	if err := registerCommands(client); err != nil {
		log.Fatalf("failed to register commands; err: %v", err)
	}

	select {}
	
}

func registerCommands(s *discordgo.Session) error {
	command := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "healthcheck",
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
	return err
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		data := i.ApplicationCommandData()
		switch data.Name {
		case "ping":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "ping",
				},
			})
		}
	}
}
