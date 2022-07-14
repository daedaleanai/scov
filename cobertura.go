package main

import (
	"encoding/xml"
	"fmt"
	"path"
	"sort"
	"strings"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + "<!DOCTYPE coverage SYSTEM 'http://cobertura.sourceforge.net/xml/coverage-04.dtd'>" + "\n"
)

type XmlPackage struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`

	LineRate   float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`

	Classes XmlWrapClasses `xml:"classes"`

	lineTotal     int
	lineCovered   int
	branchTotal   int
	branchCovered int
}

type XmlWrapMethods struct {
	Methods []XmlMethod `xml:"method"`
}

type XmlMethod struct {
	XMLName    xml.Name     `xml:"method"`
	Name       string       `xml:"name,attr"`
	Signature  string       `xml:"signature,attr"`
	LineRate   float64      `xml:"line-rate,attr"`
	BranchRate float64      `xml:"branch-rate,attr"`
	Lines      XmlWrapLines `xml:"lines"`
}

type XmlWrapClasses struct {
	Classes []XmlClass `xml:"class"`
}

type XmlClass struct {
	XMLName    xml.Name `xml:"class"`
	Name       string   `xml:"name,attr"`
	FileName   string   `xml:"filename,attr"`
	LineRate   float64  `xml:"line-rate,attr"`
	BranchRate float64  `xml:"branch-rate,attr"`
	Complexity float64  `xml:"complexity,attr"`

	Methods XmlWrapMethods `xml:"methods"`
	Lines   XmlWrapLines   `xml:"lines"`

	lines uint64
}

type XmlLine struct {
	XMLName xml.Name `xml:"line"`
	Number  uint64   `xml:"number,attr"`
	Hits    uint64   `xml:"hits,attr"`
	Branch  bool     `xml:"branch,attr"`

	ConditionCoverage string `xml:"condition-coverage,attr,omitempty"`

	Conditions []string `xml:"conditions,attr,omitempty"` // TODO ?
}

type XmlWrapLines struct {
	Lines []XmlLine `xml:"line"`
}

type XmlWrapSources struct {
	Sources []string `xml:"source"`
}

type XmlWrapPackages struct {
	Packages []*XmlPackage `xml:"package"`
}

type XmlCoverage struct {
	XMLName          xml.Name `xml:"coverage"`
	LineRate         float64  `xml:"line-rate,attr"`
	BranchRate       float64  `xml:"branch-rate,attr"`
	LinesCovered     uint64   `xml:"lines-covered,attr"`
	LinesValid       uint64   `xml:"lines-valid,attr"`
	BranchesCovered  uint64   `xml:"branches-covered,attr"`
	BranchesValid    uint64   `xml:"branches-valid,attr"`
	FunctionRate     float64  `xml:"function-rate,attr"`
	FunctionsCovered uint64   `xml:"functions-covered,attr"`
	FunctionsValid   uint64   `xml:"functions-valid,attr"`
	Complexity       float64  `xml:"complexity,attr"`
	Timestamp        uint64   `xml:"timestamp,attr"`
	Version          string   `xml:"version,attr"`

	Sources  XmlWrapSources  `xml:"sources"`
	Packages XmlWrapPackages `xml:"packages"`

	Name  string   `xml:"name"`
	Email string   `xml:"email"`
	Phone []string `xml:"phone"`
}

type ByNumber []XmlLine

func (s ByNumber) Len() int {
	return len(s)
}

func (s ByNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByNumber) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}

