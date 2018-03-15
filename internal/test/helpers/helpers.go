package helpers

import (
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected '%v' to deep equal '%v'", actual, expected)
	}
}
