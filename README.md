# e-wallet-system

## Run with Docker Compose

Start the app and PostgreSQL:

```sh
docker compose up --build
```

The app runs at `http://localhost:8080`.

PostgreSQL runs on `localhost:5432` with:

- database: `e_wallet_system`
- user: `postgres`
- password: `postgres`

Stop the containers:

```sh
docker compose down
```
