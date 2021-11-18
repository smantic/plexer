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

func Run(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	token := flags.String("token", "", "token for the discord bot")
	radarURL := flags.String("radarURL", "https://localhost:7878/api/v3", "url of radar service")
	flags.Parse(args)
	_ = radarURL

	if token == nil || *token == "" {
		log.Println("expected non empty bot token")
		return
	}

	ctx := context.Background()

	radarr := radarr.Client{}

	deps := service.Dependencies{
		Radarr: &radarr,
	}
	svc, err := service.NewService(&deps)
	if err != nil {
		log.Println(err)
		return
	}

	d := discord.NewSession(*token, &svc)
	d.Init(ctx)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("shutting down...")
}
