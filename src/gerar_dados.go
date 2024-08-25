package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func gerarCliente(id int) Cliente {
	return Cliente{
		ID:                 id,
		Idade:              rand.Intn(60) + 18,
		Genero:             rand.Intn(2),
		Localizacao:        rand.Intn(200),
		ValorTotalGasto:    rand.Float64() * 10000,
		FrequenciaCompras:  rand.Intn(100) + 1,
		TipoProduto:        rand.Intn(10),
		DiasDesdeUltCompra: rand.Intn(365),
	}
}

func gerarDadosClientes(numeroClientes int, nomeArquivo string) error {
	file, err := os.Create(nomeArquivo)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Idade", "Genero", "Localizacao", "ValorTotalGasto", "FrequenciaCompras", "TipoProduto", "DiasDesdeUltCompra"})

	for i := 1; i <= numeroClientes; i++ {
		cliente := gerarCliente(i)
		writer.Write([]string{
			strconv.Itoa(cliente.ID),
			strconv.Itoa(cliente.Idade),
			strconv.Itoa(cliente.Genero),
			strconv.Itoa(cliente.Localizacao),
			strconv.FormatFloat(cliente.ValorTotalGasto, 'f', 2, 64),
			strconv.Itoa(cliente.FrequenciaCompras),
			strconv.Itoa(cliente.TipoProduto),
			strconv.Itoa(cliente.DiasDesdeUltCompra),
		})
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numeroClientes := 1000
	nomeArquivo := "../data/clientes.csv"

	if err := gerarDadosClientes(numeroClientes, nomeArquivo); err != nil {
		fmt.Println("Erro ao gerar dados:", err)
		return
	}

	fmt.Println("Dados gerados com sucesso e salvos em", nomeArquivo)
}
