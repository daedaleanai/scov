package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/stone.code/scov/internal/tool"
)

var (
	tmpl1 = template.New("html").Funcs(template.FuncMap{"htmlSafe": htmlSafe})
	_     = template.Must(tmpl1.New("sparkbar").Parse(
		`<div class="sparkbar">{{if gt .P 99.0}}<div class="fill {{.Rating}}" style="width:100%"></div>{{else}}<div class="fill {{.Rating}}" style="width:{{printf "%.1f" .P}}%"></div><div class="empty" style="width:{{printf "%.1f" .Q}}%"></div>{{end}}</div>`,
	))
	_ = template.Must(tmpl1.New("head").Parse(
		`<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<meta name="description" content="Code coverage report">
<meta name="generator" content="https://gitlab.com/stone.code/scov">
<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
{{htmlSafe "<!--[if lte IE 8]>"}}
	<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-old-ie-min.css">
{{htmlSafe "<![endif]-->"}}
{{htmlSafe "<!--[if gt IE 8]><!-->"}}
	<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-min.css">
{{htmlSafe "<!--<![endif]-->"}}
{{ if .ProjectURL -}}
<link rel="project" href="{{.ProjectURL}}">
{{ end -}}
<style>
html { padding:1em; }
body { max-width:70em; margin:auto; }
table { margin-bottom: 1em; }
.coverage { min-width:100%; }
.coverage td:nth-child(2), .coverage th:nth-child(2) { text-align:center; }
.coverage td:nth-child(3), .coverage th:nth-child(3) { text-align:center; }
.coverage td:nth-child(4), .coverage th:nth-child(4) { text-align:center; }
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
{{ if .Script -}}
th .reveal { float:right; transition: opacity 0.5s; opacity: 0.1; }
th .reveal .pure-button { padding: 0 0.5em; }
th:hover .reveal { opacity: 1; }
{{ end -}}
footer { border-top: 1px solid rgb(203, 203, 203); margin-top: 1em; background: #e0e0e0; padding: .5em 1em; }
@media screen and (min-width: 48em) {
	.pure-gutter-md > div { box-sizing: border-box; padding: 0 0.5em; }
	.pure-gutter-md > div:first-child { padding-left: 0; }
	.pure-gutter-md > div:last-child { padding-right: 0; }
}
@media screen and (max-width: 48em) {
	.table-md td, .table-md th { padding: 0.5em; }
}
</style>
{{ if .Script -}}
<script src="index.js" ></script>
{{ end -}}
</head>
`,
	))
	_ = template.Must(tmpl1.New("h1").Parse(
		`<div class="pure-g"><h1 class="pure-u">{{.Title}}</h1></div>`,
	))
	_ = template.Must(tmpl1.New("coverageRow").Parse(
		`<td>{{.Hits}}</td><td>{{.Total}}</td><td>{{printf "%.1f" .P}}%</td>`,
	))
	_ = template.Must(tmpl1.New("coverage").Parse(
		`<table class="pure-table pure-table-horizontal coverage">
<thead><tr><th></th><th>Hits</th><th>Total</th><th>Coverage</th></tr></thead>
<tbody>
<tr><td>Lines:</td>{{template "coverageRow" .LCoverage}}</tr>
{{ if .FCoverage.Valid -}}<tr><td>Functions:</td>{{template "coverageRow" .FCoverage}}</tr>
{{ end -}}
{{ if .BCoverage.Valid -}}<tr><td>Branches:</td>{{template "coverageRow" .BCoverage}}</tr>
{{ end -}}
{{ if .RCoverage.Valid -}}<tr><td>Regions:</td>{{template "coverageRow" .RCoverage}}</tr>
{{ end -}}
</tbody>
</table>`,
	))
	_ = template.Must(tmpl1.New("coverageDetail").Parse(
		`{{- if .Valid -}}
<td>{{template "sparkbar" .}}</td><td>{{.Hits}}/{{.Total}}</td><td>{{printf "%.1f" .P}}%</td>
{{- else -}}
<td colspan="3">No data</td>
{{- end -}}`,
	))
	_ = template.Must(tmpl1.New("metadata").Parse(
		`{{ if or .SrcID .TestID .Filename -}}
<table class="pure-table pure-table-horizontal">
<tr><td>Date:</td><td>{{.Date}}</td></tr>
{{ if .SrcID -}}
<tr><td>Source&nbsp;ID:</td><td>{{.SrcID}}</td></tr>
{{ end -}}
{{ if .TestID -}}
<tr><td>Test&nbsp;ID:</td><td>{{.TestID}}</td></tr>
{{ end -}}
{{ if .Filename -}}
<tr><td>Filename:</td><td>{{.Filename}}</td></tr>
{{ end -}}
</table>
{{- else -}}
<p>Date: {{.Date}}</p>
{{ end -}}`,
	))
	_ = template.Must(tmpl1.New("footer").Parse(
		`<footer>Generated by <a href="https://gitlab.com/stone.code/scov">SCov</a>.</footer>`,
	))
	tmpl = template.Must(tmpl1.Parse(
		`<!DOCTYPE html>
<html>
{{template "head" . -}}
<body>
{{template "h1" .}}
<div class="pure-g pure-gutter-md"><div class="pure-u-1 pure-u-md-1-2">
<h2>Metadata</h2>
{{ template "metadata" . -}}
</div><div class="pure-u-1 pure-u-md-1-2">
<h2>Coverage Summary</h2>
{{ template "coverage" . -}}
</div></div>
<div class="pure-g"><div class="pure-u-1">
<h2>By File</h2>
<table class="pure-table pure-table-bordered table-md" style="width:100%">
{{ $useFunc := .FCoverage.Valid -}}
{{ $useBranch := .BCoverage.Valid -}}
{{ $useRegion := .RCoverage.Valid -}}
<thead><tr><th{{if .Script}} data-sort="text"{{end}}>Filename</th><th colspan="3"{{if .Script}} data-sort="perc"{{end}}>Line Coverage</th>{{if $useFunc}}<th colspan="3"{{if .Script}} data-sort="perc"{{end}}>Function Coverage</th>{{end}}{{if $useBranch}}<th colspan="3"{{if .Script}} data-sort="perc"{{end}}>Branch Coverage</th>{{end}}{{if $useRegion}}<th colspan="3"{{if .Script}} data-sort="perc"{{end}}>Region Coverage</th>{{end}}</tr></thead>
<tbody>
{{range $ndx, $data := .Files -}}
<tr><td><a href="{{.Name}}.html">{{.Name}}</a></td>{{template "coverageDetail" .LCoverage}}
{{- if $useFunc -}}{{ template "coverageDetail" .FCoverage }}{{- end -}}
{{- if $useBranch -}}{{ template "coverageDetail" .BCoverage }}{{- end -}}
{{- if $useRegion -}}{{ template "coverageDetail" .RCoverage }}{{- end -}}
</tr>
{{end -}}
</tbody>
</table>
</div></div>
{{if $useFunc -}}
<div class="pure-g"><div class="pure-u-1">
<h2>By Function</h2>
<table class="pure-table pure-table-bordered table-md" style="width:100%">
<thead><tr><th{{if .Script}} data-sort="text"{{end}}>Function</th><th{{if .Script}} data-sort="perc"{{end}}>Hits</th></tr></thead>
<tbody>
{{range $ndx, $data := .Funcs -}}
<tr><td><a href="{{.Filename}}.html#L{{.StartLine}}">{{.Name}}</a></td><td>{{.HitCount}}</td></tr>
{{- end -}}
</tbody>
</table>
</div></div>
{{- end}}
{{ template "footer" . }}
</body>
</html>`,
	))
	tmplSource1 = template.Must(tmpl1.New("sourcePrefix").Parse(
		`<!DOCTYPE html>
<html>
{{template "head" . -}}
<body>
{{template "h1" .}}
<div class="pure-g pure-gutter-md"><div class="pure-u-1 pure-u-md-1-2">
<h2>Metadata</h2>
{{ template "metadata" . }}
</div><div class="pure-u-1 pure-u-md-1-2">
<h2>Coverage</h2>
{{ template "coverage" . }}
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
{{ template "footer" . }}
</body></html>`,
	))
)

