package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/converte-moedas", handleConverteMoedas)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func handleConverteMoedas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body map[string]float64
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valor, ok := body["valor"]
	if !ok {
		http.Error(w, "Missing 'valor' field", http.StatusBadRequest)
		return
	}

	cotacoes, err := resgataCotacaoEuroEDolar()
	if err != nil {
		http.Error(w, "Error fetching currency rates", http.StatusInternalServerError)
		return
	}

	resultado := map[string]float64{
		"real":    valor,
		"dolar":   calculaValorEmReal(valor, cotacoes["cotacao_dolar"]),
		"euro":    calculaValorEmReal(valor, cotacoes["cotacao_euro"]),
		"maquina": 10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func calculaValorEmReal(valor, cotacao float64) float64 {
	return valor * cotacao
}

func resgataCotacaoEuroEDolar() (map[string]float64, error) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL,EUR-BRL")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseJSON map[string]map[string]string
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		return nil, err
	}

	cotacaoDolar := responseJSON["USDBRL"]["high"]
	cotacaoEuro := responseJSON["EURBRL"]["high"]

	floatCotacaoDolar, _ := strconv.ParseFloat(strings.TrimSpace(cotacaoDolar), 64)
	floatCotacaoEuro, _ := strconv.ParseFloat(strings.TrimSpace(cotacaoEuro), 64)

	return map[string]float64{
		"cotacao_dolar": floatCotacaoDolar,
		"cotacao_euro":  floatCotacaoEuro,
	}, nil
}
