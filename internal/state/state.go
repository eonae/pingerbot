package state

import (
	"database/sql"
	"pingerbot/pkg/helpers"
	"pingerbot/pkg/telegram"
)

type State struct {
	db *sql.DB
}

const (
	DuplicateGroup = helpers.Error("DuplicateGroup")
	GroupNotFound  = helpers.Error("GroupNotFound")
)

func New(db *sql.DB) State {
	return State{db}
}

func (s *State) groupExists(groupId int64) (bool, error) {
	rows, err := s.db.Query(`SELECT id FROM groups WHERE id = $1 LIMIT 1`, groupId)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (s *State) RememberGroup(c telegram.Chat) (err error) {
	exists, err := s.groupExists(c.Id)
	if err != nil {
		return err
	}

	if exists {
		return DuplicateGroup
	}

	_, err = s.db.Exec(`INSERT INTO groups (id, name) VALUES ($1, $2)`, c.Id, c.Title)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s *State) ForgetGroup(c telegram.Chat) (err error) {
	_, err = s.db.Exec(`DELETE FROM groups WHERE id = $1`, c.Id)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s *State) RememberMember(groupId int64, u telegram.User) (err error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return err
	}

	if !exists {
		return GroupNotFound
	}

	_, err = s.db.Exec(`
		INSERT INTO members (username, group_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, u.Username, groupId)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s *State) ForgetMember(groupId int64, u telegram.User) (err error) {
	_, err = s.db.Exec(`DELETE FROM members WHERE username = $1 AND group_id = $2`, u.Username, groupId)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s State) GetKnownMembers(groupId int64) ([]string, error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, GroupNotFound
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
