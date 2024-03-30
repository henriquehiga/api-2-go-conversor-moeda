## Tutorial de como iniciar o programa de converter de moedas;

1. Instale o GO na versão 1.22.1;
2. Na pasta do projeto instale as dependências com comando:
```
go install
```
3. Inicie o programa com o comando:
```cmd
go run index.go
```
-----------------------------------------------------------------
1. Para chamar a API:
Com o programa rodando faça uma chamada HTTP com método POST para a rota '/converte-moedas' enviando um objeto como o exemplo abaixo:
```json
{
    "valor": 100
}
```