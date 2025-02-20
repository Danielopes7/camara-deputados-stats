package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Danielopes7/camara-deputados-stats/service"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	data, error := service.GetData()

	if error != nil {
		fmt.Println("Erro ao buscar dados:", error)
		return
	}
	getParticipation(service.Deputados{Dados: data})

	fmt.Printf("Encontrados %d deputados\n", len(data))
}

type Resultado struct {
	Indice int
	Dados  json.RawMessage
}

func getParticipation(data service.Deputados) {
	defer measureTime()()

	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)
	for index, deputado := range data.Dados {
		if index > 50 {
			break
		}
		url := fmt.Sprintf("https://www.camara.leg.br/deputados/%d?ano=2024", deputado.ID)
		wg.Add(1)
		sem <- struct{}{}

		go makeRequest(url, &data.Dados[index], &wg, sem)
	}

	wg.Wait()

	for _, deputado := range data.Dados[:10] {
		fmt.Printf("ID: %d\n", deputado.ID)
		fmt.Printf("Nome: %s\n", deputado.Nome)
		fmt.Printf("Partido: %s\n", deputado.SiglaPartido)
		fmt.Printf("UF: %s\n", deputado.SiglaUf)

		fmt.Println("📊 Presença e Ausência do Deputado:")
		fmt.Printf("Ausencia Nao Justificada: %d\n", deputado.AusenciaNaoJustificada)
		fmt.Printf("Ausencia Justificada: %d\n", deputado.AusenciaNaoJustificada)
		fmt.Printf("Presença: %d\n", deputado.Presenca)
		fmt.Println("------------------------------------------------------------")
	}
}

func makeRequest(url string, deputado *service.Deputado, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	defer func() { <-sem }()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal("Erro ao acessar a página:", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Erro ao acessar a página: Status %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal("Erro ao processar HTML:", err)
	}

	doc.Find(".presencas__section").Each(func(i int, s *goquery.Selection) {
		s.Find(".presencas__data.presencas__mobile-outras").Each(func(j int, item *goquery.Selection) {
			label := strings.TrimSpace(item.Find(".presencas__label").Text())
			valor := strings.TrimSpace(item.Find(".presencas__qtd").Text())
			re := regexp.MustCompile(`(\d+)`)
			match := re.FindString(valor)
			if match != "" {
				num, err := strconv.Atoi(match)
				if err == nil {
					if label == "Ausências não justificadas" {
						deputado.AusenciaNaoJustificada = num
					}
					if label == "Ausências justificadas" {
						deputado.AusenciaJustificada = num
					}
					if label == "Presenças na Câmara" {
						deputado.Presenca = num
					}
				}
			}
		})
	})
}

func measureTime() func() {
	start := time.Now()
	return func() {
		fmt.Println("Execution time:", time.Since(start))
	}
}
