package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	manipulador := NewWeatherHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	manipulador.HealthCheck(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestExtrairCEP(t *testing.T) {
	testCases := []struct {
		nome     string
		caminho  string
		esperado string
	}{
		{
			nome:     "CEP sem formatação",
			caminho:  "/weather/01310930",
			esperado: "01310930",
		},
		{
			nome:     "CEP com hífen",
			caminho:  "/weather/01310-930",
			esperado: "01310-930",
		},
		{
			nome:     "Caminho vazio",
			caminho:  "/weather/",
			esperado: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := extrairCEP(tc.caminho)
			assert.Equal(t, tc.esperado, resultado)
		})
	}
}

func TestGetWeatherByCEP_CEPInvalido(t *testing.T) {
	manipulador := NewWeatherHandler()

	casosDeTeste := []struct {
		nome string
		cep  string
	}{
		{nome: "Muito curto", cep: "1234567"},
		{nome: "Muito longo", cep: "123456789"},
		{nome: "Com letras", cep: "12345ABC"},
		{nome: "Vazio", cep: ""},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			// Valida diretamente o CEP
			err := manipulador.cepService.ValidarCEP(tc.cep)
			assert.Error(t, err, "Deveria retornar erro para CEP inválido")
		})
	}
}

func TestGetWeatherByCEP_FormatoCEPValido(t *testing.T) {
	// Este teste verifica que CEPs com formato válido não retornam 422
	// Testamos diretamente com a lógica de validação de CEP
	manipulador := NewWeatherHandler()

	// Testa que 00000000 passa na validação
	err := manipulador.cepService.ValidarCEP("00000000")
	assert.NoError(t, err, "00000000 deve ser um formato de CEP válido")

	// Testa que 12121212 passa na validação
	err = manipulador.cepService.ValidarCEP("12121212")
	assert.NoError(t, err, "12121212 deve ser um formato de CEP válido")
}

func TestErrorResponseFormat(t *testing.T) {
	manipulador := NewWeatherHandler()

	// Testa validação direta de CEP inválido
	err := manipulador.cepService.ValidarCEP("abc")
	assert.Error(t, err)
	assert.Equal(t, "invalid zipcode", err.Error())
}
