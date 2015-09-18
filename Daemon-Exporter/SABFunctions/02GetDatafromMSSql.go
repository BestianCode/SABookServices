package SABFunctions

import (
	"fmt"
	"log"
	"strings"
	"time"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

// MS SQL
	_ "github.com/denisenkom/go-mssqldb"

// Translit
	"github.com/BestianRU/SABookServices/ForeignModules/unidecode"

	"github.com/BestianRU/SABookServices/SABModules"
	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func MSSQL_to_PG(conf *SABModules.Config_STR, pg_minsert int) int {

	var	(
		ckl				int

		ckl_servers			int
		num_servers			int

		ckl_que				int
		num_que				int

		queryx				string

		pg_Query_Create		=	[]string{`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), uid bytea, idparent bytea, name character varying(255), nametr character varying(255),
									primary key (uid));
						`,`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), uid bytea, idparent bytea, name character varying(255), nametr character varying(255), idorg bytea,
									primary key (uid));
						`,`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), uid bytea, idparent bytea,
									nfr character varying(255), nfir character varying(5), nlr character varying(255), nmr character varying(255), nmir character varying(5),
									nft character varying(255), nfit character varying(5), nlt character varying(255), nmt character varying(255), nmit character varying(5),
									tab character varying(255), pos character varying(255), idorg bytea, contract integer NOT NULL,
									primary key (uid));
						`}

		pg_Query_Create_Status	=	string(`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), status character varying(255),
									primary key (server));
						`)

		return_result		=	int(0)

	)

	log.Printf("MS SQL export to PG...")

	num_servers	= len(conf.MSSQL_DSN)
	num_que		= len(SABDefine.MSSQL_QUE)

	dbpg, err := sql.Open("postgres", conf.PG_DSN)

	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}
	defer dbpg.Close()

	queryx = strings.Replace(pg_Query_Create_Status, "XYZWorkTableZYX", SABDefine.PG_Table_MSSQL_Status, -1)
//	log.Printf("%s\n", queryx)
	_, err = dbpg.Query(queryx)
	if err != nil {
		log.Printf("%s\n", queryx)
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 11
	}

	for ckl_servers=0; ckl_servers<num_servers; ckl_servers++{

		db, err := sql.Open("mssql", conf.MSSQL_DSN[ckl_servers][0])
		if err != nil {
			log.Printf("MS SQL::Open() error: %v\n", err)
			continue
		}else{
			defer db.Close()


			for ckl_que=0; ckl_que<num_que; ckl_que++{

				log.Printf("\t\tServer %2d of %2d / Pass %2d of %2d / Server name: %s\n", ckl_servers+1, num_servers, ckl_que+1, num_que, conf.MSSQL_DSN[ckl_servers][1])

				queryx = strings.Replace(pg_Query_Create[ckl_que], "XYZWorkTableZYX", SABDefine.PG_Table_MSSQL[ckl_que], -1)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Create table error: %v\n", err)
					return 12
				}

				rows, err := db.Query(SABDefine.MSSQL_QUE[ckl_que])
				if err != nil {
					log.Printf("MS SQL::Query() error: %v\n", err)
					break
				}else{

					timenow:=time.Now().Format("2006.01.02 15:04:05")

					queryx = fmt.Sprintf("INSERT INTO %s (server, status) select '%s', '%s' where not exists (select server from %s where server='%s'); update %s set status='%s' where server='%s'; ", SABDefine.PG_Table_MSSQL_Status, conf.MSSQL_DSN[ckl_servers][1], timenow, SABDefine.PG_Table_MSSQL_Status, conf.MSSQL_DSN[ckl_servers][1], SABDefine.PG_Table_MSSQL_Status, timenow, conf.MSSQL_DSN[ckl_servers][1])
//					log.Printf("%s\n", queryx)
					_, err = dbpg.Query(queryx)
					if err != nil {
						log.Printf("%s\n", queryx)
						log.Printf("PG::Query() Create table error: %v\n", err)
						return 13
					}

					_, err = dbpg.Query(fmt.Sprintf("delete from %s where server='%s';", SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1]))
					if err != nil {
						log.Printf("PG::Query() Clean table error: %v\n", err)
						return 14
					}

				}

				ckl=0

				for rows.Next() {

					if ckl<1 {
						queryx = ""
					}

					return_result=94

					switch ckl_que+1 {
						case 1:
							var	(
								xid		[]byte
								xname		string
								xidparent	[]byte
								xnameq		string
							)
							rows.Scan(&xid,&xname,&xidparent)
							xname	 = SABModules.TextMutation(xname)
							xnameq	 = SABModules.TransMutation(xname, conf)
							xnametr	:= unidecode.Unidecode(xnameq)
							queryx = fmt.Sprintf("%sINSERT INTO %s (server, uid, idparent, name, nametr) select '%s', '%v', '%v','%s','%s' where not exists (select uid from %s where uid='%v'); ", queryx, SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1],xid,xidparent,xname,xnametr,SABDefine.PG_Table_MSSQL[ckl_que],xid)
						case 2:
							var	(
								xid		[]byte
								xidorg		[]byte
								xidparent	[]byte
								xname		string
								xnameq		string
							)
							rows.Scan(&xid,&xidorg,&xidparent,&xname)
							xname	 = SABModules.TextMutation(xname)
							xnameq	 = SABModules.TransMutation(xname, conf)
							xnametr := unidecode.Unidecode(xnameq)
							if fmt.Sprintf("%v", xidparent) == "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]" {
								queryx = fmt.Sprintf("%sINSERT INTO %s (server, uid, idparent, name, nametr, idorg) VALUES ('%s', '%v', '%v','%s','%s','%v'); ", queryx, SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1],xid,xidorg,xname,xnametr,xidorg)
							}else{
								queryx = fmt.Sprintf("%sINSERT INTO %s (server, uid, idparent, name, nametr, idorg) VALUES ('%s', '%v', '%v','%s','%s','%v'); ", queryx, SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1],xid,xidparent,xname,xnametr,xidorg)
							}
						case 3:
							var	(
								xid		[]byte
								xname		string
								xfio		[]string
								xfiotr		[]string
								xtab		string
								xidorg		[]byte
								xidparent	[]byte
								xpos		string
								xcontract	int
							)
							rows.Scan(&xid,&xname,&xtab,&xidorg,&xidparent,&xpos,&xcontract)
							xfio=SABModules.PeopleMutation(xname, "RUS")
							xfiotr=SABModules.PeopleMutation(unidecode.Unidecode(xname), "LAT")
							xtab=SABModules.TextMutation(xtab)
							xpos=SABModules.TextMutation(xpos)
							if len(xfio)<3 {
								queryx = fmt.Sprintf("%sINSERT INTO %s (server, uid, idparent, nlr, nfr, nfir, nmr, nmir, nlt, nft, nfit, nmt, nmit, tab, pos, idorg, contract) VALUES  ('%s', '%v', '%v','%s','%s','%s','','','%s','%s','%s','','','%s','%s','%v',%d); ", queryx, SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1], xid, xidparent, xfio[0],xfio[1],xfio[1][:2],xfiotr[0],xfiotr[1],xfiotr[1][:1],xtab, xpos, xidorg, xcontract)
							}else{
								queryx = fmt.Sprintf("%sINSERT INTO %s (server, uid, idparent, nlr, nfr, nfir, nmr, nmir, nlt, nft, nfit, nmt, nmit, tab, pos, idorg, contract) VALUES  ('%s', '%v', '%v','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%v',%d); ", queryx, SABDefine.PG_Table_MSSQL[ckl_que], conf.MSSQL_DSN[ckl_servers][1], xid, xidparent, xfio[0],xfio[1],xfio[1][:2],xfio[2],xfio[2][:2], xfiotr[0],xfiotr[1],xfiotr[1][:1],xfiotr[2],xfiotr[2][:1], xtab, xpos, xidorg, xcontract)
							}
						default:
							break
					}

					if ckl>=pg_minsert-1 {
//						log.Printf("%s\n\n", queryx)
						_, err = dbpg.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("PG::Query() Insert into table error: %v\n", err)
						}
						queryx=""
						ckl=0
					}else{
						ckl++
					}
				}

//				log.Printf("%s\n\n", queryx)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Insert into table error: %v\n", err)
				}
			}
		}
	}

	log.Printf("\tComplete")

	return return_result

}


