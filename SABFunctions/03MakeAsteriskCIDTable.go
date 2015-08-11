package SABFunctions

import (
	"fmt"
	"log"
	"strings"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABDefine"
)

func MakeAsteriskCIDTable(conf *SABDefine.Config_STR) int {

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 11
	}

	defer db.Close()
	log.Printf("MakeAsteriskCIDTable GO...\n")

	que := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (name character varying(255), number character varying(15));", SABDefine.AsteriskCIDTable)

	_, err = db.Exec(que)
	if err != nil {
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 12
	}
	que = fmt.Sprintf("%s", strings.Replace(strings.Replace(strings.Replace(SABDefine.PG_QUE, "XYZTempTableZYX", SABDefine.AsteriskCIDTableTemp, -1), "XYZWorkTableZYX", SABDefine.AsteriskCIDTable, -1), "XYZOraclePersTableZYX", SABDefine.PG_Table_Oracle[2], -1));
//	log.Printf("%s\n", que)

	_, err = db.Exec(que)
	if err != nil {
		log.Printf("PG::Query() Query executing error: %v\n", err)
		return 13
	}


	log.Printf("MakeAsteriskCIDTable FINISHED\n")

	return 94

}

