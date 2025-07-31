package rocketdomain

import (
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

var (
	ErrInvalidMissionProvided = errutil.NewError("invalid rocket mission provided")
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
