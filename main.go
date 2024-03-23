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

type ApiResponse struct {
	AddressData
	Source       string
	ResponseTime time.Duration
}

func fetchUrl(url, source string, ch chan<- ApiResponse) {
	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição para %s: %v\n", source, err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler a resposta da %s: %v\n", source, err)
		return
	}

	var address AddressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		fmt.Printf("Erro ao decodificar a resposta da %s: %v\n", source, err)
		return
	}

	timeResponse := time.Since(startTime)
	ch <- ApiResponse{AddressData: address, Source: source, ResponseTime: timeResponse}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanln(&cep)

	brasilApi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	results := make(chan ApiResponse)

	go fetchUrl(brasilApi, "Brasil API", results)
	go fetchUrl(viaCep, "Via CEP", results)

	timeout := time.After(1 * time.Second)

	select {
	case res := <-results:
		fmt.Printf("Resposta mais rápida de %s em %v\n", res.Source, res.ResponseTime)
		fmt.Println("Endereço:")
		fmt.Printf("Logradouro: %s\nBairro: %s\nCidade: %s\nEstado: %s\nCEP: %s\n",
			res.Logradouro, res.Bairro, res.Localidade, res.UF, res.CEP)
	case <-timeout:
		fmt.Println("Erro de timeout. Nenhuma resposta foi recebida em 1 segundo.")
	}
}
