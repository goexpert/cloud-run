package usecase

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	lab "github.com/goexpert/labobservabilidade"
)

func GetLogradouro(ctx context.Context, cep lab.CepDto, client *http.Client) (*lab.LogradouroDto, error) {

	_url := "https://viacep.com.br/ws/" + cep.Cep + "/json/"
	// _url := "https://opencep.com/v1/" + cep.Cep + ".json"

	var logradouro lab.LogradouroDto

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	response, err := netClient.Get(_url)
	if err != nil {
		slog.Error("erro na requisição", "error", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &logradouro)
	if err != nil {
		slog.Error("erro no unmarshal do cep", "error", err.Error())
	}

	// wClient, err := lab.NewWebclient(ctx, client, http.MethodGet, _url, nil)
	// if err != nil {
	// 	slog.Error("falha na req para o OpenCep", "error", err.Error())
	// 	return nil, err
	// }

	// err = wClient.Do(func(p []byte) error {
	// 	err = json.Unmarshal(p, &logradouro)
	// 	if err != nil {
	// 		slog.Error("erro no unmarshal do cep", "error", err.Error())
	// 	}
	// 	return err
	// })
	// if err != nil {
	// 	slog.Error("executa webclient", "error", err.Error())

	// }

	// if logradouro.Erro != "" {
	// 	return nil, errors.New("cep não encontrado")
	// }

	return &logradouro, err
}
