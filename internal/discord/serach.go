package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) Search(s *discordgo.Session, i *discordgo.InteractionCreate) {

	//ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	query := i.ApplicationCommandData().Options[0].StringValue()
	d.service.Search(context.TODO(), query)
}
