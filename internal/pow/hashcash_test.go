package pow

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashFromString(t *testing.T) {
	hashResult := hashFromString("aura")
	hashResultHex := hex.EncodeToString(hashResult)
	assert.Equal(t, "2c2a89c1c45d2977cfce538e475902f5520f4382", hashResultHex)
}

func TestBase6464EncodeBytes(t *testing.T) {
	encodeResult := base64EncodeBytes([]byte("aura"))
	assert.Equal(t, "YXVyYQ==", encodeResult)
}

func TestBase6464EncodeInt(t *testing.T) {
	encodeResult := base64EncodeInt(43389)
	assert.Equal(t, "NDMzODk=", encodeResult)
}

func TestComputeHashCash(t *testing.T) {
	t.Parallel()

	t.Run("4 zeros", func(t *testing.T) {
		hashcash := New("aura", "client", 4, 8, "")

		hashcashResult, err := hashcash.ComputeHashcash(-1)
		require.NoError(t, err)

		assert.Equal(t, true, hashcashResult.Check())
	})

	t.Run("5 zeros", func(t *testing.T) {
		hashcash := New("aura", "client", 5, 8, "")

		hashcashResult, err := hashcash.ComputeHashcash(-1)
		require.NoError(t, err)

		assert.Equal(t, true, hashcashResult.Check())
	})

	t.Run("impossible case", func(t *testing.T) {
		hashcash := New("aura", "client", 10, 8, "")

		hashcashResult, err := hashcash.ComputeHashcash(1)
		require.Error(t, err)

		assert.Equal(t, "maximum iterations exceeded", err.Error())
		assert.Equal(t, 2, hashcashResult.Counter)
	})
}
