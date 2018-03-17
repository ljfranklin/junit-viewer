package main

import (
	"fmt"
	"os"

	junit "github.com/ljfranklin/junit-viewer"
	"github.com/olekukonko/tablewriter"
)

func main() {
	junitPath := os.Args[1]
	results, err := junit.LoadFile(junitPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("## Summary\n\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tests", "Passed", "Failed", "Skipped"})

	for _, suite := range results {
		table.Append([]string{
			fmt.Sprintf("%d", suite.Tests),
			fmt.Sprintf("%d", suite.Successes),
			fmt.Sprintf("%d", suite.Failures),
			fmt.Sprintf("%d", suite.Skips),
		})
	}

	// markdown table
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.Render()
}
