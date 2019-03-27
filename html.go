package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var (
	tmpl1 = template.New("html").Funcs(template.FuncMap{"htmlSafe": htmlSafe})
	tmpl2 = template.Must(tmpl1.New("sparkbar").Parse(
		`<div class="sparkbar">{{if gt .P 99.0}}<div class="fill {{.Rating}}" style="width:100%"></div>{{else}}<div class="fill {{.Rating}}" style="width:{{printf "%.1f" .P}}%"></div><div class="empty" style="width:{{printf "%.1f" .Q}}%"></div>{{end}}</div>`,
	))
	tmplHead = template.Must(tmpl1.New("head").Parse(
		`<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta name="description" content="Code coverage report">
<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
{{htmlSafe "<!--[if lte IE 8]>"}}
	<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-old-ie-min.css">
{{htmlSafe "<![endif]-->"}}
{{htmlSafe "<!--[if gt IE 8]><!-->"}}
	<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-min.css">
{{htmlSafe "<!--<![endif]-->"}}
<link rel="generator" href="https://gitlab.com/stone.code/gcovhtml">
<style>
html { padding:1em; }
body { max-width:70em; margin:auto; }
table { margin-bottom: 1em; }
.coverage { margin-left:auto;margin-right:0; }
.coverage td:nth-child(2) { text-align:center; }
.coverage td:nth-child(3) { text-align:center; }
.coverage td:nth-child(4) { text-align:center; }
.sparkbar { border: 1px solid black; border-radius:1px; min-width:50px; height:1em; }
.sparkbar .fill { display: inline-block; height: 100%; }
.sparkbar .high { background-color:lightgreen; }
.sparkbar .medium { background-color:yellow; }
.sparkbar .low { background-color:red; }
.sparkbar .empty { display: inline-block; height: 1em; background-color: white; }
{{ if .Source -}}
.source { font-family: monospace; width:100%; margin:0; }
.source th { padding: .1em .5em; text-align:left; border-bottom: 1px solid black; }
.source td { padding: .1em .5em; white-space: pre; }
.source .hit { background:lightblue; }
.source .miss { background:LightCoral; }
.source td:nth-child(1), .source th:nth-child(1) { background:PaleGoldenrod; text-align:right; }
{{ if .BCoverage.Valid -}}
.source td:nth-child(2), .source th:nth-child(2) { background:#f2edbf; text-align:right; }
.source td:nth-child(3), .source th:nth-child(3) { background:#f6f3d4; text-align:right; }
{{ else -}}
.source td:nth-child(2), .source th:nth-child(2) { background:#f6f3d4; text-align:right; }
{{ end -}}
{{ end -}}
</style>
</head>
`,
	))
	tmplH1 = template.Must(tmpl1.New("h1").Parse(
		`<div class="pure-g"><h1 class="pure-u">{{.Title}}</h1></div>`,
	))
	tmplCoverageRow = template.Must(tmpl1.New("coverageRow").Parse(
		`<td>{{.Hits}}</td><td>{{.Total}}</td><td>{{printf "%.1f" .P}}%</td>`,
	))
	tmplCoverage = template.Must(tmpl1.New("coverage").Parse(
		`<table class="pure-table pure-table-horizontal coverage">
<thead><tr><th></th><th>Hit</th><th>Total</th><th>Coverage</th><tr></thead>
<tbody>
<tr><td>Lines:</td>{{template "coverageRow" .LCoverage}}</tr>
<tr><td>Functions:</td>{{template "coverageRow" .FCoverage}}</tr>
{{ if .BCoverage.Valid -}}
<tr><td>Branches:</td>{{template "coverageRow" .BCoverage}}</tr>
{{ end -}}
</tbody>
</table>`,
	))
	tmpl = template.Must(tmpl1.Parse(
		`<!DOCTYPE html>
<html>
{{template "head" . -}}
<body>
{{template "h1" .}}
<div class="pure-g"><div class="pure-u">
<h2>Overall</h2>
</div></div>
<div class="pure-g"><div class="pure-u-1 pure-u-md-1-2">
{{ if .ID -}}
<table class="pure-table pure-table-horizontal">
<tr><td>Date:</td><td>{{.Date}}</td></tr>
<tr><td>Source&nbsp;ID:</td><td>{{.ID}}</td></tr>
</table>
{{ else -}}
<p>Date: {{.Date}}</p>
{{ end -}}
</div><div class="pure-u-1 pure-u-md-1-2">
{{template "coverage" .}}
</div></div>
<div class="pure-g"><div class="pure-u-1">
<h2>By File</h2>
<table class="pure-table pure-table-bordered" style="width:100%">
{{ $useBranch := .BCoverage.Valid -}}
<thead><tr><th>Filename</th><th colspan="3">Line Coverage</th><th colspan="3">Function Coverage</th>{{if $useBranch}}<th colspan="3">Branch Coverage</th>{{end}}</tr></thead>
<tbody>
{{range $ndx, $data := .Files -}}
<tr><td><a href="{{.Name}}.html">{{.Name}}</a></td><td>{{template "sparkbar" .LCoverage}}</td><td>{{.LCoverage.Hits}}/{{.LCoverage.Total}}</td><td>{{printf "%.1f" .LCoverage.P}}%</td><td>{{template "sparkbar" .FCoverage}}</td><td>{{.FCoverage.Hits}}/{{.FCoverage.Total}}</td><td>{{printf "%.1f" .FCoverage.P}}%</td>
{{- if $useBranch -}}
{{- if .BCoverage.Valid -}}
<td>{{template "sparkbar" .BCoverage}}</td><td>{{.BCoverage.Hits}}/{{.BCoverage.Total}}</td><td>{{printf "%.1f" .BCoverage.P}}%</td>
{{- else -}}
<td colspan="3">No data</td>
{{- end -}}
{{- end -}}
</tr>
{{end -}}
</tbody>
</table>
</div></div>
</body>
</html>`,
	))
	tmplSource1 = template.Must(tmpl1.New("sourcePrefix").Parse(
		`<!DOCTYPE html>
<html>
{{template "head" . -}}
<body>
{{template "h1" .}}
<div class="pure-g"><div class="pure-u">
<h2>Overall</h2>
</div></div>
<div class="pure-g"><div class="pure-u-1 pure-u-md-1-2">
<table class="pure-table pure-table-horizontal">
<tr><td>Date:</td><td>{{.Date}}</td></tr>
<tr><td>Filename:</td><td>{{.Filename}}</td></tr>
</table>
</div><div class="pure-u-1 pure-u-md-1-2">
{{template "coverage" .}}
</div></div>
<div class="pure-g"><div class="pure-u">
<h2>File Listing</h2>
<table class="source"><thead>
<tr><th>Line #</th>{{if .BCoverage.Valid}}<th>Branches</th>{{end}}<th>Hit count</th><th>Source code</th></tr>
</thead><tbody>
`,
	))
	tmplSource2 = template.Must(tmpl1.New("sourcePostfix").Parse(
		`</tbody></table>
</div></div>
</body></html>`,
	))
)

