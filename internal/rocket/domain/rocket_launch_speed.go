package rocketdomain

import (
	"strconv"

	domainvalidation "github.com/soulcodex/rockets-message-processor/pkg/domain/validation"
)

const (
	minRocketLaunchSpeed = 1
)

var (
	ErrInvalidRocketLaunchSpeedProvided = domainvalidation.NewError("invalid rocket launch speed provided")
)

type LaunchSpeed uint64

func NewLaunchSpeed(speed uint64) (LaunchSpeed, error) {
	launchSpeed := LaunchSpeed(speed)

	validation := domainvalidation.NewValidator(
		domainvalidation.NotEmpty[uint64](),
		domainvalidation.Min[uint64](minRocketLaunchSpeed),
	)

	if err := validation.Validate(speed); err != nil {
		return 0, ErrInvalidRocketLaunchSpeedProvided.Wrap(err)
	}

	return launchSpeed, nil
}

func (s LaunchSpeed) Value() uint64 {
	return uint64(s)
}

func (s LaunchSpeed) String() string {
	return strconv.FormatUint(uint64(s), 10)
}
