# Golang Wire

This repository contains an sample project using Google Wire for dependency injection in Go.

## Getting Started

These instructions will help you set up and run the project on your local machine.

### Prerequisites

- [Go](https://golang.org/dl/) 1.16 or higher
- [Wire](https://github.com/google/wire)

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
3. Install dependencies:
    ```sh
    go mod download
    ```

### Running the Project

1. In order to generate the Wire code:
    ```sh
    wire
    ```

2. Run the project:
    ```sh
    go run ./cmd/api/main.go
    ```
