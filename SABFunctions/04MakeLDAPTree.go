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
		log.Printf("PG::Query() Que 1 error: %v\n", err)
		return 12
	}

	rows, err := db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[0], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 2 error: %v\n", err)
		return 13
	}

	rows.Next()
	rows.Scan(&GlobaParent)

	log.Printf("%s", GlobaParent)

	rows, err = db.Query(strings.Replace(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[1], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1), "XYZGlbDNXYZ", GlobaParent, -1))
	if err != nil {
		log.Printf("PG::Query() Que 3 error: %v\n", err)
		return 14
	}

	rows, err = db.Query(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[2], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1))
	if err != nil {
		log.Printf("PG::Query() Que 4 error: %v\n", err)
		return 15
	}

	_, err = db.Query(strings.Replace(strings.Replace(SABDefine.PG_QUE_LDAP_ORGS2[3], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1), "XYZGlbDNXYZ", GlobaParent, -1))
	if err != nil {
		log.Printf("PG::Query() Que 5 error: %v\n", err)
		return 16
	}

	que = strings.Replace(SABDefine.PG_QUE_LDAP_ORGS11, "XYZGlbParXYZ", SABDefine.GlobaParentId, -1)
//	log.Printf("%s\n", que)
	_, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 6 error: %v\n", err)
		return 17
	}

	que = strings.Replace(SABDefine.PG_QUE_LDAP_ORGS12[0], "XYZGlbParXYZ", SABDefine.GlobaParentId, -1)
//	log.Printf("%s\n", que)
	_, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 7 error: %v\n", err)
		return 18
	}

	que = SABDefine.PG_QUE_LDAP_ORGS12[1]
//	log.Printf("%s\n", que)
	rows, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 8 error: %v\n", err)
		return 19
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("%6d / %6d", lastx, lasty)

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
		log.Printf("%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
//			log.Printf("Good")
		}else{
//			log.Printf("Stop")
			break
		}
	}

	que = SABDefine.PG_QUE_LDAP_ORGS12[3]
//	log.Printf("%s\n", que)
	rows, err = db.Query(que)
	if err != nil {
		log.Printf("PG::Query() Que 112 error: %v\n", err)
		return 112
	}

	log.Printf("LDAP deps and orgs... Finish\n")
	log.Printf("LDAP peoples...\n")
	log.Printf("LDAP peoples... Finish\n")

	return 94

}

