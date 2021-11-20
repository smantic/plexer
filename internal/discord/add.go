package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/pkg/radarr"
)

func combineOpts(opts []*discordgo.ApplicationCommandInteractionDataOption) string {

	b := strings.Builder{}
	for _, o := range opts {
		b.WriteString(o.StringValue())
	}
	return b.String()
}

func (d *Discord) Add(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var response discordgo.InteractionResponse

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		movie := i.ApplicationCommandData().Options[0].Value
		log.Printf("add :%v\n", movie)

		data := discordgo.InteractionResponseData{
			Content: "adding: " + i.ApplicationCommandData().Options[0].StringValue(),
		}
		response = discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &data,
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to search for auto completes: %w", err)
		}
		choices := getAutoCompleteChoicesFrom(results)

		response = discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
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
