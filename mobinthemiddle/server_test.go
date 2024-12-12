package mobinthemiddle_test

import (
	"protohackers/mobinthemiddle"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRewriteAddresses(t *testing.T) {

	t.Run("invalid address", func(t *testing.T) {
		input := "hi bob, please send money 7F1u3wSD5RbOHQmupo9nx4TnhQ"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "hi bob, please send money 7F1u3wSD5RbOHQmupo9nx4TnhQ", out)
	})

	t.Run("valid case 2", func(t *testing.T) {
		input := "hi bob, please send money 7YWHMfk9JZe0LM0g1ZauHuiSxhI"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "hi bob, please send money 7YWHMfk9JZe0LM0g1ZauHuiSxhI", out)
	})

	t.Run("valid case 2", func(t *testing.T) {
		input := "hi bob, please send money 7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "hi bob, please send money 7YWHMfk9JZe0LM0g1ZauHuiSxhI", out)
	})

	t.Run("valid case 3", func(t *testing.T) {
		input := "hi bob, please send money 7LOrwbDlS8NujgjddyogWgIM93MV5N2VR"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "hi bob, please send money 7YWHMfk9JZe0LM0g1ZauHuiSxhI", out)
	})

	t.Run("valid case 4", func(t *testing.T) {
		input := "hi bob, please send money 7adNeSwJkMakpEcln9HEtthSRtxdmEHOT8T"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "hi bob, please send money 7YWHMfk9JZe0LM0g1ZauHuiSxhI", out)
	})

	t.Run("valid case 5", func(t *testing.T) {
		input := "[MadHunter204] Please send the payment of 750 Boguscoins to 7GG2YdcmPzcHM7jS9Ha8VjzcB71m4fcF"
		out := mobinthemiddle.RewriteAddresses(input)

		assert.Equal(t, "[MadHunter204] Please send the payment of 750 Boguscoins to 7YWHMfk9JZe0LM0g1ZauHuiSxhI", out)
	})
}
