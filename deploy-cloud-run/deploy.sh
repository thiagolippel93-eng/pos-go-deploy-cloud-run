#!/bin/bash

# Script de deploy para a API Weather CEP no Google Cloud Run
# Este script compila, envia para o Container Registry e faz deploy no Cloud Run

set -e

# Configuração
PROJECT_ID="${GCP_PROJECT_ID:-}"
REGION="southamerica-east1"
SERVICE_NAME="weather-cep"

# Verifica se o ID do projeto está definido
if [ -z "$PROJECT_ID" ]; then
    echo "Erro: Variável de ambiente GCP_PROJECT_ID não está definida"
    echo "Por favor, defina com: export GCP_PROJECT_ID=seu-projeto-id"
    exit 1
fi

# Verifica se a chave da WeatherAPI está definida
if [ -z "$WEATHER_API_KEY" ]; then
    echo "Erro: Variável de ambiente WEATHER_API_KEY não está definida"
    echo "Por favor, defina com: export WEATHER_API_KEY=sua-chave-da-weather-api"
    exit 1
fi

echo "Compilando imagem via Cloud Build..."
gcloud builds submit \
    --tag "gcr.io/${PROJECT_ID}/${SERVICE_NAME}" \
    --project "${PROJECT_ID}"

# Faz deploy no Cloud Run
echo "Fazendo deploy no Cloud Run..."
gcloud run deploy "${SERVICE_NAME}" \
    --image "gcr.io/${PROJECT_ID}/${SERVICE_NAME}" \
    --region "${REGION}" \
    --platform managed \
    --allow-unauthenticated \
    --set-env-vars "WEATHER_API_KEY=${WEATHER_API_KEY}" \
    --memory 256Mi \
    --cpu 1 \
    --min-instances 0 \
    --max-instances 10 \
    --timeout 30s \
    --project "${PROJECT_ID}"

# Obtém a URL do serviço
SERVICE_URL=$(gcloud run services describe "${SERVICE_NAME}" --region "${REGION}" --format "value(status.url)")
echo ""
echo "=========================================="
echo "Deploy concluído!"
echo "URL do Serviço: ${SERVICE_URL}"
echo "Verificação de saúde: ${SERVICE_URL}/health"
echo "Exemplo de chamada: ${SERVICE_URL}/weather/01310930"
echo "=========================================="
