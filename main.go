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
	"path/filepath"
	"regexp"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index").Parse(indexHTML))
		tmpl.Execute(w, nil)
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
        body { margin: 0; background-color: #1a1a1a; color: #ccc; font-family: sans-serif; }
        #graph { width: 100vw; height: 100vh; }
        .controls { position: absolute; top: 10px; left: 10px; z-index: 10; background: rgba(0,0,0,0.5); padding: 10px; border-radius: 5px; pointer-events: none; }
    </style>
</head>
<body>
    <div class="controls">
        <h3>Project Graph View</h3>
        <p id="status">Ready</p>
        <p style="font-size: 0.8em; color: #888;">Blue: Folders | Green: Files</p>
    </div>
    <div id="graph"></div>

    <script>
        fetch('/data')
            .then(res => res.json())
            .then(data => {
                document.getElementById('status').innerText = 'Loaded ' + data.nodes.length + ' nodes';
                const Graph = ForceGraph()
                (document.getElementById('graph'))
                    .graphData(data)
                    .nodeId('id')
                    .nodeLabel(node => node.id)
                    .nodeAutoColorBy('group')
                    .linkDirectionalParticles(1)
                    .linkDirectionalParticleSpeed(0.005)
                    .nodeCanvasObject((node, ctx, globalScale) => {
                        const label = node.name;
                        const fontSize = 12/globalScale;
                        ctx.font = fontSize + 'px Sans-Serif';
                        
                        // Draw circle
                        const radius = node.group === 2 ? 4 : 2;
                        ctx.fillStyle = node.group === 2 ? '#54a0ff' : '#1dd1a1';
                        ctx.beginPath();
                        ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
                        ctx.fill();

                        // Draw text
                        if (globalScale > 1.5) {
                            ctx.textAlign = 'center';
                            ctx.textBaseline = 'middle';
                            ctx.fillStyle = '#eee';
                            ctx.fillText(label, node.x, node.y + radius + 4);
                        }
                    });
            });
    </script>
</body>
</html>
`
