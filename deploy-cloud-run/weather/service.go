package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponse struct {
	Location LocationData `json:"location"`
	Current  CurrentData  `json:"current"`
}

type LocationData struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type CurrentData struct {
	TempC      float64          `json:"temp_c"`
	TempF      float64          `json:"temp_f"`
	TempK      float64          `json:"temp_k"`
	IsDay      int              `json:"is_day"`
	Condition  WeatherCondition `json:"condition"`
	WindKph    float64          `json:"wind_kph"`
	WindMph    float64          `json:"wind_mph"`
	Humidity   int              `json:"humidity"`
	Cloud      int              `json:"cloud"`
	FeelslikeC float64          `json:"feelslike_c"`
	FeelslikeF float64          `json:"feelslike_f"`
	UV         float64          `json:"uv"`
}

type WeatherCondition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Service struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewService() *Service {
	return &Service{
		httpClient: &http.Client{},
		baseURL:    "https://api.weatherapi.com/v1",
		apiKey:     os.Getenv("WEATHER_API_KEY"),
	}
}

func (s *Service) GetCurrentWeather(cidade string) (*CurrentData, error) {
	requestURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", s.baseURL, s.apiKey, url.QueryEscape(cidade))

	resp, err := s.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar dados do clima: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler corpo da resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código de status inesperado: %d, body: %s", resp.StatusCode, string(body))
	}

	var weatherResp WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("falha ao analisar resposta: %w", err)
	}

	return &weatherResp.Current, nil
}

func (s *Service) GetTemperatureCelsius(cidade string) (float64, error) {
	dadosClima, err := s.GetCurrentWeather(cidade)
	if err != nil {
		return 0, err
	}
	return dadosClima.TempC, nil
}
