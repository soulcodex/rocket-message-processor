package rocketdomain

import (
	"errors"

	"github.com/soulcodex/rockets-message-processor/pkg/domain"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

const rocketNotExistsErrorMessage = "rocket doesn't exists."

type RocketNotFoundError struct {
	domain.BaseError
}

func NewRocketNotFoundError(id RocketID) *RocketNotFoundError {
	return &RocketNotFoundError{
		BaseError: domain.NewError(
			rocketNotExistsErrorMessage,
			errutil.WithMetadataKeyValue("rocket.id", id.String()),
		),
	}
}

func (rnf *RocketNotFoundError) Error() string {
	return rocketNotExistsErrorMessage
}

func IsRocketNotFoundError(err error) bool {
	var self *RocketNotFoundError
	return errors.As(err, &self)
}
