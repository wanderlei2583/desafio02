package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type response struct {
	body         string
	source       string
	responseTime time.Duration
}

func fetchUrl(url string, source string, ch chan<- response, startTime time.Time) {
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

	timeResponse := time.Since(startTime)

	ch <- response{body: string(body), source: source, responseTime: timeResponse}
}

func main() {
	cep := "01001000"
	brasilApi := fmt.Sprintf("https://brsilapi.com.br/api/cep/v1/%s", cep)
	viaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	ch := make(chan response)
	startTime := time.Now()

	go fetchUrl(brasilApi, "Brasil API", ch, startTime)
	go fetchUrl(viaCep, "Via CEP", ch, startTime)

	timeout := time.After(1 * time.Second)

	select {
	case res := <-ch:
		fmt.Printf("Saída recebida da API: %s em %v: %s\n", res.source, res.responseTime, res.body)
	case <-timeout:
		fmt.Println("Erro de TimeOut: Nenhuma resposta foi recebida em até 1 segundo.")
	}
}
