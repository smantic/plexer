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
	session *discordgo.Session
	service *service.Service
}

func NewSession(token string, svc *service.Service) (Discord, error) {

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return Discord{}, err
	}

	err = discord.Open()
	if err != nil {
		return Discord{}, err
	}

	return Discord{
		token:   token,
		session: discord,
		service: svc,
	}, nil
}

func (d *Discord) Close() {
	d.session.Close()
}

// Init starts the discord service and adds handlers
func (d *Discord) Init(ctx context.Context) {

	d.session.AddHandler(d.HandleInteraction)
	d.session.AddHandler(d.Connected)

	for _, v := range commands {
		_, err := d.session.ApplicationCommandCreate(d.session.State.User.ID, "", v)
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
