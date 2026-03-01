package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var err error

const insertCotacaoSQL = `INSERT INTO cotacao (code, bid, datetime) VALUES (?, ?, ?)`

type (
	PriceResponse struct {
		USDBRL struct {
			Code       string `json:"code"`
			Codein     string `json:"codein"`
			Name       string `json:"name"`
			High       string `json:"high"`
			Low        string `json:"low"`
			VarBid     string `json:"varBid"`
			PctChange  string `json:"pctChange"`
			Bid        string `json:"bid"`
			Ask        string `json:"ask"`
			Timestamp  string `json:"timestamp"`
			CreateDate string `json:"create_date"`
		} `json:"USDBRL"`
	}

	CotacaoResponse struct {
		Bid string `json:"bid"`
	}
)

func main() {
	initDatabase()
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handleCotacao)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", mux)
}

func initDatabase() {
	db, err = sql.Open("sqlite", "client-server-api.db")
	if err != nil {
		panic(err.Error())
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS cotacao (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        code VARCHAR(3),
        bid VARCHAR(10),
        datetime TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Banco conectado e tabela criada com sucesso!")
}

func handleCotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cotacao, err := getCotacao(ctx)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cotacao)
}

func getCotacao(ctx context.Context) (CotacaoResponse, error) {
	reqCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return CotacaoResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return CotacaoResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return CotacaoResponse{}, fmt.Errorf("status inválido: %d", res.StatusCode)
	}

	var pr PriceResponse
	err = json.NewDecoder(res.Body).Decode(&pr)
	if err != nil {
		fmt.Println("Erro ao decodificar resposta:", err)
		return CotacaoResponse{}, err
	}

	insertIntoCotacao(ctx, pr.USDBRL.Code, pr.USDBRL.Bid, time.Now().Format(time.RFC3339))

	return CotacaoResponse{
		Bid: pr.USDBRL.Bid,
	}, nil

}

func insertIntoCotacao(ctx context.Context, code string, bid string, datetime string) {
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	stmt, err := db.PrepareContext(dbCtx, insertCotacaoSQL)
	if err != nil {
		fmt.Println("Erro ao preparar statement:", err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(dbCtx, code, bid, datetime)
	if err != nil {
		fmt.Println("Erro ao executar statement:", err)
		return
	}
	fmt.Printf("Cotação inserida %s com sucesso!\n", code)
}
