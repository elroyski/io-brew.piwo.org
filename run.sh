#!/bin/bash

# Sprawdź czy Go jest zainstalowany
if ! command -v go &> /dev/null; then
    echo "Go nie jest zainstalowany. Instaluję..."
    sudo apt update
    sudo apt install golang-go
fi

# Sprawdź czy PostgreSQL jest zainstalowany
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL nie jest zainstalowany. Instaluję..."
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
fi

# Sprawdź czy baza danych i użytkownik istnieją
if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='ispindel'" | grep -q 1; then
    echo "Tworzę użytkownika bazy danych..."
    sudo -u postgres psql -c "CREATE USER ispindel WITH PASSWORD 'Kochanapysia1';"
fi

if ! sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw ispindel; then
    echo "Tworzę bazę danych..."
    sudo -u postgres psql -c "CREATE DATABASE ispindel;"
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ispindel TO ispindel;"
    sudo -u postgres psql -d ispindel -c "GRANT ALL ON SCHEMA public TO ispindel;"
fi

# Pobierz zależności
echo "Pobieram zależności..."
go mod download

# Uruchom aplikację
echo "Uruchamiam aplikację..."
go run cmd/server/main.go 