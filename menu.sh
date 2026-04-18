#!/bin/bash

# Obsidian Star Graph - Interactive Launcher

clear
echo -e "\033[0;36m====================================="
echo -e "   Obsidian Star Graph Launcher"
echo -e "=====================================\033[0m"

EXE="./folder_star"
if [ ! -f "$EXE" ]; then
    echo -e "\033[0;31mError: folder_star not found. Please run build.sh first.\033[0m"
    exit 1
fi

echo -e "\nChoose a directory to scan:"
echo -e " [0] Current Directory (.)"

SUBDIRS=($(find . -maxdepth 1 -type d ! -name ".*" ! -name "."))
IDX=1
for dir in "${SUBDIRS[@]}"; do
    echo " [$IDX] ${dir#./}"
    IDX=$((IDX+1))
done
echo " [C] Enter Custom Path"

read -p $'\nSelection (default 0): ' CHOICE
if [[ -z "$CHOICE" ]]; then CHOICE="0"; fi

SELECTED_DIR="."
if [[ "$CHOICE" == "C" || "$CHOICE" == "c" ]]; then
    read -p "Enter custom path: " SELECTED_DIR
elif [[ "$CHOICE" =~ ^[0-9]+$ ]]; then
    if [[ "$CHOICE" -gt 0 && "$CHOICE" -lt "$IDX" ]]; then
        SELECTED_DIR="${SUBDIRS[$((CHOICE-1))]}"
    fi
fi

read -p "Enter port (default 8080): " PORT
if [[ -z "$PORT" ]]; then PORT="8080"; fi

echo -e "\n\033[0;32mStarting server for: $SELECTED_DIR on port $PORT...\033[0m"
$EXE -dir "$SELECTED_DIR" -port "$PORT"
