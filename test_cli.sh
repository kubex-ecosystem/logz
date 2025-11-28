#!/bin/bash

# Script de teste do Logz CLI
# Demonstra os principais recursos da ferramenta

LOGZ="./dist/logz_linux_amd64"

echo "================================================"
echo "  TESTANDO LOGZ CLI"
echo "================================================"
echo ""

# Verificar se o binário existe
if [ ! -f "$LOGZ" ]; then
    echo "❌ Binário não encontrado. Execute 'make build-dev' primeiro."
    exit 1
fi

echo "🔹 Teste 1: Diferentes níveis de log com formato texto"
echo "----------------------------------------------------"
$LOGZ -l debug -e "Mensagem de debug para diagnóstico"
$LOGZ -l info -e "Aplicação iniciada com sucesso"
$LOGZ -l success -e "Operação concluída com êxito"
$LOGZ -l warn -e "Memória acima de 80%"
$LOGZ -l error -e "Falha ao conectar ao banco de dados"
echo ""

echo "🔹 Teste 2: Formato JSON (útil para integração com sistemas)"
echo "----------------------------------------------------"
$LOGZ -l info -f json -e "Request processado"
$LOGZ -l error -f json -e "Timeout na requisição externa"
echo ""

echo "🔹 Teste 3: Simulando logs de uma aplicação real"
echo "----------------------------------------------------"
$LOGZ -l info -e "Servidor HTTP iniciado na porta 8080"
$LOGZ -l debug -e "Conectando ao Redis em localhost:6379"
$LOGZ -l success -e "Pool de conexões criado: 10 conexões ativas"
$LOGZ -l warn -e "Taxa de requisições acima do normal: 1500 req/s"
$LOGZ -l error -e "Falha ao processar pagamento: gateway timeout"
echo ""

echo "🔹 Teste 4: JSON formatado para análise detalhada"
echo "----------------------------------------------------"
$LOGZ -l info -f json -e "Autenticação realizada"
$LOGZ -l warn -f json -e "Tentativa de acesso negado"
$LOGZ -l success -f json -e "Deploy finalizado"
echo ""

echo "================================================"
echo "  TESTES CONCLUÍDOS"
echo "================================================"
echo ""
echo "💡 Dicas:"
echo "  - Use -l para definir o nível (debug, info, warn, error, etc)"
echo "  - Use -f para definir o formato (text ou json)"
echo "  - Use -e para a mensagem de log"
echo "  - Execute '$LOGZ --help' para mais opções"
echo ""
