package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"weather-cep/cep"
	"weather-cep/temperature"
	"weather-cep/weather"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type WeatherHandler struct {
	cepService     *cep.Service
	weatherService *weather.Service
	conversorTemp  *temperature.Converter
}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{
		cepService:     cep.NewService(),
		weatherService: weather.NewService(),
		conversorTemp:  temperature.NewConverter(),
	}
}

func extrairCEP(caminho string) string {
	// Padrão para extrair o CEP do caminho /weather/{cep}
	re := regexp.MustCompile(`/weather/(\d{5}-?\d{3}?)`)
	resultado := re.FindStringSubmatch(caminho)
	if len(resultado) > 1 {
		return resultado[1]
	}

	// Tenta extrair apenas os dígitos
	reDigitos := regexp.MustCompile(`/weather/(\d{5}\d{0,4})`)
	resultadoDigitos := reDigitos.FindStringSubmatch(caminho)
	if len(resultadoDigitos) > 1 {
		return resultadoDigitos[1]
	}

	return ""
}

func (h *WeatherHandler) GetWeatherByCEP(w http.ResponseWriter, r *http.Request) {
	cepParam := extrairCEP(r.URL.Path)

	if err := h.cepService.ValidarCEP(cepParam); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "CEP inválido"})
		return
	}

	// Obtém a cidade a partir do CEP
	cidade, err := h.cepService.GetCity(cepParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "can not find") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "CEP não encontrado"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
		}
		return
	}

	dadosClima, err := h.weatherService.GetCurrentWeather(cidade)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
		log.Printf("Erro ao buscar clima para a cidade %s: %v", cidade, err)
		return
	}

	respostaTemp := h.conversorTemp.Convert(dadosClima.TempC)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respostaTemp)
}

func (h *WeatherHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	if os.Getenv("WEATHER_API_KEY") == "" {
		log.Println("AVISO: Variável de ambiente WEATHER_API_KEY não definida")
	}

	manipulador := NewWeatherHandler()

	http.HandleFunc("/health", manipulador.HealthCheck)
	http.HandleFunc("/weather/", manipulador.GetWeatherByCEP)

	porta := strings.TrimSpace(os.Getenv("PORT"))
	if porta == "" {
		porta = "8080"
	}

	log.Printf("Iniciando servidor na porta %s", porta)
	if err := http.ListenAndServe(":"+porta, nil); err != nil {
		log.Fatalf("Falha ao iniciar servidor: %v", err)
	}
}
