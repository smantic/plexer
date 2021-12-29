package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/internal/service"
)

func (d *Discord) Add(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	// start with a hidden resposne
	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		},
	}

	defer func() {
		// TODO record discord error in context
		err := s.InteractionRespond(i.Interaction, &response)
		if err != nil {
			log.Printf("discord: %s\n", err.Error())
		}
	}()

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		title := strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue())
		content, err := d.service.Search(ctx, "", title, 0)
		if err != nil {
			return fmt.Errorf("failed to search for movie to add: %w", err)
		}

		var found service.ContentInfo
		for _, c := range content {
			if c.Title == title {
				found = c
				break
			}
		}

		if found.Title == "" {
			response.Data.Content = "couldn't find content to add like: " + title
			return nil
		}

		if !found.Added.IsZero() {
			response.Data.Content = title + " is already added! "
			return nil
		}

		err = d.service.Add(ctx, found)
		if err != nil {
			response.Data.Content = fmt.Sprintf("failed to add title: %s", err.Error())
			return err
		}

		response.Data.Content = "added: " + title

	case discordgo.InteractionApplicationCommandAutocomplete:
		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, "", query, 0)
		if err != nil {
			return fmt.Errorf("failed to search for auto completes: %w", err)
		}

		choices := getAutoCompleteChoicesFrom(results)
		response.Type = discordgo.InteractionApplicationCommandAutocompleteResult
		response.Data = &discordgo.InteractionResponseData{
			Choices: choices,
		}
	}

	return nil
}

func getAutoCompleteChoicesFrom(content []service.ContentInfo) []*discordgo.ApplicationCommandOptionChoice {

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(content))
	for i, c := range content {
		choice := discordgo.ApplicationCommandOptionChoice{
			Name:  c.Title,
			Value: c.Title,
		}
		choices = append(choices, &choice)

		if i >= 7 {
			break
		}
	}
	return choices
}
