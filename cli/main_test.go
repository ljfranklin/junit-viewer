package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ljfranklin/junit-viewer/internal/test/helpers"
)

var (
	mainPath string
)

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "junit-viewer-main")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	mainPath = buildMain(tmpDir)

	os.Exit(m.Run())
}

func TestPassFail(t *testing.T) {
	t.Parallel()

	successPath := filepath.Join(projectRoot(), "fixtures/success.xml")
	failurePath := filepath.Join(projectRoot(), "fixtures/failures.xml")
	cmd := exec.Command(mainPath, "--output-type", "pass-fail", successPath, failurePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd failed: %s, %s", err, string(output))
	}

	helpers.AssertEquals(t, string(output), `## Summary of last 2 run(s)

| TESTS |   PASSED   |   FAILED   | SKIPPED  | TIME  |           WHEN            |
|-------|------------|------------|----------|-------|---------------------------|
|     2 | 2 (100.0%) | 0 (0.0%)   | 0 (0.0%) | 9.837 | 2018-03-15T14:22:46+07:00 |
|     2 | 0 (0.0%)   | 2 (100.0%) | 0 (0.0%) | 0.003 | 2018-03-14T10:12:34+07:00 |
`)
}

func TestFrequentFailures(t *testing.T) {
	t.Parallel()

	successPath := filepath.Join(projectRoot(), "fixtures/success.xml")
	failurePath := filepath.Join(projectRoot(), "fixtures/failures.xml")
	mixedPath := filepath.Join(projectRoot(), "fixtures/mixed-results.xml")
	cmd := exec.Command(mainPath, "--output-type", "frequent-failures", successPath, failurePath, mixedPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cmd failed: %s, %s", err, string(output))
	}

	helpers.AssertEquals(t, string(output), `## Most frequent failures in last 3 run(s)

|        TEST         |  FAILED   |        LAST FAILED        |         LAST RAN          |
|---------------------|-----------|---------------------------|---------------------------|
| TestS3Get           | 2 (66.7%) | 2018-03-14T10:12:34+07:00 | 2018-03-15T14:22:46+07:00 |
| TestS3CompatibleGet | 1 (33.3%) | 2018-03-14T10:12:34+07:00 | 2018-03-15T14:22:46+07:00 |
`)
}

func buildMain(tmpDir string) string {
	mainPath := filepath.Join(tmpDir, "viewer")
	cmd := exec.Command("go", "build", "-o", mainPath, "github.com/ljfranklin/junit-viewer/cli")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to build main.go: %s, %s", err, string(output)))
	}

	return mainPath
}

func projectRoot() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "..")
}
