package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) Search(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {

	//ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	query := i.ApplicationCommandData().Options[0].StringValue()
	movies, err := d.service.Search(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}

	b := strings.Builder{}
	for i, m := range movies {
		str := fmt.Sprintf("%d - %s %d\n", i+1, m.Title, m.Year)
		b.WriteString(str)
	}

	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: b.String()},
	}
	err = s.InteractionRespond(i.Interaction, &response)
	if err != nil {
		log.Println(err)
	}
}