
# Voucher Redemption API

A simple REST API for voucher redemption with race condition handling using database-level locking and atomic transaction. 



## Architecture
This project is an implementation of [Uncle Bob’s Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), which is a design pattern that emphasizes separation of concerns, testability, and independence from frameworks or databases. The main idea: business rules should not depend on external systems, like databases, HTTP, or UI; though it still has coupling to SQL in usecase for handling transaction for simplicity's sake.


## Tech Stack

This project keeps the stack simple and vanilla, using Go’s standard library as much as possible. Check go.mod for a complete list of packages.

- **Go 1.23** 
- **PostgreSQL 16** 
- **net/http** 
- [**golang-migrate**](https://github.com/golang-migrate/migrate]=) - for migration
- [**go-playground/validator**](https://github.com/go-playground/validator) - for validation
- [**Viper**](https://github.com/spf13/viper) - for .environment/config handling
- [**lib/pq**](https://github.com/lib/pq) - PostgreSQL's driver (consider using [jackc/pgx](https://github.com/jackc/pgx) as it is actively maintained)


## API Reference
Pretty straightforward. The `api/http.json` file contains the API client collection for tools like Postman or Hoppscotch.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/vouchers` | Create new voucher |
| `GET` | `/api/vouchers/:code` | Get voucher details |
| `POST` | `/api/vouchers/:code/claim` | Claim a voucher |
| `GET` | `/api/vouchers/:code/redemptions` | Get redemption history |

#### Health Check

```http
GET /health
```

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| -         | -    | No parameters required |

#### Create Voucher

```http
POST /api/vouchers
```

| Parameter | Type      | Description                                   |
| :-------- | :----     | :-------------------------------              |
| code      | string    | **Required** - The voucher code you want to create    |
| quota     | int       | **Required** - Number of times this voucher can be redeemed         |
| valid_until_days | int | **Required** - Number of days the voucher is valid from creation |

#### Get Voucher Details

```http
GET /api/vouchers/:code
```

| Parameter | Type   | Description                       |
| :-------- | :---- | :-------------------------------- |
| -      | - | - |

#### Claim Voucher

```http
POST /api/vouchers/:code/claim
```

| Parameter | Type   | Description                       |
| :-------- | :---- | :-------------------------------- |
| X-User-ID      | Header | **Required** - Random User ID to mock JWT |


#### Get Redemption History

```http
GET /api/vouchers/:code/redemptions
```

| Parameter | Type   | Description                         |
| :-------- | :---- | :---------------------------------- |
| -      | - | - |


## Project Structure
```
voucher-redemption-api/
├── api/
│   ├── http.json               # API client collection
│   └── README.md
├── cmd/
│   └── main.go                 # entry point
├── internal/
│   ├── delivery/
│   │   └── handler/            # HTTP handlers (Controller)
│   ├── usecase/                # Business logic
│   ├── repository/             # Database layer
│   ├── entity/                 # Domain models
│   ├── model/                  # DTOs (Request/Response)
│   └── config/                 # Configuration
├── db/
│   └── migrations/             # SQL migrations
├── .env.example
├── Makefile
└── README.md
```
## Database Schema
```sql
CREATE TABLE vouchers (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    quota BIGINT NOT NULL,
    valid_until TIMESTAMP NOT NULL
);

CREATE TABLE redemption_history (
    id BIGSERIAL PRIMARY KEY,
    voucher_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    redeemed_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_voucher FOREIGN KEY(voucher_id) REFERENCES vouchers(id)
);

CREATE UNIQUE INDEX idx_unique_claim ON redemption_history (voucher_id, user_id);
```
## Run Locally
The database is seeded once you run the migration with test vouchers for different scenarios:
| Code                                      | Quota | Valid Until |
| ----------------------------------------- | ----- | ----------- |
| DISKON100 | 10    | +7 days     |
| FLASH50                                   | 100   | +1 day      |
| SOLDOUT                                   | 0     | +5 days     |
| EXPIRED1                                  | 50    | \-1 day     |
| CLAIMED                                   | 50    | +10 days    |


To run the project locally, make sure the following are installed on your device:
- **Go 1.23** 
- **PostgreSQL 16** (for running with your own PostgreSQL)
- **Docker** 
- [**golang-migrate**](https://github.com/golang-migrate/migrate]=) - for migration
- **make**

Clone the project

```bash
  git clone https://github.com/dhiyazhar/voucher-redemption-api.git
  cd voucher-redemption-api
```

Setup environment

```bash
  cp .env.example .env
  # Edit .env with your database credentials
```

Start database using Docker
```bash
  docker-compose up -d postgres
```
or use local PostgreSQL
```bash
  psql -U postgres -c "CREATE DATABASE voucher_db;"
```

Run migrations

```bash
  make migrate-up
```

Start server 
```bash
  make run
```
Server will start at `http://localhost:8080/`


## License

[MIT](https://choosealicense.com/licenses/mit/)

