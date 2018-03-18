package junit_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ljfranklin/junit-viewer"
	"github.com/ljfranklin/junit-viewer/internal/test/helpers"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "success.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 2)
	helpers.AssertEquals(t, results[0].Failures, 0)
	helpers.AssertEquals(t, results[0].Errors, 0)
	helpers.AssertEquals(t, results[0].Skips, 0)
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
	helpers.AssertEquals(t, results[0].TestCases[0].Passed, true)
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[0].Time), "2.560")
	helpers.AssertEquals(t, results[0].TestCases[1].Name, "TestS3CompatibleGet")
	helpers.AssertEquals(t, results[0].TestCases[1].Passed, true)
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[1].Time), "1.860")
}

func TestNoTestSuitesRoot(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "no-testsuites-root.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)
	helpers.AssertEquals(t, results[0].Tests, 2)
}

func TestFailureMessages(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "failures.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 0)
	helpers.AssertEquals(t, results[0].Failures, 2)
	helpers.AssertEquals(t, results[0].Errors, 0)
	helpers.AssertEquals(t, results[0].Skips, 0)

	helpers.AssertEquals(t, len(results[0].TestCases), 2)
	helpers.AssertEquals(t, results[0].TestCases[0].Name, "TestS3Get")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[0].Time), "0.000")
	helpers.AssertEquals(t, results[0].TestCases[0].Failed, true)
	helpers.AssertEquals(t, results[0].TestCases[0].FailureMessage, "Failed")
	helpers.AssertEquals(t, results[0].TestCases[0].FailureType, "integration")
	helpers.AssertEquals(t, results[0].TestCases[0].FailureOutput, "s3_test.go:102: AWS_ACCESS_KEY must be set")
	helpers.AssertEquals(t, results[0].TestCases[1].Name, "TestS3CompatibleGet")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[1].Time), "0.000")
	helpers.AssertEquals(t, results[0].TestCases[1].Failed, true)
	helpers.AssertEquals(t, results[0].TestCases[1].FailureMessage, "Failed")
	helpers.AssertEquals(t, results[0].TestCases[1].FailureType, "integration")
	helpers.AssertEquals(t, results[0].TestCases[1].FailureOutput, "s3_test.go:110: S3_COMPATIBLE_ACCESS_KEY must be set")
}

func TestErrorMessages(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "errors.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 0)
	helpers.AssertEquals(t, results[0].Failures, 0)
	helpers.AssertEquals(t, results[0].Errors, 2)
	helpers.AssertEquals(t, results[0].Skips, 0)

	helpers.AssertEquals(t, len(results[0].TestCases), 2)
	helpers.AssertEquals(t, results[0].TestCases[0].Name, "TestS3Get")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[0].Time), "0.000")
	helpers.AssertEquals(t, results[0].TestCases[0].Errored, true)
	helpers.AssertEquals(t, results[0].TestCases[0].ErrorMessage, "Failed")
	helpers.AssertEquals(t, results[0].TestCases[0].ErrorType, "integration")
	helpers.AssertEquals(t, results[0].TestCases[0].ErrorOutput, "s3_test.go:102: AWS_ACCESS_KEY must be set")
	helpers.AssertEquals(t, results[0].TestCases[1].Name, "TestS3CompatibleGet")
	helpers.AssertEquals(t, fmt.Sprintf("%.3f", results[0].TestCases[1].Time), "0.000")
	helpers.AssertEquals(t, results[0].TestCases[1].Errored, true)
	helpers.AssertEquals(t, results[0].TestCases[1].ErrorMessage, "Failed")
	helpers.AssertEquals(t, results[0].TestCases[1].ErrorType, "integration")
	helpers.AssertEquals(t, results[0].TestCases[1].ErrorOutput, "s3_test.go:110: S3_COMPATIBLE_ACCESS_KEY must be set")
}

func TestSkips(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "skips.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 0)
	helpers.AssertEquals(t, results[0].Failures, 0)
	helpers.AssertEquals(t, results[0].Errors, 0)
	helpers.AssertEquals(t, results[0].Skips, 2)

	helpers.AssertEquals(t, len(results[0].TestCases), 2)
	helpers.AssertEquals(t, results[0].TestCases[0].Name, "TestS3Get")
	helpers.AssertEquals(t, results[0].TestCases[0].Skipped, true)
	helpers.AssertEquals(t, results[0].TestCases[0].SkipMessage, "s3_test.go:101: Skipping this test")
	helpers.AssertEquals(t, results[0].TestCases[1].Name, "TestS3CompatibleGet")
	helpers.AssertEquals(t, results[0].TestCases[1].Skipped, true)
	helpers.AssertEquals(t, results[0].TestCases[1].SkipMessage, "s3_test.go:110: Skipping this test")
}

func TestSystemOutErr(t *testing.T) {
	t.Parallel()

	results, err := junit.Load(fixtureContents(t, "system-out-err.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)

	helpers.AssertEquals(t, len(results[0].TestCases), 1)
	helpers.AssertEquals(t, results[0].TestCases[0].Name, "TestS3Get")
	helpers.AssertEquals(t, results[0].TestCases[0].SystemOut, "some-stdout")
	helpers.AssertEquals(t, results[0].TestCases[0].SystemErr, "some-stderr")
}

func TestErrorInvalidXML(t *testing.T) {
	t.Parallel()

	_, err := junit.Load(fixtureContents(t, "malformed.xml"))
	if err == nil {
		t.Fatal("expected Load to fail but it succeeded")
	}
	if !strings.Contains(err.Error(), "parse XML") {
		t.Fatalf("expected '%s' to contain 'parse XML', but it did not", err.Error())
	}
}

func TestLoadFile(t *testing.T) {
	t.Parallel()

	results, err := junit.LoadFile(fixturePath("success.xml"))
	if err != nil {
		t.Fatal(err)
	}

	helpers.AssertEquals(t, len(results), 1)
	helpers.AssertEquals(t, results[0].Tests, 2)
	helpers.AssertEquals(t, results[0].Successes, 2)
}

func TestSortByTimestamp(t *testing.T) {
	results, err := junit.LoadFile(fixturePath("failures.xml"))
	if err != nil {
		t.Fatal(err)
	}
	moreResults, err := junit.LoadFile(fixturePath("success.xml"))
	if err != nil {
		t.Fatal(err)
	}

	results = append(results, moreResults...)

	results.SortByTimestamp()

	helpers.AssertEquals(t, len(results), 2)
	expectedTime, err := time.Parse(junit.TimeFormat, "2018-03-15T14:22:46+07:00")
	if err != nil {
		t.Fatal(err)
	}
	helpers.AssertEquals(t, results[0].Timestamp, expectedTime)
	expectedTime, err = time.Parse(junit.TimeFormat, "2018-03-14T10:12:34+07:00")
	if err != nil {
		t.Fatal(err)
	}
	helpers.AssertEquals(t, results[1].Timestamp, expectedTime)
}

func TestLoadFileErrorInvalidPath(t *testing.T) {
	t.Parallel()

	_, err := junit.LoadFile("some-invalid-path")
	if err == nil {
		t.Fatal("expected LoadFile to fail but it succeeded")
	}
}

func fixtureContents(t *testing.T, fixture string) []byte {
	contents, err := ioutil.ReadFile(fixturePath(fixture))
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func fixturePath(fixture string) string {
	return filepath.Join("fixtures", fixture)
}
