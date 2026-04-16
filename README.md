# Obsidian Star Graph (Go)

A lightweight Go application that visualizes your Markdown notes as an interactive, force-directed "star graph" in the browser, mimicking the popular graph view in [Obsidian](https://obsidian.md).

## Features

- **Recursive Directory Scanning**: Scans a specified folder and all its subfolders for `.md` files.
- **Wikilink Parsing**: Automatically detects links in the format `[[Filename]]` or `[[Filename|Alias]]`.
- **Interactive Web UI**: Provides a force-directed graph with:
  - Drag-and-drop nodes.
  - Interactive labels.
  - Animated link particles.
  - Color-coded groups (e.g., separating root files from subfolder files).
- **Real-time Progress**: Displays scan status and link analysis in the terminal.

## Prerequisites

- [Go](https://golang.org/dl/) (1.16 or later recommended)
- A folder containing Markdown files with wikilinks.

## Usage

1. **Clone/Download** this repository.
2. **Run** the program:
   ```bash
   go run main.go -dir "/path/to/your/vault"
   ```
3. **Open** your browser to `http://localhost:8080`.

## Configuration

| Flag    | Description                         | Default |
|---------|-------------------------------------|---------|
| `-dir`  | The directory to scan for `.md` files | `.`     |
| `-port` | The port to serve the web interface  | `8080`  |

---
Built with Go and [Force Graph](https://github.com/vasturiano/force-graph).
