package discord

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/libs/discord/imux"
	"github.com/smantic/plexer/internal/service"
)

func (d *Discord) DiskSpace(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	freeSpace, err := d.service.GetTotalFreeSpace(request.Context)
	if err != nil {
		err = fmt.Errorf("failed to get total free space: %w\n", err)
		respondWithErr(response, request, err)
		return
	}

	inMB := float64(freeSpace) / float64(1000000)
	response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Description: "",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Free Space",
							Value: strconv.FormatFloat(inMB, 'f', 3, 64) + " MB",
						},
					},
				},
			},
		},
	}

	imux.Respond(response, request)
}

func (d *Discord) Queue(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	q, err := d.service.GetQueue(request.Context)
	if err != nil {
		err = fmt.Errorf("failed to get queue: %w\n", err)
		respondWithErr(response, request, err)
	}

	if len(q) == 0 {
		response.Data.Content = "nothing in the queue!"
	}

	for _, i := range q {
		response.Data.Embeds = append(response.Data.Embeds, queueItemAsEmbed(i))
	}

	imux.Respond(response, request)
	return
}

func queueItemAsEmbed(i service.QueueItem) *discordgo.MessageEmbed {

	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Description: "",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  i.Title,
				Inline: true,
			},
			{
				Name:   "Type",
				Value:  string(i.ContentType),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  i.Status,
				Inline: true,
			},
			{
				Name:   "Size",
				Value:  strconv.Itoa(int(i.Size)),
				Inline: true,
			},
			{
				Name:   "Quality",
				Value:  strconv.Itoa(i.Quality),
				Inline: true,
			},
			{
				Name:   "Time Left",
				Value:  i.TimeLeft,
				Inline: true,
			},
		},
	}
}
