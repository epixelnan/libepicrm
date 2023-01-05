package epicrm_apiparts

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO use uuid.UUID for uid
func CanUser(permission string, org uuid.UUID, uid uuid.UUID, authdbcon *pgxpool.Pool) (bool, error) {
	var able bool

	// TODO FIXME IMPORTANT SECURITY limit tenant
	err := authdbcon.QueryRow(context.Background(),
		`SELECT EXISTS(
		  SELECT FROM user_org_role
		  INNER JOIN role_permission ON user_org_role.role=role_permission.role
		  WHERE uid=$1 AND organization=$2 AND permission=$3);`,
		uid, org, permission).Scan(&able)
	
	if err != nil {
		return false, err
	}
	
	return able, nil
}

// TODO accept uuid.UUID?
// TODO rename to IsUserAnyOf
// TODO reorder params like CanUser()
func UserIsAnyOf(authdbcon *pgxpool.Pool, uid string, org string, roles []string) (bool, error) {
	tx, err := authdbcon.Begin(context.Background())
	if err != nil {
		return false, fmt.Errorf("UserIsAnyOf(): Begin(): %s: %w", ErrServerError, err)
	}

	defer tx.Rollback(context.Background())

	// Wish I could use this
	// inListPlaceHolders := strings.Repeat("?,", len(roles))

	inListPlaceHolders := ""
	for i := 0; i < len(roles); i++ {
		if i > 0 {
			inListPlaceHolders += ","
		}
	
		inListPlaceHolders += "$" + strconv.Itoa(i + 3)
	}

	sql := `SELECT COUNT(*) FROM user_org_role WHERE
        uid=$1 AND organization=$2 and role IN (` + inListPlaceHolders + `)`

	sqlname := fmt.Sprintf("%x", sha1.Sum([]byte(sql)))

	_, err = tx.Prepare(context.Background(), sqlname, sql)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.SQLState() != "42P05" {
			return false, fmt.Errorf("UserIsAnyOf(): Prepare(): %s: %w", ErrServerError, err)
		}
	}

	// Can't simply do like this because the element types differ:
	// args := append([]interface{}{uid, org}, roles...)
	args := []interface{}{uid, org}
	for _, role := range roles {
		args = append(args, role)
	}

	var count int

	err = tx.QueryRow(context.Background(), sqlname, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("UserIsAnyOf(): QueryRow(): %s: %w", ErrServerError, err)
	}

	return count > 0, nil;
}
