package gin_extend

import (
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"strings"
)

func GetRouteTable(e *gin.Engine) string {
	routes := e.Routes()

	tableStrings := &strings.Builder{}
	table := tablewriter.NewWriter(tableStrings)

	for _, route := range routes {
		data := []string{route.Method, route.Path, route.Handler}
		table.Append(data)
	}
	table.Render()
	return tableStrings.String()

}
