# 🚀 Quick Start Guide

Get the Learning Game Platform running in 5 minutes!

## Prerequisites Check

Before starting, ensure you have:
- ✅ Docker Desktop installed (or Docker Engine + Docker Compose)
- ✅ At least 2GB free RAM
- ✅ Ports 3000, 8080, and 5432 available

## Step-by-Step Setup

### 1. Navigate to Project Directory

```bash
cd /home/santhosh/ubuntu-latestdev/games/division
```

### 2. Start Everything

```bash
docker-compose up --build
```

**What's happening?**
- Building Go backend (may take 1-2 minutes first time)
- Building React frontend (may take 2-3 minutes first time)
- Starting PostgreSQL database
- Running database migrations
- Loading seed data (multiplication game)

**Wait for these messages**:
```
backend    | Server starting on port 8080...
frontend   | webpack compiled successfully
db         | database system is ready to accept connections
```

### 3. Open Your Browser

**For Children:**
- Go to: http://localhost:3000
- Click "Add Player"
- Enter name (e.g., "Emma")
- Click "Multiplication Master"
- Start playing!

**For Parents/Teachers:**
- Go to: http://localhost:3000/admin
- View player progress
- See mastery heatmaps
- Reset rewards

### 4. Verify Everything Works

**Health Check:**
```bash
curl http://localhost:8080/health
```

Should return:
```json
{"status":"ok","service":"learning-game-api","version":"1.0.0"}
```

**Check Database:**
```bash
docker-compose exec db psql -U gameuser -d learning_game -c "SELECT * FROM games;"
```

Should show the multiplication game.

## Common Issues

### "Port 3000 is already in use"

```bash
# Find what's using it
lsof -i :3000

# Kill it
kill -9 [PID]

# Or change the port in docker-compose.yml
```

### "Cannot connect to database"

```bash
# Wait 10 seconds for database to initialize, then:
docker-compose restart backend
```

### "Frontend won't load"

```bash
# Check frontend logs
docker-compose logs frontend

# Hard refresh browser: Ctrl+Shift+R (Cmd+Shift+R on Mac)
```

## Stopping the Application

```bash
# Stop services (keeps data)
docker-compose down

# Stop and delete all data (fresh start)
docker-compose down -v
```

## Next Steps

1. **Read the architecture**: Open `ARCHITECTURE.md`
2. **Start implementing**: Open `TODO.md`
3. **Check boilerplate code**: See files in `backend/` and `frontend/`

## Development Workflow

### Making Backend Changes

1. Edit files in `backend/`
2. Watch logs: `docker-compose logs -f backend`
3. Changes auto-reload (Air hot reload)
4. Test endpoint: `curl http://localhost:8080/api/[endpoint]`

### Making Frontend Changes

1. Edit files in `frontend/src/`
2. Browser auto-refreshes
3. Check console for errors (F12)

### Database Changes

```bash
# Connect to database
docker-compose exec db psql -U gameuser -d learning_game

# View tables
\dt

# Query players
SELECT * FROM players;

# Exit
\q
```

## Testing Your Implementation

### Test Backend

```bash
cd backend
go test ./...
```

### Test Frontend

```bash
cd frontend
npm test
```

### Test Full Flow

1. Create player "TestKid"
2. Answer 5 questions incorrectly on times-7
3. Check mastery heatmap in admin
4. Verify times-7 is marked WEAK (red)

## Getting Help

- **Architecture questions**: See `ARCHITECTURE.md`
- **Implementation questions**: See `TODO.md`
- **Setup issues**: See `README.md`
- **API reference**: See `backend/internal/api/routes.go`

---

**Happy coding! 🎉**
