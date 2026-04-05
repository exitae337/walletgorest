# Wallet Service

REST API для управления кошельком.

## Запуск

```bash
docker compose up --build
```

Сервер: http://localhost:8080

## API

Создать кошелек
```bash
POST /api/v1/wallets
{"initialAmount": 1000}
```
Депозит
``` bash
POST /api/v1/wallet
{"walletId": "uuid", "operationType": "DEPOSIT", "amount": 500}
```
Вывод
```bash
POST /api/v1/wallet
{"walletId": "uuid", "operationType": "WITHDRAW", "amount": 200}
```
Баланс
```bash
GET /api/v1/wallets/{uuid}
```

## Примеры
```bash
# Создать
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"initialAmount": 1000}'

# Пополнить
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{"walletId": "ваш-uuid", "operationType": "DEPOSIT", "amount": 500}'

# Баланс
curl http://localhost:8080/api/v1/wallets/ваш-uuid
```

## Тесты
```bash
go test -v ./...
```

Для тестирования 1000 RPS использовалась утилита wrk:
```bash
wrk -t4 -c1000 -d10s --rate 1000 http://localhost:8080/api/v1/wallets/{wallet_uuid}
```

Для проверки 1000 RPS при обновлении баланса кошелька:

``` test.lua
wrk.method = "POST"
wrk.body = '{"walletId": "93f21b3a-70d7-46fb-a2a5-ae90291d5c33", "operationType": "DEPOSIT", "amount": 10}'
wrk.headers["Content-Type"] = "application/json"
```

Запуск:
``` bash
wrk -t10 -c1000 -d10s -s test.lua http://localhost:8080/api/v1/wallet 
```

Остановка
```bash
docker compose down
```