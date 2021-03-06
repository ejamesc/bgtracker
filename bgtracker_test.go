package bgtracker_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ejamesc/bgtracker"
)

// Test creating a new tracker when there's nothing
// in the DB.
func TestNewTracker_FromAPI(t *testing.T) {
	testDBName := genDBName()
	tr, err := bgtracker.NewTracker("basement-gang", testDBName)
	defer os.Remove(testDBName)

	ok(t, err)
	equals(t, tr.Orgname, "basement-gang")

	// Verify saved to DB
	db, _ := bolt.Open(testDBName, 0600, nil)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("trackerinfo"))
		assert(t, b != nil, "expect trackerinfo bucket to have been created, but was not")

		o := b.Get([]byte("Orgname"))
		equals(t, o, []byte("basement-gang"))

		lt := b.Get([]byte("LastUpdated"))
		equals(t, string(lt), tr.LastUpdated.Format(time.RFC3339))

		mb := tx.Bucket([]byte("members"))
		assert(t, mb != nil, "expect members bucket to have been created, but was not")

		c := mb.Cursor()
		members := []*bgtracker.BGMember{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			curr, err := bgtracker.BGMemberFromJSON(v)
			ok(t, err)
			members = append(members, curr)
		}

		equals(t, len(members), 9)

		return nil
	})
}

// Test creating a new tracker when there's some stuff
// already in the DB.
func TestNewTracker_FromDB(t *testing.T) {
	testDBName := genDBName()
	db, _ := bolt.Open(testDBName, 0600, nil)
	defer os.Remove(testDBName)

	tmpTime := time.Now().Format(time.RFC3339)
	db.Update(func(tx *bolt.Tx) error {
		tb, _ := tx.CreateBucketIfNotExists([]byte("trackerinfo"))
		tb.Put([]byte("Orgname"), []byte("basement-gang"))
		tb.Put([]byte("LastUpdated"), []byte(tmpTime))

		mb, _ := tx.CreateBucketIfNotExists([]byte("members"))
		mb.Put([]byte("nattsw"), bgmFixt.ToJSON())

		return nil
	})
	db.Close()

	tr, err := bgtracker.NewTracker("basement-gang", testDBName)

	ok(t, err)
	equals(t, tr.Orgname, "basement-gang")
	equals(t, 1, len(tr.Members))
	equals(t, tr.LastUpdated.Format(time.RFC3339), tmpTime)
}

// Test ability to get BGMember from API
func TestGetBGMember(t *testing.T) {
	username := "nattsw"
	bgm, err := bgtracker.GetBGMember(username)

	ok(t, err)

	assert(t, bgm.GithubID == username, "expected bgm.GithubID to be %s, got %s", username, bgm.GithubID)
	assert(t, bgm.Name == "", "expected bgm.Name to be empty, got %s", bgm.Name)
}

// BGMember fixture
var bgmFixt = &bgtracker.BGMember{
	GithubID:   "nattsw",
	Name:       "Natalie Tay",
	NoCommits:  0,
	StreakDays: 0,
}

// Test conversion of BGMember into JSON
func TestBGMember_ToJSON(t *testing.T) {
	js := bgmFixt.ToJSON()
	expectedJson := []byte(`{"GithubID":"nattsw","Name":"Natalie Tay","NoCommits":0,"StreakDays":0}`)
	equals(t, js, expectedJson)
}

// Test creating BGMember from JSON
func TestBGMemberFromJson(t *testing.T) {
	js := []byte(`{"GithubID":"nattsw","Name":"Natalie Tay","NoCommits":0,"StreakDays":0}`)
	bgm, err := bgtracker.BGMemberFromJSON(js)

	ok(t, err)
	equals(t, bgmFixt, bgm)
}

func genDBName() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("testdb_%d", rand.Int())
}
