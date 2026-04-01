package temperature

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelsiusParaFahrenheit(t *testing.T) {
	conversor := NewConverter()

	casosDeTeste := []struct {
		nome     string
		celsius  float64
		esperado float64
	}{
		{
			nome:     "Ponto de congelamento da água",
			celsius:  0,
			esperado: 32,
		},
		{
			nome:     "Ponto de ebulição da água",
			celsius:  100,
			esperado: 212,
		},
		{
			nome:     "Temperatura ambiente",
			celsius:  25,
			esperado: 77,
		},
		{
			nome:     "Temperatura negativa",
			celsius:  -10,
			esperado: 14,
		},
		{
			nome:     "Temperatura corporal",
			celsius:  37,
			esperado: 98.6,
		},
		{
			nome:     "Exemplo da especificação (28.5°C)",
			celsius:  28.5,
			esperado: 83.3,
		},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := conversor.CelsiusParaFahrenheit(tc.celsius)
			assert.InDelta(t, tc.esperado, resultado, 0.1, "Conversão para Fahrenheit incorreta")
		})
	}
}

func TestCelsiusParaKelvin(t *testing.T) {
	conversor := NewConverter()

	casosDeTeste := []struct {
		nome     string
		celsius  float64
		esperado float64
	}{
		{
			nome:     "Ponto de congelamento da água",
			celsius:  0,
			esperado: 273,
		},
		{
			nome:     "Ponto de ebulição da água",
			celsius:  100,
			esperado: 373,
		},
		{
			nome:     "Temperatura ambiente",
			celsius:  25,
			esperado: 298,
		},
		{
			nome:     "Temperatura corporal",
			celsius:  37,
			esperado: 310,
		},
		{
			nome:     "Exemplo da especificação (28.5°C)",
			celsius:  28.5,
			esperado: 301.5,
		},
		{
			nome:     "Próximo ao zero absoluto",
			celsius:  -270,
			esperado: 3,
		},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := conversor.CelsiusParaKelvin(tc.celsius)
			assert.InDelta(t, tc.esperado, resultado, 0.1, "Conversão para Kelvin incorreta")
		})
	}
}

func TestConvert(t *testing.T) {
	conversor := NewConverter()

	// Teste com o exemplo da especificação
	resultado := conversor.Convert(28.5)

	assert.Equal(t, 28.5, resultado.TempC, "Celsius deve ser 28.5")
	assert.InDelta(t, 83.3, resultado.TempF, 0.1, "Fahrenheit deve ser ~83.3")
	assert.InDelta(t, 301.5, resultado.TempK, 0.1, "Kelvin deve ser ~301.5")
}

func TestConvertArredondamento(t *testing.T) {
	conversor := NewConverter()

	casosDeTeste := []struct {
		nome      string
		entrada   float64
		esperadoC float64
		esperadoF float64
		esperadoK float64
	}{
		{
			nome:      "Já arredondado",
			entrada:   28.5,
			esperadoC: 28.5,
			esperadoF: 83.3,
			esperadoK: 301.5,
		},
		{
			nome:      "Precisa arredondar para baixo",
			entrada:   28.555,
			esperadoC: 28.56,
			esperadoF: 83.4,
			esperadoK: 301.56,
		},
		{
			nome:      "Precisa arredondar para cima",
			entrada:   28.554,
			esperadoC: 28.55,
			esperadoF: 83.4,
			esperadoK: 301.55,
		},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := conversor.Convert(tc.entrada)
			assert.Equal(t, tc.esperadoC, resultado.TempC, "Arredondamento de Celsius incorreto")
			assert.InDelta(t, tc.esperadoF, resultado.TempF, 0.1, "Arredondamento de Fahrenheit incorreto")
			assert.InDelta(t, tc.esperadoK, resultado.TempK, 0.1, "Arredondamento de Kelvin incorreto")
		})
	}
}

func TestFormulasDeTemperatura(t *testing.T) {
	// Testa se as fórmulas estão corretas de acordo com a especificação
	// F = C * 1.8 + 32
	// K = C + 273

	conversor := NewConverter()
	celsius := 28.5

	// Calcula valores esperados usando as fórmulas da especificação
	esperadoF := celsius*1.8 + 32
	esperadoK := celsius + 273

	realF := conversor.CelsiusParaFahrenheit(celsius)
	realK := conversor.CelsiusParaKelvin(celsius)

	assert.InDelta(t, esperadoF, realF, 0.001, "Fórmula F incorreta")
	assert.InDelta(t, esperadoK, realK, 0.001, "Fórmula K incorreta")
}

func TestTemperaturasNegativas(t *testing.T) {
	conversor := NewConverter()

	casosDeTeste := []struct {
		nome    string
		celsius float64
	}{
		{
			nome:    "Abaixo de zero",
			celsius: -10.5,
		},
		{
			nome:    "Negativo pequeno",
			celsius: -5,
		},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			resultado := conversor.Convert(tc.celsius)

			// Kelvin deve sempre ser positivo
			assert.True(t, resultado.TempK > 0, "Kelvin deve ser positivo para %v°C", tc.celsius)

			// F deve ser maior que C para temperaturas acima de -40°F
			assert.True(t, resultado.TempF > resultado.TempC, "F deve ser maior que C")

			// K deve ser maior que C
			assert.True(t, resultado.TempK > resultado.TempC, "K deve ser maior que C")
		})
	}
}

func TestTemperaturasAltas(t *testing.T) {
	conversor := NewConverter()

	// Teste com temperatura muito alta
	resultado := conversor.Convert(1000)

	assert.Equal(t, 1000.0, resultado.TempC)
	assert.Equal(t, 1832.0, resultado.TempF) // 1000 * 1.8 + 32 = 1832
	assert.Equal(t, 1273.0, resultado.TempK) // 1000 + 273 = 1273
}

func TestZeroDecimais(t *testing.T) {
	conversor := NewConverter()

	resultado := conversor.Convert(0)

	assert.Equal(t, 0.0, resultado.TempC)
	assert.Equal(t, 32.0, resultado.TempF)
	assert.Equal(t, 273.0, resultado.TempK)
}

// Testes de benchmark
func BenchmarkCelsiusParaFahrenheit(b *testing.B) {
	conversor := NewConverter()
	for i := 0; i < b.N; i++ {
		conversor.CelsiusParaFahrenheit(28.5)
	}
}

func BenchmarkCelsiusParaKelvin(b *testing.B) {
	conversor := NewConverter()
	for i := 0; i < b.N; i++ {
		conversor.CelsiusParaKelvin(28.5)
	}
}

func BenchmarkConvert(b *testing.B) {
	conversor := NewConverter()
	for i := 0; i < b.N; i++ {
		conversor.Convert(28.5)
	}
}

// Função auxiliar para verificar se dois floats são aproximadamente iguais
func saoAproximadamenteIguais(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
