#!/bin/bash
# Teste de carga na rota /inbox usando autocannon

npx autocannon -c 20 -d 10 -p 10 -m POST \
  -H "Content-Type: application/json" \
  -b '{
    "amount": 100,
    "paymentMethod": "PIX",
    "currencyCode": "BRL",
    "description": "Compra de teste"
  }' \
  http://localhost:8080/transaction &

wait