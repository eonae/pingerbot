package state

import (
	"database/sql"
	"pingerbot/pkg/helpers"
	"pingerbot/pkg/telegram"
)

type Member struct {
	Id   string // Because these could be really bit integers and pq seems to not support them
	Name string
}

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
	rows, err := s.db.Query(`SELECT id FROM pingerbot.groups WHERE id = $1 LIMIT 1`, groupId)
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

	_, err = s.db.Exec(`INSERT INTO pingerbot.groups (id, name) VALUES ($1, $2)`, c.Id, c.Title)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s *State) ForgetGroup(c telegram.Chat) (err error) {
	_, err = s.db.Exec(`DELETE FROM pingerbot.members WHERE group_id = $1`, c.Id)
	if err == sql.ErrNoRows {
		return nil
	}

	_, err = s.db.Exec(`DELETE FROM pingerbot.groups WHERE id = $1`, c.Id)
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
		INSERT INTO pingerbot.members (id, name, group_id) VALUES ($1, $2, $3)
		ON CONFLICT (id, group_id) DO UPDATE SET name = $2
	`, u.Id, u.Username, groupId)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s *State) ForgetMember(groupId int64, u telegram.User) (err error) {
	_, err = s.db.Exec(`DELETE FROM pingerbot.members WHERE id = $1 AND group_id = $2`, u.Id, groupId)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (s State) GetKnownMembers(groupId int64) ([]Member, error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, GroupNotFound
	}

	rows, err := s.db.Query(`SELECT id, name FROM pingerbot.members WHERE group_id = $1`, groupId)
	if err != nil {
		return nil, err
	}

	members := make([]Member, 0)

	for rows.Next() {
		var u Member
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			return nil, err
		}
		members = append(members, u)
	}

	return members, nil
}
