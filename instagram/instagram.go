package instagram

import (
	"context"
	"errors"

	"github.com/kr/pretty"
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

type AuthError struct{ error }

func (i Instagram) Ping() error {
	err := i.instagram.Account.Sync()
	errN := goinsta.ErrorN{}
	if !errors.As(err, &errN) {
		return err
	}
	if errN.Message == "login_required" {
		return AuthError{errN}
	}

	return nil
}

func (i Instagram) Followers(_ context.Context) (map[int64]instadiff.User, error) {
	pretty.Println()
	return getAll(i.instagram.Account.Followers())
}

func (i Instagram) Following(_ context.Context) (map[int64]instadiff.User, error) {
	return getAll(i.instagram.Account.Following())
}

func getAll(i *goinsta.Users) (map[int64]instadiff.User, error) {
	if err := i.Error(); err != nil {
		pretty.Println(err)
		return nil, err
	}

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