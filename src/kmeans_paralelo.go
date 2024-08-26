package main

import (
	"fmt"
	"math"
	"sync"
)

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

func kmeans(clientes []Cliente, k int, maxIteracoes int, tolerancia float64) []int {
	centroids := inicializarCentroids(clientes, k)

	for i := 0; i < maxIteracoes; i++ {
		clusters := atribuirClusters(clientes, centroids)
		novosCentroids := atualizarCentroids(clientes, clusters, k)

		if convergiram(centroids, novosCentroids, tolerancia) {
			fmt.Printf("Convergência atingida após %d iterações (Paralelizado).\n", i+1)
			break
		}

		centroids = novosCentroids
	}

	return atribuirClusters(clientes, centroids)
}
