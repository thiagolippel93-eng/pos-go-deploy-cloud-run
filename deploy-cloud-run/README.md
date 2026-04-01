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

**Resposta (CEP Inválido - 422):**
```json
{
  "message": "invalid zipcode"
}
```

## Configuração

A API requer uma chave da WeatherAPI para funcionar. Você pode configurar a variável de ambiente `WEATHER_API_KEY` das seguintes formas:


### Variável de Ambiente

**Linux/Mac:**
```bash
export WEATHER_API_KEY="c91bb5c1e81042c9aa2145223260104"
```

**Windows (CMD):**
```bash
set WEATHER_API_KEY="c91bb5c1e81042c9aa2145223260104"
```

## Execução Local

1. Instale as dependências:
```bash
go mod tidy
```

2. Defina a chave da WeatherAPI:
```bash
export WEATHER_API_KEY=c91bb5c1e81042c9aa2145223260104
```

3. Execute o servidor:
```bash
go run main.go
```

4. Teste a API:
```bash
curl http://localhost:8080/health
curl http://localhost:8080/weather/01310930
```

## Executando os Testes

```bash
go test ./... -v
```

## Docker

### Construindo a Imagem

```bash
docker build -t weather-cep .
```

### Executando com Docker

```bash
docker run -p 8080:8080 -e WEATHER_API_KEY=sua_chave_aqui weather-cep
```