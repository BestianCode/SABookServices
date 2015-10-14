package SABFunctions

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	//LDAP
	"github.com/go-ldap/ldap"

	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
	"github.com/BestianRU/SABookServices/SABModules"
)

func AD_to_PG(conf *SABModules.Config_STR, pg_minsert int) int {
	var (
		ckl_servers int
		num_servers int
		ckl         = int(0)
		queryx      string

		pg_AD_Create = string(`
			CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
				(domain character varying(255), server character varying(255),
					displayname character varying(255), cn character varying(255),
					dlogin character varying(255), login character varying(255),
					mail character varying(255),
					ph_int character varying(255),
					ph_mob character varying(255),
					ph_ip character varying(255),
					department character varying(255), title character varying(255),
					dn character varying(255),
					connected character varying(5) NOT NULL DEFAULT 'no'::character varying,
					primary key (dlogin));
			`)

		pg_AD_Create_Status = string(`
			CREATE TABLE IF NOT EXISTS XYZWorkTableZYX
				(server character varying(255), status character varying(255),
					primary key (server));
			`)
		return_result = int(0)
	)

	log.Printf("AD Export to PG...")

	rusFindRegExp := regexp.MustCompile(`[А-Яа-я]`)

	num_servers = len(conf.AD_LDAP)

	//	fmt.Printf("%d\n", num_servers)

	for ckl = 0; ckl < num_servers; ckl++ {
		conf.AD_LDAP[ckl][6] = "enabled"
	}

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	queryx = strings.Replace(pg_AD_Create, "XYZWorkTableZYX", SABDefine.PG_Table_AD, -1)
	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Exec() error: %v\n", err)
		return 11
	}
	queryx = strings.Replace(pg_AD_Create_Status, "XYZWorkTableZYX", SABDefine.PG_Table_AD_Status, -1)
	_, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Exec() error: %v\n", err)
		return 12
	}

	for ckl_servers = 0; ckl_servers < num_servers; ckl_servers++ {

		if conf.AD_LDAP[ckl_servers][6] != "enabled" {
			continue
		}

		log.Printf("\t\tServer %2d of %2d / Pass  1 of  1 / Domain: %s, Controller: %s\n", ckl_servers+1, num_servers, conf.AD_LDAP[ckl_servers][0], conf.AD_LDAP[ckl_servers][1])

		l, err := ldap.Dial("tcp", conf.AD_LDAP[ckl_servers][1])
		if err != nil {
			log.Printf("LDAP::Initialize() error: %v\n", err)
			continue
		}

		defer l.Close()
		//		l.Debug = true

		err = l.Bind(conf.AD_LDAP[ckl_servers][2], conf.AD_LDAP[ckl_servers][3])
		if err != nil {
			log.Printf("LDAP::Bind() error: %v\n", err)
			continue
		}

		search := ldap.NewSearchRequest(conf.AD_LDAP[ckl_servers][4], ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, conf.AD_LDAP[ckl_servers][5], SABDefine.AD_attr, nil)

		sr, err := l.Search(search)
		if err != nil {
			log.Printf("LDAP::Search() error: %v\n", err)
			continue
		}

		log.Printf("\t\t\t%s // %d\n", search.Filter, len(sr.Entries))

		if len(sr.Entries) > 10 {
			timenow := time.Now().Format("2006.01.02 15:04:05")

			queryx = fmt.Sprintf("INSERT INTO %s (server, status) select '%s', '%s' where not exists (select server from %s where server='%s'); update %s set status='%s' where server='%s'; ", SABDefine.PG_Table_AD_Status, conf.AD_LDAP[ckl_servers][1], timenow, SABDefine.PG_Table_AD_Status, conf.AD_LDAP[ckl_servers][1], SABDefine.PG_Table_AD_Status, timenow, conf.AD_LDAP[ckl_servers][1])
			//			log.Printf("%s\n", queryx)
			_, err = db.Query(queryx)
			if err != nil {
				log.Printf("%s\n", queryx)
				log.Printf("PG::Query() Insert error: %v\n", err)
				return 13
			}

			_, err = db.Query(fmt.Sprintf("delete from %s where server='%s';", SABDefine.PG_Table_AD, conf.AD_LDAP[ckl_servers][1]))
			if err != nil {
				log.Printf("PG::Query() Clean table error: %v\n", err)
				return 14
			}

		}

		ckl = 0
		for _, entry := range sr.Entries {
			if ckl < 1 {
				queryx = ""
			}
			fName := ""
			fCN := ""
			fLogin := ""
			fMail := ""
			fDep := ""
			fTitle := ""
			fDN := ""
			fPhone := ""
			fMobile := ""
			fIPPhone := ""

			for _, attr := range entry.Attributes {
				if attr.Name == "displayName" {
					fName = strings.Join(attr.Values, ",")
				}
				if attr.Name == "cn" {
					fCN = strings.Join(attr.Values, ",")
				}
				if attr.Name == "sAMAccountName" {
					fLogin = strings.Join(attr.Values, ",")
				}
				if attr.Name == "mail" {
					fMail = strings.Join(attr.Values, ",")
				}
				if attr.Name == "department" {
					fDep = strings.Join(attr.Values, ",")
				}
				if attr.Name == "title" {
					fTitle = strings.Join(attr.Values, ",")
				}
				if attr.Name == "distinguishedName" {
					fDN = strings.Join(attr.Values, ",")
				}
				if attr.Name == "telephoneNumber" {
					fPhone = strings.Join(attr.Values, ",")
				}
				if attr.Name == "mobile" {
					fMobile = strings.Join(attr.Values, ",")
				}
				if attr.Name == "pager" {
					fIPPhone = strings.Join(attr.Values, ",")
				}
			}
			if len(fName) > 0 && len(fLogin) > 0 && rusFindRegExp.FindString(fName) != "" {
				return_result = 94
				/*
					if ckl > 0 && state == 1 {
						queryx = fmt.Sprintf("%s , ", queryx)
					}

						fOU = ""
						for ckl1 := int(len(fOUa) - 2); ckl1 > 0; ckl1-- {
							fOU = fmt.Sprintf("%s/%s", fOU, fOUa[ckl1])
						}
						fOU = strings.Trim(strings.Trim(strings.Replace(strings.Replace(fOU, "OU=", "", -1), "O=", "", -1), "/"), " ")
						if fOU == "" {
							fOU = conf.ROOT_OU
						}
						pfName = SABModules.TextMutation(strings.Replace(fName[int(len(fName)-1)], "'", "", -1))
						pfCN = SABModules.TextMutation(strings.Replace(fCN[0], "'", "", -1))
						pfMail = strings.Replace(strings.ToLower(fMail), " ", "", -1)
						queryx = fmt.Sprintf("%sINSERT INTO %s (server, namerus, trnamerus,  namelat, ou, mail) select '%s','%s','%s','%s','%s','%s' where not exists (select mail from %s where mail='%s'); ", queryx, SABDefine.PG_Table_Domino, conf.AD_LDAP[ckl_servers][0], pfName, unidecode.Unidecode(pfName), pfCN, fOU, pfMail, SABDefine.PG_Table_Domino, pfMail)
						//log.Printf("%s", queryx)
						//state = 1
				*/
				queryx = fmt.Sprintf("%s INSERT INTO %s (domain,server,displayname,cn,login,dlogin,mail,department,title,dn,ph_int,ph_mob,ph_ip) values ('%s','%s','%s','%s','%s','%s@%s','%s','%s','%s','%s','%s','%s','%s');", queryx, SABDefine.PG_Table_AD, conf.AD_LDAP[ckl_servers][0], conf.AD_LDAP[ckl_servers][1], SABModules.NameMutation(fName), SABModules.NameMutation(fCN), SABModules.NameMutation(fLogin), SABModules.NameMutation(fLogin), conf.AD_LDAP[ckl_servers][0], fMail, SABModules.NameMutation(fDep), SABModules.NameMutation(fTitle), fDN, fPhone, fMobile, fIPPhone)
				//log.Printf("----- %s\n", queryx)
				//				fmt.Printf("%s/%s/%s/%s/%s/%s/%s/%s\n", conf.AD_LDAP[ckl_servers][0], conf.AD_LDAP[ckl_servers][1], fName, fCN, fLogin, fMail, fDep, fTitle)
			}
			if ckl >= pg_minsert-1 {
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
				}
				queryx = ""
				ckl = 0
			} else {
				ckl++
			}
			/*

					if ckl >= pg_minsert-1 {
						//				log.Printf("%s\n\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
						}
						queryx = ""
						ckl = 0
						//				state=0
					} else {
						ckl++
					}
				}
				//		log.Printf("%s\n\n", queryx)
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Insert into table error: %v\n", err)
			*/
		}

		if len(sr.Entries) > 1 {
			for ckl = 0; ckl < num_servers; ckl++ {
				if conf.AD_LDAP[ckl][0] == conf.AD_LDAP[ckl_servers][0] {
					conf.AD_LDAP[ckl][6] = "ok"
					if ckl != ckl_servers {
						queryx = fmt.Sprintf("INSERT INTO %s (server, status) select '%s', '%s' where not exists (select server from %s where server='%s'); update %s set status='%s' where server='%s'; ", SABDefine.PG_Table_AD_Status, conf.AD_LDAP[ckl][1], "SKIP", SABDefine.PG_Table_AD_Status, conf.AD_LDAP[ckl][1], SABDefine.PG_Table_AD_Status, "SKIP", conf.AD_LDAP[ckl][1])
						//			log.Printf("%s\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("PG::Query() Create table error: %v\n", err)
							return 19
						}
					}
				}
			}
		}

	}

	log.Printf("\tComplete")

	return return_result

}
