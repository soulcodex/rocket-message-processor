package rocketdomain

import (
	"strconv"

	"github.com/soulcodex/rockets-message-processor/pkg/domain"
	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
)

var (
	ErrInvalidRocketLaunchSpeedProvided = domain.NewError("invalid rocket launch speed provided")
)

type LaunchSpeed int64

func NewLaunchSpeed(speed int64) (LaunchSpeed, error) {
	launchSpeed := LaunchSpeed(speed)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[int64](),
	)

	if err := validation.Validate(speed); err != nil {
		return 0, ErrInvalidRocketLaunchSpeedProvided.Wrap(err)
	}

	return launchSpeed, nil
}

func (s LaunchSpeed) IsZero() bool {
	return int64(s) == 0
}

func (s LaunchSpeed) Value() int64 {
	return int64(s)
}

func (s LaunchSpeed) String() string {
	return strconv.FormatInt(int64(s), 10)
}
