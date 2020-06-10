package instagram

import (
	"context"

	instadiff "github.com/xabi93/instagram-diff"

	"github.com/ahmdrz/goinsta/v2"
)

func RestoreSession(path string) (*Instagram, error) {
	i, err := goinsta.Import(path)

	return &Instagram{i}, err
}

func Login(user, pass, exportPath string) (*Instagram, error) {
	i := goinsta.New(user, pass)

	if err := i.Login(); err != nil {
		return nil, err
	}

	i.Export(exportPath)

	return &Instagram{i}, nil
}

type Instagram struct{ instagram *goinsta.Instagram }

func (i Instagram) Followers(_ context.Context) (map[int64]instadiff.User, error) {
	return getAll(i.instagram.Account.Followers())
}

func (i Instagram) Following(_ context.Context) (map[int64]instadiff.User, error) {
	return getAll(i.instagram.Account.Following())
}

func getAll(i *goinsta.Users) (map[int64]instadiff.User, error) {
	users := make(map[int64]instadiff.User)
	for i.Next() {
		for _, u := range i.Users {
			users[u.ID] = instadiff.User{
				ID:         u.ID,
				Username:   u.Username,
				Fullname:   u.FullName,
				ProfilePic: u.ProfilePicURL,
			}
		}
	}

	return users, nil
}
