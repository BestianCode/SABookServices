package SABFunctions

import (
	"fmt"
	"log"
	"strings"
	"time"
//	"regexp"

// Oracle
	"gopkg.in/goracle.v1/oracle"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"

	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func Oracle_to_PG(conf *SABModules.Config_STR, pg_minsert int) int {

	var (

		ckl_servers		int
		num_servers		int

		ckl		=	int	(0)

		insert_exec	=	int	(0)

		row_comment		string
		row_tm			string
		row_number		string
		row_fname		string
		queryx			string

		pg_Query_Create=	string(`
						CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
							(
							server character varying(255),
							uid bytea, phone character varying(255),
							comment character varying(255),
							tm character varying(5), visible character varying(5), type integer, fname character varying(255));
						`)

		pg_Query_Create_Status	=	string(`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), status character varying(255),
									primary key (server));
						`)

		return_result	=	int(0)

	)

	log.Printf("Oracle export to PG...")

	num_servers = len(conf.Oracle_SRV)

	db, err := sql.Open("postgres", conf.PG_DSN)

	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	queryx = strings.Replace(pg_Query_Create_Status, "XYZWorkTableZYX", SABDefine.PG_Table_Oracle_Status, -1)

	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 11
	}

	queryx = strings.Replace(pg_Query_Create, "XYZWorkTableZYX", SABDefine.PG_Table_Oracle, -1)

	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 12
	}


	for ckl_servers=0; ckl_servers<num_servers; ckl_servers++{

		log.Printf("\t\tServer %2d of %2d / Pass  1 of  1 / Server name: %s\n", ckl_servers+1, num_servers, conf.Oracle_SRV[ckl_servers][3])

		cx, err := oracle.NewConnection(conf.Oracle_SRV[ckl_servers][0], conf.Oracle_SRV[ckl_servers][1], conf.Oracle_SRV[ckl_servers][2], false)

		if err != nil {
			log.Printf("Oracle::Connection() error: %v\n", err)
			continue
		}

		defer cx.Close()
		cu := cx.NewCursor()
		defer cu.Close()

		err = cu.Execute(SABDefine.Oracle_QUE, nil, nil)

		if err != nil {
			log.Printf("Oracle::Execute() error: %v\n", err)
			continue
		}

		rows, err := cu.FetchMany(pg_minsert)

		queryx = fmt.Sprintf("delete from %s where server='%s';", SABDefine.PG_Table_Oracle, conf.Oracle_SRV[ckl_servers][3])
//		log.Printf("%s\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() Clean table error: %v\n", err)
			return 13
		}

		timenow:=time.Now().Format("2006.01.02 15:04:05")

		queryx = fmt.Sprintf("INSERT INTO %s (server, status) select '%s', '%s' where not exists (select server from %s where server='%s'); update %s set status='%s' where server='%s'; ", SABDefine.PG_Table_Oracle_Status, conf.Oracle_SRV[ckl_servers][3], timenow, SABDefine.PG_Table_Oracle_Status, conf.Oracle_SRV[ckl_servers][3], SABDefine.PG_Table_Oracle_Status, timenow, conf.Oracle_SRV[ckl_servers][3])
//		log.Printf("%s\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Create table error: %v\n", err)
			return 14
		}

		ckl=0

		for err == nil && len(rows) > 0 {
			for _, row := range rows {

				if ckl<1 {
					queryx=""
				}

				if fmt.Sprintf("%s", row[4]) != "%!s(<nil>)" {
					row_comment=fmt.Sprintf("%s", row[4])
				}else{
					row_comment=""
				}

				if fmt.Sprintf("%s", row[5]) != "%!s(<nil>)" {
					row_tm=fmt.Sprintf("%s", row[5])
				}else{
					row_tm="X"
				}

				if fmt.Sprintf("%s", row[8]) != "%!s(<nil>)" {
					row_fname=SABModules.TextMutation(fmt.Sprintf("%s", row[8]))
				}else{
					row_fname=""
				}

				if fmt.Sprintf("%s", row[1]) != "%!s(<nil>)" && fmt.Sprintf("%s", row[2]) != "%!s(<nil>)"{
					return_result=94
					switch fmt.Sprintf("%d", row[7]) {
						case "1":
							row[1]=strings.Replace(fmt.Sprintf("%s", row[1]), "+7", "", -1)
							row_number=strings.Replace(fmt.Sprintf("+7%s%s", row[1], row[2]), " ", "", -1)
							row_number=SABModules.PhoneMutation(row_number)
						case "2":
							row[1]=strings.Replace(fmt.Sprintf("%s", row[1]), "+7", "", -1)
							row_number=strings.Replace(fmt.Sprintf("8%s%s", row[1], row[2]), " ", "", -1)
							row_number=SABModules.PhoneMutation(row_number)
						case "3":
							if fmt.Sprintf("%s", row[3]) != "%!s(<nil>)" {
								row[1]=strings.Replace(fmt.Sprintf("%s", row[1]), "+7", "", -1)
								row[1]=SABModules.PhoneMutation(fmt.Sprintf("%s", row[1]))
								row[2]=SABModules.PhoneMutation(fmt.Sprintf("%s", row[2]))
								row[3]=SABModules.PhoneMutation(fmt.Sprintf("%s", row[3]))
								row_number=strings.Replace(fmt.Sprintf("8(%s)%sдоб.%s", row[1], row[2], row[3]), " ", "", -1)
								row_number=strings.Replace(row_number, "доб.", " доб.", -1)
							}else{
								row[1]=strings.Replace(fmt.Sprintf("%s", row[1]), "+7", "", -1)
								row[1]=SABModules.PhoneMutation(fmt.Sprintf("%s", row[1]))
								row[2]=SABModules.PhoneMutation(fmt.Sprintf("%s", row[2]))
								row_number=strings.Replace(fmt.Sprintf("8(%s)%s", row[1], row[2]), " ", "", -1)
							}
						default:
							insert_exec=1
					}
				}else{
					insert_exec=1
				}
				if insert_exec==0 {
					queryx=fmt.Sprintf("%sinsert into %s (server, uid, phone, comment, tm, visible, type, fname) select '%s','%v','%s','%s','%s','%s','%d','%s' where not exists (select uid from %s where uid='%v' and phone='%s'); ", queryx, SABDefine.PG_Table_Oracle, conf.Oracle_SRV[ckl_servers][3], row[0], row_number, row_comment, row_tm, row[6], row[7], row_fname, SABDefine.PG_Table_Oracle, row[0], row_number)
//					fmt.Printf("%s\n", queryx)
					if ckl>=pg_minsert-1 {
//						log.Printf("%s\n\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
						}
						queryx=""
						ckl=0
					}else{
						ckl++
					}
				}
				insert_exec=0
			}
			rows, err = cu.FetchMany(pg_minsert)
		}
//		log.Printf("%s\n\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Insert into table error: %v\n", err)
		}
	}

	log.Printf("\tComplete")

	return return_result

}

