package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) Search(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var response discordgo.InteractionResponse

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		query := i.ApplicationCommandData().Options[0].StringValue()
		movies, err := d.service.Search(ctx, query)
		if err != nil {
			return err
		}

		content := fmt.Sprintf("%#v\n", movies)
		response = discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: content},
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		rest := combineOpts(i.ApplicationCommandData().Options)

		results, err := d.service.Search(ctx, rest)
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
