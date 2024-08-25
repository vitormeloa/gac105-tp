package main

type Cliente struct {
	ID                 int
	Idade              int
	Genero             int
	Localizacao        int
	ValorTotalGasto    float64
	FrequenciaCompras  int
	TipoProduto        int
	DiasDesdeUltCompra int
}

type Centroid struct {
	Idade              float64
	Genero             float64
	Localizacao        float64
	ValorTotalGasto    float64
	FrequenciaCompras  float64
	TipoProduto        float64
	DiasDesdeUltCompra float64
}
