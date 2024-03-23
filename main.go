package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Struct para armazenar a resposta e a origem
type response struct {
	body   string
	source string
}

func fetchURL(url string, source string, ch chan<- response) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch <- response{body: string(body), source: source}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP para consulta: ")
	fmt.Scanln(&cep)

	brasilAPI := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCEP := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ch := make(chan response, 1)

	go fetchURL(brasilAPI, "BrasilAPI", ch)
	go fetchURL(viaCEP, "ViaCEP", ch)

	timeout := time.After(1 * time.Second)

	select {
	case res := <-ch:
		fmt.Printf("Resposta recebida mais rápida de %s:\n %s\n\n", res.source, res.body)
	case <-timeout:
		fmt.Println("Erro de timeout. Nenhuma resposta foi recebida em menos de 1 segundo.")
	}
}
