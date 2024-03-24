package utils

type UIErrors map[string]string

func NewUIError() UIErrors {
	return make(map[string]string)
}

func (errs UIErrors) Add(field string, value string) {
	errs[field] = value
}
