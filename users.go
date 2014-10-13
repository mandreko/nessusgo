package nessusgo

import (
	"strconv"
	"time"
)

type User struct {
	Name      string
	IsAdmin   bool
	LastLogin time.Time
}

type UserResource struct {
	client *Client
}

// Requests the current Nessus server load and platform type.
func (u *UserResource) List() (*[]User, error) {
	if !u.client.isAuthenticated() {
		return nil, ErrNotAuthenticated
	}

	record := Record{}
	url_path := "/users/list"

	if err := u.client.do("POST", url_path, nil, nil, &record); err != nil {
		return nil, err
	}

	entity, err := convertUserDTO(record.Reply.Contents.Users.Users)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func convertUserDTO(u []*user) (*[]User, error) {
	users := make([]User, len(u))
	for index, element := range u {

		isAdmin, err := strconv.ParseBool(element.IsAdmin)
		if err != nil {
			return nil, err
		}

		t := time.Unix(int64(element.LastLogin), 0)
		users[index] = User{Name: element.Name, IsAdmin: isAdmin, LastLogin: t}
	}

	return &users, nil
}
