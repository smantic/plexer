package listen

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/smantic/plexer/internal/discord"
	"github.com/smantic/plexer/internal/service"
)

type Config struct {
	DiscordToken    string
	RefreshCommands bool
	SkipRegrister   bool

	// Debug will print response bodies to stdout
	Debug bool

	service.Config
}

func Run(args []string) {

	c := Config{}

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&c.DiscordToken, "token", "", "token for the discord bot")
	flags.BoolVar(&c.RefreshCommands, "refresh", false, "delete lingering commands, and re-add them")
	flags.BoolVar(&c.Debug, "debug", false, "print out response bodies")
	flags.BoolVar(&c.SkipRegrister, "skipRegister", false, "skip regerstering commands for faster bot startup")

	flags.StringVar(&c.Config.RadarrURL, "radarURL", "http://localhost:7878/api/v3", "url of radar service")
	flags.StringVar(&c.Config.RadarrKey, "radarrKey", "", "radarr api key")

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

	svc := service.New(&c.Config)

	s, err := discord.NewSession(c.DiscordToken, &svc)
	if err != nil {
		log.Printf("failed to get got: %v", err)
		return
	}
	err = s.Init(ctx, c.RefreshCommands, c.SkipRegrister)
	if err != nil {
		log.Println(err)
	}
	log.Printf("listening...")

	<-ctx.Done()
	log.Println("shutting down...")
	s.Close()
	cancel()
}
