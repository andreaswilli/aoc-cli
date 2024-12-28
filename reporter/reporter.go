package reporter

type Status string

const (
	StatusPassed Status = "PASSED"
	StatusFailed Status = "FAILED"
	StatusNoExp  Status = "NO EXP"
)

type Report struct {
	Actual   string
	Expected string
	Status   Status
}

func GetReport(actual string, expected string) (report Report) {
	report = Report{actual, expected, StatusFailed}

	if len(expected) == 0 {
		report.Status = StatusNoExp
	} else if actual == expected {
		report.Status = StatusPassed
	}
	return
}
