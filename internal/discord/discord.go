package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/smantic/libs/discord/imux"
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
					Type:         discordgo.ApplicationCommandOptionSubCommand,
					Name:         "movie",
					Description:  "add a movie",
					Required:     true,
					Autocomplete: true,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "title",
							Description:  "title of the movie that you want to add",
							Required:     true,
							Autocomplete: true,
						},
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "quality",
							Description:  "quality to download the content in",
							Required:     false,
							Autocomplete: true,
						},
					},
				},
				{
					Type:         discordgo.ApplicationCommandOptionSubCommand,
					Name:         "show",
					Description:  "add a show",
					Required:     true,
					Autocomplete: true,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "title",
							Description:  "title the show that you want to add",
							Required:     true,
							Autocomplete: true,
						},
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "season",
							Description:  "season to add",
							Required:     true,
							Autocomplete: true,
						},
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "quality",
							Description:  "quality to download the content in",
							Required:     false,
							Autocomplete: true,
						},
					},
				},
			},
		},
		{
			Name:        "search",
			Description: "search for a movie or show",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "title",
					Description:  "title of the content you want to find",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "ping",
			Description: "ping the bot",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "queue",
			Description: "see downloads in the queue",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "du",
			Description: "see disk usage",
			Type:        discordgo.ChatApplicationCommand,
		},
	}
)

// Discord is a service struct for handling discord commands.
type Discord struct {
	token   string
	session *discordgo.Session
	service *service.Service
	i       *imux.InteractionMux
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
func (d *Discord) Init(ctx context.Context, refresh bool, skip bool) error {

	i := imux.NewInteractionMux(ctx)
	i.Use(imux.LogInteraction)

	i.Add("add movie", d.AddMovie)
	i.Add("add show", d.AddShow)
	i.Add("search", d.Search)
	i.Add("du", d.DiskSpace)
	i.Add("queue", nil)
	i.Add("ping", Ping)

	d.session.AddHandler(i.Serve())
	d.session.AddHandler(d.Connected(ctx, refresh, skip))

	err := d.session.Open()
	if err != nil {
		return fmt.Errorf("failed to open discord ws: %w", err)
	}

	return nil
}

type readyHandler = func(s *discordgo.Session, r *discordgo.Ready)

func (d *Discord) Connected(ctx context.Context, refresh bool, skip bool) readyHandler {

	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("connected to: %s\n", r.User.String())

		if refresh {
			// if refresh then delete all of bot's existing commands

			existing, err := d.session.ApplicationCommands(d.session.State.User.ID, "")
			if err != nil {
				log.Printf("failed to get existing commands: %v", err)
				return
			}

			log.Printf("cleaning old commands")
			for _, e := range existing {
				if e.ApplicationID != d.session.State.User.ID {
					continue
				}

				err := d.session.ApplicationCommandDelete(d.session.State.User.ID, "", e.ID)
				if err != nil {
					log.Printf("failed to delete command %v: %v", e, err)
					return
				}
			}
		}

		if skip {
			return
		}

		for _, v := range commands {
			_, err := d.session.ApplicationCommandCreate(d.session.State.User.ID, "", v)
			if err != nil {
				log.Printf("failed to register command: %v: err: %v \n", v, err)
				continue
			}
			log.Printf("registered: %s\n", v.Name)
		}
	}
}

func Ping(response *discordgo.InteractionResponse, request *imux.InteractionRequest) {

	response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "pong!"},
	}

	imux.Respond(response, request)
}
