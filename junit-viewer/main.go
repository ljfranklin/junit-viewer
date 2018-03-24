package main

import (
	"fmt"
	"log"
	"os"

	junit "github.com/ljfranklin/junit-viewer"
	"github.com/ljfranklin/junit-viewer/junit-viewer/internal/output"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func main() {
	var (
		outputType string
	)

	rootCmd := &cobra.Command{
		Use:   "junit-viewer junit-xml-files...",
		Short: "View summary of JUnit XML files",
		Long:  "",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			results := junit.TestSuites{}
			for _, inputFile := range args {
				result, err := junit.LoadFile(inputFile)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, result...)
			}
			results.SortByTimestamp()

			table := tablewriter.NewWriter(os.Stdout)

			switch outputType {
			case "pass-fail":
				output.PrintPassFail(results, table)
			case "frequent-failures":
				output.PrintFrequentFailures(results, table)
			default:
				log.Fatalf("unknown output-type '%s'", outputType)
			}

			// markdown table
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")

			table.Render()
		},
	}
	rootCmd.PersistentFlags().StringVarP(&outputType, "output-type", "o", "", "(required) how results are presented; supports: 'pass-fail', 'frequent-failures'")
	rootCmd.MarkPersistentFlagRequired("output-type")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
