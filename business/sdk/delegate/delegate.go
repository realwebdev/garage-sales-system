// Package delegate provides the ability to make functions
// calls between different domain packages when an import
// is not possible
package delegate

import "github.com/realwebdev/garage-sales-system/foundation/logger"

// These types are just for documentation so we know
// what keys go where in the map.
type (
	domain string
	action string
)

// Delegate manages the set of functions to be called
// domain packages when an import is not possible.
type Delegate struct {
	log   *logger.Logger
	funcs map[domain]map[action][]Func
}

// New constructs a delegate for indirect api access.
func New(log *logger.Logger) *Delegate {
	return &Delegate{
		log:   log,
		funcs: make(map[domain]map[action][]Func),
	}
}
