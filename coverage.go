package main

type Coverage struct {
	Hits  int
	Total int
}

func (c Coverage) Percentage() float32 {
	return float32(c.Hits) * 100 / float32(c.Total)
}

func (c Coverage) P() float32 {
	return float32(c.Hits) * 100 / float32(c.Total)
}

func (c Coverage) Q() float32 {
	return 100 - float32(c.Hits)*100/float32(c.Total)
}

func (c *Coverage) Update(delta Coverage) {
	c.Hits += delta.Hits
	c.Total += delta.Total
}

type FileData struct {
	Filename string
	LineData map[int]uint64
	FuncData map[string]uint64
}

func NewFileData(filename string) FileData {
	return FileData{
		Filename: filename,
		LineData: make(map[int]uint64),
		FuncData: make(map[string]uint64),
	}
}

func (file *FileData) LineCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.LineData {
		if v != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}

func (file *FileData) FuncCoverage() Coverage {
	a, b := 0, 0

	for _, v := range file.FuncData {
		if v != 0 {
			a++
		}
		b++
	}
	return Coverage{a, b}
}
