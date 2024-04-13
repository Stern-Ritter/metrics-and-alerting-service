package errors

type InvalidMetricType struct {
	message string
	err     error
}

func (e InvalidMetricType) Error() string {
	return e.message
}

func (e InvalidMetricType) Unwrap() error {
	return e.err
}

func NewInvalidMetricType(message string, err error) error {
	return InvalidMetricType{message: message, err: err}
}

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

type UnsuccessRequestProccessing struct {
	message string
	err     error
}

func (e UnsuccessRequestProccessing) Error() string {
	return e.message
}

func (e UnsuccessRequestProccessing) Unwrap() error {
	return e.err
}

func NewUnsuccessRequestProccessing(message string, err error) error {
	return UnsuccessRequestProccessing{message: message, err: err}
}

type FileUnavailable struct {
	message string
	err     error
}

func (e FileUnavailable) Error() string {
	return e.message
}

func (e FileUnavailable) Unwrap() error {
	return e.err
}

func NewFileUnavailable(message string, err error) error {
	return FileUnavailable{message: message, err: err}
}
