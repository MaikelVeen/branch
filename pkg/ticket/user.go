package ticket

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	DisplayName string     `json:"display_name"`
	Email       string     `json:"email"`
	System      SystemType `json:"system"`
}

// SaveToDisk dumps the current user to the file system
// so that it can be retrieved later.
func (u *User) SaveToDisk() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/.branch-cli", dirname))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(dataBytes)
	if err != nil {
		return err
	}

	return f.Sync()
}

// LoadFromDisk reads the user from the home dir.
func LoadFromDisk() (*User, error) {
	u := &User{}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return u, err
	}

	filename := fmt.Sprintf("%s/.branch-cli", dirname)

	dat, err := os.ReadFile(filename)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(dat, u)
	if err != nil {
		return u, err
	}

	return u, nil
}
