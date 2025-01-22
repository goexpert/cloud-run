package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goexpert/cloud-run/internal/infra/server"
	lab "github.com/goexpert/labobservabilidade"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	webServer := lab.NewServer(os.Getenv("LO_PORT"))
	webServer.AddHandler("GET /cep/{cep}", server.GetWeatherViaCepHandler)
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- webServer.Run()
	}()

	select {
	case <-srvErr:
		slog.Info("Serviço finalizando via <CTRL>+C...")
		return
	case <-ctx.Done():
		slog.Info("Serviço finalizando via interrupção no sistema.")
		stop()
	}
}
