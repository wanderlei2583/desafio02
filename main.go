package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiResponse struct {
	Body         string
	Source       string
	ResponseTime float64
}

func fetchURL(url string, source string, ch chan<- ApiResponse, done chan<- bool) {
	startTime := time.Now()

	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		done <- true
		fmt.Println("Erro ao fazer a requisição para", source, ":", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		done <- true
		fmt.Println("Erro ao ler a resposta de", source, ":", err)
		return
	}
	body := string(bodyBytes)

	elapsed := time.Since(startTime).Seconds()

	ch <- ApiResponse{Body: body, Source: source, ResponseTime: elapsed}
	done <- true
}

func main() {
	var cep string
	fmt.Print("Digite o CEP para consulta: ")
	fmt.Scanln(&cep)

	brasilAPIURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	results := make(chan ApiResponse)
	done := make(chan bool, 2)

	go fetchURL(brasilAPIURL, "BrasilAPI", results, done)
	go fetchURL(viaCEPURL, "ViaCEP", results, done)

	select {
	case res := <-results:
		fmt.Printf("Resposta rápida recebida de %s em %.2f segundos\n%s\n", res.Source, res.ResponseTime, res.Body)
	case <-time.After(1 * time.Second):
		fmt.Println("Erro de timeout. Nenhuma resposta foi recebida em 1 segundo.")
	}

	<-done
	<-done
}
