package main

import (
	"fmt"
	"gonum.org/v1/plot/plotter"
	"math"
	"runtime"
	"time"
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

func calcularSpeedupEficiência(sequencialTimes, paraleloTimes plotter.Values, ks []int) {
	fmt.Println("\nResultados de Speedup e Eficiência:")
	for i, k := range ks {
		speedup := sequencialTimes[i] / paraleloTimes[i]
		eficiencia := speedup / float64(runtime.NumCPU())
		fmt.Printf("k = %d | Speedup: %f | Eficiência: %f\n", k, speedup, eficiencia)
	}
}

func calcularMedia(tempos []time.Duration) float64 {
	var soma float64
	for _, t := range tempos {
		soma += t.Seconds()
	}
	return soma / float64(len(tempos))
}

func calcularSpeedup(T_seq, T_par float64) float64 {
	return T_seq / T_par
}

func calcularFracaoParalelizavel(S float64, N int) float64 {
	return (float64(N) * (S - 1)) / (S * float64(N-1))
}

func calcularKarpFlattMetric(S float64, N int) float64 {
	if S <= 0 || N <= 0 {
		fmt.Println("Valores inválidos para S ou N")
		return math.NaN() // Valor indefinido
	}
	return (1/S - 1/float64(N)) / (1 - 1/float64(N))
}

