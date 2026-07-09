# Red-Black Tree Dictionary API

A high-performance English-Russian dictionary service built from scratch in Go. 
The core of the application is a custom implementation of a **self-balancing Red-Black Tree** for efficient O(log n) search, insertion, and deletion operations.

## Features

- **Custom Data Structure:** Hand-written Red-Black tree implementation with manual pointer management, node recoloring, and tree rotations.
- **Thread Safety:** Fully concurrent REST API and CLI, protected against data races using `sync.RWMutex`.
- **Database Persistence:** PostgreSQL integration with transactional data saving and conflict resolution (`ON CONFLICT DO NOTHING`).
- **Graceful Shutdown:** Proper resource cleanup and database connection handling.
- **Interactive CLI:** Built-in console menu for direct interaction, tree validation, and file I/O operations.
- **Dockerized:** Ready to deploy via Docker Compose using a lightweight Alpine-based image.

## 🛠 Tech Stack

- **Language:** Go 1.22
- **Database:** PostgreSQL 15
- **Infrastructure:** Docker, Docker Compose
- **Architecture:** Layered design (Storage Interface -> RBT/Postgres -> Dictionary API)

## 🚀 How to Run

### Using Docker (Recommended)
The easiest way to run the project with its database is using Docker Compose.

```bash
# Clone the repository
git clone [https://github.com/g0disd3ad/rbt.git](https://github.com/g0disd3ad/rbt.git)
cd rbt

# Start the application and database
docker-compose up -d --build

# Attach to the container to use the interactive CLI
docker attach rbt-api

```

### Local Setup

If you want to run it without Docker, ensure PostgreSQL is running and your `.env` file is configured.

```bash
# Export environment variables or use a tool like godotenv
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=dict_db

# Build and run
go run main.go

```

## 📡 REST API Endpoints

The service exposes a RESTful API on port `8080`.

### 1. Translate a word (GET)

Retrieves the translation(s) for a given English word.

**Request:**
`GET /translate?word=tree`

**Response:**

```json
{
  "word": "tree",
  "translations": ["дерево"]
}

```

### 2. Add a translation (POST)

Adds a new English-Russian translation pair to the dictionary.

**Request:**
`POST /translate`

```json
{
  "eng": "software engineering",
  "rus": "программная инженерия"
}

```

**Response:** `201 Created`

```json
{
  "status": "success"
}

```

## 🧪 Tree Validation

The project includes a built-in strict validator to ensure the structural integrity of the Red-Black Tree. You can trigger it via the CLI menu (Option 7) to check the current tree height and verify all Red-Black properties (e.g., black-height balance, no consecutive red nodes).

```