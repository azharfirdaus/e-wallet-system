# E Wallet System

## Overview

This project is a Go-based e-wallet API backed by PostgreSQL. It supports user wallet creation, multi-currency wallet balances, top up, payment, transfer, wallet suspension, and closing-balance rollups from ledger entries.

## Features

### Functional

- Create users and wallets
- Add multiple currencies to one wallet
- Top up wallet balance
- Pay using wallet balance
- Transfer balance between wallets
- Suspend wallets
- View wallet balances by currency and status
- Update closing balances from ledger entries

### Non Functional

- PostgreSQL-backed persistence
- Transactional wallet operations
- Ledger-based balance calculation
- Advisory locking for wallet write consistency
- Docker Compose setup for local development

## Decisions

- Amounts use two decimal digits and are stored as minor units to prevent unnecessary floating-point rounding mismatches.
- Rounding up or down is not used because it can introduce floating-point overflow or precision issues that may write incorrect ledger records.
- `BIGINT` is used for monetary values so minor units can be preserved without floating-point overflow risk.

## Run with Docker Compose

Start the app and PostgreSQL:

```sh
docker compose up --build
```

The app runs at `http://localhost:8080`.

PostgreSQL runs on `localhost:5432` with:

- database: `e_wallet`
- user: `admin`
- password: `password`

Stop the containers:

```sh
docker compose down
```

## API Curl Examples

Set the base URL:

```sh
BASE_URL=http://localhost:8080
```

Create a user:

```sh
curl -X POST "$BASE_URL/users"
```

Create a wallet:

```sh
curl -X POST "$BASE_URL/wallets" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "currency_code": "IDR"
  }'
```

Get wallet balance:

```sh
curl "$BASE_URL/wallets/1"
```

Top up wallet:

```sh
curl -X POST "$BASE_URL/wallets/1/topup" \
  -H "Content-Type: application/json" \
  -d '{
    "currency_code": "IDR",
    "amount": "10000.00"
  }'
```

Pay with wallet:

```sh
curl -X POST "$BASE_URL/wallets/1/pay" \
  -H "Content-Type: application/json" \
  -d '{
    "currency_code": "IDR",
    "amount": "2500.00"
  }'
```

Transfer to another wallet:

```sh
curl -X POST "$BASE_URL/wallets/1/transfer/IDR" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_id_destination": 2,
    "currency_code_destination": "IDR",
    "amount": "5000.00"
  }'
```

Suspend wallet:

```sh
curl -X POST "$BASE_URL/wallets/1/suspend"
```

Update closing balances:

```sh
curl -X POST "$BASE_URL/wallet/update_close_balance"
```
