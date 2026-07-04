# e-wallet-system

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
curl -X POST "$BASE_URL/wallets/1/transfer" \
  -H "Content-Type: application/json" \
  -d '{
    "to_wallet_id": 2,
    "currency_code": "IDR",
    "amount": "5000.00"
  }'
```

Suspend wallet:

```sh
curl -X POST "$BASE_URL/wallets/1/suspend"
```
