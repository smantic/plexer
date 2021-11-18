package listen

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

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

func Run(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	token := flags.String("token", "", "token for the discord bot")
	flags.Parse(args)

	if token == nil || *token == "" {
		log.Println("expected non empty bot token")
		return
	}

	discord, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Println(err)
		return
	}

	discord.AddHandler(HandleInteraction)
	discord.AddHandler(Connected)
	err = discord.Open()
	if err != nil {
		log.Println(err)
		return
	}

	defer discord.Close()
	for _, v := range commands {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			log.Printf("failed to register command: %v: err: %v \n", v, err)
		}
		log.Printf("registered: %s\n", v.Name)
	}
	log.Printf("listening...")

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("shutting down...")
}

func Connected(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("connected to: %s\n", r.User.String())
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	switch name {
	case "ping":
		Ping(s, i)
		return
	default:
		log.Printf("didn't recognize: %s\n", name)
		return
	}
}

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("pinged by %s\n", i.Member.User.ID)

	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "pong!"},
	}
	err := s.InteractionRespond(i.Interaction, &response)
	if err != nil {
		log.Println(err)
	}
}
