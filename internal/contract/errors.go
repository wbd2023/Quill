package contract

import "errors"

func ViolationsFound() (err error) {
	return errors.New("violations found")
}
