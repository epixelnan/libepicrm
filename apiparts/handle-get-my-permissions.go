package epicrm_apiparts

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO get service from r
func HandleGetMyPermissions(w http.ResponseWriter, r *http.Request, service string, activity string, dbcon *pgxpool.Pool) {
	var json string

	// TODO FIXME IMPORTANT SECURITY limit tenant
	err := dbcon.QueryRow(context.Background(),
		`SELECT COALESCE(json_agg(r), '[]'::json) FROM (
		  SELECT organization, permission FROM user_org_role
		  INNER JOIN role_permission ON user_org_role.role=role_permission.role
		  WHERE uid=$1) r;`,
		GetUidOrPanic(r)).Scan(&json)
	
	if err != nil {
		LogErrorAndSend500(w, service, activity, err)
		return
	}
	
	SendJsonStringResponse(w, json)
}
