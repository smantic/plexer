package discord

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/internal/service"
)

func (d *Discord) DiskSpace(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	}

	diskSpaceChan := d.service.GetDiskSpaceInfo(ctx)
	dSpace := <-diskSpaceChan
	if dSpace.Err != nil {
		return dSpace.Err
	}

	response.Data.Embeds = []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Description: "",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Total Space",
					Value: strconv.Itoa(dSpace.TotalCapacity),
				},
				{
					Name:  "Used Space",
					Value: strconv.Itoa(dSpace.UsedCapacity),
				},
				{
					Name:  "Free Space",
					Value: strconv.Itoa(dSpace.FreeSpace),
				},
			},
		},
	}

	return s.InteractionRespond(i.Interaction, &response)
}

func (d *Discord) Queue(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var response discordgo.InteractionResponse = discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	}

	qChan := d.service.GetQueue(ctx)
	queue := <-qChan
	if queue.Err != nil {
		return queue.Err
	}

	data := response.Data
	for _, i := range queue.Queue {
		data.Embeds = append(data.Embeds, queueItemAsEmbed(i))
	}

	return s.InteractionRespond(i.Interaction, &response)
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
				Value:  i.ContentType,
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  i.Status,
				Inline: true,
			},
			{
				Name:   "Size",
				Value:  strconv.Itoa(i.Size),
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
			{
				Name:   "Estimated Completion",
				Value:  i.EstimatedCompletionTime,
				Inline: true,
			},
		},
	}
}
