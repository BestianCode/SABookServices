package SABFunctions

import (
//	"fmt"
	"log"
	"strings"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"

	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func RemoveNoChildrenCache(conf *SABModules.Config_STR) int {

	var	(
		lastx	=	int(0)
		lasty	=	int(0)
		queryx		string
		ckl		int
	)

	log.Printf("Garbage collection...")

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	log.Printf("\t\tRemove BlackListed ORG and OU...\n")

	for ckl=0;ckl<len(conf.BlackList_OU);ckl++ {
		queryx=strings.Replace(SABDefine.PG_QUE_RemoveBlackListed	, "XYZDBOrgsXYZ", SABDefine.PG_Table_MSSQL[0], -1)
		queryx=strings.Replace(queryx								, "XYZDBDepsXYZ", SABDefine.PG_Table_MSSQL[1], -1)
		queryx=strings.Replace(queryx								, "XYZUidXYZ", conf.BlackList_OU[ckl], -1)

//		log.Printf("%s\n", queryx)
		_, err := db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() Que 11 error: %v\n", err)
			log.Printf("%s\n", queryx)
			return 11
		}
	}

	for ckl=0;ckl<len(SABDefine.PG_QUE_RemoveNoChildren);ckl+=2 {

		log.Printf("\t\tRemove NoChildren and NoParents from DB Cache. Pass %d of %d...\n", ckl/2+1, len(SABDefine.PG_QUE_RemoveNoChildren)/2)

		lasty=-1
		for {

			queryx=strings.Replace(SABDefine.PG_QUE_RemoveNoChildren[ckl]	, "XYZDBOrgsXYZ", SABDefine.PG_Table_MSSQL[0], -1)
			queryx=strings.Replace(queryx									, "XYZDBDepsXYZ", SABDefine.PG_Table_MSSQL[1], -1)
			queryx=strings.Replace(queryx									, "XYZDBPersXYZ", SABDefine.PG_Table_MSSQL[2], -1)
			queryx=strings.Replace(queryx									, "XYZDBPhonesXYZ", SABDefine.PG_Table_Oracle, -1)
//			log.Printf("%s\n", queryx)
			rows, err := db.Query(queryx)
			if err != nil {
				log.Printf("PG::Query() Que 12 error: %v\n", err)
				log.Printf("%s\n", queryx)
				return 12
			}

			rows.Next()
			rows.Scan(&lastx)
			log.Printf("\t%6d / %6d", lastx, lasty)

			if lastx<lasty || lasty==-1 {
				lasty=lastx

				queryx=strings.Replace(SABDefine.PG_QUE_RemoveNoChildren[ckl+1]	, "XYZDBOrgsXYZ", SABDefine.PG_Table_MSSQL[0], -1)
				queryx=strings.Replace(queryx									, "XYZDBDepsXYZ", SABDefine.PG_Table_MSSQL[1], -1)
				queryx=strings.Replace(queryx									, "XYZDBPersXYZ", SABDefine.PG_Table_MSSQL[2], -1)
				queryx=strings.Replace(queryx									, "XYZDBPhonesXYZ", SABDefine.PG_Table_Oracle, -1)

//				log.Printf("%s\n", queryx)
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("PG::Query() Que 13 error: %v\n", err)
					log.Printf("%s\n", queryx)
					return 13
				}
		
			}else{
				break
			}
		}

	}

	log.Printf("\tComplete")

	return 94

}
