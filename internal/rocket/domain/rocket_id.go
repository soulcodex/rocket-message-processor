package rocketdomain

import (
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

var (
	ErrInvalidRocketIDProvided = errutil.NewError("invalid rocket id provided")
)

type RocketID string

func NewRocketID(id string) (RocketID, error) {
	rocketID := RocketID(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.UUIDIdentifier(),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidRocketIDProvided.Wrap(err)
	}

	return rocketID, nil
}

func (r RocketID) String() string {
	return string(r)
}
