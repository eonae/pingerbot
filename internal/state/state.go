package state

import (
	"database/sql"
	"pingerbot/pkg/helpers"
	"pingerbot/pkg/telegram"
)

// Represents bot's memory
type State struct {
	db *sql.DB
}

const (
	DuplicateGroupErr = helpers.Error("DuplicateGroup")
	GroupNotFoundErr  = helpers.Error("GroupNotFound")
)

// Create new state instance
func New(db *sql.DB) State {
	return State{db}
}

// Remember group
func (s *State) RememberGroup(c telegram.Chat) (err error) {
	exists, err := s.groupExists(c.Id)
	if err != nil {
		return err
	}

	if exists {
		return DuplicateGroupErr
	}

	_, err = s.db.Exec(`INSERT INTO groups (id, name) VALUES ($1, $2)`, c.Id, c.Title)

	return err
}

// Drop all knowledge about group, including it's members
func (s *State) ForgetGroup(c telegram.Chat) (err error) {
	_, err = s.db.Exec(`DELETE FROM groups WHERE id = $1`, c.Id)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

// Remember that user is member of specified group
// Note that method is idempotent
func (s *State) RememberMember(groupId int64, username string) (err error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return err
	}

	if !exists {
		return GroupNotFoundErr
	}

	_, err = s.db.Exec(`
		INSERT INTO members (username, group_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, "@"+username, groupId)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

// Forget member. Is used when user
func (s *State) ForgetMember(groupId int64, u telegram.User) (err error) {
	_, err = s.db.Exec(`DELETE FROM members WHERE username = $1 AND group_id = $2`, u.Username, groupId)

	return err
}

func (s State) GetKnownMembers(groupId int64) ([]string, error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, GroupNotFoundErr
	}

	rows, err := s.db.Query(`SELECT username FROM members WHERE group_id = $1`, groupId)
	if err != nil {
		return nil, err
	}

	members := make([]string, 0)

	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		members = append(members, username)
	}

	return members, nil
}

func (s *State) groupExists(groupId int64) (bool, error) {
	rows, err := s.db.Query(`SELECT id FROM groups WHERE id = $1 LIMIT 1`, groupId)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}
