// Package storage provides access to the search storage backend.
//
// Blob storage is powered by Google Cloud Storage (GCS). All indexed
// repositories are stored in GCS.
//
// Google Cloud Datastore provides an eventually consistent index of
// repositories and their files.
package storage

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
	gcs "google.golang.org/cloud/storage"
)

// Client provides an interface to the storage backend.
type Client struct {
	ctx       context.Context
	bucket    string
	datastore *datastore.Client
}

// NewClient creates a new Client for the given project, using bucket for
// GCS storage.
func NewClient(ctx context.Context, projectID string, bucket string) (*Client, error) {
	d, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:       ctx,
		bucket:    bucket,
		datastore: d,
	}, nil
}

// PutRepository uploads the described repository, and all of its files in dir.
func (c *Client) PutRepository(r *Repository, dir string) error {
	// Install the repository, so we can get its key.
	if err := r.Put(c); err != nil {
		return err
	}

	dir = filepath.Clean(dir)

	// Upload all files to GCS and install their datastore entries.
	return filepath.Walk(dir, func(localPath string, info os.FileInfo, err error) error {
		log.Printf("Processing %q", localPath)
		defer log.Printf("Finished %q", localPath)

		if err != nil {
			// We don't even try to deal with errors.
			return err
		}

		if info.IsDir() {
			// We only need to upload real files.
			return nil
		}

		src, err := os.Open(localPath)
		if err != nil {
			return err
		}

		// We upload the file under the r.Root path.
		rel, err := filepath.Rel(dir, localPath)
		if err != nil {
			return err
		}
		name := path.Join(r.Root, rel)

		dst := gcs.NewWriter(c.ctx, c.bucket, name)

		if _, err := io.Copy(dst, src); err != nil {
			// CloseWithError always returns nil.
			dst.CloseWithError(err)
			return err
		}

		if err := dst.Close(); err != nil {
			return err
		}

		// Describe the file in datastore.
		file := &File{
			Repository: r.key,
			Bucket:     c.bucket,
			Path:       name,
		}
		fileKey := datastore.NewIncompleteKey(c.ctx, fileKind, nil)
		if _, err := c.datastore.Put(c.ctx, fileKey, file); err != nil {
			return err
		}

		return nil
	})
}
