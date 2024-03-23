package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type response struct {
	body   string
	source string
}

func fetchUrl(url string, source string, ch chan<- response) {
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
	cep := "01001000"
	brasilApi := fmt.Sprintf("https://brsilapi.com.br/api/cep/v1/%s", cep)
	viaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ch := make(chan response, 1)
	go fetchUrl(brasilApi, "Brasil API", ch)
	go fetchUrl(viaCep, "Via CEP", ch)

	timeout := time.After(1 * time.Second)

	select {
	case res := <-ch:
		fmt.Printf("Saída recebida da API: %s: %s\n", res.source, res.body)
	case <-timeout:
		fmt.Println("Tempo escedido: Nenhuma resposta foi recebida em até 1 segundo.")
	}
}
