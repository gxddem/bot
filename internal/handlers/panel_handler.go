package handlers

import (
	"fmt"
	"foxhole-bot/internal/services"

	"github.com/bwmarrin/discordgo"
)

type PanelController struct {
	service services.PanelService
	client  *discordgo.Session
}

func NewPanelController(service *services.PanelService, client *discordgo.Session) {
	handler := &PanelController{service: *service, client: client}

	handler.RegisterCommands(handler.client)
	handler.client.AddHandler(handler.InteractionCreate)
}

func (pc *PanelController) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		case "init":
			pc.InitializeNewList(s, i)

		case "append":
			pc.AppendItemToList(s, i, data)
		case "remove":
			pc.RemoveItemFromList(s, i, data)
		case "purge":
			pc.PurgeList(s,i)
		}
		
	}
}

func (pc *PanelController) RegisterCommands(s *discordgo.Session) error {
	command := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "healthcheck",
	}
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}

	command = &discordgo.ApplicationCommand{
		Name:        "init",
		Description: "Create a new list instance, that be linked to the channel",
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}

	command = &discordgo.ApplicationCommand{
		Name:        "purge",
		Description: "Removes anything that panel currently contains",
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}
	

	command = &discordgo.ApplicationCommand{
		Name:        "delete",
		Description: "Deletes channel's currently linked list",
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}

	command = &discordgo.ApplicationCommand{
		Name:        "append",
		Description: "Appends an item to the list",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "What to add",
				Required:    true,
			},
		},
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}

	command = &discordgo.ApplicationCommand{
		Name:        "remove",
		Description: "Removes an item from the list by provided index",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "What to delete",
				Required:    true,
			},
		},
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", command)
	if err != nil {
		return err
	}

	return nil
}

func (pc *PanelController) InitializeNewList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ChannelID

	panel, err := pc.service.CreateList(channelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}
	messageContent := "```\n\nItems:\n\n```"
	message, err := s.ChannelMessageSend(i.ChannelID, messageContent)
	if err != nil {
		return
	}
	panel.MessageID = message.ID

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The list has been created",
		},
	})
}

func (pc *PanelController) AppendItemToList(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	channelID := i.ChannelID

	panel, err := pc.service.GetList(channelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}

	text := data.Options[0].StringValue()
	if len(text) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintln("No data provided"),
			},
		})
		return
	}

	if err := pc.service.AppendItemToList(panel, text); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}
	_, err = s.ChannelMessageEdit(i.ChannelID, panel.MessageID, formatList(panel.Items))
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success",
		},
	})
}

func (pc *PanelController) RemoveItemFromList (s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	panel, err := pc.service.GetList(i.ChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}
	if err := pc.service.RemoveItemFromList(panel, int(data.Options[0].IntValue())); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}

	_, err = s.ChannelMessageEdit(panel.ChannelID, panel.MessageID, formatList(panel.Items))
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success",
		},
	})
}

func (pc *PanelController) PurgeList(s *discordgo.Session, i *discordgo.InteractionCreate) {

	panel, err := pc.service.GetList(i.ChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}

	if err := pc.service.PurgeList(panel); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}

	_, err = s.ChannelMessageEdit(panel.ChannelID, panel.MessageID, formatList(panel.Items))
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err),
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Success",
		},
	})
}




func formatList(items []string) string {
	text := "```Items:\n"

	for i, item := range items {
		text += fmt.Sprintf("%d | %s\n", i+1, item)
	}
	text += "```"

	return text
}
