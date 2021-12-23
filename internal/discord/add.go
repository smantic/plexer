package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/pkg/radarr"
)

func (d *Discord) Add(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	// start with a hidden resposne
	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
		},
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:

		rootFolderInfo := d.getRadarrRootFolder(ctx)

		title := i.ApplicationCommandData().Options[0].StringValue()
		title = strings.TrimSpace(title)

		movies, err := d.service.Search(ctx, title)
		if err != nil {
			return fmt.Errorf("failed to search for movie to add: %w", err)
		}

		if len(movies) == 0 {
			response.Data.Content = "couldn't find a movie to add like: " + title
			break
		}

		var m radarr.Movie
		for _, movie := range movies {
			if movie.Title == title {
				m = movie
				break
			}
		}

		if len(m.Path) > 0 {
			response.Data.Content = title + " is already added! "
			break
		}

		info := <-rootFolderInfo
		m.RootFolderPath = info.Path

		if int(info.FreeSpace) < m.SizeOnDisk {
			response.Data.Content = fmt.Sprintf("not enough space on disk!! only %d space left", info.FreeSpace)
			break
		}

		err = d.service.Add(ctx, m)
		if err != nil {
			return fmt.Errorf("failed to add title: %w", err)
		}

		response.Data.Content = "added: " + title

	case discordgo.InteractionApplicationCommandAutocomplete:
		query := i.ApplicationCommandData().Options[0].StringValue()
		results, err := d.service.Search(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to search for auto completes: %w", err)
		}

		choices := getAutoCompleteChoicesFrom(results)
		response.Type = discordgo.InteractionApplicationCommandAutocompleteResult
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
