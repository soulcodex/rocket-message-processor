package rocketdomain

import (
	"errors"

	"github.com/soulcodex/rockets-message-processor/pkg/domain"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

const rocketStoreErrorMessage = "rocket store error occurred"

type RocketStoreError struct {
	domain.BaseError
}

func NewRocketStoreError() *RocketStoreError {
	return &RocketStoreError{
		BaseError: domain.NewError(rocketStoreErrorMessage),
	}
}

func (rse *RocketStoreError) Error() string {
	return rocketStoreErrorMessage
}

func (rse *RocketStoreError) Severity() errutil.ErrorSeverity {
	return errutil.SeverityFatal
}

func IsRocketStoreError(err error) bool {
	var self *RocketStoreError
	return errors.As(err, &self)
}
