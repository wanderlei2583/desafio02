package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type response struct {
	body         string
	source       string
	responseTime time.Duration
	Address      addressData
}

type addressData struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	CEP        string `json:"cep"`
}

func fetchUrl(url string, source string, ch chan<- response, startTime time.Time) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição de", source, ":", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erroa ao ler a resposta de", source, ":", err)
		return
	}

	var address addressData
	err = json.Unmarshal(body, &address)
	if err != nil {
		fmt.Println("Erro ao fazer o unmarshal da resposta de", source, ":", err)
		return
	}

	timeResponse := time.Since(startTime)

	ch <- response{body: string(body), source: source, responseTime: timeResponse, Address: address}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanln(&cep)

	brasilApi := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ch := make(chan response)
	startTime := time.Now()

	go fetchUrl(brasilApi, "Brasil API", ch, startTime)
	go fetchUrl(viaCep, "Via CEP", ch, startTime)

	timeout := time.After(1 * time.Second)

	select {
	case res := <-ch:
		fmt.Printf("Saída recebida da API: %s em %v: %s\n", res.source, res.responseTime)
		fmt.Println("Endereço:")
		fmt.Printf("Logradouro:", res.Address.Logradouro)
		fmt.Printf("Bairro:", res.Address.Bairro)
		fmt.Printf("Cidade:", res.Address.Localidade)
		fmt.Printf("Estado:", res.Address.UF)
		fmt.Printf("CEP:", res.Address.CEP)
	case <-timeout:
		fmt.Println("Erro de TimeOut: Nenhuma resposta foi recebida em até 1 segundo.")
	}
}
