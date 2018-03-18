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

			fmt.Printf("## Summary of last %d run(s)\n\n", len(results))

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Tests", "Passed", "Failed", "Skipped", "Time", "When"})

			for _, suite := range results {
				table.Append([]string{
					fmt.Sprintf("%d", suite.Tests),
					fmt.Sprintf("%d (%.1f%%)", suite.Successes, (float64(suite.Successes)/float64(suite.Tests))*100),
					fmt.Sprintf("%d (%.1f%%)", suite.Failures, (float64(suite.Failures)/float64(suite.Tests))*100),
					fmt.Sprintf("%d (%.1f%%)", suite.Skips, (float64(suite.Skips)/float64(suite.Tests))*100),
					fmt.Sprintf("%.3f", suite.Time),
					fmt.Sprintf("%s", suite.Timestamp.Format(junit.TimeFormat)),
				})
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
