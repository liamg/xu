package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneForByte(t *testing.T) {

	cases := []struct {
		b       byte
		literal bool
		r       rune
	}{
		{b: 97, literal: true, r: 'a'},
		{b: 127, literal: false, r: '.'},
		{b: 0, literal: false, r: '.'},
		{b: 8, literal: false, r: '.'},
	}

	for _, testCase := range cases {
		t.Run(
			fmt.Sprintf("Test %d -> %c [%t]", testCase.b, testCase.r, testCase.literal),
			func(t *testing.T) {
				r, literal := getRuneForByte(testCase.b)
				assert.Equal(t, r, testCase.r)
				assert.Equal(t, literal, testCase.literal)
			},
		)
	}
}
