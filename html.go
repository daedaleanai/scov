package main

import (
	"bufio"
	"fmt"
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
.source th { padding: .1em .5em; }
.source td { padding: .1em .5em; white-space: pre; }
.source .hit { background:lightblue; }
.source .miss { background:LightCoral; }
.source .ln { background:PaleGoldenrod; text-align:right; }
.source .ld { text-align:right; }
</style>
</head>
`,
	))
	tmplH1 = template.Must(tmpl1.New("h1").Parse(
		`<div class="pure-g"><h1 class="pure-u">{{.Title}}</h1></div>`,
	))
	tmplCoverage = template.Must(tmpl1.New("coverage").Parse(
		`<table class="pure-table pure-table-bordered">
	<thead><tr><td></td><th>Coverage</th><th>Hit</th><th>Total</th><tr></thead>
	<tbody>
	<tr><td>Lines:</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td><td>{{.LCoverage.Hit}}</td><td>{{.LCoverage.Total}}</td></tr>
	<tr><td>Functions:</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td><td>{{.FCoverage.Hit}}</td><td>{{.FCoverage.Total}}</td></tr>
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
{{template "coverage" .}}
</div>
<div class="pure-g"><div class="pure-u">
<h2>By File</h2>
<table class="pure-table pure-table-bordered">
<thead><tr><th>Filename</th><th colspan="3">Line Coverage</th><th colspan="2">Function Coverage</th></tr></thead>
<tbody>
{{range $ndx, $data := .Files -}}
<tr><td><a href="{{.Name}}.html">{{.Name}}</a></td><td>{{template "sparkbar" .LCoverage}}</td><td>{{printf "%.1f" .LCoverage.Percentage}}%</td><td>{{.LCoverage.Hit}}/{{.LCoverage.Total}}</td><td>{{printf "%.1f" .FCoverage.Percentage}}%</td><td>{{.FCoverage.Hit}}/{{.FCoverage.Total}}</td></tr>
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

			buildHtmlForSource(outdir, name)
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

func buildHtmlForSource(outdir, sourcename string) error {
	filename := filepath.Join(outdir, sourcename+".html")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpa, tmpb := lineCoverageForFile(lineCountData[sourcename])
	lcoverage := Coverage{tmpa, tmpb}
	tmpa, tmpb = funcCoverageForFile(funcCountData[sourcename])
	fcoverage := Coverage{tmpa, tmpb}

	file.WriteString("<html>\n")
	tmplHead.Execute(file, map[string]interface{}{
		"Title": "Coverage - " + sourcename,
	})
	file.WriteString("<body>\n")
	tmplH1.Execute(file, map[string]interface{}{
		"Title": "Coverage - " + sourcename,
	})
	file.WriteString(`<div class="pure-g"><div class="pure-u">`)
	file.WriteString(`<h2>Overall</h2>`)
	tmplCoverage.Execute(file, map[string]interface{}{
		"Title":     "Coverage - " + sourcename,
		"LCoverage": lcoverage,
		"FCoverage": fcoverage,
	})
	file.WriteString(`</div></div>`)
	file.WriteString(`<div class="pure-g"><div class="pure-u">`)
	file.WriteString(`<table class="source"><thead>`)
	file.WriteString(`<tr><th>Line #</th><th>Line data</th><th>Source code</th></tr>`)
	file.WriteString(`</thead><tbody>`)
	buildSourceListing(file, sourcename)
	file.WriteString(`</tbody></table>`)
	file.WriteString(`</div></div>`)
	file.WriteString("</body>\n</html>\n")
	return nil
}

func buildSourceListing(out *os.File, sourcename string) error {
	filename := filepath.Join(*srcdir, sourcename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
		return err
	}
	defer file.Close()

	lineNo := 1
	println("--")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if lc, ok := lineCountData[sourcename][lineNo]; ok {
			cl := "miss"
			if lc > 0 {
				cl = "hit"
			}
			fmt.Fprintf(out, "<tr class=\"%s\"><td class=\"ln\">%d</td><td class=\"ld\">%d</td><td>%s</td></tr>\n",
				cl, lineNo, lc,
				template.HTMLEscapeString(scanner.Text()),
			)
		} else {
			fmt.Fprintf(out, "<tr><td class=\"ln\">%d</td><td></td><td>%s</td></tr>\n",
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