type FileStatistics struct {
	Name      string
	LCoverage Coverage
	FCoverage Coverage
	BCoverage Coverage
}

func createHTML(outdir string, data map[string]*FileData, date time.Time) error {
	err := os.MkdirAll(outdir, 0700)
	if err != nil {
		return err
	}

	err = createHTMLIndex(filepath.Join(outdir, "index.html"), data, date)
	if err != nil {
		return err
	}

	for name, data := range data {
		filename := filepath.Join(outdir, name+".html")
		err = createHTMLForSource(filename, name, data, date)
		if err != nil {
			return err
		}
	}

	return nil
}

func createHTMLIndex(filename string, data map[string]*FileData, date time.Time) error {
	w, err := Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeHTMLIndex(w.File(), data, date)
	w.Keep(err)
	return err
}

func writeHTMLIndex(out io.Writer, data map[string]*FileData, date time.Time) error {
	LCov := Coverage{}
	FCov := Coverage{}
	BCov := Coverage{}
	files := []FileStatistics{}
	for name, data := range data {
		stats := FileStatistics{Name: name}

		stats.LCoverage = data.LineCoverage()
		LCov.Accumulate(stats.LCoverage)
		stats.FCoverage = data.FuncCoverage()
		FCov.Accumulate(stats.FCoverage)
		stats.BCoverage = data.BranchCoverage()
		BCov.Accumulate(stats.BCoverage)

		files = append(files, stats)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	params := map[string]interface{}{
		"Title":     *title,
		"ID":        *srcid,
		"LCoverage": LCov,
		"FCoverage": FCov,
		"BCoverage": BCov,
		"Files":     files,
		"Date":      date.Format(time.UnixDate),
	}

	return tmpl.Execute(out, params)
}

func createHTMLForSource(filename string, sourcename string, data *FileData, date time.Time) error {
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return err
	}

	w, err := Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeHTMLForSource(w.File(), sourcename, data, date)
	w.Keep(err)
	return err
}

func writeHTMLForSource(out io.Writer, sourcename string, data *FileData, date time.Time) error {
	bcov := data.BranchCoverage()
	params := map[string]interface{}{
		"Title":     *title + " > " + filepath.Base(sourcename),
		"Source":    true,
		"Date":      date.Format(time.UnixDate),
		"Filename":  sourcename,
		"LCoverage": data.LineCoverage(),
		"FCoverage": data.FuncCoverage(),
		"BCoverage": bcov,
	}

	err := tmplSource1.Execute(out, params)
	if err != nil {
		return err
	}
	err = writeSourceListing(out, sourcename, data.LineData, bcov.Valid(), data.BranchData)
	if err != nil {
		return err
	}
	return tmplSource2.Execute(out, params)
}

func writeSourceListing(writer io.Writer, sourcename string, lineCountData map[int]uint64, withBranchData bool, branchData map[int][]BranchStatus) error {
	filename := filepath.Join(*srcdir, sourcename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(writer)

	lineNo := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.WriteString("<tr")
		lc, ok := lineCountData[lineNo]
		if ok {
			if lc > 0 {
				w.WriteString(` class="hit">`)
			} else {
				w.WriteString(` class="miss">`)
			}
		} else {
			w.WriteString(">")
		}
		fmt.Fprintf(w, "<td>%d</td>", lineNo)
		writeBranchDescription(w, withBranchData, branchData[lineNo])
		if ok {
			fmt.Fprintf(w, `<td>%d</td>`, lc)
		} else {
			w.WriteString(`<td></td>`)
		}
		fmt.Fprintf(w, "<td>%s</td></tr>\n", template.HTMLEscapeString(scanner.Text()))
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return w.Flush()
}

func htmlSafe(text string) template.HTML {
	return template.HTML(text)
}

func writeBranchDescription(w *bufio.Writer, withBranchData bool, data []BranchStatus) {
	if !withBranchData {
		return
	}
	if len(data) == 0 {
		w.WriteString(`<td></td>`)
		return
	}
	if data[0] == BranchNotExec {
		w.WriteString(`<td>[ NE ]</td>`)
		return
	}

	w.WriteString(`<td>[`)
	for _, v := range data {
		if v == BranchTaken {
			w.WriteString(" +")
		} else {
			w.WriteString(" -")
		}
	}
	w.WriteString(` ]</td>`)
}
