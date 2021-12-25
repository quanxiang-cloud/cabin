package error

const (
	// Unknown unknown error
	Unknown int64 = -1
	// Internal internal server error
	Internal int64 = -2

	// Success
	Success int64 = 0
	// ErrParams parameter error
	ErrParams int64 = 1
)

type Table map[int64]string

var CodeTable Table

var baseCode = map[int64]string{
	Unknown:   "unknown err",
	Internal:  "internal server error",
	Success:   "success",
	ErrParams: "parameter error",
}

// Translation translation code to message
func Translation(code int64) string {
	if CodeTable != nil {
		if text, ok := CodeTable[code]; ok {
			return text
		}
	}
	if text, ok := baseCode[code]; ok {
		return text
	}
	return "unknown code."
}
