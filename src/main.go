package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func carregarDados(nomeArquivo string) ([]Cliente, error) {
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var clientes []Cliente
	for _, record := range records[1:] {
		idade, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter Idade: %v", err)
		}
		genero, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter Genero: %v", err)
		}
		localizacao, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter Localizacao: %v", err)
		}
		valorTotalGasto, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter ValorTotalGasto: %v", err)
		}
		frequenciaCompras, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter FrequenciaCompras: %v", err)
		}
		tipoProduto, err := strconv.Atoi(record[6])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter TipoProduto: %v", err)
		}
		diasDesdeUltCompra, err := strconv.Atoi(record[7])
		if err != nil {
			return nil, fmt.Errorf("erro ao converter DiasDesdeUltCompra: %v", err)
		}

		cliente := Cliente{
			Idade:              idade,
			Genero:             genero,
			Localizacao:        localizacao,
			ValorTotalGasto:    valorTotalGasto,
			FrequenciaCompras:  frequenciaCompras,
			TipoProduto:        tipoProduto,
			DiasDesdeUltCompra: diasDesdeUltCompra,
		}
		clientes = append(clientes, cliente)
	}
	return clientes, nil
}

// Função para inicializar os centróides
func inicializarCentroids(clientes []Cliente, k int) []Centroid {
	rand.Seed(time.Now().UnixNano())
	centroids := make([]Centroid, k)
	for i := 0; i < k; i++ {
		c := clientes[rand.Intn(len(clientes))]
		centroids[i] = Centroid{
			Idade:              float64(c.Idade),
			Genero:             float64(c.Genero),
			Localizacao:        float64(c.Localizacao),
			ValorTotalGasto:    c.ValorTotalGasto,
			FrequenciaCompras:  float64(c.FrequenciaCompras),
			TipoProduto:        float64(c.TipoProduto),
			DiasDesdeUltCompra: float64(c.DiasDesdeUltCompra),
		}
	}
	return centroids
}

func main() {
	clientes, err := carregarDados("data/clientes.csv")
	if err != nil {
		log.Fatalf("Erro ao carregar dados: %v", err)
	}

	ks := []int{3, 5, 8}
	maxIteracoes := 100
	tolerancia := 0.0001

	for _, k := range ks {
		fmt.Printf("Teste com k = %d\n", k)

		startSequencial := time.Now()
		kmeansSequencial(clientes, k, maxIteracoes, tolerancia)
		elapsedSequencial := time.Since(startSequencial)
		fmt.Printf("Tempo de execução (Sequencial): %s\n", elapsedSequencial)

		startParalelizado := time.Now()
		kmeans(clientes, k, maxIteracoes, tolerancia)
		elapsedParalelizado := time.Since(startParalelizado)
		fmt.Printf("Tempo de execução (Paralelizado): %s\n", elapsedParalelizado)
	}
}
