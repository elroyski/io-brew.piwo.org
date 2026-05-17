#!/bin/bash
# Watchdog / restart skrypt dla io-brew.piwo.org
# Uruchamiaj przez cron co 5 minut:
# */5 * * * * /bin/bash /home/piwo/domains/io-brew.piwo.org/restart.sh >> /home/piwo/domains/io-brew.piwo.org/watchdog.log 2>&1

APP_DIR="/home/piwo/domains/io-brew.piwo.org/ispindel-app"
BINARY="${APP_DIR}/ispindel-app"
PID_FILE="${APP_DIR}/app.pid"
LOG_FILE="${APP_DIR}/app.log"
PORT=5414
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

log() {
    echo "[${TIMESTAMP}] $1"
}

# Sprawdź czy aplikacja odpowiada na HTTP
app_responds() {
    curl -sf --max-time 5 "http://127.0.0.1:${PORT}/" > /dev/null 2>&1
    return $?
}

# Zabij stary proces po PID
kill_old_process() {
    if [ -f "$PID_FILE" ]; then
        OLD_PID=$(cat "$PID_FILE")
        if [ -n "$OLD_PID" ] && kill -0 "$OLD_PID" 2>/dev/null; then
            log "Zabijam stary proces PID=${OLD_PID}"
            kill "$OLD_PID" 2>/dev/null
            sleep 2
            # Wymuś jeśli nadal żyje
            if kill -0 "$OLD_PID" 2>/dev/null; then
                kill -9 "$OLD_PID" 2>/dev/null
                log "Wymuszono kill -9 PID=${OLD_PID}"
            fi
        fi
        rm -f "$PID_FILE"
    fi

    # Dodatkowe: zabij wszystkie procesy o nazwie binarki
    pkill -f "${BINARY}" 2>/dev/null
    sleep 1
}

# Uruchom aplikację
start_app() {
    cd "$APP_DIR" || { log "BŁĄD: nie mogę wejść do $APP_DIR"; exit 1; }

    if [ ! -f "$BINARY" ]; then
        log "BŁĄD: brak pliku binarnego $BINARY"
        exit 1
    fi

    if [ ! -x "$BINARY" ]; then
        chmod +x "$BINARY"
    fi

    nohup "${BINARY}" >> "${LOG_FILE}" 2>&1 &
    NEW_PID=$!
    echo "$NEW_PID" > "$PID_FILE"
    log "Uruchomiono aplikację PID=${NEW_PID}"

    # Poczekaj chwilę i sprawdź czy działa
    sleep 5
    if app_responds; then
        log "OK: aplikacja odpowiada na porcie ${PORT}"
    elif kill -0 "$NEW_PID" 2>/dev/null; then
        log "UWAGA: proces żyje (PID=${NEW_PID}) ale HTTP jeszcze nie odpowiada - może się inicjalizuje"
    else
        log "BŁĄD: aplikacja nie wystartowała! Sprawdź ${LOG_FILE}"
    fi
}

# --- Główna logika ---

if app_responds; then
    log "OK: aplikacja działa na porcie ${PORT} - nic nie robię"
    exit 0
fi

log "UWAGA: aplikacja nie odpowiada - restartuję..."
kill_old_process
start_app
