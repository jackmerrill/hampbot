package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/commands/fun"
	studentlife "github.com/jackmerrill/hampbot/internal/commands/studentlife"
	util "github.com/jackmerrill/hampbot/internal/commands/util"
	"github.com/zekroTJA/shireikan"

	"github.com/charmbracelet/log"
)

func main() {
	token := os.Getenv("TOKEN")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	log.Info("Starting bot...")

	err = session.Open()
	if err != nil {
		panic(err)
	}

	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}()

	handler := shireikan.New(&shireikan.Config{
		GeneralPrefix:         ">",
		AllowBots:             false,
		AllowDM:               true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: true,
		OnError: func(ctx shireikan.Context, typ shireikan.ErrorType, err error) {
			log.Error(err)
		},
	})

	log.Info("Initializing Random")
	rand.Seed(time.Now().UnixNano())

	log.Info("Registering commands...")

	handler.Register(&util.Ping{})
	log.Debug("Registered ping command")

	handler.Register(&studentlife.Laundry{})
	log.Debug("Registered laundry command")

	handler.Register(&util.Steal{})
	log.Debug("Registered steal command")

	handler.Register(&fun.AI{})
	log.Debug("Registered ai command")

	handler.Register(&fun.XKCD{})
	log.Debug("Registered xkcd command")

	log.Info("Registered all commands")

	handler.Setup(session)

	log.Info("Bot is now running. Press CTRL-C to exit.")
}
