package epicrm_apiparts

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetSysOrgId(appname string, authdbcon *pgxpool.Pool) uuid.UUID {
	var sysOrgId uuid.UUID

	err := authdbcon.QueryRow(context.Background(),
		"SELECT sys_org_id FROM sysinfo;").Scan(&sysOrgId)
	
	if err != nil {
		log.Fatalf("Error: getSysOrgId(appname = %s): %s", appname, err.Error())
	}

	return sysOrgId
}

// TODO FIXME this is a workaround.
// Calling GetSysOrgId() at startup is problematic if the DB is empty
func GetSysOrgIdOr500(appname string, authdbcon *pgxpool.Pool, w http.ResponseWriter) (uuid.UUID, bool) {
	var sysOrgId uuid.UUID

	err := authdbcon.QueryRow(context.Background(),
		"SELECT sys_org_id FROM sysinfo;").Scan(&sysOrgId)
	
	if err != nil {
		LogErrorAndSend500(w, appname, "getSysOrgId", err)
		return sysOrgId, false
	}

	return sysOrgId, true
}
