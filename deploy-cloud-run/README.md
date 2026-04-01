# API Weather CEP

Microsserviço em Go que recebe um CEP brasileiro (código postal), identifica a cidade correspondente e retorna a temperatura atual do clima em Celsius, Fahrenheit e Kelvin.

### GET /weather/{cep}
Obtém o clima pelo CEP.

**Parâmetros de Caminho:**
- `cep`: Código postal brasileiro (8 dígitos)

**Resposta (Sucesso - 200):**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

## Execução Local

1. Instale as dependências:
```bash
go mod tidy
```

2. Execute o servidor:
```bash
$env:WEATHER_API_KEY="c91bb5c1e81042c9aa2145223260104"; go run main.go
```

3. Teste a API:
```bash
curl http://localhost:8080/health
curl http://localhost:8080/weather/01310930
```

## URL de Acesso Google Cloud Run

https://deploy-weather-156587497249.southamerica-east1.run.app/weather/89031070

## Executando os Testes

```bash
go test ./... -v
```