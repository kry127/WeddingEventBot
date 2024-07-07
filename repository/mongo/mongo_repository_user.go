package mongorepo

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/bson"
	"kry127.ru/weddingbot/repository"
)

const (
	MongoRepositoryUserDB                  = "db"
	MongoRepositoryUserRememberCollection  = "remember"
	MongoRepositoryUserSubscribeCollection = "subscribe"
)

type MongoID struct {
	ID any `bson:"_id"`
}

func MakeMongoID(id any) MongoID {
	return MongoID{ID: id}
}

type MongoUser struct {
	ID   int64 `bson:"_id"`
	User *telego.User
}

func MakeMongoUser(user *telego.User) MongoUser {
	return MongoUser{
		ID:   user.ID,
		User: user,
	}
}

type MongoRepositoryUser struct {
	mongoClient *mongo.Client
}

func NewMongoRepositoryUser(mongoClient *mongo.Client) *MongoRepositoryUser {
	return &MongoRepositoryUser{
		mongoClient: mongoClient,
	}
}

func (r *MongoRepositoryUser) rememberCollection() *mongo.Collection {
	return r.mongoClient.Database(MongoRepositoryUserDB).Collection(MongoRepositoryUserRememberCollection)
}

func (r *MongoRepositoryUser) subscribeCollection() *mongo.Collection {
	return r.mongoClient.Database(MongoRepositoryUserDB).Collection(MongoRepositoryUserSubscribeCollection)
}

func (r *MongoRepositoryUser) Remember(ctx context.Context, user *telego.User) error {
	_, err := r.rememberCollection().InsertOne(ctx, MakeMongoUser(user))
	return err
}
func (r *MongoRepositoryUser) GetRememberedUser(ctx context.Context, userID int64) (*telego.User, error) {
	singleResult := r.rememberCollection().FindOne(ctx, MakeMongoID(userID))
	var user MongoUser
	if err := singleResult.Decode(&user); err != nil {
		return nil, fmt.Errorf("cannot decode single result: %w", err)
	}
	return user.User, nil
}
func (r *MongoRepositoryUser) ListRememberedUser(ctx context.Context) ([]*telego.User, error) {
	cursor, err := r.rememberCollection().Find(ctx, []bson.D{})
	if err != nil {
		return nil, fmt.Errorf("cannot open mongo cursor: %w", err)
	}
	defer cursor.Close(ctx)

	var result []*telego.User
	for cursor.Next(ctx) {
		var user MongoUser
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("cannot decode user: %w", err)
		}
		result = append(result, user.User)
	}
	if cursor.Err() != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return result, nil
}

func (r *MongoRepositoryUser) Subscribe(ctx context.Context, userID int64) error {
	_, err := r.rememberCollection().InsertOne(ctx, MakeMongoID(userID))
	if err != nil {
		return fmt.Errorf("cannot insert id for subscription: %w", err)
	}
	return nil
}
func (r *MongoRepositoryUser) Subscribed(ctx context.Context, userID int64) (bool, error) {
	findOneRes := r.rememberCollection().FindOne(ctx, MakeMongoID(userID))
	if findOneRes.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if findOneRes.Err() != nil {
		return false, fmt.Errorf("error finding id in subscription: %w", findOneRes.Err())
	}
	return true, nil
}
func (r *MongoRepositoryUser) Unsubscribe(ctx context.Context, userID int64) error {
	_, err := r.rememberCollection().DeleteOne(ctx, MakeMongoID(userID))
	if err != nil {
		return fmt.Errorf("cannot delete id for subscription: %w", err)
	}
	return nil
}

func (r *MongoRepositoryUser) ScheduleMessage(ctx context.Context, messageSchedule repository.ScheduledMessage) error {
	panic("implement me")
}
func (r *MongoRepositoryUser) ListScheduledMessages(ctx context.Context) ([]repository.MessageDelivery, error) {
	panic("implement me")
}
