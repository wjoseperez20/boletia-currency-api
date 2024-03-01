# boletia-currency-api

## Overview

This app represents a tech challenge from Boletia. It is a RESTful API for historical currency rate values, developed in
Go, incorporating features such as JWT Authentication, rate limiting, Swagger documentation, caching with Redis, and
database operations through GORM. The application utilizes the Gin Gonic web framework and is containerized using
Docker.

Additionally, it includes a daemon that periodically consults the currency API's "/latest" endpoint every few minutes
and populates the database with the obtained data.

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [Docker + Docker compose](https://docs.docker.com/engine/install/)

### Installation

1. Clone the repository

```bash
git clone https://github.com/wjoseperez20/boletia-currency-api.git
```

2. Navigate to the directory

```bash
cd boletia-currency-api
```

3. Build and run the Docker containers

```bash
make setup && make build && make up
```

### Environment Variables

Local: You need to put some important information into a file called .env on your computer.

- `POSTGRES_HOST`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `JWT_SECRET`
- `API_SECRET_KEY`
- `DAEMON_WAKEUP`
- `CURRENCY_API_ENDPOINT`
- `CURRENCY_API_KEY`
- `CURRENCY_API_TIMEOUT`

In .env.sample you can find an example of the .env file.

### API Documentation

The API is documented using Swagger and can be accessed at:

```
http://localhost:8001/swagger/index.html
```

## Usage

### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

```bash
curl -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8001/api/v1/currencies
```
