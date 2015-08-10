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

	que := fmt.Sprintf("%s", strings.Replace(SABDefine.PG_QUE, "XYZZYX", SABDefine.AsteriskCIDTableTemp, -1));
//	log.Printf("%s\n", que)

	_, err = db.Exec(que)
	if err != nil {
		log.Printf("PG::Query() Query executing error: %v\n", err)
		return 12
	}


	log.Printf("MakeAsteriskCIDTable FINISHED\n")

	return 94

}
