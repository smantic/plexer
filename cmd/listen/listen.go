package listen

import (
	"flag"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Run(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	token := flags.String("token", "", "token for the discord bot")
	flags.Parse(args)

	if token == nil || *token == "" {
		log.Println("expected non empty bot token")
		return
	}

	var (
		commands = []*discordgo.ApplicationCommand{
			{
				Name:        "seach",
				Description: "search for a movie or show",
				Type:        discordgo.ChatApplicationCommand,
				// TODO: auto suggestions
				//Options:     []*discordgo.ApplicationCommand{},
			},
			{
				Name:        "ping",
				Description: "ping the bot",
				Type:        discordgo.ChatApplicationCommand,
				// TODO: auto suggestions
				//Options:     []*discordgo.ApplicationCommand{},
			},
		}
	)

	_ = commands

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Println(err)
		return
	}

	discord.AddHandler(Ping)
}

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("pinged by %s", i.User.ID)

	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "pong!"},
	}
	err := s.InteractionRespond(i.Interaction, &response)
	if err != nil {
		log.Println(err)
	}
}
