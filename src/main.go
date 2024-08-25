package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
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

func calcularDistancia(c Cliente, centroide Centroid) float64 {
	return math.Sqrt(
		math.Pow(float64(c.Idade)-centroide.Idade, 2) +
			math.Pow(float64(c.Genero)-centroide.Genero, 2) +
			math.Pow(float64(c.Localizacao)-centroide.Localizacao, 2) +
			math.Pow(c.ValorTotalGasto-centroide.ValorTotalGasto, 2) +
			math.Pow(float64(c.FrequenciaCompras)-centroide.FrequenciaCompras, 2) +
			math.Pow(float64(c.TipoProduto)-centroide.TipoProduto, 2) +
			math.Pow(float64(c.DiasDesdeUltCompra)-centroide.DiasDesdeUltCompra, 2),
	)
}

func atribuirClusters(clientes []Cliente, centroids []Centroid) []int {
	clusters := make([]int, len(clientes))
	var wg sync.WaitGroup

	for i := range clientes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			menorDist := math.MaxFloat64
			clusterIndex := 0
			for j, centroide := range centroids {
				dist := calcularDistancia(clientes[i], centroide)
				if dist < menorDist {
					menorDist = dist
					clusterIndex = j
				}
			}
			clusters[i] = clusterIndex
		}(i)
	}
	wg.Wait()
	return clusters
}

func atualizarCentroids(clientes []Cliente, clusters []int, k int) []Centroid {
	centroids := make([]Centroid, k)
	contadores := make([]int, k)

	for i, cliente := range clientes {
		cluster := clusters[i]
		centroids[cluster].Idade += float64(cliente.Idade)
		centroids[cluster].Genero += float64(cliente.Genero)
		centroids[cluster].Localizacao += float64(cliente.Localizacao)
		centroids[cluster].ValorTotalGasto += cliente.ValorTotalGasto
		centroids[cluster].FrequenciaCompras += float64(cliente.FrequenciaCompras)
		centroids[cluster].TipoProduto += float64(cliente.TipoProduto)
		centroids[cluster].DiasDesdeUltCompra += float64(cliente.DiasDesdeUltCompra)
		contadores[cluster]++
	}

	for i := 0; i < k; i++ {
		if contadores[i] > 0 {
			centroids[i].Idade /= float64(contadores[i])
			centroids[i].Genero /= float64(contadores[i])
			centroids[i].Localizacao /= float64(contadores[i])
			centroids[i].ValorTotalGasto /= float64(contadores[i])
			centroids[i].FrequenciaCompras /= float64(contadores[i])
			centroids[i].TipoProduto /= float64(contadores[i])
			centroids[i].DiasDesdeUltCompra /= float64(contadores[i])
		}
	}

	return centroids
}

func convergiram(centroidsAntigos, centroidsNovos []Centroid, tolerancia float64) bool {
	for i := range centroidsAntigos {
		if calcularDistanciaEntreCentroids(centroidsAntigos[i], centroidsNovos[i]) > tolerancia {
			return false
		}
	}
	return true
}

func calcularDistanciaEntreCentroids(a, b Centroid) float64 {
	return math.Sqrt(
		math.Pow(a.Idade-b.Idade, 2) +
			math.Pow(a.Genero-b.Genero, 2) +
			math.Pow(a.Localizacao-b.Localizacao, 2) +
			math.Pow(a.ValorTotalGasto-b.ValorTotalGasto, 2) +
			math.Pow(a.FrequenciaCompras-b.FrequenciaCompras, 2) +
			math.Pow(a.TipoProduto-b.TipoProduto, 2) +
			math.Pow(a.DiasDesdeUltCompra-b.DiasDesdeUltCompra, 2),
	)
}

func kmeans(clientes []Cliente, k int, maxIteracoes int, tolerancia float64) []int {
	centroids := inicializarCentroids(clientes, k)

	for i := 0; i < maxIteracoes; i++ {
		clusters := atribuirClusters(clientes, centroids)
		novosCentroids := atualizarCentroids(clientes, clusters, k)

		if convergiram(centroids, novosCentroids, tolerancia) {
			fmt.Printf("Convergência atingida após %d iterações.\n", i+1)
			break
		}

		centroids = novosCentroids
	}

	return atribuirClusters(clientes, centroids)
}

func kmeansSequencial(clientes []Cliente, k int, maxIteracoes int, tolerancia float64) []int {
	centroids := inicializarCentroids(clientes, k)

	for i := 0; i < maxIteracoes; i++ {
		clusters := atribuirClustersSequencial(clientes, centroids)
		novosCentroids := atualizarCentroidsSequencial(clientes, clusters, k)

		if convergiram(centroids, novosCentroids, tolerancia) {
			fmt.Printf("Convergência atingida após %d iterações.\n", i+1)
			break
		}

		centroids = novosCentroids
	}

	return atribuirClustersSequencial(clientes, centroids)
}

func atribuirClustersSequencial(clientes []Cliente, centroids []Centroid) []int {
	clusters := make([]int, len(clientes))
	for i, cliente := range clientes {
		menorDist := math.MaxFloat64
		clusterIndex := 0
		for j, centroide := range centroids {
			dist := calcularDistancia(cliente, centroide)
			if dist < menorDist {
				menorDist = dist
				clusterIndex = j
			}
		}
		clusters[i] = clusterIndex
	}
	return clusters
}

func atualizarCentroidsSequencial(clientes []Cliente, clusters []int, k int) []Centroid {
	centroids := make([]Centroid, k)
	contadores := make([]int, k)

	for i, cliente := range clientes {
		cluster := clusters[i]
		centroids[cluster].Idade += float64(cliente.Idade)
		centroids[cluster].Genero += float64(cliente.Genero)
		centroids[cluster].Localizacao += float64(cliente.Localizacao)
		centroids[cluster].ValorTotalGasto += cliente.ValorTotalGasto
		centroids[cluster].FrequenciaCompras += float64(cliente.FrequenciaCompras)
		centroids[cluster].TipoProduto += float64(cliente.TipoProduto)
		centroids[cluster].DiasDesdeUltCompra += float64(cliente.DiasDesdeUltCompra)
		contadores[cluster]++
	}

	for i := 0; i < k; i++ {
		if contadores[i] > 0 {
			centroids[i].Idade /= float64(contadores[i])
			centroids[i].Genero /= float64(contadores[i])
			centroids[i].Localizacao /= float64(contadores[i])
			centroids[i].ValorTotalGasto /= float64(contadores[i])
			centroids[i].FrequenciaCompras /= float64(contadores[i])
			centroids[i].TipoProduto /= float64(contadores[i])
			centroids[i].DiasDesdeUltCompra /= float64(contadores[i])
		}
	}

	return centroids
}

func main() {
	clientes, err := carregarDados("../data/clientes.csv")
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
