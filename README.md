# Segmentação de Clientes com K-Means

## Descrição do Problema

A segmentação de clientes é uma técnica fundamental para dividir uma base de clientes em grupos com características semelhantes. Isso permite que empresas personalizem suas estratégias de marketing, melhorem a oferta de produtos e aumentem a satisfação do cliente. A segmentação se baseia em variáveis como idade, gênero, localização, comportamento de compra, entre outras.

## Aplicação do K-Means

O algoritmo K-Means é utilizado para agrupar clientes em K clusters, onde cada cluster contém clientes com características semelhantes. O processo envolve:

1. **Inicialização**: Seleção inicial de centroids.
2. **Atribuição**: Cada cliente é atribuído ao cluster mais próximo.
3. **Atualização**: Centroids são recalculados com base nos clientes atribuídos, repetindo o processo até que os clusters se estabilizem.

## Melhorias Implementadas

Foram realizadas as seguintes melhorias no projeto:

- **Comparação de Desempenho**: Implementação de uma versão paralelizada do K-Means para comparação de desempenho com a versão sequencial.
- **Análises Estatísticas**: Implementação de cálculos de tempo de execução, speedup, eficiência e a métrica de Karp-Flatt para cada valor de K.
- **Geração de Gráficos**: Criação de gráficos para visualização dos tempos de execução.

## Como Executar

1. **Clonar o Repositório**:
   ```bash
   git clone git@github.com:vitormeloa/gac105-tp.git
   cd gac105-tp

2. **Rodar o Script**:
   ```bash
   sudo chmod +x run_kmeans.sh
    ./run_kmeans.sh
   ```
   
3. **Visualizar os Resultados**:
    - Os resultados são armazenados no diretório `data/`.

## Autores
    - Eduardo Oliveira Gomes
    - Marco Antônio Martins Ribeiro de Jesus
    - Vitor Melo Assunção