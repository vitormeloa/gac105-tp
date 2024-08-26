package main

import (
	"fmt"
	"math"
)

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

func kmeansSequencial(clientes []Cliente, k int, maxIteracoes int, tolerancia float64) []int {
	centroids := inicializarCentroids(clientes, k)

	for i := 0; i < maxIteracoes; i++ {
		clusters := atribuirClustersSequencial(clientes, centroids)
		novosCentroids := atualizarCentroidsSequencial(clientes, clusters, k)

		if convergiram(centroids, novosCentroids, tolerancia) {
			fmt.Printf("Convergência atingida após %d iterações (Sequencial).\n", i+1)
			break
		}

		centroids = novosCentroids
	}

	return atribuirClustersSequencial(clientes, centroids)
}
