package SABFunctions

import (
	"fmt"
	"log"
	"strings"
	"time"


// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

//LDAP
	"github.com/go-ldap/ldap"

// Translit
	"github.com/BestianRU/gounidecode/unidecode"

	"github.com/BestianRU/SABookServices/SABModules"
	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func LDAP_to_PG(conf *SABModules.Config_STR, pg_minsert int) int {

	var (

		ckl_servers		int
		num_servers		int

		fName			[]string
		fCN			[]string
		fOUa			[]string
		fOU			string
		fMail			string

		pfName			string
		pfCN			string
		pfMail			string

		ckl		=	int	(0)
		state		=	int	(0)

		queryx			string

		pg_Domino_Create=	string(`
						CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
							(server character varying(255),
							namerus character varying(255), trnamerus character varying(255), namelat character varying(255),
							ou character varying(255), mail character varying(255),
								primary key (mail));
						`)

		pg_Query_Create_Status	=	string(`
							CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
								(server character varying(255), status character varying(255),
									primary key (server));
						`)

		return_result		=	int(0)

	)

	log.Printf("LDAP Export to PG...")

	num_servers = len(conf.LDAP_URL)

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	queryx = strings.Replace(pg_Domino_Create, "XYZWorkTableZYX", SABDefine.PG_Table_Domino, -1)
	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 11
	}

	queryx = strings.Replace(pg_Query_Create_Status, "XYZWorkTableZYX", SABDefine.PG_Table_Domino_Status, -1)
	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("%s\n", queryx)
		log.Printf("PG::Query() Create table error: %v\n", err)
		return 12
	}


	for ckl_servers=0; ckl_servers<num_servers; ckl_servers++{

		log.Printf("\t\tServer %2d of %2d / Pass  1 of  1 / Server name: %s\n", ckl_servers+1, num_servers, conf.LDAP_URL[ckl_servers][0])

		l, err := ldap.Dial("tcp", conf.LDAP_URL[ckl_servers][0])
		if err != nil {
			log.Printf("LDAP::Initialize() error: %v\n", err)
			continue
		}

		defer l.Close()
//		l.Debug = true

		err = l.Bind(conf.LDAP_URL[ckl_servers][1], conf.LDAP_URL[ckl_servers][2])
		if err != nil {
			log.Printf("LDAP::Bind() error: %v\n", err)
			continue
		}

		search := ldap.NewSearchRequest(conf.LDAP_URL[ckl_servers][3], ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, conf.LDAP_URL[ckl_servers][4], SABDefine.LDAP_attr, nil)

		sr, err := l.Search(search)
		if err != nil {
			log.Printf("LDAP::Search() error: %v\n", err)
			continue
		}

		log.Printf("\t\t\t%s // %d\n", search.Filter, len(sr.Entries))

		if len(sr.Entries)>10 {
			timenow:=time.Now().Format("2006.01.02 15:04:05")

			queryx = fmt.Sprintf("INSERT INTO %s (server, status) select '%s', '%s' where not exists (select server from %s where server='%s'); update %s set status='%s' where server='%s'; ", SABDefine.PG_Table_Domino_Status, conf.LDAP_URL[ckl_servers][0], timenow, SABDefine.PG_Table_Domino_Status, conf.LDAP_URL[ckl_servers][0], SABDefine.PG_Table_Domino_Status, timenow, conf.LDAP_URL[ckl_servers][0])
//			log.Printf("%s\n", queryx)
			_, err = db.Query(queryx)
			if err != nil {
				log.Printf("%s\n", queryx)
				log.Printf("PG::Query() Create table error: %v\n", err)
				return 14
			}

			_, err = db.Query(fmt.Sprintf("delete from %s where server='%s';", SABDefine.PG_Table_Domino, conf.LDAP_URL[ckl_servers][0]))
			if err != nil {
				log.Printf("PG::Query() Clean table error: %v\n", err)
				return 14
			}

		}

		ckl=0

		for _, entry := range sr.Entries {

			if ckl<1 {
				queryx = ""
			}


			for _, attr := range entry.Attributes {
				if attr.Name == "altfullname" {
					x  := strings.Join(attr.Values, ",")
					fOUa = strings.Split(x, ",")
					fName = strings.Split(fOUa[0], "=")
				}
				if attr.Name == "cn" {
					x  := strings.Join(attr.Values, ",")
					fCN = strings.Split(x, ",")
				}
				if attr.Name == "mail" {
					fMail = strings.Join(attr.Values, ",")
				}
			}

			if len(fName)>0 && len(fCN)>0 {

				return_result=94

				if ckl>0 && state == 1{
					queryx = fmt.Sprintf("%s , ", queryx)
				}
				fOU=""
				for ckl1:=int(len(fOUa)-2);ckl1>0;ckl1-- {
					fOU=fmt.Sprintf("%s/%s", fOU, fOUa[ckl1])
				}
				fOU=strings.Trim(strings.Trim(strings.Replace(strings.Replace(fOU, "OU=", "", -1), "O=", "", -1), "/"), " ")
				if fOU == "" {
					fOU=conf.ROOT_OU
				}
				pfName=SABModules.TextMutation(strings.Replace(fName[int(len(fName)-1)], "'", "", -1))
				pfCN=SABModules.TextMutation(strings.Replace(fCN[0], "'", "", -1))
				pfMail=strings.Replace(strings.ToLower(fMail), " ", "", -1)
				queryx = fmt.Sprintf("%sINSERT INTO %s (server, namerus, trnamerus,  namelat, ou, mail) select '%s','%s','%s','%s','%s','%s' where not exists (select mail from %s where mail='%s'); ", queryx, SABDefine.PG_Table_Domino, conf.LDAP_URL[ckl_servers][0], pfName, unidecode.Unidecode(pfName), pfCN, fOU, pfMail, SABDefine.PG_Table_Domino, pfMail)
//				log.Printf("%s", queryx)
//				state = 1
			}

			if ckl>=pg_minsert-1 {
//				log.Printf("%s\n\n", queryx)
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
				}
				queryx=""
				ckl=0
//				state=0
			}else{
				ckl++
			}
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

