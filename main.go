package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url1 := "https://cdn.apicep.com/file/apicep/"
	end_url1 := ".json"

	url2 := "http://viacep.com.br/ws/"
	end_url2 := "/json/"

	cep := "13423-100"

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		resp := CallCepApi(url1, cep, end_url1)
		c1 <- "ApiCEP.com respondeu mais rápido com o resultado: " + resp + "\n"
	}()

	go func() {
		resp := CallCepApi(url2, cep, end_url2)
		c2 <- "ViaCEP respondeu mais rápido com o resultado: " + resp + "\n"
	}()

	select {
	case msg := <-c1: // APICep.com
		fmt.Printf(msg)

	case msg := <-c2: // ViaCEP
		fmt.Printf(msg)

	case <-time.After(time.Second):
		println("Timeout!")
	}
}

func CallCepApi(url string, cep string, end_url string) string {
	req, err := http.Get(url + cep + end_url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	return string(res)
}
