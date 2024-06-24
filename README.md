# Bookstore Service

This repository contains a simple bookstore service written in Go, using the Gin framework, sqlx for database operations, and MySQL as the database.

## Getting Started

### Prerequisites

- Go (>= 1.22)
- MySQL

### Installation

1. **Clone the repository and navigate to the directory:**
   clone this repo
   cd bookstore
   cp .env.example .env

2. **Start the MySQL Docker container:**
   in different terminal run this to start the docker
    ```sh
    docker-compose up -d 
    ```

3. **Install dependencies:**
   run : go mod tidy

4. **Initialize the database schema:**
   run : make init-db

5. **Seed the database:**
   run : make seed-db

6. **Running Unit Test:**
   run : make test

7. **Start the service:**
   run : make run

8. **Trying some endpoints:**
   You can find the Postman collection in the database folder. Import it into Postman to test the API endpoints.


