# TierMaster

TierMaster is a leaderboard microservice project built with Go and Redis. I decided to do this to learn about Redis (and to go without using a relational database for once). Plus, microservices are really cool too :)

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
.
├── cmd
│   └── main.go             # Entry point
├── docker-compose.yml      # Docker configuration for running the Redis ins
├── Dockerfile              # For building and deploying the microservice
├── docs                    # OpenAPI spec
├── go.mod                  # Contains metadata for 3p packages                
├── go.sum  
├── internal
│   ├── config
│   │   └── config.go       # Util for loading config from .env files with Viper
│   ├── repositories        
│   └── service
└── README.md

````

# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* Go (version 1.20+)
* Docker & Docker Compose (for local development with Redis)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/tiermaster.git](https://github.com/your-username/tiermaster.git)
    cd tiermaster
    ```

2.  **Set up Environment Variables:**
    Create a `.env` file in the root directory based on `.env.example`.
    ```dotenv
    # Example .env file
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=
    REDIS_DB=0
    SERVICE_PORT=8080
    ```

3.  **Run with Docker Compose (Recommended for local dev):**
    This will start both the Redis instance and the Go microservice.
    ```bash
    docker-compose -f docker/docker-compose.yml up --build
    ```
    The service should be accessible at `http://localhost:8080`.

4.  **Run Natively (without Docker for the Go service):**
    * **Start a Redis instance:** Ensure you have a Redis server running (e.g., `redis-server` if installed locally, or via a Docker container).
    * **Build and run the Go application:**
        ```bash
        go build -o bin/tiermaster-service ./cmd/tiermaster-service
        ./bin/tiermaster-service
        ```

# API Endpoints

Once the service is running, you can interact with it via the following HTTP endpoints. (Assuming `http://localhost:8080` as the base URL).

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

# License

This project is licensed under the MIT License - see the `LICENSE` file for details.

# Acknowledgments

  * Redis documentation and community for excellent resources.
  * Go community for amazing libraries and support.

