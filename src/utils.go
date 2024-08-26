package main

import (
	"math"
)

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
