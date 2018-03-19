package main

import (
	"fmt"
	"os"

	junit "github.com/ljfranklin/junit-viewer"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func main() {
	var (
		outputType string
	)

	rootCmd := &cobra.Command{
		Use:   "junit-viewer",
		Short: "TODO",
		Long:  "TODO",
		// TODO: verify at least one arg
		Run: func(cmd *cobra.Command, args []string) {
			results := junit.TestSuites{}
			for _, inputFile := range args {
				result, err := junit.LoadFile(inputFile)
				if err != nil {
					panic(err)
				}
				results = append(results, result...)
			}
			results.SortByTimestamp()

			table := tablewriter.NewWriter(os.Stdout)

			switch outputType {
			case "pass-fail":
				printPassFail(results, table)
			case "frequent-failures":
				printFrequentFailures(results, table)
			}

			// markdown table
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")

			table.Render()
		},
	}
	rootCmd.PersistentFlags().StringVarP(&outputType, "output-type", "o", "", "TODO")
	rootCmd.MarkFlagRequired("output-type")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
