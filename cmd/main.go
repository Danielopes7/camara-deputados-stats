package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

func getParticipation(data service.Deputados) {

	for index, deputado := range data.Dados {
		if index > 5 {
			break
		}
		url := fmt.Sprintf("https://www.camara.leg.br/deputados/%d?ano=2024", deputado.ID)

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

		presencas := make(map[string]string)

		doc.Find(".presencas__section").Each(func(i int, s *goquery.Selection) {
			s.Find(".presencas__data.presencas__mobile-outras").Each(func(j int, item *goquery.Selection) {
				label := strings.TrimSpace(item.Find(".presencas__label").Text())
				valor := strings.TrimSpace(item.Find(".presencas__qtd").Text())
				chave := fmt.Sprintf("[%d] -   %s", i, label)

				presencas[chave] = valor
			})
		})

		fmt.Printf("ID: %d\n", deputado.ID)
		fmt.Printf("Nome: %s\n", deputado.Nome)
		fmt.Printf("Partido: %s\n", deputado.SiglaPartido)
		fmt.Printf("UF: %s\n", deputado.SiglaUf)

		fmt.Println("📊 Presença e Ausência do Deputado:")
		for chave, valor := range presencas {
			fmt.Printf("%s: %s\n", chave, valor)
		}
		fmt.Println("------------------------------------------------------------")

	}

}
