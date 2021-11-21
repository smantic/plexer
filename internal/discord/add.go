package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/pkg/radarr"
)

func (d *Discord) Add(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		},
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		title := i.ApplicationCommandData().Options[0].StringValue()

		results, err := d.service.Search(ctx, title)
		if err != nil {
			return fmt.Errorf("failed to search for movie to add: %w", err)
		}

		if len(results) == 0 {
			data := discordgo.InteractionResponseData{
				Content: "couldn't find a movie to add like: " + title,
			}
			response = discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &data,
			}
			break
		}

		movie := results[0]
		if movie.Added != "" {
			response.Data.Content = title + " is already added! "
			break
		}

		err = d.service.Add(ctx, results[0])
		if err != nil {
			return fmt.Errorf("failed to add title: %w", err)
		}

		response.Data.Content = title

	case discordgo.InteractionApplicationCommandAutocomplete:
		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to search for auto completes: %w", err)
		}

		choices := getAutoCompleteChoicesFrom(results)
		response.Data = &discordgo.InteractionResponseData{
			Choices: choices,
		}
	}

	return s.InteractionRespond(i.Interaction, &response)
}

func getAutoCompleteChoicesFrom(movies []radarr.Movie) []*discordgo.ApplicationCommandOptionChoice {

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(movies))
	for i, m := range movies {
		c := discordgo.ApplicationCommandOptionChoice{
			Name:  m.Title,
			Value: m.Title,
		}
		choices = append(choices, &c)

		if i >= 7 {
			break
		}
	}
	return choices
}
