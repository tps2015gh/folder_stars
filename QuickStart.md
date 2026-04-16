# Quick Start Guide

Get your Obsidian-style graph up and running in less than a minute.

## 1. Prepare your environment
Ensure you have Go installed on your system. You can check by running:
```powershell
go version
```

## 2. Initialize and Build
If you haven't already, navigate to the project directory and initialize the module:
```powershell
go mod tidy
go build -o obsidian-graph.exe main.go
```

## 3. Run the Scan
Run the application against a folder. For example, if you want to scan a folder named `MyNotes`:
```powershell
./obsidian-graph.exe -dir "./MyNotes"
```

## 4. View the Graph
1. Once the terminal says `Server started at http://localhost:8080`, open your web browser.
2. Go to: **[http://localhost:8080](http://localhost:8080)**
3. You will see your notes as nodes. Hover over them to see filenames, and click/drag to interact!

## 5. Keyboard/Mouse Controls
- **Scroll**: Zoom in and out.
- **Left Click + Drag**: Rotate/Move the view.
- **Click Node**: Highlight or interact (expandable features).
- **Drag Node**: Pin the node to a specific location.

---
**Tip**: Use the `-port` flag if 8080 is already in use:
```powershell
./obsidian-graph.exe -dir "." -port 9000
```
