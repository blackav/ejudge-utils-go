package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/blackav/ejudge-utils-go/internal/app"
	"github.com/blackav/ejudge-utils-go/pkg/ejclient"
	"github.com/blackav/ejudge-utils-go/pkg/ejudge"
	"github.com/blackav/ejudge-utils-go/pkg/sber"
	"github.com/blackav/ejudge-utils-go/pkg/slogt"
	"github.com/ilyakaznacheev/cleanenv"
)

type EjudgeConfig struct {
	URL       string `env:"URL" env-required:"true"`
	ContestID int32  `env:"CONTEST_ID" env-required:"true"`
	Token     string `env:"TOKEN" env-required:"true"`
}

type SberConfig struct {
	AuthURL  string `env:"AUTH_URL" env-default:"https://ngw.devices.sberbank.ru:9443/api/v2/oauth"`
	GenURL   string `env:"GEN_URL" env-default:"https://api.giga.chat"`
	AuthKey  string `env:"AUTH_KEY"`
	CertFile string `env:"CERT_FILE"`
	Model    string `env:"MODEL" env-default:"GigaChat-3-Ultra"`
}

type Config struct {
	Ejudge       EjudgeConfig `env-prefix:"EJUDGE_"`
	Sber         SberConfig   `env-prefix:"SBER_"`
	TemplateFile string       `env:"TEMPLATE_FILE" env-default:"template.mustache"`
}

func main() {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("failed to parse environment", slogt.Error(err))
		os.Exit(1)
	}
	flag.Parse()

	ejc, err := ejudge.New(ejclient.Config{
		URL:       cfg.Ejudge.URL,
		ContestID: cfg.Ejudge.ContestID,
		Token:     cfg.Ejudge.Token,
	})
	if err != nil {
		slog.Error("failed to initialize ejudge client", slogt.Error(err))
		os.Exit(1)
	}

	ctx := context.Background()

	gen, err := sber.New(ctx, &sber.Config{
		AuthURL:  cfg.Sber.AuthURL,
		GenURL:   cfg.Sber.GenURL,
		AuthKey:  cfg.Sber.AuthKey,
		CertFile: cfg.Sber.CertFile,
		Model:    cfg.Sber.Model,
	})
	if err != nil {
		slog.Error("failed to initialize sber client", slogt.Error(err))
		os.Exit(1)
	}

	// we check that the connection works by counting the available models
	models, err := gen.ListModels(ctx)
	if err != nil {
		slog.Error("failed to get the list of models", slogt.Error(err))
		os.Exit(1)
	}
	chatCount := 0
	for i := range models {
		if models[i].Type == "chat" {
			chatCount++
		}
	}
	if chatCount == 0 {
		slog.Error("no chat models found")
		os.Exit(1)
	}
	slog.Info("chat models found", slog.Int("count", chatCount))

	a := app.New(app.Config{
		Ejudge:       ejc,
		Generator:    gen,
		ID:           "The ejudge reviewer",
		TemplateFile: cfg.TemplateFile,
	})

	a.Run(ctx)
}
