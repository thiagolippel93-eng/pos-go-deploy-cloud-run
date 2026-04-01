package temperature

// Converter manipula conversões de temperatura
type Converter struct{}

// NewConverter cria um novo conversor de temperatura
func NewConverter() *Converter {
	return &Converter{}
}

// CelsiusParaFahrenheit converte Celsius para Fahrenheit
// Fórmula: F = C * 1.8 + 32
func (c *Converter) CelsiusParaFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// CelsiusParaKelvin converte Celsius para Kelvin
// Fórmula: K = C + 273
func (c *Converter) CelsiusParaKelvin(celsius float64) float64 {
	return celsius + 273
}

// TemperatureResponse contém dados de temperatura em todas as escalas
type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// Convert retorna temperaturas em todas as escalas a partir de Celsius
func (c *Converter) Convert(celsius float64) TemperatureResponse {
	return TemperatureResponse{
		TempC: arredondarParaDoisDecimais(celsius),
		TempF: arredondarParaDoisDecimais(c.CelsiusParaFahrenheit(celsius)),
		TempK: arredondarParaDoisDecimais(c.CelsiusParaKelvin(celsius)),
	}
}

// arredondarParaDoisDecimais arredonda um float para duas casas decimais
func arredondarParaDoisDecimais(valor float64) float64 {
	return float64(int(valor*100+0.5)) / 100
}
