package main

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
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
	numExecucoes := 10

	sequencialTimes := make(plotter.Values, len(ks))
	paraleloTimes := make(plotter.Values, len(ks))

	outFile, err := os.Create("data/analises.txt")

	if err != nil {
		log.Fatalf("Erro ao criar arquivo de saída: %v", err)
	}

	defer outFile.Close()

	for i, k := range ks {
		fmt.Fprintf(outFile, "Teste com k = %d\n", k)

		var temposSequenciais []time.Duration
		var temposParalelos []time.Duration

		for j := 0; j < numExecucoes; j++ {
			startSequencial := time.Now()
			kmeansSequencial(clientes, k, maxIteracoes, tolerancia)
			elapsedSequencial := time.Since(startSequencial)
			temposSequenciais = append(temposSequenciais, elapsedSequencial)

			startParalelizado := time.Now()
			kmeans(clientes, k, maxIteracoes, tolerancia)
			elapsedParalelizado := time.Since(startParalelizado)
			temposParalelos = append(temposParalelos, elapsedParalelizado)
		}

		mediaSequencial := calcularMedia(temposSequenciais)
		mediaParalelo := calcularMedia(temposParalelos)

		sequencialTimes[i] = mediaSequencial
		paraleloTimes[i] = mediaParalelo

		fmt.Fprintf(outFile, "Tempo médio de execução (Sequencial): %.6f segundos\n", mediaSequencial)
		fmt.Fprintf(outFile, "Tempo médio de execução (Paralelizado): %.6f segundos\n", mediaParalelo)

		speedup := calcularSpeedup(mediaSequencial, mediaParalelo)
		fmt.Fprintf(outFile, "Speedup: %.2f\n", speedup)

		fracaoParalelizavel := calcularFracaoParalelizavel(speedup, 4)
		fmt.Fprintf(outFile, "Fração Paralelizável (p): %.4f\n", fracaoParalelizavel)

		karpFlatt := calcularKarpFlattMetric(speedup, 4)
		fmt.Fprintf(outFile, "Métrica de Karp-Flatt: %.4f\n", karpFlatt)

		fmt.Fprintf(outFile, "---------------------------\n")
	}

	err = gerarGrafico(sequencialTimes, paraleloTimes, ks)

	if err != nil {
		log.Fatalf("Erro ao gerar gráfico: %v", err)
	}

	fmt.Println("Execução concluída. Resultados salvos em data/analises.txt")
}

func gerarGrafico(sequencialTimes, paraleloTimes plotter.Values, ks []int) error {
	p := plot.New()

	p.Title.Text = "Comparação de Tempo de Execução"
	p.X.Label.Text = "Número de Clusters (k)"
	p.Y.Label.Text = "Tempo de Execução (segundos)"

	barWidth := vg.Points(10)
	spaceBetweenBars := vg.Points(5)
	sequencialBars, err := plotter.NewBarChart(sequencialTimes, barWidth)

	if err != nil {
		return err
	}

	sequencialBars.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	sequencialBars.Offset = -spaceBetweenBars

	paraleloBars, err := plotter.NewBarChart(paraleloTimes, barWidth)

	if err != nil {
		return err
	}

	paraleloBars.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	paraleloBars.Offset = spaceBetweenBars
	p.Add(sequencialBars, paraleloBars)

	labels := make([]string, len(ks))

	for i, k := range ks {
		labels[i] = fmt.Sprintf("%d", k)
	}

	p.NominalX(labels...)

	p.Legend.Add("Sequencial", sequencialBars)
	p.Legend.Add("Paralelizado", paraleloBars)
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(8*vg.Inch, 6*vg.Inch, "data/comparacao_kmeans.png"); err != nil {
		return err
	}

	fmt.Println("Gráfico gerado com sucesso: comparacao_kmeans.png")
	return nil
}
