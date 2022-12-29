package data

import (
	"context"
	"fmt"

	"github.com/utagai/look/datum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	sortStage = bson.M{
		"$sort": bson.M{
			"_id": 1,
		},
	}
)

type MongoDBData struct {
	client    *mongo.Client
	resultSet *indexableResult
	cache     *mongoDBDataCache
}

var _ Data = (*MongoDBData)(nil)

// NewMongoDBData takes a URI and datums and provides a MongoDB-backed data.
// Note that as per the official MongoDB driver behavior, the database &
// collection specified in the URI are not respected, and therefore this
// constructor simply uses the namespace {dbName}.{collName}.
func NewMongoDBData(uri, dbName, collName string, datums []datum.Datum) (*MongoDBData, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	db := client.Database(dbName)
	coll := db.Collection(collName)
	if err := loadDataIntoMongoDB(ctx, coll, datums); err != nil {
		return nil, fmt.Errorf("failed to insert the datums into MongoDB: %w", err)
	}
	cache := newMongoDBDataCache(db, coll, []bson.M{sortStage})
	// The initial query is just the source collection.
	resultSet, err := cache.runQuery(ctx, "[]")
	if err != nil {
		return nil, fmt.Errorf("failed to acquire the result set: %w", err)
	}

	return &MongoDBData{
		client:    client,
		resultSet: resultSet,
		cache:     cache,
	}, nil
}

func loadDataIntoMongoDB(ctx context.Context, coll *mongo.Collection, datums []datum.Datum) error {
	if err := coll.Drop(ctx); err != nil {
		return fmt.Errorf("failed to initially drop collection: %w", err)
	}

	datumInterfaces := make([]interface{}, len(datums))
	for i, datum := range datums {
		datumInterfaces[i] = datum
	}
	_, err := coll.InsertMany(ctx, datumInterfaces)
	if err != nil {
		return err
	}

	return nil
}

func (md *MongoDBData) At(ctx context.Context, index int) (datum.Datum, error) {
	if datum, err := md.resultSet.At(ctx, index); err != nil {
		return nil, fmt.Errorf("failed to retrieve datum from the cache: %v", err)
	} else {
		return datum, nil
	}
}

func (md *MongoDBData) Find(ctx context.Context, q string) (Data, error) {
	// Consider the empty string to be equivalent to the query "[]" for
	// convenience.
	if q == "" {
		q = "[]"
	}
	// Be cheeky and don't actually run any queries. Just return a Data that will
	// return the requested data lazily.
	// This is helpful because it lets us avoid spawning excessive goroutines
	// for partially typed yet valid queries.
	resultSet, err := md.cache.runQuery(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get the result set for the new query: %w", err)
	}
	return &MongoDBData{
		client:    md.client,
		resultSet: resultSet,
		cache:     md.cache,
	}, nil
}

func (md *MongoDBData) Length(ctx context.Context) (int, error) {
	return md.resultSet.Length(ctx)
}
