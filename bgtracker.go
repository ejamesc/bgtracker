package bgtracker

import "github.com/google/go-github/github"

type BGMember struct {
	GithubID   string
	Name       string
	NoCommits  int
	StreakDays int
}

func NewBGMember(username string) (*BGMember, error) {
	client := github.NewClient(nil)
	user, _, err := client.Users.Get(username)
	if err != nil {
		return nil, err
	}

	var name string
	if user.Name != nil {
		name = *user.Name
	} else {
		name = ""
	}

	bgm := &BGMember{
		GithubID:   *user.Login,
		Name:       name,
		NoCommits:  0,
		StreakDays: 0,
	}

	return bgm, nil
}
