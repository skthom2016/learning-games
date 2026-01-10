# Docker Setup Guide

This document explains how to run the Learning Game Platform using Docker Compose.

## Prerequisites

- Docker Desktop installed and running
- Docker Compose (comes with Docker Desktop)

## Quick Start

### First Time Setup (Initial Database Creation)

```bash
# From the project root directory
docker-compose up -d
```

**What happens on first startup:**
1. PostgreSQL container starts with empty volume
2. Migration scripts run automatically in order:
   - `001_init_schema.sql` - Creates tables
   - `002_seed_multiplication_game.sql` - Adds multiplication game
   - `003_seed_division_game.sql` - Adds division game
3. Backend container starts and connects to database
4. Frontend container starts and connects to backend

**This takes about 30-60 seconds on first run.**

### Accessing the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Subsequent Starts

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

**Important**: Your data is persisted! The `postgres_data` volume survives:
- `docker-compose down`
- Container restarts
- System reboots

### Starting Fresh (Complete Reset)

⚠️ **WARNING: This deletes ALL data including players, progress, and rewards!**

```bash
# Stop and remove containers, networks, and volumes
docker-compose down -v

# Remove dangling images (optional)
docker system prune -a

# Start fresh (will run migrations again)
docker-compose up -d
```

## Database Persistence Explained

### How It Works

1. **Named Volume**: `postgres_data` is a Docker-managed volume
2. **Mount Point**: `/var/lib/postgresql/data` in the container
3. **Persistence**: Data survives all container operations
4. **Migration Detection**: Migrations only run when volume is empty

### Checking Your Data

```bash
# View all volumes
docker volume ls

# Inspect postgres_data volume
docker volume inspect division_postgres_data

# Access database directly
docker exec -it learning-game-db psql -U gameuser -d learning_game

# List all tables
\dt

# View players
SELECT * FROM players;

# Exit psql
\q
```

## Development Features

### Hot Reload

Both frontend and backend support hot-reload:

- **Backend**: Changes to Go files trigger automatic rebuild (using Air)
- **Frontend**: Changes to React files trigger automatic browser refresh

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db
```

### Entering Containers

```bash
# Backend shell
docker exec -it learning-game-backend sh

# Frontend shell
docker exec -it learning-game-frontend sh

# Database shell
docker exec -it learning-game-db psql -U gameuser -d learning_game
```

## Troubleshooting

### Database Connection Issues

```bash
# Check if database is healthy
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Port Conflicts

If ports 3000, 8080, or 5432 are already in use:

Edit `docker-compose.yml` and change the port mappings:
```yaml
ports:
  - "3001:3000"  # Use 3001 instead of 3000
```

### Migration Issues

If migrations didn't run properly:

```bash
# Check if migrations exist
ls backend/migrations/

# Verify they're mounted in container
docker exec learning-game-db ls /docker-entrypoint-initdb.d/

# Reset and start fresh
docker-compose down -v
docker-compose up -d
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Rebuild without cache
docker-compose build --no-cache
docker-compose up -d
```

## Architecture

```
┌─────────────────┐
│   Browser       │
│  localhost:3000 │
└────────┬────────┘
         │
┌────────▼────────┐
│   Frontend      │
│  (React/Nginx)  │
│  Port: 3000     │
└────────┬────────┘
         │ HTTP
┌────────▼────────┐
│   Backend       │
│  (Go/Gin)       │
│  Port: 8080     │
└────────┬────────┘
         │ TCP
┌────────▼────────┐
│  PostgreSQL     │
│   Database      │
│  Port: 5432     │
└─────────────────┘
    postgres_data
    (Volume)
```

## Environment Variables

### Backend
- `DB_HOST`: Database hostname (db)
- `DB_PORT`: Database port (5432)
- `DB_USER`: Database user (gameuser)
- `DB_PASSWORD`: Database password (gamepass123)
- `DB_NAME`: Database name (learning_game)
- `SERVER_PORT`: Server port (8080)

### Frontend
- `REACT_APP_API_URL`: Backend API URL (http://localhost:8080)

## Security Notes

⚠️ **This is a development setup!** For production:

1. Change all default passwords
2. Use environment file (`.env`) instead of hardcoded values
3. Remove database port exposure (`5432:5432`)
4. Use proper secrets management
5. Enable SSL/TLS
6. Implement authentication/authorization
7. Use production-ready images and builds

## Next Steps

After Docker setup is complete:

1. Visit http://localhost:3000
2. Create a player
3. Start playing multiplication or division games
4. View progress and stars earned

## Support

For issues or questions:
- Check logs: `docker-compose logs -f`
- Verify container status: `docker-compose ps`
- Check documentation: `CLAUDE.md`
