package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/plexer/internal/service"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "add",
			Description: "add a movie or a show",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "title",
					Description:  "title of movie you want to find",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "search",
			Description: "search for a movie or show",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "title of movie you want to find",
					Required:    true,
				},
			},
		},
		{
			Name:        "ping",
			Description: "ping the bot",
			Type:        discordgo.ChatApplicationCommand,
		},
	}
)

// Discord is a service struct for handling discord commands.
type Discord struct {
	token   string
	session *discordgo.Session
	service *service.Service
}

// NewSession creates a new session.
func NewSession(token string, svc *service.Service) (Discord, error) {

	discord, err := discordgo.New("Bot " + token)
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
	err := d.session.Close()
	if err != nil {
		log.Println(err)
	}
}

// Init starts the discord service and adds handlers.
// if refresh is true we will delete all the old commands and re-add them.
func (d *Discord) Init(ctx context.Context, refresh bool) error {

	d.session.AddHandler(d.HandleInteraction)
	d.session.AddHandler(d.Connected)

	err := d.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord ws: %w", err)
	}

	if refresh {
		existing, err := d.session.ApplicationCommands(d.session.State.User.ID, "")
		if err != nil {
			return fmt.Errorf("failed to get existing commands: %w", err)
		}
		log.Printf("cleaning old commands")
		for _, e := range existing {
			err := d.session.ApplicationCommandDelete(d.session.State.User.ID, "", e.ID)
			if err != nil {
				return fmt.Errorf("failed to delete command %v: %w", e, err)
			}
		}
	}

	for _, v := range commands {
		_, err := d.session.ApplicationCommandCreate(d.session.State.User.ID, "", v)
		if err != nil {
			return fmt.Errorf("failed to register command: %v: err: %w \n", v, err)
		}
		log.Printf("registered: %s\n", v.Name)
	}

	return nil
}

func (d *Discord) Connected(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("connected to: %s\n", r.User.String())
}

func (d *Discord) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	ctx := context.Background()
	name := i.ApplicationCommandData().Name
	log.Printf("received command: %s\n", name)
	switch name {
	case "add":
		err := d.Add(ctx, s, i)
		if err != nil {
			log.Println(err)
		}
		return
	case "search":
		err := d.Search(ctx, s, i)
		if err != nil {
			log.Println(err)
		}
		return
	case "ping":
		d.Ping(s, i)
		return
	default:
		log.Printf("didn't recognize command: %s\n", name)
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
