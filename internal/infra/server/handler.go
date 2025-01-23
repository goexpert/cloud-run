package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/goexpert/cloud-run/internal/entity"
	"github.com/goexpert/cloud-run/internal/usecase"
	lab "github.com/goexpert/labobservabilidade"
)

func GetWeatherViaCepHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	w.Header().Add("Content-Type", "application/json")

	httpClient := http.DefaultClient

	cepDto, err := lab.NewCep(r.PathValue("cep"))
	if err != nil {

		code := http.StatusInternalServerError
		if err.Error() == "cep deve ter 8 dígitos numéricos" {
			code = http.StatusNotFound
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&lab.DtoError{Message: err.Error()})
		return
	}

	addressDto, err := usecase.GetLogradouro(ctx, *cepDto, httpClient)
	if err != nil {

		code := http.StatusInternalServerError
		message := err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			code = http.StatusNotFound
			message = "CEP não encontrado"
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&lab.DtoError{Message: message})
		return
	}

	weatherDto, err := usecase.GetWeather(ctx, *addressDto, httpClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&lab.DtoError{Message: err.Error()})
		return
	}

	localeWeatherDto, err := entity.NewLocaleWeather(addressDto.Localidade, weatherDto.Current.TempC)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&lab.DtoError{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(localeWeatherDto)
}
