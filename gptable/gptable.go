package gptable

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

var t table.Writer

func Init() {
	t = table.NewWriter()
	
	t.SetOutputMirror(os.Stdout)

	t.SetStyle(table.StyleLight)
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
}

func SetHeader(headers ...string) {
    row := make(table.Row, len(headers))
    for i, h := range headers {
        row[i] = h
    }

    t.AppendHeader(row)
}

func AppendRow(values ...string) {
    row := make(table.Row, len(values))
    for i, v := range values {
        row[i] = v
    }

	t.AppendRow(row)
}

func Render() {
	t.Render()
}
