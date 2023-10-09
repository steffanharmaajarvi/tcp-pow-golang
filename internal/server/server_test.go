package server

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/protocol"
	storage "wordofwisdom/internal/storage"
)

func TestProcessRequest(t *testing.T) {
	storageInst := getStorage()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "storage", storageInst)
	ctx = context.WithValue(ctx, "storageExpiration", int64(3600))

	clientInfo := "129.23.23.23"

	t.Parallel()

	t.Run("Unknown message format", func(t *testing.T) {
		message := "<>wqe2"
		_, err := processRequest(ctx, message, clientInfo)

		require.Error(t, err)
		assert.Equal(t, "cannot parse header: <>wqe2", err.Error())

	})

	t.Run("Unknown header", func(t *testing.T) {
		message := "9|"
		_, err := processRequest(ctx, message, clientInfo)

		require.Error(t, err)
		assert.Equal(t, "unknown header", err.Error())

	})

	t.Run("Exit case", func(t *testing.T) {
		message := "0|"
		_, err := processRequest(ctx, message, clientInfo)

		require.Error(t, err)
		assert.Equal(t, "close connection", err.Error())

	})

	t.Run("Request Challenge case", func(t *testing.T) {
		message := "1|"
		requestResult, err := processRequest(ctx, message, clientInfo)

		require.NoError(t, err)

		var hashcash pow.Hashcash
		err = json.Unmarshal([]byte(requestResult.Payload), &hashcash)
		require.NoError(t, err)

		assert.Equal(t, protocol.ResponseChallenge, requestResult.Header)
		assert.Equal(t, 0, hashcash.Counter)
		assert.Equal(t, clientInfo, hashcash.Client)
		assert.Equal(t, uint(3), hashcash.Bits)
		assert.Equal(t, "", hashcash.Extra)
		assert.NotEmpty(t, hashcash.BaseValue)

	})

	t.Run("Request Resource case", func(t *testing.T) {
		message := "3|{\"Bits\":3,\"Zeros\":1,\"SaltLen\":8,\"Counter\":1,\"Datetime\":1696816036,\"BaseValue\":\"6129484611666145821\",\"Extra\":\"\",\"Client\":\"172.20.0.3:41004\"}"

		storageInst.Add("6129484611666145821", 3600)

		requestResult, err := processRequest(ctx, message, "172.20.0.3:41004")

		require.NoError(t, err)

		assert.Equal(t, protocol.ResponseResource, requestResult.Header)
		assert.Contains(t, storage.Quotes, requestResult.Payload)
	})

	t.Run("Request Resource with no storage found", func(t *testing.T) {
		message := "3|{\"Bits\":3,\"Zeros\":1,\"SaltLen\":8,\"Counter\":1,\"Datetime\":1696816036,\"BaseValue\":\"6129484611666145821\",\"Extra\":\"\",\"Client\":\"172.20.0.3:41004\"}"

		_, err := processRequest(ctx, message, "172.20.0.3:41004")

		require.Error(t, err)

		assert.Equal(t, "challenge expired or not sent", err.Error())
	})

	t.Run("Request Resource with invalid hashcash sum", func(t *testing.T) {
		message := "3|{\"Bits\":3,\"Zeros\":1,\"SaltLen\":8,\"Counter\":1,\"Datetime\":1696816036,\"BaseValue\":\"61294841666145821\",\"Extra\":\"\",\"Client\":\"172.20.0.3:41004\"}"
		storageInst.Add("61294841666145821", 3600)

		_, err := processRequest(ctx, message, "172.20.0.3:41004")

		require.Error(t, err)

		assert.Equal(t, "invalid hashcash", err.Error())
	})

	t.Run("Request Resource with different client info", func(t *testing.T) {
		message := "3|{\"Bits\":3,\"Zeros\":1,\"SaltLen\":8,\"Counter\":1,\"Datetime\":1696816036,\"BaseValue\":\"61294841666145821\",\"Extra\":\"\",\"Client\":\"172.20.0.3:41004\"}"
		storageInst.Add("6129484611666145821", 3600)

		_, err := processRequest(ctx, message, clientInfo)

		require.Error(t, err)

		assert.Equal(t, "hashcash client mismatch", err.Error())
	})

}
