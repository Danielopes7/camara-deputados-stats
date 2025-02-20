package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct para os dados da API externa
type Deputado struct {
	ID                     int    `json:"id"`
	URI                    string `json:"uri"`
	Nome                   string `json:"nome"`
	SiglaPartido           string `json:"siglaPartido"`
	URIPartido             string `json:"uriPartido"`
	SiglaUf                string `json:"siglaUf"`
	IDLegislatura          int    `json:"idLegislatura"`
	URLFoto                string `json:"urlFoto"`
	Email                  string `json:"email"`
	AusenciaJustificada    int
	AusenciaNaoJustificada int
	Presenca               int
}

type Deputados struct {
	Dados []Deputado `json:"dados"`
}

// GetData busca dados da API externa
func GetData() ([]Deputado, error) {
	resp, err := http.Get("https://dadosabertos.camara.leg.br/api/v2/deputados")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados: %w", err)
	}
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	var resposta Deputados

	if err := json.NewDecoder(resp.Body).Decode(&resposta); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return resposta.Dados, nil
}
