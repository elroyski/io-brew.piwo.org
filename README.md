# iSpindel.piwo.org

Aplikacja webowa do zarządzania iSpindel - urządzeniem do monitorowania fermentacji piwa.

## Wymagania

- Go 1.23 lub nowszy
- PostgreSQL 14 lub nowszy

## Instalacja

1. Sklonuj repozytorium:
```bash
git clone https://github.com/twojego-username/ispindel.piwo.org.git
cd ispindel.piwo.org
```

2. Zainstaluj zależności:
```bash
go mod download
```

3. Skonfiguruj bazę danych PostgreSQL:
```bash
sudo -u postgres psql
CREATE DATABASE ispindel;
CREATE USER ispindel WITH PASSWORD 'twoje_haslo';
GRANT ALL PRIVILEGES ON DATABASE ispindel TO ispindel;
```

4. Uruchom aplikację:
```bash
go run cmd/server/main.go
```

## Struktura projektu

- `cmd/server/` - główny plik uruchomieniowy serwera
- `internal/` - wewnętrzny kod aplikacji
  - `handlers/` - obsługa żądań HTTP
  - `models/` - modele danych
  - `services/` - logika biznesowa
- `pkg/` - współdzielone pakiety
  - `auth/` - autentykacja i autoryzacja
  - `database/` - konfiguracja bazy danych
- `web/templates/` - szablony HTML 