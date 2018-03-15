package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func main() {
	fmt.Printf("## Failures matching regex\n\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Matches", "Regex", "Last Failure"})

	// RFC822Z
	data := [][]string{
		[]string{"10", `\w502\w`, "2006-01-02T15:04:05Z07:00"},
		[]string{"3", `InsufficientResources`, "2006-01-02T15:04:05Z07:00"},
	}
	for _, v := range data {
		table.Append(v)
	}

	// markdown table
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
}
