package SABFunctions

import (
//	"fmt"
	"log"
//	"strings"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABDefine"
)

func RemoveNoChildrenORA(conf *SABDefine.Config_STR) int {

	var	(
		lastx	=	int(0)
		lasty	=	int(0)
	)


	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	log.Printf("Remove NoChildren from Ora ORGS...\n")

	rows, err := db.Query(SABDefine.PG_QUE_RemoveNoChildren[0])
	if err != nil {
		log.Printf("PG::Query() Que 11 error: %v\n", err)
		return 11
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		_, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[1])
		if err != nil {
			log.Printf("PG::Query() Que 12 error: %v\n", err)
			return 12
		}

		rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[0])
		if err != nil {
			log.Printf("PG::Query() Que 13 error: %v\n", err)
			return 13
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
		}else{
			break
		}
	}

	log.Printf("Remove NoChildren from Ora DEPS...\n")

	rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[2])
	if err != nil {
		log.Printf("PG::Query() Que 11 error: %v\n", err)
		return 11
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		_, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[3])
		if err != nil {
			log.Printf("PG::Query() Que 12 error: %v\n", err)
			return 12
		}

		rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[2])
		if err != nil {
			log.Printf("PG::Query() Que 13 error: %v\n", err)
			return 13
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
		}else{
			break
		}
	}

	log.Printf("\t\tFINISHED\n")

	return 94

}

func RemoveNoChildrenLDAP(conf *SABDefine.Config_STR) int {

	var	(
		lastx	=	int(0)
		lasty	=	int(0)
	)


	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	log.Printf("Remove NoChildren from Ora LDAP 1...\n")

	rows, err := db.Query(SABDefine.PG_QUE_RemoveNoChildren[4])
	if err != nil {
		log.Printf("PG::Query() Que 11 error: %v\n", err)
		return 11
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		_, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[5])
		if err != nil {
			log.Printf("PG::Query() Que 12 error: %v\n", err)
			return 12
		}

		rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[4])
		if err != nil {
			log.Printf("PG::Query() Que 13 error: %v\n", err)
			return 13
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
		}else{
			break
		}
	}

	log.Printf("Remove NoChildren from Ora LDAP 2...\n")

	rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[6])
	if err != nil {
		log.Printf("PG::Query() Que 111 error: %v\n", err)
		return 111
	}

	rows.Next()
	rows.Scan(&lastx)

	lasty=lastx
	log.Printf("\t%6d / %6d", lastx, lasty)

	for {
		_, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[7])
		if err != nil {
			log.Printf("PG::Query() Que 112 error: %v\n", err)
			return 112
		}

		rows, err = db.Query(SABDefine.PG_QUE_RemoveNoChildren[6])
		if err != nil {
			log.Printf("PG::Query() Que 113 error: %v\n", err)
			return 113
		}
		rows.Next()
		rows.Scan(&lastx)
		log.Printf("\t%6d / %6d", lastx, lasty)
		if lastx<lasty {
			lasty=lastx
		}else{
			break
		}
	}

	log.Printf("\t\tFINISHED\n")

	return 94

}

