package listen

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/smantic/plexer/internal/discord"
	"github.com/smantic/plexer/internal/service"
	"github.com/smantic/plexer/pkg/radarr"
	"github.com/webtor-io/go-jackett"
)

type Config struct {
	DiscordToken string
	RadarrURL    string
	RadarrKey    string
	JackettURL   string
	JackettKey   string
}

func Run(args []string) {

	c := Config{}

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&c.DiscordToken, "token", "", "token for the discord bot")
	flags.StringVar(&c.RadarrURL, "radarURL", "https://localhost:7878/api/v3", "url of radar service")
	flags.StringVar(&c.RadarrKey, "radarrKey", "", "radarr api key")
	flags.StringVar(&c.JackettURL, "jackett", "", "url of jacket service")
	flags.StringVar(&c.JackettKey, "jackettKey", "", "jackett api key")

	flags.Parse(args)

	if c.DiscordToken == "" {
		log.Println("expected non empty bot token")
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	svc := service.Service{
		Radarr: &radarr.Client{
			BaseURL: c.RadarrURL,
			Apikey:  c.RadarrKey,
		},
		Jackett: jackett.NewJackett(&jackett.Settings{
			ApiURL: c.JackettURL,
			ApiKey: c.JackettKey,
		}),
	}

	s, err := discord.NewSession(c.DiscordToken, &svc)
	if err != nil {
		log.Printf("failed to get bot: %w", err)
		return
	}
	s.Init(ctx)

	<-ctx.Done()
	log.Println("shutting down...")

	s.Close()

	cancel()
}
