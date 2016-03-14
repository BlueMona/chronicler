package chronicler

type StateReporter interface {
	ReportState() (stateCode, stateDescription string)
}
