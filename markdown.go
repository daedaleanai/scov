package main

import (
	"text/template"

	"gitlab.com/stone.code/scov/internal/tool"
)

var (
	mdtmpltop = template.New("markdown")
	_         = template.Must(mdtmpltop.New("coverageRow").Parse(
		`{{.Hits}} | {{.Total}} | {{printf "%.1f" .P}}%`,
	))
	_ = template.Must(mdtmpltop.New("coverageDetail").Parse(
		`{{if .Valid}} {{.Hits}}/{{.Total}} ({{printf "%.1f" .P}}%) {{else}} No Data {{end}}`,
	))
	_ = template.Must(mdtmpltop.New("footer").Parse(
		`***
Generated by [SCov](https://gitlab.com/stone.code/scov).
`,
	))
	mdtmpl = template.Must(mdtmpltop.New("markdown").Parse(
		`# {{.Title}}

## Metadata

{{ if or .SrcID .TestID -}}
|       |       |
| :---- | :---- |
| Date: | {{.UnixDate }} |
{{ if .SrcID -}}
| Source ID: | {{.SrcID}} |
{{ end -}}
{{ if .TestID -}}
| Test ID: | {{.TestID}} |
{{ end -}}
{{- else -}}
Date: {{.UnixDate}}
{{ end }}

## Coverage Summary

|        | Hits   | Total  | Coverage |
| :----- | :----: | :----: | :------: |
| Lines: | {{template "coverageRow" .LCoverage}} |
{{ if .FCoverage.Valid -}}
| Functions: | {{template "coverageRow" .FCoverage}} |
{{ end -}}
{{ if .BCoverage.Valid -}}
| Branches: | {{template "coverageRow" .BCoverage}} |
{{ end -}}
{{ if .RCoverage.Valid -}}
| Regions: | {{template "coverageRow" .RCoverage}} |
{{ end }}

## By File

{{ $useFunc := .FCoverage.Valid -}}
{{ $useBranch := .BCoverage.Valid -}}
{{ $useRegion := .RCoverage.Valid -}}
| Filename | Line Coverage |{{if $useFunc }} Function Coverage |{{end}}{{if $useBranch}} Branch Coverage |{{end}}
| :------- | :-----------: |{{if $useFunc }} :---------------: |{{end}}{{if $useBranch}} :-------------: |{{end}}
{{range $ndx, $data := .Files -}}
| {{.Name}} |{{template "coverageDetail" .LCoverage}} 
{{- if $useFunc -}}
|{{template "coverageDetail" .FCoverage}} 
{{- end -}}
{{- if $useBranch -}}
|{{template "coverageDetail" .BCoverage}} 
{{- end -}}
{{- if $useRegion -}}
|{{template "coverageDetail" .RCoverage}} 
{{- end -}}
|
{{ end }}

## By Function

| Function | Hits |
| :------- | :--: |
{{range $ndx, $data := .Funcs -}}
| {{.Name}} | {{.HitCount }} |
{{ end }}

{{ template "footer" }}
`,
	))
)

func createMarkdownReport(filename string, report *Report) error {
	w, err := tool.Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	err = mdtmpl.Execute(w.File(), report)
	w.Keep(err)
	return err
}
