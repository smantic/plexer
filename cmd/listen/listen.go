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
	flags.StringVar(&c.RadarrURL, "radarURL", "http://localhost:7878/api/v3", "url of radar service")
	flags.StringVar(&c.RadarrKey, "radarrKey", "", "radarr api key")

	err := flags.Parse(args)
	if err != nil {
		log.Println(err)
		return
	}

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
	}

	s, err := discord.NewSession(c.DiscordToken, &svc)
	if err != nil {
		log.Printf("failed to get got: %v", err)
		return
	}
	err = s.Init(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Printf("listening...")

	<-ctx.Done()
	log.Println("shutting down...")

	//s.Close()

	cancel()
}
