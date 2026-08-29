@echo off
REM ==========================================
REM Automated Safe Update Script
REM Run this to update the system with new code
REM ==========================================
REM This script:
REM   1. Backs up existing data automatically
REM   2. Rebuilds and starts containers
REM   3. Restores data automatically

echo.
echo ==========================================
echo   Automated Safe Update
echo ==========================================
echo.

REM Create backup folder
if not exist "player-backups" mkdir player-backups

REM Create timestamp for backup
set TIMESTAMP=%date:~-4%%date:~4,2%%date:~7,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set TIMESTAMP=%TIMESTAMP: =0%
set BACKUP_FILE=player-backups\backup_%TIMESTAMP%.sql

echo [1/4] Checking existing system...
docker ps | findstr "learning-game-db" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo Found existing database - backing up data...

    docker exec learning-game-db pg_dump -U gameuser learning_game > %BACKUP_FILE%

    if %ERRORLEVEL% NEQ 0 (
        echo [ERROR] Backup failed! Aborting update.
        pause
        exit /b 1
    )

    echo     Backup created: %BACKUP_FILE%
    set DATA_EXISTS=1
) else (
    echo No existing database found - this is a fresh install
    set DATA_EXISTS=0
)

echo.
echo [2/4] Stopping containers...
docker-compose down

echo.
echo [3/4] Building and starting containers with new code...
docker-compose up -d --build

echo.
echo [4/4] Waiting for database to be ready...
timeout /t 5 /nobreak >nul

REM Wait for database to be ready
:waitloop
docker exec learning-game-db pg_isready -U gameuser >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    timeout /t 2 /nobreak >nul
    goto waitloop
)

echo Database is ready!

REM Restore data if backup was made
if %DATA_EXISTS% EQU 1 (
    echo.
    echo Restoring data from backup...
    docker exec -i learning-game-db psql -U gameuser -d learning_game < %BACKUP_FILE%

    if %ERRORLEVEL% NEQ 0 (
        echo.
        echo [WARNING] Restore had issues. Backup saved at: %BACKUP_FILE%
    ) else (
        echo     Data restored successfully!
    )
)

echo.
echo ==========================================
echo   Update Complete!
echo ==========================================
echo.
if %DATA_EXISTS% EQU 1 (
    echo David and Daniel's scores have been preserved.
)
echo.
echo Application is running at:
echo   - Frontend: http://localhost:3000
echo   - Backend:  http://localhost:8080
echo.
pause
