package compareoffers

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidEscalator is returned in case of incorrectly defined escalator.
	ErrInvalidEscalator = fmt.Errorf("invalid escalator")

	// ErrInvalidAdvance is returned in case of incorrectly defined advance.
	ErrInvalidAdvance = fmt.Errorf("invalid advance")

	// ErrInvalidRate is returned in case of incorrectly defined royalty rate.
	ErrInvalidRate = fmt.Errorf("invalid royalty rate")

	// ErrInvalidCopies is returned in case of incorrectly defined number of copies.
	ErrInvalidCopies = fmt.Errorf("invalid number of copies")
)

// Offer represents an offer comprised of an advance and a royalty progression.
type Offer struct {
	Advance   float64
	Escalator []Step
}

// Step represents a step in a royalty progression.
type Step struct {
	Rate   float64
	Copies int
}

// CalcEarnings calculates royalty earnings for a given price, a number of copies.
func (o *Offer) CalcEarnings(p float64, c int) float64 {
	var e float64
	for _, step := range o.Escalator {
		// the last step is a step where either:
		// - copies left is less than the next threshold, or
		// - the threshold is 0
		if c < step.Copies || step.Copies == 0 {
			e += float64(c) * p * step.Rate / 100
			break
		}
		e += float64(step.Copies) * p * step.Rate / 100
		c -= step.Copies
	}
	// if the advance has not earned out earnings = advance
	if e < o.Advance {
		return o.Advance
	}
	return e
}

// NewOfferFromString parses a string representation and returns a new offer.
func NewOfferFromString(raw string) (*Offer, error) {
	ae := strings.Split(raw, "__")

	adv, err := strconv.ParseFloat(ae[0], 64)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", ae[0], ErrInvalidAdvance)
	}

	rawRoys := strings.Split(ae[1], "_")
	offer := &Offer{
		Advance:   adv,
		Escalator: make([]Step, len(rawRoys)),
	}

	var lastUntil int = 0
	var lastRate float64 = -1
	for i, r := range rawRoys {
		rc := strings.Split(r, "-")

		// parse the rate
		rate, err := strconv.ParseFloat(rc[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", rc[0], ErrInvalidRate)
		}
		if rate < lastRate {
			return nil, fmt.Errorf(
				"rate %.2f is lower than the previous rate %.2f: %w",
				rate,
				lastRate,
				ErrInvalidEscalator,
			)
		}
		step := Step{Rate: rate}
		lastRate = rate

		// if this is the last step we don't care about copies
		if i == len(rawRoys)-1 {
			offer.Escalator[i] = step
			break
		}

		// parse step copies
		until, err := strconv.Atoi(rc[1])
		if err != nil {
			return nil, fmt.Errorf(
				"can't convert %q to integer: %w",
				rc[1],
				ErrInvalidCopies,
			)
		}
		if lastUntil >= until {
			return nil, fmt.Errorf(
				"previous threshold %d is not lower than new threshold %d: %w",
				lastUntil,
				until,
				ErrInvalidEscalator,
			)
		}
		step.Copies = until - lastUntil
		lastUntil = until

		// add step
		offer.Escalator[i] = step
	}
	return offer, nil
}
