package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApiCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	fmt.Println("Input cep: ")
	var cepIn string
	_, err := fmt.Scanln(&cepIn)
	if err != nil {
		log.Fatal(err)
	}

	chanelViaCep := make(chan ViaCep)
	chanelBrasilApiCep := make(chan BrasilApiCep)

	var viaCepReturn ViaCep
	var BrasilApiCepReturn BrasilApiCep

	go goRoutinesGetCep("http://viacep.com.br/ws/"+cepIn+"/json", viaCepReturn, chanelViaCep)
	go goRoutinesGetCep("http://brasilapi.com.br/api/cep/v1/"+cepIn, BrasilApiCepReturn, chanelBrasilApiCep)

	select {
	case returnViaCep := <-chanelViaCep:
		fmt.Printf("Return from ViaCep: \n %v\n", returnViaCep)

	case returnBrasilApiCep := <-chanelBrasilApiCep:
		if returnBrasilApiCep.Cep != "" {
			fmt.Printf("Return from BrasilApiCep: \n %v\n", returnBrasilApiCep)
		} else {
			fmt.Printf("Return from BrasilApiCep: All CEP services returned an error.")
		}

	case <-time.After(time.Second):
		panic("Timeout error!")
	}
}

func goRoutinesGetCep[T any](url string, returnApiStruct T, chanelApiCep chan T) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error in request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	err = json.Unmarshal(body, &returnApiStruct)
	if err != nil {
		log.Printf("Error parsing json: %v", err)
	}
	chanelApiCep <- returnApiStruct
}
