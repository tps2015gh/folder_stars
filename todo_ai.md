# TODO for Next AI

This document outlines potential next steps and improvements for the "Obsidian Star Graph" project.

## Current State:
The project is a Go application that visualizes a directory structure as an interactive force-directed graph in the browser.
It includes:
- Recursive scanning of multiple file types (`.md`, `.php`, `.js`, `.css`, etc.).
- Parsing of `[[wikilinks]]` in `.md` and `.txt` files.
- Automatic "Parent-Child" linking for a structured graph view.
- `.ignore` file support for excluding directories/files.
- Interactive web UI with:
    - Theme toggle (Dark/Light mode).
    - Search functionality with Obsidian-style highlighting (dimming non-matches, highlighting matches).
    - Node context menu (right-click) for:
        - "Show in Explorer" (OS file explorer).
        - "Open Terminal" (PowerShell/Terminal).
        - "Copy Path" (relative and full).
        - "Focus/Unfocus" (isolate node and its neighbors).

## Next Steps / Potential Enhancements:

1.  **Advanced Link Parsing**:
    - Implement parsing for language-specific import/include statements (e.g., `import` in Go/JS, `require`/`use` in PHP) to establish more meaningful links within codebases.
    - Allow users to define custom regex patterns for links.

2.  **Graph Customization via UI**:
    - Add UI controls for filtering nodes by file type, folder, or specific properties.
    - Implement options to toggle different link types (e.g., show/hide parent-child links, show/hide wikilinks).

3.  **Persistent Layout**:
    - Implement saving and loading of graph node positions to maintain a preferred layout across sessions. This could involve storing data in local storage or a small config file.

4.  **Node Information Panel**:
    - When a node is clicked, display a side panel with more details about the file (e.g., file size, last modified date, first few lines of content, list of links).

5.  **Performance Optimization**:
    - For very large projects, consider optimizations for graph rendering and data processing (e.g., lazy loading nodes, WebWorkers for heavy computation).

6.  **Bidirectional Link Representation**:
    - Visually distinguish between incoming and outgoing links for better graph understanding.

7.  **Customizable Themes/Styles**:
    - Expand theme options beyond dark/light mode, allowing more user customization (e.g., node colors, link thickness, font sizes).