func createHTML(outdir string, data map[string]*FileData, report *Report) error {
	err := os.MkdirAll(outdir, 0700)
	if err != nil {
		return err
	}

	err = createHTMLIndex(filepath.Join(outdir, "index.html"), report)
	if err != nil {
		return err
	}

	if report.AllowHTMLScripting {
		err = createJS(filepath.Join(outdir, "index.js"))
		if err != nil {
			return err
		}
	}

	for name, data := range data {
		filename := filepath.Join(outdir, name+".html")
		err = createHTMLForSource(filename, name, data, report)
		if err != nil {
			return err
		}
	}

	return nil
}

func createHTMLIndex(filename string, report *Report) error {
	w, err := tool.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeHTMLIndex(w.File(), report)
	w.Keep(err)
	return err
}

func writeHTMLIndex(out io.Writer, report *Report) error {
	params := map[string]interface{}{
		"Title":      report.Title,
		"SrcID":      report.SrcID,
		"TestID":     report.TestID,
		"ProjectURL": report.ProjectURL,
		"LCoverage":  report.LCoverage,
		"FCoverage":  report.FCoverage,
		"BCoverage":  report.BCoverage,
		"RCoverage":  report.RCoverage,
		"Files":      report.Files,
		"Funcs":      report.Funcs,
		"Date":       report.UnixDate(),
		"Script":     report.AllowHTMLScripting,
	}

	return tmpl.Execute(out, params)
}

func createJS(filename string) error {
	w, err := tool.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeJS(w.File())
	w.Keep(err)
	return err
}

