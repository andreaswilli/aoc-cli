package trigger_test

import (
	"aoc-cli/trigger"
	"testing"
)

func TestOneShotTrigger_Listen(t *testing.T) {
	trigger := trigger.OneShotTrigger{}

	triggeredTimes := 0

	for range trigger.Listen() {
		triggeredTimes += 1

		if triggeredTimes > 1 {
			t.Fatalf(
				"OneShotTrigger should trigger once, triggered %d times",
				triggeredTimes,
			)
		}
	}

  if triggeredTimes == 0 {
    t.Fatalf("OneShotTrigger should trigger once, triggered 0 times")
  }
}