func createCoberturaReport(filename string, data FileDataSet, report *Report) error {
	w, err := Open(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	xmlReport := XmlCoverage{
		LineRate:         float64(report.LCoverage.P()) / 100.0,
		BranchRate:       float64(report.BCoverage.P()) / 100.0,
		LinesCovered:     uint64(report.LCoverage.Hits),
		LinesValid:       uint64(report.LCoverage.Total),
		BranchesCovered:  uint64(report.BCoverage.Hits),
		BranchesValid:    uint64(report.BCoverage.Total),
		FunctionRate:     float64(report.FCoverage.P()) / 100.0,
		FunctionsCovered: uint64(report.FCoverage.Hits),
		FunctionsValid:   uint64(report.FCoverage.Total),
		Complexity:       0,
		Timestamp:        uint64(report.Date.Unix()),
		Version:          "2.0.3",
		Sources:          XmlWrapSources{Sources: []string{"."}},
		Packages:         XmlWrapPackages{Packages: []*XmlPackage{}},
	}

	packages := map[string]*XmlPackage{}

	for _, file := range report.Files {
		parent := path.Dir(file.Name)
		pkgName := strings.ReplaceAll(parent, "/", ".")

		pkg, ok := packages[pkgName]

		if !ok {
			pkg = &XmlPackage{
				Name: pkgName,
			}
			packages[pkgName] = pkg
			xmlReport.Packages.Packages = append(xmlReport.Packages.Packages, pkg)
		}

		pkg.lineTotal = pkg.lineTotal + data.FileData(file.Name).LineCoverage().Total
		pkg.lineCovered = pkg.lineCovered + data.FileData(file.Name).LineCoverage().Hits
		pkg.branchTotal = pkg.branchTotal + data.FileData(file.Name).BranchCoverage().Total
		pkg.branchCovered = pkg.branchCovered + data.FileData(file.Name).BranchCoverage().Hits

		cls := XmlClass{
			Name:       strings.ReplaceAll(file.Name, "/", "."),
			FileName:   file.Name,
			LineRate:   float64(file.LCoverage.P()) / 100.0,
			BranchRate: float64(file.BCoverage.P()) / 100.0,
			Complexity: 0,
		}

		for name, finfo := range data.FileData(file.Name).FuncData {
			rate := 0
			if finfo.HitCount > 0 {
				rate = 1
			}

			cls.Methods.Methods = append(cls.Methods.Methods, XmlMethod{
				Name:       name,
				BranchRate: float64(rate),
				LineRate:   float64(rate),
				Signature:  "",
				Lines: XmlWrapLines{
					Lines: []XmlLine{
						{
							Number: uint64(finfo.StartLine),
							Hits:   finfo.HitCount,
							Branch: false,
						},
					},
				},
			})
		}

		for idx, hits := range data.FileData(file.Name).LineData {
			branch := len(data.FileData(file.Name).BranchData[idx]) > 1
			branchCoverage := ""
			if branch {
				covered := 0
				for _, b := range data.FileData(file.Name).BranchData[idx] {
					if b == BranchTaken {
						covered = covered + 1
					}
				}
				branchCoverage = fmt.Sprintf("%d%% (%d/%d)", covered*100/len(data.FileData(file.Name).BranchData[idx]), covered, len(data.FileData(file.Name).BranchData[idx]))
			}
			cls.Lines.Lines = append(cls.Lines.Lines, XmlLine{
				Number:            uint64(idx),
				Hits:              hits,
				Branch:            branch,
				ConditionCoverage: branchCoverage,
			})
		}

		sort.Sort(ByNumber(cls.Lines.Lines))

		cls.lines = uint64(data.FileData(file.Name).LineCoverage().Total)

		pkg.Classes.Classes = append(pkg.Classes.Classes, cls)
	}

	for _, pkg := range packages {
		pkg.LineRate = float64(pkg.lineCovered) * 100.0 / float64(pkg.lineTotal)
		pkg.BranchRate = float64(pkg.branchCovered) * 100.0 / float64(pkg.branchTotal)
	}

	_, err = w.File().Write([]byte(Header))
	if err != nil {
		return err
	}

	xmldata, err := xml.Marshal(xmlReport)
	if err != nil {
		return err
	}
	_, err = w.File().Write(xmldata)

	w.Keep(err)
	return err
}
