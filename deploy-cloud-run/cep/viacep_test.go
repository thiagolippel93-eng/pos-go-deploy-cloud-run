package cep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidarCEP(t *testing.T) {
	servico := NewService()

	casosDeTeste := []struct {
		nome         string
		cep          string
		esperaErro   bool
		mensagemErro string
	}{
		{
			nome:       "CEP válido sem formatação",
			cep:        "01310930",
			esperaErro: false,
		},
		{
			nome:       "CEP válido com hífen",
			cep:        "01310-930",
			esperaErro: false,
		},
		{
			nome:       "CEP válido com ponto",
			cep:        "01310.930",
			esperaErro: false,
		},
		{
			nome:       "CEP válido com espaços",
			cep:        "01310 930",
			esperaErro: false,
		},
		{
			nome:         "CEP muito curto",
			cep:          "1234567",
			esperaErro:   true,
			mensagemErro: "invalid zipcode",
		},
		{
			nome:         "CEP muito longo",
			cep:          "123456789",
			esperaErro:   true,
			mensagemErro: "invalid zipcode",
		},
		{
			nome:         "CEP com letras",
			cep:          "12345ABC",
			esperaErro:   true,
			mensagemErro: "invalid zipcode",
		},
		{
			nome:         "CEP com caracteres especiais",
			cep:          "1234@567",
			esperaErro:   true,
			mensagemErro: "invalid zipcode",
		},
		{
			nome:         "CEP vazio",
			cep:          "",
			esperaErro:   true,
			mensagemErro: "invalid zipcode",
		},
		{
			nome:       "CEP com 8 dígitos em formato misto",
			cep:        "01.310-930",
			esperaErro: false,
		},
	}

	for _, tc := range casosDeTeste {
		t.Run(tc.nome, func(t *testing.T) {
			err := servico.ValidarCEP(tc.cep)
			if tc.esperaErro {
				assert.Error(t, err)
				assert.Equal(t, tc.mensagemErro, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidarCEPCasosExtremos(t *testing.T) {
	servico := NewService()

	// Testa todos os zeros (deve ser formato válido)
	err := servico.ValidarCEP("00000000")
	assert.NoError(t, err)

	// Testa todos os noves (deve ser formato válido)
	err = servico.ValidarCEP("99999999")
	assert.NoError(t, err)

	// Testa dígito único
	err = servico.ValidarCEP("1")
	assert.Error(t, err)

	// Testa padrão repetido
	err = servico.ValidarCEP("12121212")
	assert.NoError(t, err)
}
