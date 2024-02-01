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
	"github.com/jackmerrill/hampbot/internal/listeners"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
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
		GeneralPrefix:         config.BotPrefix,
		AllowBots:             false,
		AllowDM:               true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: true,
		OnError: func(ctx shireikan.Context, typ shireikan.ErrorType, err error) {
			if typ != shireikan.ErrTypCommandNotFound {
				ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
					Embed:     embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription(err.Error()).MessageEmbed,
					Reference: ctx.GetMessage().Reference(),
				})
				log.Error(err)
			}
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

	handler.Register(&fun.Dalle{})
	log.Debug("Registered dalle command")

	handler.Register(&studentlife.PVTA{})
	log.Debug("Registered PVTA command")

	handler.Register(&studentlife.Abbreviation{})
	log.Debug("Registered abbreviation command")

	err = studentlife.InitAll(session)
	if err != nil {
		log.Error("Failed to initialize laundry notify- skipping.")
	} else {
		handler.Register(&studentlife.LaundryNotify{})
	}

	log.Debug("Registered laundry notify command")

	handler.Register(&studentlife.Where{})
	log.Debug("Registered where command")

	handler.Register(&fun.MetricTime{})
	log.Debug("Registered metrictime command")

	handler.Register(&util.VerifyCommand{})
	go util.StartWebserver(session)
	log.Debug("Registered verify command")

	log.Info("Registered all commands")

	log.Info("Setting up activities...")

	activities := config.Statuses

	go func() {
		for {
			activity := activities[rand.Intn(len(activities))]
			session.UpdateStatusComplex(discordgo.UpdateStatusData{
				Activities: []*discordgo.Activity{&activity},
				Status:     "online",
			})
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Info("Setting up listeners...")

	deleteHandler := &listeners.MessageDeleteListener{}
	editHandler := &listeners.MessageEditListener{}
	reactAddHandler := &listeners.MessageReactAddListener{}
	reactRemoveHandler := &listeners.MessageReactRemoveListener{}

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		config.MessageLog[e.ID] = *e.Message
	})
	session.AddHandler(deleteHandler.Exec)
	session.AddHandler(editHandler.Exec)
	session.AddHandler(reactAddHandler.Exec)
	session.AddHandler(reactRemoveHandler.Exec)

	handler.Setup(session)

	log.Info("Bot is now running. Press CTRL-C to exit.")
}
