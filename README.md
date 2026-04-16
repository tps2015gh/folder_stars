# Obsidian Star Graph (Go)

This project is a Go-based application designed to bring the powerful visualization capabilities of Obsidian's graph view to any directory structure. It transforms your files and folders into an interactive, force-directed graph displayed in your web browser. Think of it as a dynamic map of your digital workspace, revealing connections and hierarchies in a way that plain file explorers cannot.

## How it Works

The application operates on a simple yet effective client-server model:

1.  **Go Backend**: When you run the Go program, it first scans a specified directory (or the current directory by default).
    *   It identifies and indexes all supported files (like `.md`, `.php`, `.js`, `.html`, etc.) and folders, creating nodes for each.
    *   It automatically establishes "parent-child" links between files/folders and their containing directories, forming a fundamental tree structure.
    *   For Markdown (`.md`) and text (`.txt`) files, it parses wikilinks (`[[link]]`) to identify explicit connections between notes.
    *   It respects an `.ignore` file to exclude common noisy directories (like `vendor`, `node_modules`) for a cleaner graph.
    *   It then starts a local web server, ready to serve the graph data and the interactive frontend.

2.  **Web Frontend**: Your browser connects to the Go server.
    *   It fetches the processed graph data (nodes and links) as JSON.
    *   Using the `force-graph` JavaScript library, it renders an interactive graph.
    - You can zoom, pan, and drag nodes.
    - Search functionality allows you to quickly find and highlight files.
    - **One-Click Copy**: Left-click any node to instantly copy its full absolute path to your clipboard.
    - A right-click context menu provides advanced actions like opening files in Explorer/Terminal and focus mode.
    - A theme toggle allows switching between dark and light interfaces for comfortable viewing.

## Purpose

The goal is to provide a clear, visual overview of your project's or notes' structure. Whether you're navigating a code repository, exploring a personal knowledge base, or understanding the interconnections of your documents, this tool helps you see the bigger picture and discover relationships you might otherwise miss. It aims to be a lightweight, fast, and user-friendly tool for gaining spatial awareness of your digital assets.

## Features

- **Recursive Directory Scanning**: Scans a specified folder and all its subfolders for supported files.
- **Comprehensive File Support**: Indexes `.md`, `.php`, `.js`, `.css`, `.html`, `.json`, `.txt`, and more.
- **Wikilink Parsing**: Automatically detects and visualizes links in `[[Link]]` format within Markdown and text files.
- **Structural Linking**: Creates implicit links between files/folders and their parent directories, forming a hierarchical "star" graph.
- **Interactive Web UI**: Provides a force-directed graph with:
  - Drag-and-drop nodes for manual arrangement.
  - Search functionality with Obsidian-style highlighting: matches are highlighted with a glow effect, while non-matches are dimmed for clarity.
  - **One-Click Copy**: Left-click any node to instantly copy its full absolute path to the clipboard.
  - A dynamic context menu (right-click) offering:
    - **Show in Explorer**: Opens the OS file explorer and selects the item.
    - **Open Terminal**: Opens PowerShell (Windows) or Terminal (macOS/Linux) in the node's directory.
    - **Copy Path**: Copy Relative or Full (absolute) paths to clipboard with visual confirmation.
    - **Focus/Unfocus**: Isolate a node and its immediate connections for detailed inspection.
  - Theme toggle (Dark/Light mode).
- **`.ignore` File Support**: Excludes specified files/directories (e.g., `vendor`, `node_modules`, `.git`, `*.exe`, `*.log`) from the graph for a cleaner visualization.

## Privacy & Technical Design

This application is designed to be **local-only** and **privacy-first**:

-   **No AI/Embeddings**: The graph is built using traditional file system traversal and regex-based parsing. It does **not** send your data to any external AI services or use any embedding APIs.
-   **Local Processing**: All directory scanning and wikilink analysis happen entirely on your machine.
-   **Static Analysis**: Relationships are established purely through file hierarchy and explicit `[[wikilinks]]` found in your documents.
-   **Browser-Based Rendering**: The visualization is rendered in your browser using local graph data, ensuring your digital workspace remains private and secure.

## Prerequisites

- [Go](https://golang.org/dl/) (1.16 or later recommended)
- A directory containing files to visualize.

## Usage

1.  **Clone/Download** this repository.
2.  **Navigate** to the project directory in your terminal.
3.  **Build** the executable:
    ```bash
    go build -o file_stars.exe main.go
    ```
4.  **Run** the application, specifying the target directory:
    ```bash
    ./file_stars.exe -dir "/path/to/your/directory"
    ```
    (If no directory is specified, it defaults to the current directory.)
5.  **Open** your web browser to the address provided in the terminal (e.g., `http://localhost:8080`).

## Configuration

| Flag    | Description                         | Default |
|---------|-------------------------------------|---------|
| `-dir`  | The directory to scan for files     | `.`     |
| `-port` | The port to serve the web interface | `8080`  |

---
Built with Go, `force-graph`, and inspired by Obsidian's graph view.
