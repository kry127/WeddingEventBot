package mongo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"kry127.ru/weddingbot/config"
)

func TestConnection(t *testing.T) {
	client, err := MongoClientFromConfig(new(config.Config))
	require.NoError(t, err)
	require.NotNil(t, client)
}
