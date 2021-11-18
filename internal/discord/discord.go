package discord

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/internal/service"
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

// Discord is a service struct for handling discord commands
type Discord struct {
	token   string
	service *service.Service
}

func NewSession(token string, svc *service.Service) Discord {
	return Discord{
		token:   token,
		service: svc,
	}
}

// Init starts the discord service and adds handlers
func (d *Discord) Init(ctx context.Context) {

	discord, err := discordgo.New("Bot " + d.token)
	if err != nil {
		log.Println(err)
		return
	}

	discord.AddHandler(d.HandleInteraction)
	discord.AddHandler(d.Connected)
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
}

func (d *Discord) Connected(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("connected to: %s\n", r.User.String())
}

func (d *Discord) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	switch name {
	case "search":
		d.Search(s, i)
	case "ping":
		d.Ping(s, i)
		return
	default:
		log.Printf("didn't recognize: %s\n", name)
		return
	}
}

func (d *Discord) Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
