package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AddressData struct {
	Logradouro string `json:"logradouro,omitempty"`
	Bairro     string `json:"bairro,omitempty"`
	Localidade string `json:"localidade,omitempty"`
	UF         string `json:"uf,omitempty"`
	CEP        string `json:"cep,omitempty"`
}

func fetchUrl(url string, source string) (AddressData, time.Duration, error) {
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return AddressData{}, 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AddressData{}, 0, err
	}

	var address AddressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		return AddressData{}, 0, err
	}

	timeResponse := time.Since(startTime)

	return address, timeResponse, nil
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanln(&cep)

	brasilApi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	results := make(chan string, 2)

	fetchAndProcess := func(url string, source string) {
		address, timeResponse, err := fetchUrl(url, source)
		if err != nil {
			results <- fmt.Sprintf("Erro ao buscar o CEP na API %s: %s", source, err)
			return
		}

		result := fmt.Sprintf("Resposta de %s em %v\nLogradouro: %s\nBairro: %s\nCidade: %s\nEstado: %s\nCEP: %s\n",
			source, timeResponse, address.Logradouro, address.Bairro, address.Localidade, address.UF, address.CEP)
		results <- result
	}

	go fetchAndProcess(brasilApi, "Brasil API")
	go fetchAndProcess(viaCep, "Via CEP")

	timeout := time.After(1 * time.Second)

	for i := 0; i < 2; i++ {
		select {
		case res := <-results:
			fmt.Println(res)
		case <-timeout:
			fmt.Println("Erro de TimeOut, nenhuma resposta foi recebida em até 1 segundo.")
			return
		}
	}
}
