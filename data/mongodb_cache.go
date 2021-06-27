package data

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/utagai/look/datum"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDBDataCache is a MongoDB-backed cache for query result set data.
// This is unfortunately necessary because certain kinds of MongoDB queries do
// not guarantee a deterministic ordering of their result documents. An example
// of this is $group.
// This cache works by taking the hex-encoded string of a pipeline, and
// creating a new collection via $out of the result of running it against the
// source collection. These resulting collections form the 'cache'.
type mongoDBDataCache struct {
	sourceDB        *mongo.Database
	sourceColl      *mongo.Collection
	prefilter       []bson.M
	queryToCollName map[string]string
}

func newMongoDBDataCache(sourceDB *mongo.Database, sourceColl *mongo.Collection, prefilter []bson.M) *mongoDBDataCache {
	return &mongoDBDataCache{
		sourceDB:        sourceDB,
		sourceColl:      sourceColl,
		prefilter:       prefilter,
		queryToCollName: make(map[string]string),
	}
}

// newQuery creates a new query to be cached.
func (m *mongoDBDataCache) newQuery(ctx context.Context, q string) error {
	pipeline := []bson.M{}
	if err := bson.UnmarshalExtJSON([]byte(q), true, &pipeline); err != nil {
		return ErrUnableToParseQuery
	}
	pipelineHex := hex.EncodeToString([]byte(q))
	outPipeline := append(pipeline, bson.M{"$out": pipelineHex})
	// $out produces an empty cursor, just make sure this did not error:
	if _, err := m.sourceColl.Aggregate(ctx, outPipeline); err != nil {
		return fmt.Errorf("failed to cache the results of the pipeline (%q): %w", q, err)
	}

	m.queryToCollName[q] = pipelineHex

	return nil
}

// runQuery runs the given query and/or returns an indexable result set with
// deterministic ordering.
func (m *mongoDBDataCache) runQuery(ctx context.Context, q string) (*indexableResult, error) {
	if collNameForQuery, ok := m.queryToCollName[q]; ok {
		return &indexableResult{
			prefilter: m.prefilter,
			coll:      m.sourceDB.Collection(collNameForQuery),
		}, nil
	} else {
		if err := m.newQuery(ctx, q); err != nil {
			return nil, fmt.Errorf("tried caching the results on fetch, but failed: %w", err)
		}

		return m.runQuery(ctx, q)
	}
}

type indexableResult struct {
	prefilter []bson.M
	coll      *mongo.Collection
}

func (ir *indexableResult) At(ctx context.Context, index int) (datum.Datum, error) {
	// $skip cannot take a negative value.
	if index < 0 {
		return nil, ErrOutOfBounds
	}
	pipeline := []bson.M{
		{
			"$skip": index,
		},
		{
			"$limit": 1,
		},
	}
	pipeline = append(ir.prefilter, pipeline...)
	// Use $natural to ensure a deterministic ordering of the result set.
	// We don't use indexes for these collections anyways, so we don't lose
	// anything here, we will always be doing a COLLSCAN regardless.
	// $natural is the default sort order regardless, but being explicit
	// doesn't hurt, I think.
	cursor, err := ir.coll.Aggregate(ctx, pipeline, options.Aggregate().SetHint(bson.M{"$natural": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to get a cursor for the cached collection: %w", err)
	}

	if !cursor.Next(ctx) {
		return nil, ErrOutOfBounds
	}

	var datum datum.Datum
	if err := cursor.Decode(&datum); err != nil {
		return nil, fmt.Errorf("failed to decode datum: %w", err)
	}

	return datum, nil
}

func (ir *indexableResult) Length(ctx context.Context) (int, error) {
	count64, err := ir.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to get the count for the result set: %w", err)
	}

	return int(count64), nil
}
