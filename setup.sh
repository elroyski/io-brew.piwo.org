#!/bin/bash

# Konfiguracja zmiennych
BASE_DIR="/home/piwo/domains/io-brew.piwo.org"
APP_DIR="${BASE_DIR}/io-brew"
LOG_FILE="${BASE_DIR}/setup.log"

echo "Rozpoczynam instalację w $APP_DIR" | tee -a $LOG_FILE

# Sprawdź, czy aktualny katalog to $BASE_DIR
if [ "$(pwd)" != "$BASE_DIR" ]; then
    echo "UWAGA: Uruchamiasz skrypt poza katalogiem $BASE_DIR"
    echo "Przechodzę do katalogu $BASE_DIR..."
    cd $BASE_DIR || { echo "Nie mogę przejść do katalogu $BASE_DIR. Kończę."; exit 1; }
fi

# Tworzenie backupu pliku .env jeśli istnieje
if [ -f "${APP_DIR}/.env" ]; then
    echo "Tworzę backup istniejącego pliku .env" | tee -a $LOG_FILE
    cp "${APP_DIR}/.env" "${BASE_DIR}/.env.backup" || { echo "Nie mogę utworzyć backupu .env. Kończę."; exit 1; }
    echo "Backup pliku .env utworzony w ${BASE_DIR}/.env.backup" | tee -a $LOG_FILE
fi

# Usunięcie starego katalogu aplikacji (jeśli istnieje)
if [ -d "$APP_DIR" ]; then
    echo "Usuwanie istniejącego katalogu aplikacji..." | tee -a $LOG_FILE
    rm -rf "$APP_DIR" || { echo "Nie mogę usunąć katalogu $APP_DIR. Kończę."; exit 1; }
fi

# Klonowanie repozytorium
echo "Klonowanie repozytorium do $APP_DIR..." | tee -a $LOG_FILE
git clone https://github.com/elroyski/io-brew.piwo.org.git "$APP_DIR" || { echo "Klonowanie repozytorium nie powiodło się. Kończę."; exit 1; }

# Przejście do katalogu aplikacji
cd "$APP_DIR" || { echo "Nie mogę przejść do katalogu $APP_DIR. Kończę."; exit 1; }

# Pobieranie zależności
echo "Pobieranie zależności Go..." | tee -a $LOG_FILE
go mod download || { echo "Nie mogę pobrać zależności. Kończę."; exit 1; }
echo "Zależności pobrane pomyślnie." | tee -a $LOG_FILE

# Kompilacja aplikacji
echo "Kompilacja aplikacji..." | tee -a $LOG_FILE
go build -v -o io-brew ./cmd/server || { echo "Kompilacja nie powiodła się. Kończę."; exit 1; }
echo "Aplikacja skompilowana pomyślnie." | tee -a $LOG_FILE

# Ustawienie uprawnień
echo "Ustawianie uprawnień wykonywania..." | tee -a $LOG_FILE
chmod +x io-brew || { echo "Nie mogę ustawić uprawnień. Kończę."; exit 1; }
chmod +x restart.sh || { echo "Nie mogę ustawić uprawnień dla restart.sh. Upewnij się, że plik istnieje."; }

echo "Konfiguracja zakończona! Aplikacja jest gotowa do uruchomienia." | tee -a $LOG_FILE
echo "Aby uruchomić aplikację, wykonaj: cd $APP_DIR && ./restart.sh" | tee -a $LOG_FILE 
