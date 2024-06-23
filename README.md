# SourceScout

SourceScout é uma ferramenta simples escrita em Go para verificar o código-fonte HTML de uma URL especificada. Ele busca informações interessantes, como comentários, scripts externos e metadados, e também pode extrair e analisar subdomínios.

## Funcionalidades

- Verifica o código-fonte HTML de uma URL.
- Busca comentários, scripts externos e metadados.
- Extrai e analisa subdomínios.
- Inclui intervalos aleatórios entre as requisições para reduzir a detecção.

## Requisitos

- Go 1.16+ (https://golang.org/dl/)

## Como usar
go run sourcescout.go
Digite a URL que você quer verificar: https://www.exemplo.com

## Resposta:
Domínio base extraído: www.exemplo.com
Fazendo requisição para: https://www.exemplo.com
Resposta recebida, status code: 200
Resposta lida, tamanho: 12345
Verificando conteúdo interessante...
Comentário encontrado: <!-- Este é um comentário -->
Script externo encontrado: https://www.exemplo.com/script.js
Metadado encontrado: description=Este é um site de exemplo
Extraindo subdomínios...
Nenhum subdomínio encontrado.

by brunod3vs
