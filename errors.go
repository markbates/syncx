package syncx

const (
	ErrBrokenRange = stringError("range was broken")
	ErrNilMap      = stringError("map is nil")
)

type stringError string

func (s stringError) Error() string {
	return string(s)
}
