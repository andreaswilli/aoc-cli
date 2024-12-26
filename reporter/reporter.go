package reporter

type Status string

const (
	StatusPassed Status = "PASSED"
	StatusFailed Status = "FAILED"
	StatusNoExp  Status = "NO EXP"
)

func Report(actual string, expected string) (status Status) {
	if len(expected) == 0 {
		return StatusNoExp
	}
	if actual == expected {
		return StatusPassed
	}
	return StatusFailed
}
