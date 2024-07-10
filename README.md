# Golang Wire

This repository contains an sample project using Google Wire for dependency injection in Go.

## Getting Started

These instructions will help you set up and run the project on your local machine.

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher
- [Docker and Docker compose](https://docs.docker.com/engine/install/)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/susilnem/golang-wire.git
    cd golang-wire
    ```

2. Create env file and update the database
    ```sh
    cp .env.sample .env
    ```
3. Setup
    ```
     $ docker-compose build
     $ docker-compose up
     $ docker-compose run --rm migrate
    ```

### Setup Environment
```
DB_HOST=db
DB_NAME=postgres
DB_USER=postgres
DB_PORT=5432
DB_PASSWORD=postgres
```
