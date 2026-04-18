#!/bin/bash
echo "Building folder_star..."
go build -o folder_star main.go
if [ $? -eq 0 ]; then
    chmod +x folder_star
    echo "Build successful: folder_star"
else
    echo "Build failed!"
    exit 1
fi
