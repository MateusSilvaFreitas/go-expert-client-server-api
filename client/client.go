package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro na requisição:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Erro: status code", res.StatusCode)
		return
	}

	var c Cotacao
	json.NewDecoder(res.Body).Decode(&c)

	file, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("Dólar: %v", c.Bid))

}
