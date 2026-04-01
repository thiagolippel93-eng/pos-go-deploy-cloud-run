package cep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// Service fornece funcionalidade de consulta de CEP
type Service struct {
	httpClient *http.Client
	baseURL    string
}

// NewService cria um novo serviço de CEP
func NewService() *Service {
	return &Service{
		httpClient: &http.Client{},
		baseURL:    "https://viacep.com.br/ws",
	}
}

// ValidarCEP verifica se o CEP é válido (8 dígitos)
func (s *Service) ValidarCEP(cep string) error {
	// Remove qualquer formatação (pontos, hifens, espaços)
	cepLimpo := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(cep, ".", ""), "-", ""), " ", "")

	// Verifica se tem exatamente 8 dígitos
	if len(cepLimpo) != 8 {
		return fmt.Errorf("invalid zipcode")
	}

	// Verifica se todos os caracteres são dígitos
	valido, _ := regexp.MatchString(`^\d{8}$`, cepLimpo)
	if !valido {
		return fmt.Errorf("invalid zipcode")
	}

	return nil
}

// GetLocation busca dados de localização para um CEP fornecido
func (s *Service) GetLocation(cep string) (*ViaCEPResponse, error) {
	// Limpa o CEP
	cepLimpo := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(cep, ".", ""), "-", ""), " ", "")

	// Valida o formato do CEP
	if err := s.ValidarCEP(cepLimpo); err != nil {
		return nil, err
	}

	// Constrói a URL
	url := fmt.Sprintf("%s/%s/json/", s.baseURL, cepLimpo)

	// Faz a requisição
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar dados do CEP: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler corpo da resposta: %w", err)
	}

	// Verifica erros HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("código de status inesperado: %d", resp.StatusCode)
	}

	// Analisa a resposta
	var viaCEPResp ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEPResp); err != nil {
		return nil, fmt.Errorf("falha ao analisar resposta: %w", err)
	}

	// Verifica se o CEP foi encontrado (ViaCEP retorna "erro": true quando não encontrado)
	if viaCEPResp.Cep == "" || strings.Contains(string(body), `"erro": true`) {
		return nil, fmt.Errorf("can not find zipcode")
	}

	return &viaCEPResp, nil
}

// GetCity retorna apenas o nome da cidade para um CEP fornecido
func (s *Service) GetCity(cep string) (string, error) {
	localizacao, err := s.GetLocation(cep)
	if err != nil {
		return "", err
	}
	return localizacao.Localidade, nil
}
