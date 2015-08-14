package SABFunctions

import (
//	"fmt"
	"log"
	"strings"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABDefine"
)

func GetORGs(conf *SABDefine.Config_STR) int {

	var	(
		GlobaParent	string
		que		string
		lastx	=	int(0)
		lasty	=	int(0)
	)

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 11
	}

	defer db.Close()
	log.Printf("LDAP deps and orgs...\n")

	_, err = db.Exec(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS1, "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 12 error: %v\n", err)
		return 12
	}

	rows, err := db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[0], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 13 error: %v\n", err)
		return 13
	}

	rows.Next()
	rows.Scan(&GlobaParent)

	log.Printf("\tGlobal prent: %s", GlobaParent)

	rows, err = db.Query(strings.Replace(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[1], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1), "XYZGlbDNXYZ", GlobaParent, -1))
	if err != nil {
		log.Printf("PG::Query() Que 14 error: %v\n", err)
		return 14
	}

	rows, err = db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[2], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 15 error: %v\n", err)
		return 15
	}

	_, err = db.Query(strings.Replace(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[3], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1), "XYZGlbDNXYZ", GlobaParent, -1))
	if err != nil {
		log.Printf("PG::Query() Que 16 error: %v\n", err)
		return 16
	}

	que = strings.Replace(SABDefine.PG_QUE_LDAP_ORGS11, "XYZGlbParXYZ", SABDefine.GlobaParentId, -1)
//	log.Printf("%s\n", que)
	_, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 17 error: %v\n", err)
		return 17
	}

	que = strings.Replace(SABDefine.PG_QUE_LDAP_ORGS12[0], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1)
//	log.Printf("%s\n", que)
	_, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 18 error: %v\n", err)
		return 18
	}

	que = SABDefine.PG_QUE_LDAP_ORGS12[1]
//	log.Printf("%s\n", que)
	rows, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 19 error: %v\n", err)
		return 19
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		que = SABDefine.PG_QUE_LDAP_ORGS12[2]
//		log.Printf("%s\n", que)
		_, err = db.Query(que)
		if err != nil {
			log.Printf("PG::Query() Que 110 error: %v\n", err)
			return 110
		}

		que = SABDefine.PG_QUE_LDAP_ORGS12[1]
//		log.Printf("%s\n", que)
		rows, err = db.Query(que)
		if err != nil {
			log.Printf("PG::Query() Que 111 error: %v\n", err)
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

	rows, err = db.Query(SABDefine.PG_QUE_LDAP_ORGS12[3])
	if err != nil {
		log.Printf("PG::Query() Que 112 error: %v\n", err)
		return 112
	}

	log.Printf("\t\tFINISHED\n")

	log.Printf("LDAP peoples...\n")
	_, err = db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS21, "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 210 error: %v\n", err)
		return 210
	}

	_, err = db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS22[0], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 211 error: %v\n", err)
		return 211
	}

	rows, err = db.Query(SABDefine.PG_QUE_LDAP_ORGS22[1])
	if err != nil {
		log.Printf("PG::Query() Que 212 error: %v\n", err)
		return 212
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)
	for {

		que = SABDefine.PG_QUE_LDAP_ORGS22[2]
//		log.Printf("%s\n", que)
		_, err = db.Query(que)
		if err != nil {
			log.Printf("PG::Query() Que 214 error: %v\n", err)
			return 213
		}

		rows, err = db.Query(SABDefine.PG_QUE_LDAP_ORGS22[1])
		if err != nil {
			log.Printf("PG::Query() Que 214 error: %v\n", err)
			return 214
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

	rows, err = db.Query(SABDefine.PG_QUE_LDAP_ORGS22[3])
	if err != nil {
		log.Printf("PG::Query() Que 215 error: %v\n", err)
		return 215
	}

	log.Printf("\t\tFINISHED\n")

	log.Printf("LDAP phones...\n")

	rows, err = db.Query(SABDefine.PG_QUE_LDAP_ORGS31)
	if err != nil {
		log.Printf("PG::Query() Que 311 error: %v\n", err)
		return 311
	}

	log.Printf("\t\tFINISHED\n")

	return 94

}

