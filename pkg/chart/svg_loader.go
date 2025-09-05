package chart

import (
	"embed"
	"path"
	"strings"
)

//go:embed svg_paths/*.svg
var svgFiles embed.FS

// SVGPaths contains all loaded SVG symbol paths
var SVGPaths map[string]string

// init loads all SVG files into memory
func init() {
	SVGPaths = make(map[string]string)
	loadSVGPaths()
}

// loadSVGPaths reads all SVG files from the embedded filesystem
func loadSVGPaths() {
	files, err := svgFiles.ReadDir("svg_paths")
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".svg") {
			content, err := svgFiles.ReadFile(path.Join("svg_paths", file.Name()))
			if err != nil {
				continue
			}

			// Remove .svg extension to get the symbol name
			name := strings.TrimSuffix(file.Name(), ".svg")
			SVGPaths[name] = string(content)
		}
	}
}

// GetSVGPath returns the SVG path content for a given symbol name
func GetSVGPath(name string) string {
	if content, exists := SVGPaths[name]; exists {
		return content
	}
	return ""
}
