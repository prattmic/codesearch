package storage

import (
	"google.golang.org/cloud/datastore"
)

// Datastore object Kinds.

// NOTE: due to a lack of consistency guarantees between GCS and Datastore,
// there may be differences between the Datastore objects and what is
// actually found in GCS.

const (
	repositoryKind = "Repository"
	fileKind       = "File"
)

// Repository describes a cloned repository.
type Repository struct {
	// URL is the source URL of this repository.
	URL string

	// VCS is the version control system of the repository.
	// Currently, it must be "git".
	VCS string

	// Head is the current HEAD commit ref of the stored files.
	Head string

	// Bucket is the name of the GCS bucket this repository is stored in.
	Bucket string

	// Root is the path to the root of this repository in the GCS bucket.
	Root string

	// The key of this repository in Datastore. Not valid until Put is
	// called.
	key *datastore.Key
}

// Put installs the Repository in the Client's datastore.
func (r *Repository) Put(c *Client) error {
	key := datastore.NewIncompleteKey(c.ctx, repositoryKind, nil)

	key, err := c.datastore.Put(c.ctx, key, r)
	if err != nil {
		return err
	}

	r.key = key

	return nil
}

// File describes a single static file.
type File struct {
	// Repository points to the repository this file originated from.
	Repository *datastore.Key

	// Bucket is the name of the GCS bucket this file is stored in.
	Bucket string

	// Path is the path this file in the GCS bucket.
	Path string
}
