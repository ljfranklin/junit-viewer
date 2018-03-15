package junit_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ljfranklin/junit-viewer/internal/junit"
	"github.com/ljfranklin/junit-viewer/internal/test/helpers"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixturePath("success.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 2)
	helpers.AssertEquals(t, results[0].Failures, 0)
	expectedTime, err := time.Parse(junit.TimeFormat, "2018-03-15T14:22:46+07:00")
	if err != nil {
		t.Fatal(err)
	}
	helpers.AssertEquals(t, results[0].Timestamp, expectedTime)
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].Time), "9.837")
	helpers.AssertEquals(t, results[0].Properties, map[string]string{
		"go.version": "go1.9.2",
	})

	helpers.AssertEquals(t, len(results[0].TestCases), 2)
	helpers.AssertEquals(t, results[0].TestCases[0].Name, "TestS3Get")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[0].Time), "2.560")
	helpers.AssertEquals(t, results[0].TestCases[1].Name, "TestS3CompatibleGet")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[1].Time), "1.860")
}

func fixturePath(fixture string) string {
	return filepath.Join("..", "..", "fixtures", "junit", fixture)
}
