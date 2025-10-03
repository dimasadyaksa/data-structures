package heap

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNegativeCap     = Error("heap: capacity cannot be negative")
	ErrZeroCap        = Error("heap: capacity cannot be zero")
	ErrCapacityReached = Error("heap: capacity reached and cannot grow")
)