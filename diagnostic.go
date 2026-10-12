package wsdl

// Severity classifies a validation diagnostic.
type Severity string

const (
	// SeverityError makes a retained diagnostic contribute to Diagnostics.Err.
	SeverityError Severity = "error"
	// SeverityWarning records an issue without by itself making Diagnostics.Err non-nil.
	SeverityWarning Severity = "warning"
)

// Diagnostic describes one bounded parsing, validation, or compilation issue.
type Diagnostic struct {
	// Code identifies the diagnostic or protocol fault code in the owning vocabulary.
	Code string
	// Severity controls whether this retained diagnostic contributes to HasErrors.
	Severity Severity
	// Message contains inspectable input-derived detail excluded from Diagnostics.Error.
	Message string
	// Path identifies the component/diagnostic boundary and may contain caller data.
	Path string
	// Location records decoder source-position metadata, not a trusted resource boundary.
	Location Location
}

// Diagnostics is an ordered collection of issues.
type Diagnostics []Diagnostic

// HasErrors reports whether at least one error diagnostic is present.
func (d Diagnostics) HasErrors() bool {
	for _, diagnostic := range d {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Error returns a category without input-derived diagnostic details.
// Inspect the diagnostics explicitly for messages, paths and locations.
func (d Diagnostics) Error() string {
	if len(d) == 0 {
		return ""
	}
	return "wsdl: validation failed"
}

// Err returns diagnostics as an error when any error diagnostic is present.
func (d Diagnostics) Err() error {
	if !d.HasErrors() {
		return nil
	}
	return d
}
