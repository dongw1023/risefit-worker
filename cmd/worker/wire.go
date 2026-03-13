//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/risefit/email-worker/pkg/config"
	"github.com/risefit/email-worker/pkg/email"
	"github.com/risefit/email-worker/pkg/handler"
)

type Server struct {
	Handler *handler.EmailHandler
	Config  *config.Config
}

func ProvideEmailProvider(cfg *config.Config) email.Provider {
	return email.NewSendGridProvider(cfg.EmailProviderAPIKey, cfg.FromEmail)
}

func InitializeServer() (*Server, error) {
	wire.Build(
		config.Load,
		ProvideEmailProvider,
		handler.NewEmailHandler,
		wire.Struct(new(Server), "*"),
	)
	return nil, nil
}
