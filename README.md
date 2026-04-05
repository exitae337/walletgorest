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

Остановка
```bash
docker compose down
```