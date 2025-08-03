package rocketdomain

import (
	"iter"
)

type RocketCollection struct {
	rockets []*Rocket
}

func NewRocketCollection(rockets ...*Rocket) RocketCollection {
	return RocketCollection{rockets: rockets}
}

func (rc RocketCollection) All() []*Rocket {
	return rc.rockets
}

func (rc RocketCollection) Primitives() iter.Seq[RocketPrimitives] {
	return func(yield func(RocketPrimitives) bool) {
		for _, c := range rc.rockets {
			if !yield(c.Primitives()) {
				return
			}
		}
	}
}
