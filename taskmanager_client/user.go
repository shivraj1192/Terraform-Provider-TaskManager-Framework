package taskmanager_client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateUser(user User) (*UserModel, error) {
	url := "/register"
	body, err := c.Post(url, user, UserLockName)
	if err != nil {
		return nil, err
	}

	var userModel UserModel
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (c *Client) GetUser(userID int64) (*User, error) {
	url := fmt.Sprintf("/users/%d", userID)
	body, err := c.Get(url, UserLockName)
	if err != nil {
		return nil, err
	}

	var userModel UserModel
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		return nil, err
	}

	return &userModel.User, nil
}

func (c *Client) UpdateUser(user User) (*UserModel, error) {
	url := fmt.Sprintf("/users/%d", user.ID)
	body, err := c.Put(url, user, UserLockName)
	if err != nil {
		return nil, err
	}

	var userModel UserModel
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (c *Client) DeleteUser(userID int64) error {
	url := fmt.Sprintf("/users/%d", userID)
	err := c.Delete(url, UserLockName)
	if err != nil {
		return err
	}

	return nil
}
