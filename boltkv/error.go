package boltkv

import "errors"

var (
	TargetKeyNotFoundError = errors.New("target key not found")
	TargetKeyExpiredError   = errors.New("target key expired")
)
