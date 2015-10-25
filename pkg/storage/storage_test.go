package storage

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	gcs "google.golang.org/cloud/storage"
)

var (
	keyFile   = flag.String("key", os.Getenv("GCLOUD_KEY"), "Google Cloud JSON key file")
	projectID = flag.String("project", os.Getenv("GCLOUD_PROJECT_ID"), "Google Cloud Project ID")
	bucket    = flag.String("bucket", os.Getenv("GCLOUD_BUCKET"), "GCS test bucket")
)

func testContext() (context.Context, error) {
	if *keyFile == "" {
		return nil, errors.New("-key must be set")
	}

	if *projectID == "" {
		return nil, errors.New("-project must be set")
	}

	key, err := ioutil.ReadFile(*keyFile)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(key, gcs.ScopeFullControl, datastore.ScopeDatastore)
	if err != nil {
		return nil, err
	}

	return cloud.NewContext(*projectID, conf.Client(oauth2.NoContext)), nil
}

// Can we even create a client?
func TestClient(t *testing.T) {
	ctx, err := testContext()
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	if _, err := NewClient(ctx, *projectID, *bucket); err != nil {
		t.Errorf("Failed to create client: %v", err)
	}
}

func TestPutRepository(t *testing.T) {
	ctx, err := testContext()
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	c, err := NewClient(ctx, *projectID, *bucket)
	if err != nil {
		t.Errorf("Failed to create client: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	// We assume this executes within pkg/storage.
	repoPath := path.Join(cwd, "../..")

	r := &Repository{
		URL:    "github.com/prattmic/codesearch",
		VCS:    "git",
		Head:   "testHEAD",
		Bucket: *bucket,
		Root:   "github.com/prattmic/codesearch",
	}

	if err := c.PutRepository(r, repoPath); err != nil {
		t.Errorf("PutRepository failed: %v", err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
