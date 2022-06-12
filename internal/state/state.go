package state

import (
	"database/sql"
	"fmt"
	"pingerbot/pkg/helpers"
	"pingerbot/pkg/telegram"
	"strings"
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
func (s *State) RememberMember(groupId int64, username string, tags []string) (err error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return err
	}

	if !exists {
		return GroupNotFoundErr
	}

	if len(tags) == 0 {
		_, err = s.db.Exec(`
			INSERT INTO members (group_id, username, tag) VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, groupId, "@"+username, "-")

		return err
	}

	i := 0
	values := make([]string, 0, len(tags))
	parameters := make([]interface{}, 0, len(tags)*3)

	for _, tag := range tags {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", i+1, i+2, i+3))
		i += 3
		parameters = append(parameters, groupId, "@"+username, tag)
	}

	query := fmt.Sprintf(`
		INSERT INTO members (group_id, username, tag) VALUES %s
		ON CONFLICT DO NOTHING
	`, strings.Join(values, ", "))

	fmt.Println(query, parameters)

	_, err = s.db.Exec(query, parameters...)

	return err
}

// Forget member. Is used when user
func (s *State) ForgetMember(groupId int64, username string, tags []string) (err error) {
	if len(tags) == 0 {
		_, err = s.db.Exec(`DELETE FROM members WHERE username = $1 AND group_id = $2 AND tag = $3`, username, groupId, "-")

		return err
	}

	i := 2
	placeholders := make([]string, 0, len(tags))
	parameters := []interface{}{username, groupId}
	for _, tag := range tags {
		i++
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		parameters = append(parameters, tag)
	}

	_, err = s.db.Exec(
		fmt.Sprintf(
			`DELETE FROM members WHERE username = $1 AND group_id = $2 AND tag IN (%s)`,
			strings.Join(placeholders, ", "),
		),
		parameters...,
	)

	return err
}

func (s State) GetKnownMembers(groupId int64, tags []string) ([]string, error) {
	exists, err := s.groupExists(groupId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, GroupNotFoundErr
	}

	var query string
	var parameters []interface{}

	if len(tags) == 0 {
		query = `SELECT username FROM members WHERE group_id = $1 AND tag = $2`
		parameters = []interface{}{groupId, "-"}
	} else {
		placeholders := make([]string, 0, len(tags))
		i := 1
		for range tags {
			i++
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		}
		query = fmt.Sprintf(
			`SELECT username FROM members WHERE group_id = $1 AND tag IN (%s)`,
			strings.Join(placeholders, ", "),
		)
		parameters = make([]interface{}, 0, len(tags)+1)
		parameters = append(parameters, groupId)
		for _, tag := range tags {
			parameters = append(parameters, tag)
		}
	}

	fmt.Println(query, parameters)

	rows, err := s.db.Query(query, parameters...)
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
