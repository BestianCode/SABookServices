package SABFunctions

import (
	"fmt"
	"log"
	"strings"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"

	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func LDAP_Make(conf *SABModules.Config_STR) int {

	var	(
		GlobalParentInsert	=	string("")

		GlobalParent			string
		GlobalParentID			int

		ckl				int

		lastx			=	int(0)
		lasty			=	int(0)

		queryx				string

		ldap_table_check	=	int(1)

		ldap_que_check_tables	=	string ("SELECT count(tablename) FROM pg_catalog.pg_tables where tablename like 'ldap%';")

	)

	log.Printf(".")
	log.Printf("..")
	log.Printf("...")
	log.Printf("Building LDAP Tree...")

	for ckl=0;ckl<len(conf.ROOT_DN);ckl++ {
		GlobalParentInsert=fmt.Sprintf("%sINSERT INTO ldap_entries VALUES (%d,'%s',%d,%d,%d,'%d','%d',%d); ", GlobalParentInsert, ckl+1, conf.ROOT_DN[ckl][0], 3, ckl, ckl+1, ckl+1, ckl, 0)
		GlobalParentInsert=fmt.Sprintf("%sINSERT INTO ldapx_institutes VALUES (%d,'%s','%d','%d',%d); ", GlobalParentInsert, ckl+1, conf.ROOT_DN[ckl][1], ckl+1, ckl, 0)
	}
	GlobalParentID=ckl
	GlobalParent=conf.ROOT_DN[ckl-1][0]
//	log.Printf("%d %s\n\n%s\n", GlobalParentID, GlobalParent, GlobalParentInsert)

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	log.Printf("\tChecking LDAP DB consistency...")

	rows, err := db.Query(ldap_que_check_tables)
	if err != nil {
		log.Printf("PG::Query() Check LDAP tables error: %v\n", err)
		return 11
	}

	rows.Next()
	rows.Scan(&ldap_table_check)

	if ldap_table_check<SABDefine.LDAP_Tables_am {
		log.Printf("\t\tHmmm...  %d tables instead of %d... Rebuilding LDAP DB!", ldap_table_check, SABDefine.LDAP_Tables_am)
		queryx=strings.Replace(SABDefine.LDAP_Scheme_create, "XYZInsertIntoXYZ", GlobalParentInsert, -1)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() Create LDAP scheme error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 12
		}
		log.Printf("\tLDAP DB Rebuilded!")
	}else{
		log.Printf("\tLDAP DB Good!")
	}

	log.Printf("\t\tUpdate LDAP orgs and deps...")

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1	, "XYZGlbParXYZ",	fmt.Sprintf("%d", GlobalParentID), -1)
	queryx=strings.Replace(queryx				, "XYZDBOrgsXYZ",	SABDefine.PG_Table_MSSQL[0], -1)
	queryx=strings.Replace(queryx				, "XYZDBDepsXYZ",	SABDefine.PG_Table_MSSQL[1], -1)
	queryx=strings.Replace(queryx				, "XYZGlbDNXYZ",	GlobalParent, -1)
//	log.Printf("%s\n", queryx)
	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() ORGS1 Execute error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 13
	}

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1X_GET	, "XYZDBDepsXYZ",	SABDefine.PG_Table_MSSQL[1], -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() ORGS1X_GET error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 14
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		queryx=strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1X_PUT	, "XYZDBDepsXYZ",	SABDefine.PG_Table_MSSQL[1], -1)
//		log.Printf("%s\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() ORGS1X_PUT in FOR error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 110
		}

		queryx=strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1X_GET	, "XYZDBDepsXYZ",	SABDefine.PG_Table_MSSQL[1], -1)
//		log.Printf("%s\n", queryx)
		rows, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() ORGS1X_GET in FOR error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 111
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
//			log.Printf("Good")
		}else{
//			log.Printf("Stop")
			break
		}
	}

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1_END	, "XYZDBDepsXYZ",	SABDefine.PG_Table_MSSQL[1], -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() ORGS1_END error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 15
	}

	log.Printf("\t\tUpdate LDAP persons...")

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PERS1	, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
	queryx=strings.Replace(queryx				, "XYZGlbParXYZ",	fmt.Sprintf("%d", GlobalParentID), -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() PERS1 error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 16
	}

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PERS1X_GET	, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() PERS1X_GET error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 17
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PERS1X_PUT	, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
//		log.Printf("%s\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() PERS1X_PUT in FOR error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 112
		}

		queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PERS1X_GET	, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
//		log.Printf("%s\n", queryx)
		rows, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() PERS1X_GET error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 113
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
//			log.Printf("Good")
		}else{
//			log.Printf("Stop")
			break
		}
	}

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PERS1_END	, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() PERS1_END error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 18
	}

	log.Printf("\t\tUpdate Phone and Mail...")

	queryx=strings.Replace(SABDefine.PG_QUE_LDAP_PHONES	, "XYZDBPhonesXYZ",	SABDefine.PG_Table_Oracle, -1)
	queryx=strings.Replace(queryx				, "XYZDBMailXYZ",	SABDefine.PG_Table_Domino, -1)
	queryx=strings.Replace(queryx				, "XYZDBPersXYZ",	SABDefine.PG_Table_MSSQL[2], -1)
//	log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() PHONES error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 18
	}

	log.Printf("\tComplete")

	return 94

}

