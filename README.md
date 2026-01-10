# Learning Game Platform

A kid-friendly adaptive learning game platform for children aged 6-9, initially focused on multiplication practice.

## 🎯 Project Overview

This platform provides:
- **Adaptive Learning**: System automatically focuses on weak areas
- **Reward System**: Stars and achievements motivate learning
- **Child-Safe Design**: No timers, no punishment, encouraging feedback
- **Parent Analytics**: Detailed progress tracking and insights

See [`ARCHITECTURE.md`](./ARCHITECTURE.md) for complete system design.

## 🏗️ Technology Stack

- **Frontend**: React 18 (browser-based UI)
- **Backend**: Golang 1.21 (REST API)
- **Database**: PostgreSQL 15
- **Deployment**: Docker Compose (localhost)

## 📁 Project Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/          # Main application entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers and routes
│   │   ├── database/        # Database connection
│   │   ├── models/          # Data models
│   │   └── services/        # Business logic
│   │       ├── platform/    # Player management
│   │       ├── game/        # Game logic (multiplication)
│   │       ├── learning/    # Adaptive learning engine
│   │       └── rewards/     # Reward system
│   ├── migrations/          # SQL schema files
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── api/             # API client
│   │   ├── screens/
│   │   │   ├── child/       # Child UI screens
│   │   │   └── admin/       # Admin UI screens
│   │   ├── App.js
│   │   └── index.js
│   ├── Dockerfile
│   └── package.json
│
├── docker-compose.yml
├── ARCHITECTURE.md          # Complete system architecture
├── README.md                # This file
└── TODO.md                  # Implementation checklist
```

## 🚀 Quick Start

### Prerequisites

- Docker Desktop or Docker Engine + Docker Compose
- At least 2GB free RAM
- Ports 3000, 8080, and 5432 available

### Setup & Run

1. **Clone or navigate to the project directory**:
   ```bash
   cd /path/to/learning-game-platform
   ```

2. **Start all services**:
   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the backend (Golang API)
   - Build the frontend (React app)
   - Start PostgreSQL database
   - Run database migrations
   - Start all services

3. **Access the application**:
   - **Child UI**: http://localhost:3000
   - **Admin UI**: http://localhost:3000/admin
   - **Backend API**: http://localhost:8080
   - **Health Check**: http://localhost:8080/health

4. **Stop services**:
   ```bash
   docker-compose down
   ```

5. **Reset database** (delete all data):
   ```bash
   docker-compose down -v
   docker-compose up --build
   ```

## 🧑‍💻 Development

### Backend Development

1. Navigate to backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run locally (without Docker):
   ```bash
   # Set environment variables
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=gameuser
   export DB_PASSWORD=gamepass123
   export DB_NAME=learning_game
   export SERVER_PORT=8080

   # Run server
   go run cmd/server/main.go
   ```

4. Hot reload is enabled in Docker via Air (changes auto-reload)

### Frontend Development

1. Navigate to frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Run locally (without Docker):
   ```bash
   npm start
   ```

4. Hot reload is enabled in Docker (changes auto-reload)

## 📝 Implementation Status

This is a **boilerplate project with TODOs**. See [`TODO.md`](./TODO.md) for detailed implementation checklist.

**Current State**:
- ✅ Project structure created
- ✅ Docker setup complete
- ✅ Database schema defined
- ✅ API routes scaffolded
- ✅ React screens scaffolded
- ⚠️ **Business logic NOT implemented** (marked with TODOs)
- ⚠️ **Database queries NOT implemented** (marked with TODOs)
- ⚠️ **Frontend API calls return mock data** (marked with TODOs)

**To Complete MVP**: See [`TODO.md`](./TODO.md)

## 🎮 How to Use (Once Implemented)

### For Children

1. Open http://localhost:3000
2. Select your name (or add new player)
3. Click "Multiplication Master"
4. Answer questions
5. Earn stars!

### For Parents/Teachers

1. Open http://localhost:3000/admin
2. View player progress
3. See mastery heatmap (which tables need practice)
4. Reset rewards when giving real-world gifts

## 🗄️ Database

### Connection Details

- **Host**: localhost
- **Port**: 5432
- **Database**: learning_game
- **User**: gameuser
- **Password**: gamepass123

### Schema Migrations

Located in `backend/migrations/`:
- `001_init_schema.sql` - Creates all tables
- `002_seed_multiplication_game.sql` - Adds multiplication game data

Migrations run automatically on first startup via Docker.

## 🧪 Testing

```bash
# Backend tests (once implemented)
cd backend
go test ./...

# Frontend tests (once implemented)
cd frontend
npm test
```

## 📚 Key Documentation

- **[ARCHITECTURE.md](./ARCHITECTURE.md)**: Complete system design (READ THIS FIRST)
- **[TODO.md](./TODO.md)**: Implementation checklist for developers
- **API Documentation**: See `backend/internal/api/routes.go` for endpoint list

## 🔍 Troubleshooting

### Port Already in Use

If you see "port already allocated" error:
```bash
# Check what's using the port
sudo lsof -i :3000  # or :8080, :5432

# Kill the process or change ports in docker-compose.yml
```

### Database Connection Failed

```bash
# Ensure PostgreSQL container is healthy
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Frontend Not Loading

```bash
# View frontend logs
docker-compose logs frontend

# Rebuild frontend
docker-compose up --build frontend
```

### Backend API Not Responding

```bash
# View backend logs
docker-compose logs backend

# Check if backend is running
curl http://localhost:8080/health

# Rebuild backend
docker-compose up --build backend
```

## 🤝 Contributing

This is a boilerplate for implementation. To contribute:

1. Read `ARCHITECTURE.md` thoroughly
2. Check `TODO.md` for open tasks
3. Pick a TODO item
4. Implement following the architecture specification
5. Test your changes
6. Update TODO.md to mark item complete

## 📄 License

[Add your license here]

## 👥 Authors

[Add your name/team here]

## 🙏 Acknowledgments

- Architecture designed with Claude Sonnet 4.5
- Built for children aged 6-9
- Focus on learning through play
