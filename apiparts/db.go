package epicrm_apiparts

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Expected by this package; apps should use this variable as the minimum
// version when calling GetDBConn() for authdb if the connection handle is
// meant to be passed to functions from this package.
var DBVER_AUTHDB = 1

func pgxTryExec(dbcon *pgxpool.Pool, sql string, args ...interface{}) {
	_, err := dbcon.Exec(context.Background(), sql, args...)
	if err != nil {
		log.Print("pgxTryExec(): " + err.Error())
	}
}

func pgQuote(conn *pgxpool.Pool, id string, fmt string) string {
	var ret string

	err := conn.QueryRow(context.Background(),
		"SELECT FORMAT('" + fmt + "', $1::VARCHAR)", id).Scan(&ret)

	if err != nil {
		log.Fatal("pgQuote(fmt): " + err.Error())
	}
	
	return ret
}

func qI(conn *pgxpool.Pool, id string) string {
	return pgQuote(conn, id, "%I")
}

func qL(conn *pgxpool.Pool, id string) string {
	return pgQuote(conn, id, "%L")
}

// TODO FIXME SECURITY convenient, but causes passwords to be logged if something goes wrong
// TODO log on success
func CreateForeignDataWrapper(localdb string, foreigndb string, table string, conn *pgxpool.Pool) {
	envPrefixForeign := strings.ToUpper(foreigndb)
	envPrefixLocal   := strings.ToUpper(localdb)
	
	foreignServer := foreigndb + "_remote"

	pgxTryExec(conn, "CREATE EXTENSION IF NOT EXISTS postgres_fdw;")
	
	pgxTryExec(conn, fmt.Sprintf(`CREATE SERVER IF NOT EXISTS %s
        FOREIGN DATA WRAPPER postgres_fdw
        OPTIONS (host %s, port %s, dbname %s);`,
        qI(conn, foreignServer),
        qL(conn, os.Getenv(envPrefixForeign + "_POSTGRES_HOST")),
        qL(conn, os.Getenv(envPrefixForeign + "_POSTGRES_PORT")),
        qL(conn, os.Getenv(envPrefixForeign + "_POSTGRES_NAME"))))

	pgxTryExec(conn, fmt.Sprintf(`CREATE USER MAPPING IF NOT EXISTS FOR %s
        SERVER %s
        OPTIONS (user %s, password %s);`,
				qI(conn, os.Getenv(envPrefixLocal + "_POSTGRES_USER")),
        qI(conn, foreignServer),
        qL(conn, os.Getenv(envPrefixForeign + "_POSTGRES_USER")),
        qL(conn, os.Getenv(envPrefixForeign + "_POSTGRES_PASSWORD"))))

	pgxTryExec(conn, fmt.Sprintf(`GRANT USAGE ON FOREIGN SERVER %s TO %s;`,
        qI(conn, foreignServer),
        qI(conn, os.Getenv(envPrefixLocal + "_POSTGRES_USER"))))

	// TODO don't if already exists
	pgxTryExec(conn, fmt.Sprintf(`IMPORT FOREIGN SCHEMA public
        LIMIT TO (%s) FROM SERVER %s INTO public;`,
        qI(conn, table),
        qI(conn, foreignServer)))
}

func AssertDBIsUp(dbname string, appname string) {
	conn := GetDBConn(dbname, appname, 0)
	if conn == nil { // Shouldn't reach here
		panic("conn = nil")
	}
	conn.Close()
}

func assertDBVersion(dbname string, appname string, dbverMinimum int, dbcon *pgxpool.Pool) {
	var dbverActual int

	err := dbcon.QueryRow(context.Background(),
		"SELECT ival FROM dbmeta WHERE key='dbversion';").Scan(&dbverActual)
	
	if err != nil {
		log.Fatalf("Error: assertDBVersion(dbname = %s, appname = %s): %s",
			dbname, appname, err.Error())
	}
	
	if dbverActual < dbverMinimum {
		log.Fatalf(
			"Failure: assertDBVersion(dbname = %s, appname = %s): dbverActual = %d (needed %d)",
			dbname, appname, dbverActual, dbverMinimum)
	}

	log.Printf(
		"Success: assertDBVersion(dbname = %s, appname = %s, dbverMinimum %d); dbverActual = %d",
		dbname, appname, dbverMinimum, dbverActual)
}

func GetDBConn(dbname string, appname string, dbversion int) *pgxpool.Pool {
	DBCON_RETRY := 10
	DBCON_INTERV := 6 * time.Second

	var dbcon *pgxpool.Pool = nil
	var dberr error

	envprefix := strings.ToUpper(dbname)
	logprefix := "Database connection from app " + appname + " to " + dbname

	for i := 1; dbcon == nil && i <= DBCON_RETRY; i++ {
		dbcon, dberr = pgxpool.Connect( context.Background(),
			// TODO FIXME connects without password (check both pgbouncer/pgpool and the actual db node)?
			"user=" + os.Getenv(envprefix + "_POSTGRES_USER") +
				" password=" + os.Getenv(envprefix + "_POSTGRES_PASSWORD") +
				" host=" + os.Getenv(envprefix + "_POSTGRES_HOST") +
				" port=" + os.Getenv(envprefix + "_POSTGRES_PORT") +
				" dbname=" + os.Getenv(envprefix + "_POSTGRES_NAME") )

		if dberr != nil {
		
			log.Print(logprefix + ": " + dberr.Error())
			
			dbcon = nil
			
			time.Sleep(DBCON_INTERV)
		} else {
			log.Print(logprefix + ": Success.")
		}
	}

	if dberr != nil {
		log.Fatal(dberr)
	}
	
	assertDBVersion(dbname, appname, dbversion, dbcon)
	
	return dbcon
}
