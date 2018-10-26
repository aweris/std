package repository

import (
	"context"
	"encoding/json"
	"github.com/aweris/stp/internal/inventory"
	"github.com/aweris/stp/internal/models"
	"github.com/aweris/stp/storage"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
)

const (
	bucketCategory        = "inv_category"
	bucketCategoryMeta    = "_meta"
	bucketCategoryIdx     = "index"
	bucketCategoryIdxName = "idx_category_name"
)

type boltDBCategoryRepository struct {
	db *storage.BoltDB
}

func (bcr *boltDBCategoryRepository) init() error {
	return bcr.db.Update(func(tx *bolt.Tx) error {
		cb, err := tx.CreateBucketIfNotExists([]byte(bucketCategory))
		if err != nil {
			return err
		}

		mt, err := cb.CreateBucketIfNotExists([]byte(bucketCategoryMeta))
		if err != nil {
			return err
		}

		ib, err := mt.CreateBucketIfNotExists([]byte(bucketCategoryIdx))
		if err != nil {
			return err
		}

		_, err = ib.CreateBucketIfNotExists([]byte(bucketCategoryIdxName))
		if err != nil {
			return err
		}

		return nil
	})
}

// NewBoltDBCategoryRepository creates repository for bolt db
func NewBoltDBCategoryRepository(db *storage.BoltDB) inventory.CategoryRepository {
	bcr := &boltDBCategoryRepository{db}

	if err := bcr.init(); err != nil {
		log.Fatalln(err)
	}

	return bcr
}

// AddOrUpdateCategory adding or updating category and related indexes without checking existing value.
func (bcr *boltDBCategoryRepository) AddOrUpdateCategory(ctx context.Context, cat *models.Category) (*models.Category, error) {
	err := bcr.db.Update(func(tx *bolt.Tx) error {
		tb := tx.Bucket([]byte(bucketCategory))

		// put tenant to bucket
		data, err := json.Marshal(cat)
		if err != nil {
			return err
		}
		err = tb.Put(cat.Id.Bytes(), data)
		if err != nil {
			return err
		}

		// getting tenant name index bucket
		mb := tb.Bucket([]byte(bucketCategoryMeta))
		ib := mb.Bucket([]byte(bucketCategoryIdx))
		idx := ib.Bucket([]byte(bucketCategoryIdxName))

		// adding tenant name to bucket
		err = idx.Put([]byte(strings.ToLower(cat.Name)), cat.Id.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
	return cat, err
}