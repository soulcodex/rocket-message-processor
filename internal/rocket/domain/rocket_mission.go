package rocketdomain

import (
	"github.com/soulcodex/rockets-message-processor/pkg/domain"
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
)

var (
	ErrInvalidMissionProvided = domain.NewError("invalid rocket mission provided")
)

type Mission string

func NewMission(id string) (Mission, error) {
	rocketID := Mission(id)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[string](),
	)

	if err := validation.Validate(id); err != nil {
		return "", ErrInvalidMissionProvided.Wrap(err)
	}

	return rocketID, nil
}

func (m Mission) String() string {
	return string(m)
}
