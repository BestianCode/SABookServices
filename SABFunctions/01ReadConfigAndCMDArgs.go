package SABFunctions

import (
	"log"
//	"fmt"
	"os"
	"flag"
//	"time"

	"encoding/json"

// PostgreSQL
//	"database/sql"
//	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABDefine"

)

func ParseCommandLine(config_file string) string {
	ConfigPtr := flag.String("config", config_file, "Path to Configuration file")
	flag.Parse()
	config_file=*ConfigPtr
	log.Printf("path to Configuration file: %s", config_file)
	return config_file
}


func ReadConfigFile (config_file string){
	conf_file, err := os.Open(config_file)
	if err != nil {
		log.Fatalf("Error open Configuration file %s %v\n", config_file, err)
	}

	conf_decoder := json.NewDecoder(conf_file)
	err = conf_decoder.Decode(&SABDefine.Conf)
	if err != nil {
		log.Fatalf("Error read SABDefine.Configuration file %s %v\n", config_file, err)
	}

	conf_file.Close()
}

/*func CheckState(conf *SABDefine.Config_STR, message string) {

	var	cs_table	=	string ("CheckState")

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
//		return 11
	}

	defer db.Close()

	que := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (dtime character varying(255), state character varying(15));", cs_table)

	_, err = db.Exec(que)
	if err != nil {
		log.Printf("PG::Query() Create table error: %v\n", err)
//		return 12
	}
	que = fmt.Sprintf("INSERT INTO %s (dtime,state) VALUES ('%s','%s');", cs_table, time.Now(), message)
//	log.Printf("%s\n", que)

	_, err = db.Exec(que)
	if err != nil {
		log.Printf("PG::Query() Query executing error: %v\n", err)
//		return 13
	}

//	return 94

}

*/