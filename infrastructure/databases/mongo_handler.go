package database

import (
	"backend-agent-demo/adapter/repository"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoHandler struct {
	db     *mongo.Database
	client *mongo.Client
}

func NewMongoHandler(cfg *config) (*mongoHandler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ctxTimeout)
	defer cancel()

	URI := fmt.Sprintf("mongodb://%s:%s@%s",
		cfg.user,
		cfg.password,
		cfg.host,
	)

	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &mongoHandler{
		db:     client.Database(cfg.database),
		client: client,
	}, nil
}

func (mg mongoHandler) Store(ctx context.Context, collection string, data interface{}) error {
	if _, err := mg.db.Collection(collection).InsertOne(ctx, data); err != nil {
		return err
	}

	return nil
}

func (mg mongoHandler) Update(ctx context.Context, collection string, query, update interface{}) error {
	if _, err := mg.db.Collection(collection).UpdateOne(ctx, query, update); err != nil {
		return err
	}

	return nil
}

func (mg mongoHandler) FindAll(ctx context.Context, collection string, query, result interface{}) error {
	cur, err := mg.db.Collection(collection).Find(ctx, query)
	if err != nil {
		return err
	}

	defer cur.Close(ctx)
	if err = cur.All(ctx, result); err != nil {
		return err
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

func (mg mongoHandler) FindOne(
	ctx context.Context,
	collection string,
	query,
	projection,
	result interface{},
) error {
	err := mg.db.Collection(collection).FindOne(ctx, query, options.FindOne().SetProjection(projection)).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

// mongo session transaction
func (mg *mongoHandler) StartSession() (repository.Session, error) {
	session, err := mg.client.StartSession()
	if err != nil {
		log.Fatal(err)
	}

	return newMongoHandlerSession(session), nil
}

type mongoDBSession struct {
	session mongo.Session
}

func newMongoHandlerSession(session mongo.Session) *mongoDBSession {
	return &mongoDBSession{session: session}
}

func (ms *mongoDBSession) WithTransaction(ctx context.Context, fnc func(context.Context) error) error {
	cb := func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := fnc(sessCtx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err := ms.session.WithTransaction(ctx, cb)
	if err != nil {
		return err
	}

	return nil
}

func (ms *mongoDBSession) EndSession(ctx context.Context) {
	ms.session.EndSession(ctx)
}
