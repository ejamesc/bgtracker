package bgtracker

import "github.com/google/go-github/github"

type BGMember struct {
	GithubID   string
	Name       string
	NoCommits  int
	StreakDays int
}

var client = github.NewClient(nil)

func GetBGMember(username string) (*BGMember, error) {
	user, _, err := client.Users.Get(username)
	if err != nil {
		return nil, err
	}

	bgm := memberFromUser(*user)

	return bgm, nil
}

func GetAllBGMembers(orgname string) ([]*BGMember, error) {
	memberList, _, err := client.Organizations.ListMembers(orgname, &github.ListMembersOptions{})
	if err != nil {
		return nil, err
	}

	var bgMembers = []*BGMember{}
	for _, u := range memberList {
		bgm := memberFromUser(u)
		bgMembers = append(bgMembers, bgm)
	}
	return bgMembers, nil
}

func memberFromUser(user github.User) *BGMember {
	name := ""
	if user.Name != nil {
		name = *user.Name
	}

	bgm := &BGMember{
		GithubID:   *user.Login,
		Name:       name,
		NoCommits:  0,
		StreakDays: 0,
	}

	return bgm
}
