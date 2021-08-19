package outputx

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/xunull/goc/commonx"
	"strings"
)

func OutputStringMap(dm map[string]interface{}) {
	tableStrings := &strings.Builder{}
	table := tablewriter.NewWriter(tableStrings)
	for k, item := range dm {
		data := []string{k, commonx.JsonString(item)}
		table.Append(data)
	}
	table.Render()
	fmt.Println(tableStrings.String())
}
