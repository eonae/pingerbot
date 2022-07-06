package state

import (
	"context"
	"database/sql"
	"fmt"
	"pingerbot/pkg/helpers"
	"pingerbot/pkg/telegram"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Represents bot's memory
type State struct {
	db *pgxpool.Pool
}

const (
	DuplicateGroupErr = helpers.Error("DuplicateGroup")
	GroupNotFoundErr  = helpers.Error("GroupNotFound")
	defaultTag        = "#all"
)

// Create new state instance
func New(db *pgxpool.Pool) State {
	return State{db}
}

// Remember group
func (s *State) RememberGroup(c telegram.Chat) (err error) {
	chatId := GroupId(c.Id).String()
	exists, err := s.groupExists(chatId)
	if err != nil {
		return err
	}

	if exists {
		return DuplicateGroupErr
	}

	ctx := context.Background()

	query := `INSERT INTO groups (id, name) VALUES ($1, $2)`
	_, err = s.db.Exec(ctx, query, chatId, c.Title)

	return err
}

// Drop all knowledge about group, including it's members
func (s *State) ForgetGroup(c telegram.Chat) (err error) {
	ctx := context.Background()
	query := `INSERT INTO groups (id, name) VALUES ($1, $2)`

	_, err = s.db.Exec(ctx, query, c.Id)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

// Remember that user is member of specified group
// Note that method is idempotent
func (s *State) RememberMember(groupId GroupId, username string, tags []string) (err error) {
	exists, err := s.groupExists(groupId.String())
	if err != nil {
		return err
	}

	if !exists {
		return GroupNotFoundErr
	}

	ctx := context.Background()

	query := `
		INSERT INTO members (group_id, username, tag) VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`

	if len(tags) == 0 {
		_, err = s.db.Exec(ctx, query, groupId, "@"+username, defaultTag)

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

	query = fmt.Sprintf(`
		INSERT INTO members (group_id, username, tag) VALUES %s
		ON CONFLICT DO NOTHING
	`, strings.Join(values, ", "))

	fmt.Println(query, parameters)

	_, err = s.db.Exec(ctx, query, parameters...)

	return err
}

// Forget member. Is used when user
func (s *State) ForgetMember(groupId GroupId, username string, tags []string) (err error) {
	ctx := context.Background()
	query := `DELETE FROM members WHERE username = $1 AND group_id = $2 AND tag = ANY($3)`

	if len(tags) == 0 {
		tags = []string{defaultTag}
	}

	_, err = s.db.Exec(ctx, query, username, groupId.String(), tags)

	return err
}

type GroupId int64

func (id GroupId) String() string {
	return strconv.Itoa(int(id))
}

func (s State) GetKnownMembers(groupId GroupId, tags []string) ([]string, error) {
	exists, err := s.groupExists(groupId.String())
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, GroupNotFoundErr
	}

	if len(tags) == 0 {
		tags = []string{defaultTag}
	}

	query := `SELECT username FROM members WHERE group_id = $1 AND tag = ANY($2)`

	ctx := context.Background()
	rows, err := s.db.Query(ctx, query, groupId.String(), tags)
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

type GroupList map[string][]string

func (l GroupList) String() string {
	sb := strings.Builder{}

	for tag, members := range l {
		sb.WriteString(tag)
		sb.WriteString("\n")
		for _, username := range members {
			stripped := username[1:]
			sb.WriteString(fmt.Sprintf("🠖 %s\n", stripped))
		}
	}

	return sb.String()
}

func (s *State) ListGroupMembers(groupId GroupId, tags []string) (GroupList, error) {
	ctx := context.Background()

	tagsFilter := ""
	parameters := []interface{}{groupId.String()}

	if len(tags) > 0 {
		tagsFilter = "AND tag = ANY($2)"
		parameters = append(parameters, tags)
	}

	query := fmt.Sprintf(`
		SELECT
			tag,
			array_agg(username) as members
		FROM members
		WHERE group_id = $1
		%s
		GROUP BY tag;
	`, tagsFilter)

	rows, err := s.db.Query(ctx, query, parameters...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(GroupList)

	for rows.Next() {
		var tag string
		var members []string

		err = rows.Scan(&tag, &members)
		if err != nil {
			return nil, err
		}

		fmt.Println("usernames:", members)

		if len(tags) == 0 || helpers.Includes(tags, tag) {
			result[tag] = members
		}
	}

	fmt.Printf("%+v", result)

	return result, nil
}

func (s *State) groupExists(groupId string) (bool, error) {
	ctx := context.Background()
	query := `SELECT id FROM groups WHERE id = $1 LIMIT 1`

	rows, err := s.db.Query(ctx, query, groupId)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}
