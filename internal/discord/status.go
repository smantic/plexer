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

	freeSpace, err := d.service.GetFreeSpace(ctx)
	if err != nil {
		return err
	}

	response.Data.Embeds = []*discordgo.MessageEmbed{
		{
			Type:        discordgo.EmbedTypeRich,
			Description: "",
			Fields: []*discordgo.MessageEmbedField{
				//{
				//	Name:  "Total Space",
				//	Value: strconv.Itoa(dSpace.TotalCapacity),
				//},
				//{
				//	Name:  "Used Space",
				//	Value: strconv.Itoa(dSpace.UsedCapacity),
				//},
				{
					Name:  "Free Space",
					Value: strconv.Itoa(int(freeSpace.FreeSpace)),
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

	q, err := d.service.GetQueue(ctx)
	if err != nil {
		return err
	}

	data := response.Data
	if len(q) == 0 {
		data.Content = "nothing in the queue!"
	}

	for _, i := range q {
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
