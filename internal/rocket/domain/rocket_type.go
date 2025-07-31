package rocketdomain

import (
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

const (
	minRocketTypeLength = 5
)

var (
	ErrInvalidRocketTypeProvided = errutil.NewError("invalid rocket type provided")
)

type RocketType string

func NewRocketType(id string) (RocketType, error) {
	rocketType := RocketType(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
		domainvalidation.MinLength(minRocketTypeLength),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidRocketTypeProvided.Wrap(err)
	}

	return rocketType, nil
}

func (r RocketType) String() string {
	return string(r)
}
