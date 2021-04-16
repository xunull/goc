package base_api

import (
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"strings"
)

func GetRoutes(c *gin.Context) {
	routes := Engine.Routes()
	tableStrings := &strings.Builder{}
	table := tablewriter.NewWriter(tableStrings)

	for _, route := range routes {
		data := []string{route.Method, route.Path, route.Handler}
		table.Append(data)
	}
	table.Render()
	c.String(http.StatusOK, tableStrings.String())
}
