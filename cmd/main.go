package main

import (
	"api-gov/service"
	"fmt"
)

func main() {
	data, error := service.GetData()

	if error != nil {
		fmt.Println("Erro ao buscar dados:", error)
		return
	}
	for _, deputado := range data {
		fmt.Printf("Nome: %d\n", deputado.ID)
		fmt.Printf("Partido: %s\n", deputado.SiglaPartido)
		fmt.Printf("UF: %s\n", deputado.SiglaUf)
		fmt.Println("---------------")
	}
	fmt.Printf("Encontrados %d deputados\n", len(data))
}
