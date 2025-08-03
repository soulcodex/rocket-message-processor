package rocketdomain

import (
	"github.com/soulcodex/rockets-message-processor/pkg/domain"
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
)

const (
	minRocketTypeLength = 5
)

var (
	ErrInvalidRocketTypeProvided = domain.NewError("invalid rocket type provided")
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
