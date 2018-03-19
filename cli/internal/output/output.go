package output

import (
	"fmt"
	"sort"
	"time"

	junit "github.com/ljfranklin/junit-viewer"
	"github.com/olekukonko/tablewriter"
)

func PrintPassFail(results junit.TestSuites, table *tablewriter.Table) {
	fmt.Printf("## Summary of last %d run(s)\n\n", len(results))

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
}

type testWithCount struct {
	TestCase     junit.TestCase
	TotalCount   int
	FailureCount int
	LastRan      time.Time
	LastFailed   time.Time
}

func PrintFrequentFailures(results junit.TestSuites, table *tablewriter.Table) {
	testCounts := map[string]*testWithCount{}
	for _, ts := range results {
		for _, tc := range ts.TestCases {
			if _, ok := testCounts[tc.Name]; !ok {
				testCounts[tc.Name] = &testWithCount{
					TestCase:     tc,
					TotalCount:   0,
					FailureCount: 0,
					LastRan:      ts.Timestamp,
					LastFailed:   time.Time{},
				}
			}
			testCounts[tc.Name].TotalCount++
			if testCounts[tc.Name].LastRan.Before(ts.Timestamp) {
				testCounts[tc.Name].LastRan = ts.Timestamp
			}
			if tc.Failed {
				testCounts[tc.Name].FailureCount++
				if testCounts[tc.Name].LastFailed.Before(ts.Timestamp) {
					testCounts[tc.Name].LastFailed = ts.Timestamp
				}
			}
		}
	}
	outputList := []testWithCount{}
	for _, counts := range testCounts {
		if counts.FailureCount > 0 {
			outputList = append(outputList, *counts)
		}
	}
	sort.Slice(outputList, func(i, j int) bool {
		return float64(outputList[i].FailureCount)/float64(outputList[i].TotalCount) > float64(outputList[j].FailureCount)/float64(outputList[j].TotalCount)
	})

	fmt.Printf("## Most frequent failures in last %d run(s)\n\n", len(results))

	table.SetHeader([]string{"Test", "Failed", "Last Failed", "Last Ran"})

	for _, tc := range outputList {
		table.Append([]string{
			tc.TestCase.Name,
			fmt.Sprintf("%d (%.1f%%)", tc.FailureCount, (float64(tc.FailureCount)/float64(tc.TotalCount))*100),
			tc.LastFailed.Format(junit.TimeFormat),
			tc.LastRan.Format(junit.TimeFormat),
		})
	}
}
