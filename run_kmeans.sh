#!/bin/bash

# verifica se o Go está instalado
if ! command -v go &> /dev/null
then
    echo "Go não está instalado. Instalando Go..."
    sudo apt-get update
    sudo apt-get install -y golang-go

    if ! command -v go &> /dev/null
    then
        echo "Falha ao instalar Go. Por favor, instale-o manualmente."
        exit 1
    fi
fi

# compila e roda o script que gera a base de dados
go run src/gerar_dados.go src/models.go
if [ $? -ne 0 ]; then
    echo "Erro ao executar o script de geração de dados."
    exit 1
fi

# compila e roda o programa principal
go run src/main.go src/kmeans_paralelo.go src/kmeans_sequencial.go src/utils.go src/models.go
if [ $? -ne 0 ]; then
    echo "Erro ao executar o programa principal."
    exit 1
fi

echo "Execução concluída com sucesso."