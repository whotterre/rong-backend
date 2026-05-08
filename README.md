# TierMaster

TierMaster is a leaderboard service built with Go and Redis. The project has been consolidated into a single monolith service (formerly separated into a gateway + leaderboard).

# Tech Stack
* **Go:** The primary language for building the microservice, chosen for its performance, concurrency features, and strong tooling.
* **Redis:** Utilized as the primary data store for leaderboard data, specifically leveraging Redis Sorted Sets for efficient ranking and retrieval.

# Features

* **Add/Update Score:** Atomically add or update a user's score on the leaderboard.
* **Get Top N:** Retrieve the top `N` users from the leaderboard.
* **Get User Rank & Score:** Fetch a specific user's current rank and score.
* **Remove User:** Remove a user entirely from the leaderboard.
* **Scalable:** Designed to be horizontally scalable to handle high request volumes.
* **Lightweight:** Built with a minimal footprint, typical of Go services.

# Project Structure

The project follows a clean, modular structure to ensure maintainability and clear separation of concerns:

```
tiermaster/
├── README.md
├── app.env
├── cmd
│   ├── main.go
│   └── setupRoutes.go
├── Dockerfile
├── go.mod
├── go.sum
└── internal
    ├── config
    │   └── config.go
    ├── conn
    │   └── conn.go
    ├── handlers
    │   └── lb.go
    ├── models
    │   └── user.go
    ├── repositories
    │   └── lb.go
    └── services
        └── lb.go
````

# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.20+)
- Docker & Docker Compose (for local development with Redis)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/tiermaster.git](https://github.com/your-username/tiermaster.git)
    cd tiermaster
    ```

2.  **Set up Environment Variables:**
    Copy or edit `app.env` in the repository root and adjust any values (Redis address, ports, etc.).

3.  **Run with Docker (Recommended for local dev):**
    Build and run the monolith container from the repository root.
    ```bash
    docker build -t tiermaster .
    docker run --env-file app.env -p 3001:3001 tiermaster
    ```
    The service should be accessible at `http://localhost:3001` (or the port defined in `app.env`).

4.  **Run Natively (without Docker for the Go service):**
    * **Start a Redis instance:** Ensure you have a Redis server running locally or via Docker.
    * **Build and run the Go application:**
        ```bash
        go build -o bin/tiermaster ./cmd
        ./bin/tiermaster
        ```

# API Endpoints

Once the service is running, you can interact with it via the following HTTP endpoints. (Assuming `http://localhost:3001` as the base URL).

* **Add or Update User Score**
    * `POST /scores`
    * **Request Body:**
        ```json
        {
            "user_id": "user123",
            "score": 1500
        }
        ```
    * **Response:** `200 OK` (or `400 Bad Request` if invalid input)

* **Get Top N Leaderboard Entries**
    * `GET /leaderboard?limit=10` (default limit, e.g., 10 or 100)
    * **Query Parameters:**
        * `limit`: (Optional) The number of top entries to return.
    * **Response:**
        ```json
        [
            {"user_id": "userABC", "score": 2500, "rank": 1},
            {"user_id": "userXYZ", "score": 2450, "rank": 2}
        ]
        ```

* **Get User's Rank and Score**
    * `GET /leaderboard/{user_id}`
    * **Example:** `GET /leaderboard/user123`
    * **Response:**
        ```json
        {
            "user_id": "user123",
            "score": 1500,
            "rank": 5
        }
        ```
        (Returns `404 Not Found` if user not on leaderboard)

* **Remove User from Leaderboard**
    * `DELETE /leaderboard/{user_id}`
    * **Example:** `DELETE /leaderboard/user123`
    * **Response:** `200 OK`

* **Health Check**
    * `GET /health`
    * **Response:** `200 OK` if the service is running and can connect to Redis.


# Contributing

Feel free to open issues or pull requests. Contributions are welcome\!

