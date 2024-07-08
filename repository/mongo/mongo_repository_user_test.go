package mongorepo

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/mymmrac/telego"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"kry127.ru/weddingbot/config"
	dbaasmongo "kry127.ru/weddingbot/dbaas/mongo"
)

func TestRememberUserAndSubscribe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// cfg, err := config.LoadConfig()
	// require.NoError(t, err)
	cfg := new(config.Config)

	client, err := dbaasmongo.MongoClientFromConfig(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)

	dbName := "testdb"
	names, err := client.Database(dbName).ListCollectionNames(ctx, bson.D{})
	require.NoError(t, err)
	require.Empty(t, names, "the test database should be empty")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		client.Database(dbName).Drop(ctx)
	}()

	repo := NewMongoRepositoryUser(dbName, client)
	user1 := telego.User{
		ID:        1,
		FirstName: "Vasya",
		LastName:  "Pupkin",
	}
	user2 := telego.User{
		ID:        2,
		FirstName: "Gena",
		LastName:  "Bukin",
	}

	repo.Remember(ctx, &user1)

	user1db, err := repo.GetRememberedUser(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, user1, user1db)

	user2dbFake, err := repo.GetRememberedUser(ctx, 2)
	require.Nil(t, user2dbFake)

	repo.Remember(ctx, &user2)

	user2db, err := repo.GetRememberedUser(ctx, 2)
	require.NoError(t, err)
	require.Equal(t, user2, user2db)

	users, err := repo.ListRememberedUser(ctx)
	require.NoError(t, err)
	require.Len(t, users, 2)
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	require.Equal(t, user1, users[0])
	require.Equal(t, user2, users[1])

	subscribed1fake, err := repo.Subscribed(ctx, 1)
	require.NoError(t, err)
	require.False(t, subscribed1fake)
	subscribed2fake, err := repo.Subscribed(ctx, 2)
	require.NoError(t, err)
	require.False(t, subscribed2fake)

	require.NoError(t, repo.Subscribe(ctx, 1))
	require.NoError(t, repo.Subscribe(ctx, 2))

	subscribed1, err := repo.Subscribed(ctx, 1)
	require.NoError(t, err)
	require.True(t, subscribed1)
	subscribed2, err := repo.Subscribed(ctx, 2)
	require.NoError(t, err)
	require.True(t, subscribed2)

	require.NoError(t, repo.Unsubscribe(ctx, 1))
	require.NoError(t, repo.Unsubscribe(ctx, 2))

	subscribed1fake4, err := repo.Subscribed(ctx, 1)
	require.NoError(t, err)
	require.False(t, subscribed1fake4)
	subscribed2fake4, err := repo.Subscribed(ctx, 2)
	require.NoError(t, err)
	require.False(t, subscribed2fake4)

}
