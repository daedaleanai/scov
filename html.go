package main

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	tmpl1 = template.New("html")
	tmpl2 = template.Must(tmpl1.New("sparkbar").Parse(
		`<div class="sparkbar"><div class="fill" style="width:{{.P}}px"></div><div class="empty" style="width:{{.Q}}px"></div></div>`,
	))
	tmpl = template.Must(tmpl1.Parse(
		`<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
<!--[if lte IE 8]>
    <link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-old-ie-min.css">
<![endif]-->
<!--[if gt IE 8]><!-->
    <link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-min.css">
<!--<![endif]-->
<style>
body { max-width:70em; margin:auto; }
.sparkbar { border: 1px solid black; border-radius:1px; }
.sparkbar .fill { display: inline-block; height: 1em; background-color:lightgreen; }
.sparkbar .empty { display: inline-block; height: 1em; background-color: white; }
</style>
</head>
<body>
<div class="pure-g">
<h1 class="pure-u">{{.Title}}</h1>
</div>
<div class="pure-g"><div class="pure-u">
<h2>Overall</h2>
<table class="pure-table pure-table-bordered">
<thead><tr><td></td><th>Coverage</th><th>Hit</th><th>Total</th><tr></thead>
<tbody>
<tr><td>Lines:</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td><td>{{.LCoverage.Hit}}</td><td>{{.LCoverage.Total}}</td></tr>
<tr><td>Functions:</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td><td>{{.FCoverage.Hit}}</td><td>{{.FCoverage.Total}}</td></tr>
</tbody>
</table>
</div>
<div class="pure-g"><div class="pure-u">
<h2>By File</h2>
<table class="pure-table pure-table-bordered">
<thead><tr><th>Filename</th><th colspan="3">Line Coverage</th><th colspan="2">Function Coverage</th></tr></thead>
<tbody>
{{range $ndx, $data := .Files -}}
<tr><td>{{.Name}}</td><td>{{template "sparkbar" .LCoverage}}</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td><td>{{.LCoverage.Hit}}/{{.LCoverage.Total}}</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td><td>{{.FCoverage.Hit}}/{{.FCoverage.Total}}</td></tr>
{{end}}
</tbody>
</table>
</div></div>
</body>
</html>`,
	))
)

type Coverage struct {
	Hit   int
	Total int
}

func (c Coverage) Percentage() float32 {
	return float32(c.Hit) * 100 / float32(c.Total)
}

func (c Coverage) P() float32 {
	return float32(c.Hit) * 100 / float32(c.Total)
}

func (c Coverage) Q() float32 {
	return 100 - float32(c.Hit)*100/float32(c.Total)
}

func (c *Coverage) Update(a, b int) {
	c.Hit += a
	c.Total += b
}

type FileStatistics struct {
	Name      string
	LCoverage Coverage
	FCoverage Coverage
}

func buildHtml(outdir string) error {
	filename := filepath.Join(outdir, "index.html")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lineCoverage := Coverage{}
	funcCoverage := Coverage{}
	files := []FileStatistics{}
	for name, data := range lineCountData {
		if *external || !strings.HasPrefix(name, "/") {
			stats := FileStatistics{Name: name}
			tmpa, tmpb := lineCoverageForFile(data)
			stats.LCoverage = Coverage{tmpa, tmpb}
			lineCoverage.Update(tmpa, tmpb)

			tmpa, tmpb = funcCoverageForFile(funcCountData[name])
			stats.FCoverage = Coverage{tmpa, tmpb}
			funcCoverage.Update(tmpa, tmpb)

			files = append(files, stats)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	params := map[string]interface{}{
		"Title":     "LCovGo",
		"LCoverage": lineCoverage,
		"FCoverage": funcCoverage,
		"Files":     files,
	}

	return tmpl.Execute(file, params)
}
