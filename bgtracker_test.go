package bgtracker_test

import (
	"testing"

	"github.com/ejamesc/bgtracker"
)

func TestGetBGMember(t *testing.T) {
	username := "nattsw"
	bgm, err := bgtracker.NewBGMember(username)

	ok(t, err)

	assert(t, bgm.GithubID == username, "expected bgm.GithubID to be %s, got %s", username, bgm.GithubID)
	assert(t, bgm.Name == "", "expected bgm.Name to be empty, got %s", bgm.Name)
	assert(t, bgm.NoCommits == 0, "expected bgm.NoCommits to be 0, got %v", bgm.NoCommits)
	assert(t, bgm.StreakDays == 0, "expected bgm.StreakDays to be 0, got %v", bgm.StreakDays)
}