func writeJS(out io.Writer) error {
	const js = `function sortTable( table, cmp ) {
	var tbody = table.getElementsByTagName('tbody')[0]

	var list = tbody.childNodes;
	var arr = [];
	for ( var i in list ) {
		if ( list[i].nodeType==1 ) arr.push( list[i].cloneNode(true) )
	}
	arr.sort( cmp )
	while ( tbody.firstChild ) tbody.removeChild(tbody.firstChild )
	for ( var i in arr ) {
		tbody.appendChild( arr[i] )
	}
}

var cmpTable = {
	'text' : function(a) {
		return a.textContent
	},
	'perc' : function(a) {
		return parseFloat(a.textContent)
	}
}

function makeOnClickCallback(column,cmp,dir) {
	return function() {
		var table = this.parentElement.parentElement.parentElement.parentElement.parentElement
		sortTable( table, function(a,b) {
			var v1 = cmp(a.childNodes[column])
			var v2 = cmp(b.childNodes[column])
			return v1 < v2 ? -dir : v1 > v2 ? +dir : 0
		})
	}
}

function cellIndex(e) {
	var thead = e.parentElement
	var arr = thead.getElementsByTagName('th')
	var cellindex = 0

	for ( var i of arr ) {
		if ( i == e ) {
			return cellindex
		}
		cellindex += i.colSpan
	}

	return -1
}

window.onload = function() {
	var elems = document.getElementsByTagName('th')
	for ( let elem of elems ) {
		var attr = elem.getAttribute('data-sort')
		if ( !attr ) continue

		var cellindex = cellIndex(elem)
		var colspan = elem.colSpan
		elem.innerHTML = elem.innerHTML + ' <span class="reveal"><a class="pure-button">▲</a><a class="pure-button">▼</a></span>'
		elem.getElementsByTagName('a')[0].onclick = makeOnClickCallback(cellindex+(colspan-1), cmpTable[attr], +1)
		elem.getElementsByTagName('a')[1].onclick = makeOnClickCallback(cellindex+(colspan-1), cmpTable[attr], -1)
	}
}`

	_, err := out.Write([]byte(js))
	return err
}

func createHTMLForSource(filename string, sourcename string, data *FileData, report *Report) error {
	err := os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		return err
	}

	w, err := tool.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = writeHTMLForSource(w.File(), sourcename, data, report)
	w.Keep(err)
	return err
}

func writeHTMLForSource(out io.Writer, sourcename string, data *FileData, report *Report) error {
	bcov := data.BranchCoverage()
	params := map[string]interface{}{
		"Title":     report.Title + " > " + filepath.Base(sourcename),
		"SrcID":     report.SrcID,
		"TestID":    report.TestID,
		"Source":    true,
		"Date":      report.UnixDate(),
		"Filename":  sourcename,
		"LCoverage": data.LineCoverage(),
		"FCoverage": data.FuncCoverage(),
		"BCoverage": bcov,
		"RCoverage": data.RegionCoverage(),
	}

	err := tmplSource1.Execute(out, params)
	if err != nil {
		return err
	}
	err = writeSourceListing(out, filepath.Join(report.SrcDir, sourcename), data.LineData, bcov.Valid(), data.BranchData)
	if err != nil {
		return err
	}
	return tmplSource2.Execute(out, params)
}

func rowClassAttribute(hitCount uint64, ok bool) string {
	if !ok {
		return ""
	}
	if hitCount > 0 {
		return ` class="hit"`
	}
	return ` class="miss"`
}

func writeSourceListing(writer io.Writer, filename string, lineCountData map[int]uint64, withBranchData bool, branchData map[int][]BranchStatus) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(writer)

	lineNo := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hitCount, ok := lineCountData[lineNo]
		fmt.Fprintf(w, `<tr id="L%d"%s>`, lineNo, rowClassAttribute(hitCount, ok))
		fmt.Fprintf(w, "<td>%d</td>", lineNo)
		writeBranchDescription(w, withBranchData, branchData[lineNo])
		if ok {
			fmt.Fprintf(w, `<td>%d</td>`, hitCount)
		} else {
			// Ignore error.  Since we are using a bufio.Writer, it will be
			// reported when we flush.
			_, _ = w.WriteString(`<td></td>`)
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
	// Ignore write errors.
	// Since we are using a bufio.Writer, write errors will be reported when
	// we flush.

	if !withBranchData {
		return
	}
	if len(data) == 0 {
		_, _ = w.WriteString(`<td></td>`)
		return
	}
	if data[0] == BranchNotExec {
		_, _ = w.WriteString(`<td>[ NE ]</td>`)
		return
	}

	_, _ = w.WriteString(`<td>[`)
	for _, v := range data {
		if v == BranchTaken {
			_, _ = w.WriteString(" +")
		} else {
			_, _ = w.WriteString(" -")
		}
	}
	_, _ = w.WriteString(` ]</td>`)
}
