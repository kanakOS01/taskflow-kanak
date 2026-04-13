# TaskFlow API

## 1. Overview
TaskFlow provides a RESTful API for managing users, projects, and tasks. It features secure authentication, ownership-based access control, and project statistics.

**Tech Stack:**
- **Language:** Go (Golang) 1.22+
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL 15
- **Authentication:** JWT (JSON Web Tokens)
- **Containerization:** Docker & Docker Compose
- **Migrations:** golang-migrate

## 2. Architecture Decisions
- **Clean Architecture Pattern:** The project is structured into `Handler`, `Service`, and `Repository` layers to ensure separation of concerns and testability.
- **Feature-Based Folder Structure:** Code is organized by domain features (e.g., `auth`, `projects`, `tasks`) in the `internal/` directory to promote encapsulation and reduce coupling between different parts of the system.
- **Stateless Authentication:** JWT is used to avoid session management overhead and allow the API to scale horizontally easily.
- **Postgres-Specific Features:** `pgxpool` is used as the database driver to leverage PostgreSQL-specific benefits and performance optimizations.
- **Database Migrations:** Used `golang-migrate` for versioned schema management, ensuring database consistency across all environments and enabling automated migrations via Docker.
- **Ownership-Based Access:** All project and task operations are restricted based on ownership (e.g., only project owners can create tasks), ensuring data security by design.



## 3. Running Locally
Assume you have **Docker** installed. No other dependencies are required.

```bash
# 1. Clone the repository
git clone https://github.com/kanakOS01/taskflow-kanak
cd taskflow-kanak

# 2. Start the application
docker compose up --build

# 3. Access the API
# The server will be available at http://localhost:8000
```

## 4. Running Migrations
Migrations run **automatically** on startup via the Docker entrypoint script. You don't need to run any manual commands. 

If you need to check migration status manually inside the container:
```bash
docker compose exec api migrate -path /app/migrations -database $DATABASE_URL version
```

## 5. Test Credentials
Use these credentials to test the API immediately without registration:

| Email                                     | Password  |
| ----------------------------------------- | --------- |
| [kanak@email.com](mailto:kanak@email.com) | 123456789 |
| [test@email.com](mailto:test@email.com)   | 123456789 |



## 6. API Reference
A detailed breakdown of all available endpoints, request schemas, and response examples can be found in:
👉 **[API_REFERENCE.md](./API_REFERENCE.md)**

The API spec can be found at [OpenAPI Spec](./api_spec.json)

## 7. What I'd Do With More Time
- **Filtering & Search:** Add advanced query params for tasks (search by title, filter by date ranges).
- **Soft Deletes & Audit Logs:** Implement a `deleted_at` pattern and track every change made to a project or task.
- **Security & Auth:**
    - Implement better security logic (e.g., forcing users to join projects to interact with tasks).
    - Add Role-Based Access Control (RBAC) and refresh tokens.
- **Database & Performance:** Implement DB transaction support, optimize SQL queries, and add more `.env` configs (e.g., pool size).
- **Documentation & Testing:** Add Swagger/OpenAPI documentation and increase code coverage with comprehensive unit/integration tests.
- **Docker Optimization:** Further optimize Docker images for size and production security.
