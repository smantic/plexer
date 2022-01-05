package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/libs/discord/imux"
	"github.com/smantic/plexer/internal/service"
)

const (
	SEARCH_RESULTS_LIMIT int = 10
)

func (d *Discord) AddMovie(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	ctx := request.Context
	cmd := request.Interaction.ApplicationCommandData()

	switch cmd.Type() {
	case discordgo.InteractionApplicationCommand:

		var (
			title   string = cmd.Options[0].StringValue()
			quality *discordgo.ApplicationCommandInteractionDataOption
		)

		_ = quality
		if cmd.Options[1] == nil {
			quality = cmd.Options[1]
		}

		content, err := d.service.Search(ctx, service.CONTENT_MOVIE, title, SEARCH_RESULTS_LIMIT)
		if err != nil {
			err := fmt.Errorf("failed to search for movie to add: %w", err)
			respondWithErr(response, request, err)
			return
		}

		var found service.ContentInfo
		for _, c := range content {
			if c.Title == title {
				found = c
				break
			}
		}

		if len(content) == 0 {
			response.Data.Content = "no title found"
			// log error
			imux.Respond(response, request)
		}

		err = d.service.AddMovie(ctx, found, "")
		if err != nil {
			err := fmt.Errorf("failed to add movie: %w", err)
			respondWithErr(response, request, err)
			return
		}
		return
	case discordgo.InteractionApplicationCommandAutocomplete:

		switch len(cmd.Options) {
		case 0:
			return
		case 1:
			d.searchAutoCompleteAndRespond(service.CONTENT_MOVIE, cmd.Options[0].StringValue(), response, request)
		case 2:
			// todo quality auto complete
		}
	}
}

// AddShow handles the add show command.
// add show title season quality
func (d *Discord) AddShow(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	ctx := request.Context
	cmd := request.Interaction.ApplicationCommandData()

	switch cmd.Type() {
	case discordgo.InteractionApplicationCommand:

		var (
			title   string = cmd.Options[0].StringValue()
			season  string = cmd.Options[1].StringValue()
			quality *discordgo.ApplicationCommandInteractionDataOption
		)

		_ = season
		_ = quality
		if cmd.Options[2] == nil {
			quality = cmd.Options[2]
		}

		content, err := d.service.Search(ctx, service.CONTENT_SHOW, title, SEARCH_RESULTS_LIMIT)
		if err != nil {
			err := fmt.Errorf("failed to search for show to add: %w", err)
			respondWithErr(response, request, err)
			return
		}

		var found service.ContentInfo
		for _, c := range content {
			if c.Title == title {
				found = c
				break
			}
		}

		err = d.service.Add(ctx, found)
		if err != nil {
			err := fmt.Errorf("failed to add show: %w", err)
			respondWithErr(response, request, err)
			return
		}

		response.Data.Content = "added " + title
		return

	case discordgo.InteractionApplicationCommandAutocomplete:

		switch len(cmd.Options) {
		case 0:
			return
		case 1:
			d.searchAutoCompleteAndRespond(service.CONTENT_MOVIE, cmd.Options[0].StringValue(), response, request)
			return
		case 2:
			// todo season auto complete
		case 3:
			// todo  quality auto complete
		}
	}
}

func (d *Discord) searchAutoCompleteAndRespond(
	kind service.ContentType,
	title string,
	response *discordgo.InteractionResponse,
	request *imux.InteractionRequest,
) {

	if len(title) == 0 {
		return
	}

	content, err := d.service.Search(request.Context, kind, title, SEARCH_RESULTS_LIMIT)
	if err != nil {
		err := fmt.Errorf("failed to search for show to add: %w", err)
		respondWithErr(response, request, err)
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(content))
	for _, c := range content {
		choice := discordgo.ApplicationCommandOptionChoice{
			Name:  c.Title,
			Value: c.Title,
		}
		choices = append(choices, &choice)
	}

	response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}

	imux.Respond(response, request)
}

func respondWithErr(response *discordgo.InteractionResponse, request *imux.InteractionRequest, err error) {
	response.Data.Content = err.Error()
	imux.Respond(response, request)
}
