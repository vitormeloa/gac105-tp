#!/bin/bash

# compila e roda script que gera a base de dados
go run src/gerar_dados.go src/models.go

# compila e roda o programa
go run src/main.go src/kmeans_paralelo.go src/kmeans_sequencial.go src/utils.go src/models.go