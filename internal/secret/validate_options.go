package secret

// ValidateReport summarises results of a ValidateMap call.
type ValidateReport struct {
	Total   int
	Invalid int
	Errors  map[string]error
}

// NewValidateReport builds a report from the error map returned by ValidateMap.
func NewValidateReport(secrets map[string]string, errs map[string]error) ValidateReport {
	r := ValidateReport{
		Total:  len(secrets),
		Errors: errs,
	}
	if errs != nil {
		r.Invalid = len(errs)
	}
	return r
}

// OK returns true when no validation errors exist.
func (r ValidateReport) OK() bool {
	return r.Invalid == 0
}

// Summary returns a short human-readable summary string.
func (r ValidateReport) Summary() string {
	if r.OK() {
		return "all secrets valid"
	}
	return formatSummary(r.Total, r.Invalid)
}

func formatSummary(total, invalid int) string {
	return fmt.Sprintf("%d/%d secret(s) failed validation", invalid, total)
}
