package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/internal/service"
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

		query := strings.TrimSpace(
			i.ApplicationCommandData().Options[0].StringValue(),
		)
		data := response.Data

		content, err := d.service.Search(ctx, "", query, 0)
		if err != nil {
			return err
		}

		if len(content) == 0 {
			data.Content = "could not find " + query
		}

		var c service.ContentInfo
		for _, i := range content {
			if strings.EqualFold(i.Title, query) {
				c = i
				break
			}
		}

		inMB := float64(c.Size) / float64(1000000)
		data.Embeds = []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Description: "",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Title",
						Value:  fmt.Sprintf("%s (%d)\n", c.Title, c.Year),
						Inline: true,
					},
					{
						Name:   "added",
						Value:  strconv.FormatBool(!c.Added.IsZero()),
						Inline: true,
					},

					{
						Name:   "Genre",
						Inline: true,
					},
					{
						Name:   "Available",
						Value:  "true",
						Inline: true,
					},
					{
						Name:   "Imdb",
						Value:  c.ImdbID,
						Inline: true,
					},
					{
						Name:  "Overview",
						Value: strconv.FormatFloat(inMB, 'f', 3, 64),
					},
					{
						Name:  "Size",
						Value: strconv.FormatInt(c.Size, 10),
					},
				},
			},
		}

		if len(c.Genre) > 0 {
			data.Embeds[0].Fields[2].Value = c.Genre[0]
		}

		//if len(c.Images) > 0 {
		//	image := c.Images[0]
		//	data.Embeds = append(data.Embeds, &discordgo.MessageEmbed{
		//		Image: &discordgo.MessageEmbedImage{
		//			URL: image.RemoteUrl,
		//		},
		//		URL:         image.RemoteUrl,
		//		Type:        discordgo.EmbedTypeImage,
		//		Description: "",
		//	})
		//}
	case discordgo.InteractionApplicationCommandAutocomplete:

		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, "", query, 0)
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
