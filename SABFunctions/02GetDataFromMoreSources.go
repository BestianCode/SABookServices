package SABFunctions

import (
	"fmt"
	"log"
	"strings"

//LDAP
	"gopkg.in/ldap.v1"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

// Oracle
	"gopkg.in/goracle.v1/oracle"

	"github.com/BestianRU/SABookServices/SABDefine"
)

func LDAP_to_PG(conf *SABDefine.Config_STR) int {

	var (
		fName			[]string
		fCN			[]string
		fOUa			[]string
		fOU			string
		fMail			string
		queryx			string

		ckl		=	int	(0)
		state		=	int	(0)
	)


	l, err := ldap.Dial("tcp", conf.LDAP_URL)
	if err != nil {
		log.Printf("LDAP::Initialize() error: %v\n", err)
		return 21
	}

	defer l.Close()
	// l.Debug = true

	err = l.Bind(conf.LDAP_User, conf.LDAP_Pass)
	if err != nil {
		log.Printf("LDAP::Bind() error: %v\n", err)
		return 22
	}
	search := ldap.NewSearchRequest(
		conf.LDAP_BASE,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		conf.LDAP_Filter,
		SABDefine.LDAP_attr,
		nil)

	sr, err := l.Search(search)
	if err != nil {
		log.Printf("LDAP::Search() error: %v\n", err)
		return 23
	}

	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 11
	}

	defer db.Close()

	qwetabletrunc := fmt.Sprintf("truncate %s;", conf.PG_Table_LDAP)
	_, err = db.Query(qwetabletrunc)
	if err != nil {
		log.Printf("PG::Query() truncate table error: %v\n", err)
		return 12
	}

//	log.Printf("LDAP Export: %d // %s // %s // %s\n", result.Count(), result.Filter(), result.Base(), strings.Join(result.Attributes(), ", "))
	log.Printf("LDAP Export: %s // %d\n", search.Filter, len(sr.Entries))

	for _, entry := range sr.Entries {
//		log.Printf("dn=%s\n", entry.Dn())
		if ckl == 0 {
			queryx = fmt.Sprintf("INSERT INTO %s (namerus, namelat, ou, mail) VALUES ", conf.PG_Table_LDAP)
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
			queryx = fmt.Sprintf("%s ('%s','%s','%s','%s')", queryx, strings.Trim(strings.Replace(fName[int(len(fName)-1)], "'", "", -1), " "), strings.Trim(strings.Replace(fCN[0], "'", "", -1), " "), fOU, strings.Trim(strings.ToLower(fMail), " "))
			state = 1
		}
		if ckl>SABDefine.PG_MultiInsert {
			queryx = fmt.Sprintf("%s;", queryx) 
//			log.Printf("%s\n", queryx)
			if queryx != fmt.Sprintf("INSERT INTO %s (namerus, namelat, ou, mail) VALUES ;", conf.PG_Table_LDAP) {
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
				}
			}
			ckl=0
			state=0
		}else{
			ckl++
		}
	}

	log.Printf("LDAP FINISHED\n")

	return 94
}

func Oracle_to_PG(mode int, conf *SABDefine.Config_STR) int {

	var 	(

		pg_Table_Oracle		=	[]string{"dump_oracle_org", "dump_oracle_dep", "dump_oracle_pers"}
		pg_Query_Start		=	[]string{"INSERT INTO dump_oracle_org (uid, name) VALUES ", "INSERT INTO dump_oracle_dep (uid, idparent, name) VALUES ", "INSERT INTO dump_oracle_pers (uid, iddep, namefull,namelf,namelfi,position,phoneint, phonetown, phonecell, mail, idorg, idpos, sort) VALUES "}

		queryx			string
	)

	if mode<1 ||mode >3 {
		log.Fatalf("Oracle_to_PG mode select ERROR! Mode == \n", mode)
	}

	cx, err := oracle.NewConnection(conf.Oracle_user, conf.Oracle_pass, conf.Oracle_sid, false)

	if err != nil {
		log.Printf("Oracle::Connection() error: %v\n", err)
		return 1
	}

	defer cx.Close()
	cu := cx.NewCursor()
	defer cu.Close()

	log.Printf("Oracle Select for %s\n", pg_Table_Oracle[mode-1])

	err = cu.Execute(SABDefine.Oracle_QUE[mode-1], nil, nil)

	if err != nil {
		log.Printf("Oracle::Execute() error: %v\n", err)
		return 2
	}

	db, err := sql.Open("postgres", conf.PG_DSN)

	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 11
	}

	defer db.Close()

	rows, err := cu.FetchMany(SABDefine.PG_MultiInsert)

	qwetabletrunc := fmt.Sprintf("truncate %s;", pg_Table_Oracle[mode-1])

	_, err = db.Query(qwetabletrunc)
	if err != nil {
		log.Printf("PG::Query() truncate table error: %v\n", err)
		return 12
	}

	log.Printf("Oracle Export to PG %s\n", pg_Table_Oracle[mode-1])

	for err == nil && len(rows) > 0 {
		queryx = pg_Query_Start[mode-1]
		for ckl, row := range rows {
			if ckl>0 {
				queryx = fmt.Sprintf("%s, ", queryx)
			}
			switch mode {
				case 1:
					queryx = fmt.Sprintf("%s ('%v','%s')", queryx, row[0],row[1]) 
				case 2:
					queryx = fmt.Sprintf("%s ('%v','%v','%s')", queryx, row[0],row[1],row[2]) 
				case 3:
					queryx = fmt.Sprintf("%s ('%v','%v','%s','%s','%s','%s','%s','%s','%s','%s','%v','%v','%d')", queryx, row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9], row[10], row[11], row[12])
				default:
					break
			}
//			n++
		}

		queryx = fmt.Sprintf("%s;", queryx)

		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("PG::Query() insert error: %v /// %s\n", err, queryx)
		}
		rows, err = cu.FetchMany(SABDefine.PG_MultiInsert)
	}

	log.Printf("Oracle FINISHED for %s\n", pg_Table_Oracle[mode-1])

	return 94
}

