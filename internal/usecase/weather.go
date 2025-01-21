package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/goexpert/cloud-run/internal/dto"
	lab "github.com/goexpert/labobservabilidade"
)

func GetWeather(ctx context.Context, a lab.LogradouroDto, client *http.Client) (*dto.WeatherDto, error) {

	var urlQuery = map[string]string{}
	urlQuery["key"] = os.Getenv("WEATHER_TOKEN")
	urlQuery["q"] = a.Localidade
	urlQuery["aqi"] = "no"

	wClient, err := lab.NewWebclient(ctx, client, http.MethodGet, "https://api.weatherapi.com/v1/current.json", urlQuery)
	if err != nil {
		slog.Error("weatherapi client", "error", err.Error())
		return nil, err
	}

	slog.Debug("header", "Header", wClient.Request().Header)

	var clima dto.WeatherDto

	err = wClient.Do(func(p []byte) error {
		err = json.Unmarshal(p, &clima)
		if err != nil {
			slog.Error("erro body unmarshal", "error", err.Error())
		}
		return err
	})
	if err != nil {
		slog.Error("executa request", "error", err.Error())
		return nil, err

	}
	slog.Debug("reponse", "WeatherResponseDto", clima)
	return &clima, nil
}
