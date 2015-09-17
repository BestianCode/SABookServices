package main

import (
	"fmt"
//	"os/exec"
	"log"
//	"strings"
//	"regexp"
	"net"
//	"time"

	"database/sql"

// SQLite3
	_ "github.com/mattn/go-sqlite3"

// PostgreSQL
	_ "github.com/lib/pq"

// Asterisk ARI
	"code.google.com/p/gami"

	"github.com/BestianRU/SABookServices/SABModules"
)

func main() {

	const (
		pName				=	string("SABook AsteriskCIDUpdater")
		pVer				=	string("3 2015.09.10.21.10")
	)

	var	(
		def_config_file		=	string ("./AsteriskCIDUpdater.json")			// Default configuration file
		def_log_file		=	string ("/var/log/ABook/AsteriskCIDUpdater.log")	// Default log file
		def_daemon_mode		=	string ("NO")						// Default start in foreground

		sqlite_key				string
		sqlite_value			string

		pg_name					string
		pg_phone				string

		pg_array				[100000][3]string
		sq_array				[100000][3]string

		pg_array_len		=	int(0)
		sq_array_len		=	int(0)

		ckl1				=	int(0)
		ckl2				=	int(0)
		ckl_status			=	int(0)

		ast_cmd					string

		rconf					SABModules.Config_STR

		sql_mode			=	int(0)

		query 					string
	)

	fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

	rconf.LOG_File = def_log_file

	def_config_file, def_daemon_mode = SABModules.ParseCommandLine(def_config_file, def_daemon_mode)

//	log.Printf("%s %s %s", def_config_file, def_daemon_mode, os.Args[0])

	SABModules.ReadConfigFile(def_config_file, &rconf)

	sqlite_select	:= fmt.Sprintf("SELECT key, value FROM astdb where key like '%%%s%%';", rconf.AST_CID_Group)
	pg_select		:= fmt.Sprintf("select x.cid_name, y.phone from ldapx_persons x, ldapx_phones y, (select a.phone, count(a.phone) as phone_count from ldapx_phones as a, ldapx_persons as b where a.pers_id=b.uid and b.contract=0 and a.pass=2 and b.lang=1 group by a.phone order by a.phone) as subq where x.uid=y.pers_id and y.pass=2 and x.lang=1 and subq.phone=y.phone and subq.phone_count<2 and y.phone like '%s%%' and x.contract=0 group by x.cid_name, y.phone order by y.phone;", rconf.AST_Num_Start)

	SABModules.Pid_Check(&rconf)
	SABModules.Pid_ON(&rconf)

	SABModules.Log_ON(&rconf)
	
	log.Printf(".")
	log.Printf("..")
	log.Printf("...")
	log.Printf("-> %s V%s", pName, pVer)
	log.Printf("--> Go!")

	db, err := sql.Open("sqlite3", rconf.AST_SQLite_DB)
	if err != nil {
		log.Printf("SQLite3::Open() error: %v\n", err)
		return
	}

	defer db.Close()

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return
	}

	defer dbpg.Close()

	dbast, err := net.Dial("tcp", fmt.Sprintf("%s:%d", rconf.AST_ARI_Host, rconf.AST_ARI_Port))
	if err != nil {
		log.Printf(".")
		log.Printf("Asterisk ARI::Dial() error: %v", err)
		log.Printf("\tWorking in SQL SQLite mode!")
		log.Printf(".")
		sql_mode=1
	}

	if sql_mode==0 {
 		defer dbast.Close()
 	}

	ast_gami := gami.NewAsterisk(&dbast, nil)
    ast_get := make(chan gami.Message, 10000)


	if sql_mode==0 {
		err = ast_gami.Login(rconf.AST_ARI_User, rconf.AST_ARI_Pass)
		if err != nil {
			log.Printf("Asterisk ARI::Login() error: %v\n", err)
			return
		}
	}

	rows, err := db.Query(sqlite_select)
	if err != nil {
		log.Printf("SQLite3::Query() error: %v\n", err)
		return
	}
		
	for rows.Next() {

		err = rows.Scan(&sqlite_key,&sqlite_value)
		if err != nil {
			log.Printf("rows.Scan error: %v\n", err)
			return
		}

		sq_array[sq_array_len][0]=sqlite_key
		sq_array[sq_array_len][1]=sqlite_value
		sq_array[sq_array_len][2]=SABModules.PhoneMutation(sqlite_key)
		sq_array_len++

	}

	rows, err = dbpg.Query(pg_select)
	if err != nil {
		log.Printf("PG::Query() error: %v\n", err)
		return
	}

	pg_array_len=0
	for rows.Next() {

		err = rows.Scan(&pg_name, &pg_phone)
		if err != nil {
			log.Printf("rows.Scan error: %v\n", err)
			return
		}

		pg_array[pg_array_len][0]=fmt.Sprintf("/%s/%s", rconf.AST_CID_Group, pg_phone)
		pg_array[pg_array_len][1]=pg_name
		pg_array[pg_array_len][2]=SABModules.PhoneMutation(pg_phone)
		pg_array_len++

	}

	for ckl1=0;ckl1<sq_array_len;ckl1++ {
		ckl_status=0
		for ckl2=0;ckl2<pg_array_len;ckl2++ {
			if sq_array[ckl1][0]==pg_array[ckl2][0] && sq_array[ckl1][1]==pg_array[ckl2][1] {
				ckl_status=1
				break
			}
		}
		if ckl_status==0 {
			if sql_mode==0 {
				ast_cmd = fmt.Sprintf("database del %s %s", rconf.AST_CID_Group, sq_array[ckl1][2])
				log.Printf ("\t- %s\n", ast_cmd)

				ast_cb := func(m gami.Message) {
					ast_get <- m
				}

				err = ast_gami.Command(ast_cmd, &ast_cb)
				if err != nil {
					log.Printf("Asterisk ARI::Command() error: %v\n", err)
					return
				}
			
				for x1, x2 := range <- ast_get {
					if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
						log.Printf("\t\t\t%s\n", x2)
					}
				}
			}else{
				query = fmt.Sprintf("delete from astdb where key='%s';", sq_array[ckl1][0])
				log.Printf("\t- %s\n", query)
				_, err := db.Exec(query)
				if err != nil {
					log.Printf("SQLite3::Query() DEL error: %v\n", err)
					return
				}
			}
		}
	}

	for ckl1=0;ckl1<pg_array_len;ckl1++ {
		ckl_status=0
		for ckl2=0;ckl2<sq_array_len;ckl2++ {
			if pg_array[ckl1][0]==sq_array[ckl2][0] && pg_array[ckl1][1]==sq_array[ckl2][1] {
				ckl_status=1
				break
			}
		}
		if ckl_status==0 {

			if sql_mode==0 {
				ast_cmd = fmt.Sprintf("database del %s %s", rconf.AST_CID_Group, pg_array[ckl1][2])
				log.Printf ("\t- %s\n", ast_cmd)

				ast_cb := func(m gami.Message) {
					ast_get <- m
				}

				err = ast_gami.Command(ast_cmd, &ast_cb)
				if err != nil {
					log.Printf("Asterisk ARI::Command() error: %v\n", err)
					return
				}

				for x1, x2 := range <- ast_get {
					if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
						log.Printf("\t\t\t%s\n", x2)
					}
				}

				ast_cmd = fmt.Sprintf("database put %s %s \"%s\"", rconf.AST_CID_Group, pg_array[ckl1][2], pg_array[ckl1][1])
				log.Printf ("\t+ %s\n", ast_cmd)

				ast_cb = func(m gami.Message) {
					ast_get <- m
				}

				err = ast_gami.Command(ast_cmd, &ast_cb)
				if err != nil {
					log.Printf("Asterisk ARI::Command() error: %v\n", err)
					return
				}
			
				for x1, x2 := range <- ast_get {
					if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
						log.Printf("\t\t\t%s\n", x2)
					}
				}
			}else{
				query = fmt.Sprintf("delete from astdb where key='%s';",pg_array[ckl1][0])
				log.Printf("\t- %s\n", query)
				_, err := db.Exec(query)
				if err != nil {
					log.Printf("SQLite3::Query() DEL error: %v\n", err)
					return
				}

				query = fmt.Sprintf("insert into astdb (key,value) values ('%s','%s');",pg_array[ckl1][0], pg_array[ckl1][1])
				log.Printf("\t+ %s\n", query)
				_, err = db.Exec(query)
				if err != nil {
					log.Printf("SQLite3::Query() INS error: %v\n", err)
					return
				}
			}
		}
	}


	log.Printf("...")
	log.Printf("..")
	log.Printf(".")

	SABModules.Log_OFF()

	SABModules.Pid_OFF(&rconf)

}
