package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

func main() {
	clientes, err := carregarDados("data/clientes.csv")
	if err != nil {
		log.Fatalf("Erro ao carregar dados: %v", err)
	}

	ks := []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	maxIteracoes := 100
	tolerancia := 0.0001

	sequencialTimes := make(plotter.Values, len(ks))
	paraleloTimes := make(plotter.Values, len(ks))

	for i, k := range ks {
		fmt.Printf("Teste com k = %d\n", k)

		startSequencial := time.Now()
		kmeansSequencial(clientes, k, maxIteracoes, tolerancia)
		elapsedSequencial := time.Since(startSequencial)
		sequencialTimes[i] = elapsedSequencial.Seconds()
		fmt.Printf("Tempo de execução (Sequencial): %s\n", elapsedSequencial)

		startParalelizado := time.Now()
		kmeans(clientes, k, maxIteracoes, tolerancia)
		elapsedParalelizado := time.Since(startParalelizado)
		paraleloTimes[i] = elapsedParalelizado.Seconds()
		fmt.Printf("Tempo de execução (Paralelizado): %s\n", elapsedParalelizado)
	}

	err = gerarGrafico(sequencialTimes, paraleloTimes, ks)
	if err != nil {
		log.Fatalf("Erro ao gerar gráfico: %v", err)
	}
}

func gerarGrafico(sequencialTimes, paraleloTimes plotter.Values, ks []int) error {
	p := plot.New()

	p.Title.Text = "Comparação de Tempo de Execução"
	p.X.Label.Text = "Número de Clusters (k)"
	p.Y.Label.Text = "Tempo de Execução (segundos)"

	sequencialBars, err := plotter.NewBarChart(sequencialTimes, vg.Points(20))
	if err != nil {
		return err
	}
	sequencialBars.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	sequencialBars.Offset = -vg.Points(10)
	paraleloBars, err := plotter.NewBarChart(paraleloTimes, vg.Points(20))
	if err != nil {
		return err
	}
	paraleloBars.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	paraleloBars.Offset = vg.Points(10)

	p.Add(sequencialBars, paraleloBars)

	p.NominalX(fmt.Sprintf("%d", ks[0]), fmt.Sprintf("%d", ks[1]), fmt.Sprintf("%d", ks[2]))

	p.Legend.Add("Sequencial", sequencialBars)
	p.Legend.Add("Paralelizado", paraleloBars)

	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "data/comparacao_kmeans.png"); err != nil {
		return err
	}

	fmt.Println("Gráfico gerado com sucesso: comparacao_kmeans.png")
	return nil
}
