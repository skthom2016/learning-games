@echo off
REM Database Reset Script - Keep only David and Daniel
REM This script connects to the running PostgreSQL database and resets it

echo ========================================
echo Learning Game Database Reset
echo ========================================
echo.
echo This will:
echo 1. Keep David and Daniel as the only players
echo 2. Delete all other players and their data
echo 3. Reset all scores and progress
echo.
echo Press Ctrl+C to cancel, or
pause

echo.
echo Connecting to database...
docker exec -i learning-game-db psql -U gameuser -d learning_game < backend\migrations\007_reset_and_seed_david_daniel.sql

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo SUCCESS! Database has been reset.
    echo ========================================
    echo Only David and Daniel remain in the database.
    echo All other data has been cleared.
) else (
    echo.
    echo ========================================
    echo ERROR! Database reset failed.
    echo ========================================
    echo Make sure Docker is running and the database container is up.
    echo Run: docker-compose up -d
)

echo.
pause
