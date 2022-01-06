package discord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/libs/discord/imux"
	"github.com/smantic/plexer/internal/service"
)

func (d *Discord) Search(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	ctx := request.Context
	cmd := request.Interaction.ApplicationCommandData()

	switch cmd.Type() {
	case discordgo.InteractionApplicationCommand:
		var title = cmd.Options[0].StringValue()

		content, err := d.service.Search(ctx, "", title, SEARCH_RESULTS_LIMIT)
		if err != nil {
			err := fmt.Errorf("failed to search content: %w", err)
			respondWithErr(response, request, err)
			return
		}

		if len(content) == 0 {
			response.Data.Content = "no results for " + title
		}

		var found service.ContentInfo
		for _, c := range content {
			if strings.EqualFold(c.Title, title) {
				found = c
				break
			}
		}

		inMB := float64(found.Size) / float64(1000000)
		response.Data.Embeds = []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Description: "",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Title",
						Value:  fmt.Sprintf("%s (%d)\n", found.Title, found.Year),
						Inline: true,
					},
					{
						Name:   "added",
						Value:  strconv.FormatBool(!found.Added.IsZero()),
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
						Value:  found.ImdbID,
						Inline: true,
					},
					{
						Name:   "Size",
						Value:  strconv.FormatFloat(inMB, 'f', 3, 64) + " MB",
						Inline: true,
					},
					{
						Name:  "Overview",
						Value: found.Overview,
					},
				},
			},
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		d.searchAutoCompleteAndRespond(service.CONTENT_MOVIE, cmd.Options[0].StringValue(), response, request)
		return
	}
}
