package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/pkg/radarr"
)

func (d *Discord) Search(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	// start with a hidden resposne
	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		},
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		query := i.ApplicationCommandData().Options[0].StringValue()
		query = strings.TrimSpace(query)
		data := response.Data

		movies, err := d.service.Search(ctx, query)
		if err != nil {
			return err
		}

		var m radarr.Movie
		for _, movie := range movies {
			if movie.Title == query {
				m = movie
			}
		}

		if m.Title == "" {
			data.Content = "could not find " + query
		}

		data.Embeds = []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Description: "",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Title",
						Value:  fmt.Sprintf("%s (%d)\n", m.Title, m.Year),
						Inline: true,
					},
					{
						Name:   "added",
						Value:  m.Added,
						Inline: true,
					},

					{
						Name:   "Genre",
						Inline: true,
					},
					{
						Name:   "Available",
						Value:  strconv.FormatBool(m.IsAvailable),
						Inline: true,
					},
					{
						Name:   "Imdb",
						Value:  m.ImdbID,
						Inline: true,
					},
					{
						Name:  "Overview",
						Value: fmt.Sprintf("%s\n", m.Overview),
					},
				},
			},
		}

		if len(m.Genres) > 0 {
			data.Embeds[0].Fields[2].Value = m.Genres[0]
		}

		if len(m.Images) > 0 {
			image := m.Images[0]
			data.Embeds = append(data.Embeds, &discordgo.MessageEmbed{
				Image: &discordgo.MessageEmbedImage{
					URL: image.RemoteUrl,
				},
				URL:         image.RemoteUrl,
				Type:        discordgo.EmbedTypeImage,
				Description: "",
			})
		}
	case discordgo.InteractionApplicationCommandAutocomplete:

		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to search for auto completes: %w", err)
		}

		response.Type = discordgo.InteractionApplicationCommandAutocompleteResult
		response.Data = &discordgo.InteractionResponseData{
			Choices: getAutoCompleteChoicesFrom(results),
		}
	}
	return s.InteractionRespond(i.Interaction, &response)
}
