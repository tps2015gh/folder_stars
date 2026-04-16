package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Node represents a file or folder in the graph
type Node struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Group int    `json:"group"` // 1: File, 2: Folder
}

// Link represents a connection (wikilink or parent-child)
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Graph represents the full graph data
type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

var (
	wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]+)?\]\]`)
	targetDir     string
	extensions    = []string{".md", ".php", ".js", ".css", ".html", ".json", ".txt"}
	ignorePatterns []string
)

func loadIgnorePatterns(root string) {
	ignoreFile := filepath.Join(root, ".ignore")
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		// Try root of project if not in target
		ignoreFile = ".ignore"
	}

	file, err := os.Open(ignoreFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			ignorePatterns = append(ignorePatterns, line)
		}
	}
}

func shouldIgnore(path string) bool {
	for _, pattern := range ignorePatterns {
		if strings.Contains(path, string(filepath.Separator)+pattern+string(filepath.Separator)) || 
		   strings.HasSuffix(path, string(filepath.Separator)+pattern) ||
		   strings.HasPrefix(path, pattern+string(filepath.Separator)) ||
		   filepath.Base(path) == pattern {
			return true
		}
	}
	return false
}

func main() {
	flag.StringVar(&targetDir, "dir", ".", "Directory to scan")
	port := flag.String("port", "8080", "Port to serve on")
	flag.Parse()

	absPath, err := filepath.Abs(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	loadIgnorePatterns(absPath)

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		graph := scanDirectory(absPath)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(graph)
	})

	http.HandleFunc("/explore", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			return
		}
		fullPath := filepath.Join(absPath, filepath.FromSlash(path))

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			// /select, highlights the file. If it's a directory, it opens it.
			cmd = exec.Command("explorer", "/select,", fullPath)
		} else if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-R", fullPath)
		} else {
			cmd = exec.Command("xdg-open", filepath.Dir(fullPath))
		}
		cmd.Run()
	})

	http.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			return
		}
		fullPath := filepath.Join(absPath, filepath.FromSlash(path))

		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			fullPath = filepath.Dir(fullPath)
		}

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "start", "powershell.exe", "-NoExit", "-Command", fmt.Sprintf("cd '%s'", fullPath))
		} else if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "Terminal", fullPath)
		} else {
			cmd = exec.Command("x-terminal-emulator", "--working-directory", fullPath)
		}
		cmd.Run()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(indexHTML))
		tmpl.Execute(w, struct{ AbsPath string }{AbsPath: filepath.ToSlash(absPath)})
	})

	fmt.Printf("Scanning directory: %s\n", absPath)
	fmt.Printf("Server started at http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func scanDirectory(root string) Graph {
	nodes := make([]Node, 0)
	links := make([]Link, 0)
	nodeMap := make(map[string]bool)

	fmt.Println("--- Starting Scan ---")

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			return nil
		}
		
		rel = filepath.ToSlash(rel)
		if shouldIgnore(rel) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Add node for file or folder
		isDir := info.IsDir()
		ext := filepath.Ext(path)
		
		isValidExt := false
		for _, e := range extensions {
			if e == ext {
				isValidExt = true
				break
			}
		}

		if isDir || isValidExt {
			group := 2 // Folder
			if !isDir {
				group = 1 // File
			}
			
			nodes = append(nodes, Node{ID: rel, Name: filepath.Base(rel), Group: group})
			nodeMap[rel] = true

			// Create parent link to form the "Star" structure
			parent := filepath.ToSlash(filepath.Dir(rel))
			if parent != "." {
				links = append(links, Link{Source: parent, Target: rel})
			}
		}

		return nil
	})

	fmt.Printf("--- Found %d nodes. Analyzing wikilinks... ---\n", len(nodes))

	// Analysis pass for wikilinks (only in text files)
	for id := range nodeMap {
		if !strings.HasSuffix(id, ".md") && !strings.HasSuffix(id, ".txt") {
			continue
		}

		path := filepath.Join(root, id)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		matches := wikilinkRegex.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			target := filepath.ToSlash(match[1])
			
			// Try exact match or base name match
			if _, exists := nodeMap[target]; exists {
				links = append(links, Link{Source: id, Target: target})
			} else {
				for existingId := range nodeMap {
					if filepath.Base(existingId) == target || strings.TrimSuffix(filepath.Base(existingId), filepath.Ext(existingId)) == target {
						links = append(links, Link{Source: id, Target: existingId})
						break
					}
				}
			}
		}
	}

	fmt.Println("--- Scan Complete ---")
	return Graph{Nodes: nodes, Links: links}
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Project Star Graph</title>
    <script src="//unpkg.com/force-graph"></script>
    <style>
        :root {
            --bg-color: #1a1a1a;
            --text-color: #ccc;
            --control-bg: rgba(0,0,0,0.8);
            --link-color: rgba(255, 255, 255, 0.3);
            --menu-bg: #2a2a2a;
            --menu-hover: #3a3a3a;
        }
        body.light-mode {
            --bg-color: #f5f5f5;
            --text-color: #333;
            --control-bg: rgba(255,255,255,0.9);
            --link-color: rgba(0, 0, 0, 0.2);
            --menu-bg: #fff;
            --menu-hover: #eee;
        }
        body { margin: 0; background-color: var(--bg-color); color: var(--text-color); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; transition: background-color 0.3s; overflow: hidden; }
        #graph { width: 100vw; height: 100vh; }
        .controls { 
            position: absolute; top: 10px; left: 10px; z-index: 10; 
            background: var(--control-bg); padding: 15px; border-radius: 8px; 
            border: 1px solid #444; pointer-events: auto;
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
        input[type="text"] {
            width: 100%; padding: 8px; margin-bottom: 10px; border-radius: 4px;
            border: 1px solid #555; background: #333; color: white; box-sizing: border-box;
        }
        body.light-mode input[type="text"] { background: white; color: black; border: 1px solid #ccc; }
        button { 
            background: #444; color: white; border: none; padding: 5px 10px; 
            border-radius: 4px; cursor: pointer; margin-top: 5px;
        }
        body.light-mode button { background: #ddd; color: #333; }
        
        #context-menu {
            position: fixed; display: none; z-index: 100;
            background: var(--menu-bg); border: 1px solid #444;
            border-radius: 4px; padding: 5px 0; min-width: 160px;
            box-shadow: 0 4px 10px rgba(0,0,0,0.5);
            color: var(--text-color);
        }
        #context-menu div {
            padding: 8px 15px; cursor: pointer; font-size: 14px;
        }
        #context-menu div:hover { background: var(--menu-hover); }
    </style>
</head>
<body>
    <div class="controls">
        <h3 style="margin:0 0 10px 0; font-size: 1.1em;">Project Graph</h3>
        <input type="text" id="search" placeholder="Search files..." oninput="handleSearch(this.value)">
        <p id="status" style="margin: 5px 0; font-size: 0.9em;">Ready</p>
        <div style="font-size: 0.85em; margin-bottom: 10px;">
            <span style="color: #54a0ff">●</span> Folders &nbsp;
            <span style="color: #1dd1a1">●</span> Files
        </div>
        <button onclick="toggleTheme()">Toggle Theme</button>
        <button onclick="resetFocus()" id="btn-reset" style="display:none">Show All</button>
    </div>

    <div id="context-menu"></div>
    <div id="graph"></div>

    <script>
        let isLightMode = false;
        let graphInstance = null;
        let allData = { nodes: [], links: [] };
        let selectedNode = null;
        let searchQuery = "";

        function toggleTheme() {
            isLightMode = !isLightMode;
            document.body.classList.toggle('light-mode');
            if (graphInstance) {
                const linkColor = isLightMode ? 'rgba(0,0,0,0.2)' : 'rgba(255,255,255,0.3)';
                graphInstance.linkColor(() => linkColor);
            }
        }

        function handleSearch(val) {
            searchQuery = val.toLowerCase();
            if (graphInstance) graphInstance.nodeCanvasObject(graphInstance.nodeCanvasObject()); // Refresh
        }

        function resetFocus() {
            selectedNode = null;
            document.getElementById('btn-reset').style.display = 'none';
            graphInstance.graphData(allData);
        }

        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(() => {
                const status = document.getElementById('status');
                const oldText = status.innerText;
                status.innerText = 'Copied: ' + text.split('/').pop();
                const oldColor = status.style.color;
                status.style.color = '#1dd1a1';
                setTimeout(() => {
                    status.innerText = oldText;
                    status.style.color = oldColor;
                }, 2000);
            }).catch(err => {
                console.error('Failed to copy: ', err);
            });
        }

        function showContextMenu(e, node) {
            e.preventDefault();
            const menu = document.getElementById('context-menu');
            menu.style.display = 'block';
            menu.style.left = e.clientX + 'px';
            menu.style.top = e.clientY + 'px';
            
            menu.innerHTML = "";
            const actions = [
                { label: 'Show in Explorer', action: () => fetch('/explore?path=' + encodeURIComponent(node.id)) },
                { label: 'Open Terminal', action: () => fetch('/terminal?path=' + encodeURIComponent(node.id)) },
                { label: 'Copy Relative Path', action: () => copyToClipboard(node.id) },
                { label: 'Copy Full Path', action: () => copyToClipboard('{{.AbsPath}}/' + node.id) },
                { label: 'Focus Node', action: () => focusNode(node) }
            ];

            actions.forEach(a => {
                const div = document.createElement('div');
                div.innerText = a.label;
                div.onclick = () => { a.action(); menu.style.display = 'none'; };
                menu.appendChild(div);
            });

            document.addEventListener('click', () => menu.style.display = 'none', { once: true });
        }

        function focusNode(node) {
            selectedNode = node;
            const neighbors = new Set();
            neighbors.add(node.id);
            allData.links.forEach(link => {
                if (link.source.id === node.id || link.source === node.id) neighbors.add(link.target.id || link.target);
                if (link.target.id === node.id || link.target === node.id) neighbors.add(link.source.id || link.source);
            });

            const filteredNodes = allData.nodes.filter(n => neighbors.has(n.id));
            const filteredLinks = allData.links.filter(l => {
                const s = l.source.id || l.source;
                const t = l.target.id || l.target;
                return neighbors.has(s) && neighbors.has(t);
            });
            
            graphInstance.graphData({ nodes: filteredNodes, links: filteredLinks });
            document.getElementById('btn-reset').style.display = 'inline-block';
        }

        fetch('/data')
            .then(res => res.json())
            .then(data => {
                allData = data;
                document.getElementById('status').innerText = 'Loaded ' + data.nodes.length + ' nodes';
                graphInstance = ForceGraph()
                (document.getElementById('graph'))
                    .graphData(data)
                    .nodeId('id')
                    .nodeLabel(node => node.id + ' (Click to copy path)')
                    .onNodeClick(node => copyToClipboard('{{.AbsPath}}/' + node.id))
                    .onNodeRightClick(showContextMenu)
                    .linkColor(() => isLightMode ? 'rgba(0,0,0,0.2)' : 'rgba(255,255,255,0.3)')
                    .linkWidth(1.5)
                    .linkDirectionalParticles(2)
                    .linkDirectionalParticleWidth(2)
                    .linkDirectionalParticleSpeed(0.005)
                    .nodeCanvasObject((node, ctx, globalScale) => {
                        const label = node.name;
                        const fontSize = 12/globalScale;
                        ctx.font = fontSize + 'px Sans-Serif';
                        
                        const isMatch = searchQuery && label.toLowerCase().includes(searchQuery);
                        const isDimmed = searchQuery && !isMatch;

                        // Draw circle
                        const radius = node.group === 2 ? 6 : 4;
                        ctx.fillStyle = node.group === 2 ? '#54a0ff' : '#1dd1a1';
                        
                        if (isDimmed) ctx.globalAlpha = 0.1;
                        if (isMatch) {
                            ctx.shadowColor = '#fff';
                            ctx.shadowBlur = 15;
                        }

                        ctx.beginPath();
                        ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
                        ctx.fill();
                        
                        ctx.shadowBlur = 0; // reset
                        ctx.globalAlpha = 1.0;

                        // Draw text on zoom or if it's a match
                        if (globalScale > 1.5 || isMatch) {
                            ctx.textAlign = 'center';
                            ctx.textBaseline = 'middle';
                            ctx.fillStyle = isLightMode ? '#333' : '#eee';
                            if (isDimmed) ctx.globalAlpha = 0.2;
                            ctx.fillText(label, node.x, node.y + radius + 5);
                            ctx.globalAlpha = 1.0;
                        }
                    });
            });
    </script>
</body>
</html>
`
