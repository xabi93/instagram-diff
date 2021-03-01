package instagram

import (
	"context"
	"fmt"
	"sort"

	"golang.org/x/sync/errgroup"
)

type client interface {
	Followers(context.Context) (map[int64]User, error)
	Following(context.Context) (map[int64]User, error)
}

func New(cli client) *Instadiff {
	return &Instadiff{cli}
}

type Instadiff struct {
	cli client
}

type DiffResult struct {
	FollowNotFollower []User
	FollowerNotFollow []User
}

func (i Instadiff) Diff(ctx context.Context) (*DiffResult, error) {
	g, _ := errgroup.WithContext(ctx)

	followers := make(map[int64]User)
	g.Go(func() error {
		var err error
		followers, err = i.cli.Followers(ctx)
		if err != nil {
			return fmt.Errorf("getting followers: %s", err.Error())
		}

		return nil
	})

	following := make(map[int64]User)
	g.Go(func() error {
		var err error
		following, err = i.cli.Following(ctx)
		if err != nil {
			return fmt.Errorf("getting followings: %s", err.Error())
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &DiffResult{
		FollowNotFollower: diffAndSort(following, followers),
		FollowerNotFollow: diffAndSort(followers, following),
	}, nil
}

func diffAndSort(src, comparator map[int64]User) []User {
	var users []User

	for id, u := range src {
		if _, ok := comparator[id]; !ok {
			users = append(users, u)
		}
	}

	sort.Slice(users, func(j, i int) bool {
		return users[j].Username <= users[i].Username
	})

	return users
}
