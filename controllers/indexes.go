package controllers

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates all required MongoDB indexes on startup.
// It is idempotent - re-running never drops existing data.
// Call this once after all Init*Collection calls in main.go.
func EnsureIndexes(client *mongo.Client, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(dbName)

	if err := ensureUniqueMetadataReady(ctx, db.Collection("election_metadata")); err != nil {
		return err
	}

	type indexSpec struct {
		collection string
		model      mongo.IndexModel
	}

	specs := []indexSpec{
		{
			collection: "voters",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("voters_email_unique"),
			},
		},
		{
			collection: "voters",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "registrations.election_address", Value: 1}},
				Options: options.Index().SetName("voters_registration_election_address"),
			},
		},
		{
			collection: "election_metadata",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "election_address", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("metadata_election_address_unique"),
			},
		},
		{
			collection: "election_metadata",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "status", Value: 1}},
				Options: options.Index().SetName("metadata_status"),
			},
		},
		{
			collection: "candidates",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "electionAddress", Value: 1}},
				Options: options.Index().SetName("candidates_election_address"),
			},
		},
		{
			collection: "candidates",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "email", Value: 1},
					{Key: "electionAddress", Value: 1},
				},
				Options: options.Index().SetName("candidates_email_election_address"),
			},
		},
		{
			collection: "otps",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetName("otps_email"),
			},
		},
		{
			collection: "otps",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "expiresAt", Value: 1}},
				Options: options.Index().SetExpireAfterSeconds(0).SetName("otps_ttl"),
			},
		},
		{
			collection: "companies",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("companies_email_unique"),
			},
		},
		{
			collection: "students",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("students_email_unique"),
			},
		},
		{
			collection: "email_jobs",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "status", Value: 1}, {Key: "createdAt", Value: 1}},
				Options: options.Index().SetName("email_jobs_status_createdAt"),
			},
		},
		{
			collection: "vote_jobs",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "status", Value: 1}, {Key: "createdAt", Value: 1}},
				Options: options.Index().SetName("vote_jobs_status_createdAt"),
			},
		},
		{
			collection: "vote_jobs",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "electionAddress", Value: 1}, {Key: "voterEmail", Value: 1}},
				Options: options.Index().SetName("vote_jobs_election_voter"),
			},
		},
	}

	created := 0
	for _, s := range specs {
		exists, err := hasEquivalentIndex(ctx, db.Collection(s.collection), s.model)
		if err != nil {
			log.Printf("[INDEX WARN] Could not check existing indexes for %s.%s: %v (skipping)", s.collection, indexName(s.model), err)
			continue
		}
		if exists {
			// index with same keys + options already present under any name — skip
			continue
		}

		if _, err := db.Collection(s.collection).Indexes().CreateOne(ctx, s.model); err != nil {
			// Code 85 IndexOptionsConflict / 86 IndexKeySpecsConflict / 68 IndexAlreadyExists
			// All mean the index already exists — safe to ignore
			log.Printf("[INDEX WARN] %s.%s: %v (index likely exists under a different name — safe to ignore)", s.collection, indexName(s.model), err)
			continue
		}
		created++
	}

	log.Printf("[OK] MongoDB indexes ensured: %d new, %d already present (total specs: %d)", created, len(specs)-created, len(specs))
	return nil
}

func indexName(m mongo.IndexModel) string {
	if m.Options != nil && m.Options.Name != nil {
		return *m.Options.Name
	}
	return "(unnamed)"
}

func ensureUniqueMetadataReady(ctx context.Context, coll *mongo.Collection) error {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$election_address"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$type", Value: "string"}}},
			{Key: "count", Value: bson.D{{Key: "$gt", Value: 1}}},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("check election_metadata duplicates: %w", err)
	}
	defer cursor.Close(ctx)

	type duplicateRow struct {
		ID    string `bson:"_id"`
		Count int32  `bson:"count"`
	}

	var duplicates []duplicateRow
	if err := cursor.All(ctx, &duplicates); err != nil {
		return fmt.Errorf("decode election_metadata duplicates: %w", err)
	}
	if len(duplicates) == 0 {
		return nil
	}

	sample := duplicates[0]
	return fmt.Errorf(
		"cannot create unique index on election_metadata.election_address: found %d duplicate address groups (example %q has %d documents); clean duplicates before starting the app",
		len(duplicates),
		sample.ID,
		sample.Count,
	)
}

func hasEquivalentIndex(ctx context.Context, coll *mongo.Collection, model mongo.IndexModel) (bool, error) {
	cursor, err := coll.Indexes().List(ctx)
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	wantKeys := normalizeIndexKeys(model.Keys)
	wantUnique := indexUnique(model.Options)
	wantTTL, wantHasTTL := indexTTL(model.Options)

	for cursor.Next(ctx) {
		var idxDoc bson.M
		if err := cursor.Decode(&idxDoc); err != nil {
			return false, err
		}

		gotKeys := normalizeIndexKeys(idxDoc["key"])
		if !reflect.DeepEqual(gotKeys, wantKeys) {
			continue
		}

		gotUnique := false
		if uniqueVal, ok := idxDoc["unique"].(bool); ok {
			gotUnique = uniqueVal
		}
		if gotUnique != wantUnique {
			continue
		}

		gotTTL, gotHasTTL := extractTTLValue(idxDoc["expireAfterSeconds"])
		if gotHasTTL != wantHasTTL {
			continue
		}
		if gotHasTTL && gotTTL != wantTTL {
			continue
		}

		return true, nil
	}

	if err := cursor.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func normalizeIndexKeys(keys interface{}) bson.D {
	switch v := keys.(type) {
	case bson.D:
		return v
	case bson.M:
		out := make(bson.D, 0, len(v))
		for key, value := range v {
			out = append(out, bson.E{Key: key, Value: value})
		}
		return out
	case map[string]interface{}:
		out := make(bson.D, 0, len(v))
		for key, value := range v {
			out = append(out, bson.E{Key: key, Value: value})
		}
		return out
	default:
		return nil
	}
}

func indexUnique(opts *options.IndexOptions) bool {
	return opts != nil && opts.Unique != nil && *opts.Unique
}

func indexTTL(opts *options.IndexOptions) (int32, bool) {
	if opts == nil || opts.ExpireAfterSeconds == nil {
		return 0, false
	}
	return *opts.ExpireAfterSeconds, true
}

func extractTTLValue(v interface{}) (int32, bool) {
	switch t := v.(type) {
	case int32:
		return t, true
	case int64:
		return int32(t), true
	case float64:
		return int32(t), true
	default:
		return 0, false
	}
}
