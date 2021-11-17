package listen

import (
	"flag"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Run(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	token := flags.String("bot_token", "", "token for bot")
	flags.Parse(args)

	if token == nil || *token == "" {
		log.Println("expected non empty bot token")
		return
	}

	discord, err := discordgo.New("Bot " + *token)
	_ = discord
	_ = err
}
