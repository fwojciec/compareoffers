package compareoffers_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fwojciec/compareoffers"
)

func TestNewOfferFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in  string
		out *compareoffers.Offer
		e   error
	}{
		{
			in: "1500.50__7.25-5000_8.50",
			out: &compareoffers.Offer{
				Advance: 1500.50,
				Escalator: []compareoffers.Step{
					{7.25, 5000},
					{8.50, 0},
				},
			},
			e: nil,
		},
		{
			in: "1500__8",
			out: &compareoffers.Offer{
				Advance: 1500,
				Escalator: []compareoffers.Step{
					{8, 0},
				},
			},
			e: nil,
		},
		{
			in: "2500__8-5000_9-10000_10",
			out: &compareoffers.Offer{
				Advance: 2500,
				Escalator: []compareoffers.Step{
					{8, 5000},
					{9, 5000},
					{10, 0},
				},
			},
			e: nil,
		},
		{
			in:  "wrong__8",
			out: nil,
			e:   compareoffers.ErrInvalidAdvance,
		},
		{
			in:  "1500__wrong",
			out: nil,
			e:   compareoffers.ErrInvalidRate,
		},
		{
			in:  "1500__8-5000_7",
			out: nil,
			e:   compareoffers.ErrInvalidEscalator,
		},
		{
			in:  "1500__7-5000_8-4000_9",
			out: nil,
			e:   compareoffers.ErrInvalidEscalator,
		},
		{
			in:  "1500__7-wrong_8",
			out: nil,
			e:   compareoffers.ErrInvalidCopies,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			res, err := compareoffers.NewOfferFromString(tc.in)
			if !errors.Is(err, tc.e) {
				t.Errorf("expected %v error value, but received %v", tc.e, err)
			}
			if !(reflect.DeepEqual(tc.out, res)) {
				t.Errorf("expected the offer to be %v, but received %v", tc.out, res)
			}
		})
	}
}
