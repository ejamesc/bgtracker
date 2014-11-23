package bgtracker_test

import (
	"testing"

	"github.com/ejamesc/bgtracker"
)

func TestGetBGMember(t *testing.T) {
	username := "nattsw"
	bgm, err := bgtracker.GetBGMember(username)

	ok(t, err)

	assert(t, bgm.GithubID == username, "expected bgm.GithubID to be %s, got %s", username, bgm.GithubID)
	assert(t, bgm.Name == "", "expected bgm.Name to be empty, got %s", bgm.Name)
	assert(t, bgm.NoCommits == 0, "expected bgm.NoCommits to be 0, got %v", bgm.NoCommits)
	assert(t, bgm.StreakDays == 0, "expected bgm.StreakDays to be 0, got %v", bgm.StreakDays)
}

func TestGetAllBGMembers(t *testing.T) {
	orgname := "basement-gang"
	members, err := bgtracker.GetAllBGMembers(orgname)

	ok(t, err)
	assert(t, len(members) == 9, "expected members to be 9, got %v instead", len(members))
}
