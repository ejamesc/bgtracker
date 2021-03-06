package bgtracker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

var membersBucketName = []byte("members")
var trackerInfoBucketName = []byte("trackerinfo")

// A Tracker is a singleton struct that is responsible for
// retrieving and storing API data every hour.
type Tracker struct {
	Orgname     string
	Members     []*BGMember
	LastUpdated time.Time
}

// Generator for new Tracker
func NewTracker(orgname, dbname string) (*Tracker, error) {
	tr := &Tracker{
		Orgname: orgname,
	}

	db, err := bolt.Open(dbname, 0600, nil)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		infoB := tx.Bucket(trackerInfoBucketName)
		membersB := tx.Bucket(membersBucketName)

		var err error
		if infoB == nil || membersB == nil {
			err = tr.loadFromAPI(tx)
		} else {
			// TODO
			err = tr.loadFromDB(infoB, membersB)

		}

		if err != nil {
			return fmt.Errorf("fail to init: %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tr, nil
}

// Load data from the DB and store in Tracker
func (t *Tracker) loadFromDB(info, members *bolt.Bucket) error {
	lastUpdated, err := time.Parse(time.RFC3339, string(info.Get([]byte("LastUpdated"))))
	if err != nil {
		return err
	}
	t.LastUpdated = lastUpdated

	c := members.Cursor()
	membersList := []*BGMember{}
	for k, v := c.First(); k != nil; k, v = c.Next() {
		curr, err := BGMemberFromJSON(v)
		if err != nil {
			return err
		}
		membersList = append(membersList, curr)
	}
	t.Members = membersList
	return nil
}

// Get all members from the Github API, store in Tracker
// and store in DB
func (t *Tracker) loadFromAPI(tx *bolt.Tx) error {
	memberList, _, err := client.Organizations.ListMembers(t.Orgname, &github.ListMembersOptions{})
	if err != nil {
		return err
	}

	var bgMembers = []*BGMember{}
	for _, u := range memberList {
		bgm := memberFromUser(u)
		bgMembers = append(bgMembers, bgm)
	}

	t.Members = bgMembers
	t.LastUpdated = time.Now()

	// Store in DB
	trackB, err := reInitBucket(tx, trackerInfoBucketName)
	if err != nil {
		return err
	}

	trackB.Put([]byte("Orgname"), []byte(t.Orgname))
	formattedTime := t.LastUpdated.Format(time.RFC3339)
	trackB.Put([]byte("LastUpdated"), []byte(formattedTime))

	memB, err := reInitBucket(tx, membersBucketName)
	if err != nil {
		return err
	}

	for _, mem := range t.Members {
		err = memB.Put([]byte(mem.GithubID), mem.ToJSON())
		if err != nil {
			return err
		}
	}

	return err
}

// A BGMember represents a member of the Basement Gang
type BGMember struct {
	GithubID   string
	Name       string
	NoCommits  int
	StreakDays int
}

func BGMemberFromJSON(js []byte) (*BGMember, error) {
	var bgm BGMember
	err := json.Unmarshal(js, &bgm)
	return &bgm, err
}

func (b *BGMember) ToJSON() []byte {
	j, _ := json.Marshal(b)
	return j
}

// TODO: remove, because hardcoded
var client = github.NewClient(nil)

// Get one BGMember from the Github API
func GetBGMember(username string) (*BGMember, error) {
	user, _, err := client.Users.Get(username)
	if err != nil {
		return nil, err
	}

	bgm := memberFromUser(*user)
	return bgm, nil
}

// Helpers
// Creates a BGMember from a Github API User
func memberFromUser(user github.User) *BGMember {
	name := ""
	if user.Name != nil {
		name = *user.Name
	}

	return &BGMember{
		GithubID:   *user.Login,
		Name:       name,
		NoCommits:  0,
		StreakDays: 0,
	}
}

// Recreates a bucket
func reInitBucket(tx *bolt.Tx, bucketName []byte) (*bolt.Bucket, error) {
	b := tx.Bucket(bucketName)
	if b != nil {
		tx.DeleteBucket(bucketName)
	}
	return tx.CreateBucket(bucketName)
}
