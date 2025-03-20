#!/bin/bash

# Sprawdź czy są jakieś zmiany do commitowania
if ! git status | grep -q "nothing to commit"; then
    # Pokaż status zmian
    echo "Status zmian:"
    git status
    
    # Dodaj wszystkie zmiany
    git add .
    
    # Zapytaj o nazwę commita
    echo -n "Podaj nazwę commita: "
    read commit_message
    
    # Jeśli nazwa commita jest pusta, użyj domyślnej
    if [ -z "$commit_message" ]; then
        commit_message="Aktualizacja projektu"
    fi
    
    # Wykonaj commit
    git commit -m "$commit_message"
    
    # Wysyłaj zmiany na GitHub
    echo "Wysyłam zmiany na GitHub..."
    git push origin main
    
    if [ $? -eq 0 ]; then
        echo "Zmiany zostały pomyślnie wysłane na GitHub!"
    else
        echo "Wystąpił błąd podczas wysyłania zmian na GitHub."
    fi
else
    echo "Brak zmian do commitowania."
fi 