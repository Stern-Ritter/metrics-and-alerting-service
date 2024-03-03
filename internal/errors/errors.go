package errors

type InvalidMetricName struct {
	message string
	err     error
}

func (e InvalidMetricName) Error() string {
	return e.message
}

func (e InvalidMetricName) Unwrap() error {
	return e.err
}

func NewInvalidMetricName(message string, err error) error {
	return InvalidMetricName{message: message, err: err}
}

type InvalidMetricValue struct {
	message string
	err     error
}

func (e InvalidMetricValue) Error() string {
	return e.message
}

func (e InvalidMetricValue) Unwrap() error {
	return e.err
}

func NewInvalidMetricValue(message string, err error) error {
	return InvalidMetricValue{message: message, err: err}
}
