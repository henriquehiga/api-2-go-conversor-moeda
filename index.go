package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/converte-moedas", handleConverteMoedas)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

type Conversao struct {
	Dolar   float64 `json:"dolar"`
	Euro    float64 `json:"euro"`
	Maquina string  `json:"maquina"`
	Real    float64 `json:"real"`
}

type Resposta struct {
	Conversao Conversao `json:"conversao"`
}

func handleConverteMoedas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body map[string]float64
	json.NewDecoder(r.Body).Decode(&body)

	valor, ok := body["valor"]
	if !ok {
		http.Error(w, "É obrigatório enviar um valor", http.StatusBadRequest)
		return
	}

	cotacoes, err := resgataCotacaoEuroEDolar()
	if err != nil {
		http.Error(w, "Erro ao resgatar cotação de euro e dolar", http.StatusInternalServerError)
		return
	}

	hostname, _ := os.Hostname()
	resultado := Resposta{
		Conversao: Conversao{
			Real:    valor,
			Dolar:   valor / cotacoes["cotacao_dolar"],
			Euro:    valor / cotacoes["cotacao_euro"],
			Maquina: hostname,
		},
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

	cotacaoDolar, _ := strconv.ParseFloat(strings.TrimSpace(responseJSON["USDBRL"]["high"]), 64)
	cotacaoEuro, _ := strconv.ParseFloat(strings.TrimSpace(responseJSON["EURBRL"]["high"]), 64)

	return map[string]float64{
		"cotacao_dolar": cotacaoDolar,
		"cotacao_euro":  cotacaoEuro,
	}, nil
}
