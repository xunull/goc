package wrap_tablewriter

import (
	"github.com/olekukonko/tablewriter"
	"strconv"
	"strings"
)

func GetSimpleTable(data []string) string {
	tableStr := &strings.Builder{}
	table := tablewriter.NewWriter(tableStr)
	table.SetHeader([]string{"Id", "Data"})
	for i, name := range data {
		data := []string{
			strconv.Itoa(i + 1), name,
		}
		table.Append(data)
	}
	table.Render()
	return tableStr.String()
}

func GetSimpleListTable(data [][]string) string {
	tableStr := &strings.Builder{}
	table := tablewriter.NewWriter(tableStr)
	for i, list := range data {
		d := make([]string, 0, len(list)+1)
		d = append(d, strconv.Itoa(i+1))
		d = append(d, list...)
		table.Append(d)
	}
	table.Render()
	return tableStr.String()
}
