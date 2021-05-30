package errorx

import "strings"

type ErrorList []error

func (el ErrorList) Error() string {
	if len(el) == 0 {
		return ""
	}
	var r strings.Builder
	r.WriteString("Errors:\n")
	for _, err := range el {
		r.WriteString(err.Error() + "\n")
	}
	return r.String()
}

func (el ErrorList) Append(errs ...error) ErrorList {
	var errList ErrorList
	if len(el) > 0 {
		errList = el
	}
	for _, err := range errs {
		if err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) == 0 {
		return nil
	}
	return errList
}
