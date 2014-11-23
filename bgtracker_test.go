package bgtracker_test

import (
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ejamesc/bgtracker"
)

func TestGetBGMember(t *testing.T) {
	username := "nattsw"
	bgm, err := bgtracker.GetBGMember(username)

	ok(t, err)

	assert(t, bgm.GithubID == username, "expected bgm.GithubID to be %s, got %s", username, bgm.GithubID)
	assert(t, bgm.Name == "", "expected bgm.Name to be empty, got %s", bgm.Name)
}

func TestNewTracker_FromAPI(t *testing.T) {
	tr, err := bgtracker.NewTracker("basement-gang")

	ok(t, err)
	equals(t, tr.Orgname, "basement-gang")

	// Verify saved to DB
	db, _ := bolt.Open("bgtracker.db", 0600, nil)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("trackerinfo"))
		assert(t, b != nil, "expect members bucket to have been created, but was not")

		o := b.Get([]byte("Orgname"))
		equals(t, o, []byte("basement-gang"))

		lt := b.Get([]byte("LastUpdated"))
		equals(t, string(lt), tr.LastUpdated.Format(time.RFC3339))

		return nil
	})
}
