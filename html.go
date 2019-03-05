package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	tmpl1 = template.New("html")
	tmpl2 = template.Must(tmpl1.New("sparkbar").Parse(
		`<div class="sparkbar"><div class="fill" style="width:{{.P}}px"></div><div class="empty" style="width:{{.Q}}px"></div></div>`,
	))
	tmplHead = template.Must(tmpl1.New("head").Parse(
		`<head>
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
.source { font-family: monospace; width:100%; margin:3em 0; }
.source th { padding: .1em .5em; text-align:left; border-bottom: 1px solid black; }
.source td { padding: .1em .5em; white-space: pre; }
.source .hit { background:lightblue; }
.source .miss { background:LightCoral; }
.source .ln { background:PaleGoldenrod; text-align:right; }
.source .ld { background:#f6f3d4; text-align:right; }
</style>
</head>
`,
	))
	tmplH1 = template.Must(tmpl1.New("h1").Parse(
		`<div class="pure-g"><h1 class="pure-u">{{.Title}}</h1></div>`,
	))
	tmplCoverage = template.Must(tmpl1.New("coverage").Parse(
		`<table class="pure-table pure-table-bordered" style="margin-left:auto;margin-right:0">
	<thead><tr><td></td><th>Hit</th><th>Total</th><th>Coverage</th><tr></thead>
	<tbody>
	<tr><td>Lines:</td><td>{{.LCoverage.Hits}}</td><td>{{.LCoverage.Total}}</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td></tr>
	<tr><td>Functions:</td><td>{{.FCoverage.Hits}}</td><td>{{.FCoverage.Total}}</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td></tr>
	</tbody>
	</table>`,
	))
	tmpl = template.Must(tmpl1.Parse(
		`<html>
{{template "head" .}}
<body>
{{template "h1" .}}
<div class="pure-g"><div class="pure-u">
<h2>Overall</h2>
</div></div>
<div class="pure-g"><div class="pure-u-1 pure-u-md-1-2">
<p>Coverage generated on: {{.Date}}</p>
</div><div class="pure-u-1 pure-u-md-1-2">
{{template "coverage" .}}
</div></div>
<div class="pure-g"><div class="pure-u">
<h2>By File</h2>
<table class="pure-table pure-table-bordered">
<thead><tr><th>Filename</th><th colspan="3">Line Coverage</th><th colspan="2">Function Coverage</th></tr></thead>
<tbody>
{{range $ndx, $data := .Files -}}
<tr><td><a href="{{.Name}}.html">{{.Name}}</a></td><td>{{template "sparkbar" .LCoverage}}</td><td>{{.LCoverage.Hits}}/{{.LCoverage.Total}}</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td><td>{{.FCoverage.Hits}}/{{.FCoverage.Total}}</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td></tr>
{{end}}
</tbody>
</table>
</div></div>
</body>
</html>`,
	))
)

type FileStatistics struct {
	Name      string
	LCoverage Coverage
	FCoverage Coverage
}

func buildHtml(outdir string, data map[string]FileData) error {
	err := os.MkdirAll(outdir, 0700)
	if err != nil {
		return err
	}

	filename := filepath.Join(outdir, "index.html")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	LCov := Coverage{}
	FCov := Coverage{}
	files := []FileStatistics{}
	for name, data := range data {
		if *external || !strings.HasPrefix(name, "/") {
			stats := FileStatistics{Name: name}

			stats.LCoverage = data.LineCoverage()
			LCov.Update(stats.LCoverage)
			stats.FCoverage = data.FuncCoverage()
			FCov.Update(stats.FCoverage)

			files = append(files, stats)

			buildHtmlForSource(outdir, name, data)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	params := map[string]interface{}{
		"Title":     *title,
		"LCoverage": LCov,
		"FCoverage": FCov,
		"Files":     files,
		"Date":      time.Now().UTC().Format(time.UnixDate),
	}

	return tmpl.Execute(file, params)
}

func buildHtmlForSource(outdir, sourcename string, data FileData) error {
	filename := filepath.Join(outdir, sourcename+".html")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lcov := data.LineCoverage()
	fcov := data.FuncCoverage()

	params := map[string]interface{}{
		"Title":     *title + " > " + sourcename,
		"LCoverage": lcov,
		"FCoverage": fcov,
	}

	file.WriteString("<html>\n")
	tmplHead.Execute(file, params)
	file.WriteString("<body>\n")
	tmplH1.Execute(file, params)
	file.WriteString(`<div class="pure-g"><div class="pure-u">`)
	file.WriteString(`<h2>Overall</h2>`)
	tmplCoverage.Execute(file, params)
	file.WriteString(`</div></div>`)
	file.WriteString(`<div class="pure-g"><div class="pure-u">`)
	file.WriteString(`<table class="source"><thead>`)
	file.WriteString(`<tr><th class="ln">Line #</th><th class="ld">Hit count</th><th>Source code</th></tr>`)
	file.WriteString(`</thead><tbody>`)
	buildSourceListing(file, sourcename, data.LineData)
	file.WriteString(`</tbody></table>`)
	file.WriteString(`</div></div>`)
	file.WriteString("</body>\n</html>\n")
	return nil
}

func buildSourceListing(out *os.File, sourcename string, lineCountData map[int]uint64) error {
	filename := filepath.Join(*srcdir, sourcename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	lineNo := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if lc, ok := lineCountData[lineNo]; ok {
			cl := "miss"
			if lc > 0 {
				cl = "hit"
			}
			fmt.Fprintf(out, "<tr class=\"%s\"><td class=\"ln\">%d</td><td class=\"ld\">%d</td><td>%s</td></tr>\n",
				cl, lineNo, lc,
				template.HTMLEscapeString(scanner.Text()),
			)
		} else {
			fmt.Fprintf(out, "<tr><td class=\"ln\">%d</td><td class=\"ld\"></td><td>%s</td></tr>\n",
				lineNo, template.HTMLEscapeString(scanner.Text()),
			)
		}
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
