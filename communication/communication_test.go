package communication_test

import (
	"fmt"
	"testing"
)

// Test helper

func LogTestError(have, want interface{}, t *testing.T, extraInfo ...string) {
	msg := fmt.Sprintf("%s: Test %s failed\n\tHave: %+v\n\tWant: %+v\n", "core_test", t.Name(), have, want)
	if len(extraInfo) > 0 {
		msg += "\n"
		msg += fmt.Sprintln(extraInfo)
	}
	t.Errorf(msg)
}
