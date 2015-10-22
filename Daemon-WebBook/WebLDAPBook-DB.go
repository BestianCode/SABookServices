package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	// PostgreSQL
	_ "github.com/lib/pq"
	// SQLite
	_ "github.com/mattn/go-sqlite3"

	//LDAP
	"github.com/go-ldap/ldap"
)

var (
	db   *sql.DB
	dbpg *sql.DB
)

func pgInit() {
	var (
		err error

		query_create = string(`
CREATE TABLE IF NOT EXISTS wb_auth_session (
	username varchar(255),
	sessname varchar(255)  PRIMARY KEY,
	exptime integer
);
		`)
	)

	dbpg, err = sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	//defer dbpg.Close()
	/*
		err = db.Ping()
		if err != nil {
			log.Fatalf("SQLiteINIT::Ping() error: %v\n", err)
		}

		_, err = db.Begin()
		if err != nil {
			log.Fatalf("SQLiteINIT::Begin() error: %v\n", err)
		}
	*/
	_, err = dbpg.Exec(query_create)
	if err != nil {
		log.Fatalf("PG_INIT::Exec() create tables error: %v\n", err)
	}

}

func sqliteInit() {
	var (
		err error

		query_create = string(`
PRAGMA journal_mode=WAL;
CREATE TABLE IF NOT EXISTS ldap (
	FullName varchar(255),
	FirstName varchar(255),
	LastName varchar(255),
	DN varchar(1024) PRIMARY KEY
);
		`)
	)

	db, err = sql.Open("sqlite3", rconf.WLB_SQLite_DB)
	if err != nil {
		log.Fatalf("SQLiteINIT::Open() error: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("SQLiteINIT::Ping() error: %v\n", err)
	}

	_, err = db.Begin()
	if err != nil {
		log.Fatalf("SQLiteINIT::Begin() error: %v\n", err)
	}

	_, err = db.Exec(query_create)
	if err != nil {
		log.Fatalf("SQLiteINIT::Exec() create tables error: %v\n", err)
	}

}

func sqliteUpdate() {
	var (
		ckl1, ckl2 int
		err        error
		l          *ldap.Conn

		query_create = string(`
CREATE temp TABLE IF NOT EXISTS ldap_cache (
	FullName varchar(255),
	FirstName varchar(255),
	LastName varchar(255),
	DN varchar(1024) PRIMARY KEY
);
		`)
		query_clean = string(`
delete from ldap_cache;
		`)
		query_update = string(`
delete from ldap where DN not in (select DN from ldap_cache);
insert into ldap (FullName, LastName, FirstName, DN)
select FullName, LastName, FirstName, DN from ldap_cache where DN not in (select DN from ldap);
		`)
	)

	sleepTime = rconf.Sleep_Time

	ldap_Attr := make([]string, 4)
	userdb := make(map[string]string, 4)

	log.Printf("SQLite Update ***** Starting...\n")

	_, err = db.Exec(query_create)
	if err != nil {
		log.Printf("SQLiteUPD::Exec() error: %v\n", err)
		sleepTime = 60
		return
	}

	if initLDAPConnector() == "error" {
		log.Printf("SQLite Update ***** Error connecting to LDAP !!!")
		sleepTime = 60
		return
	}

	l, err = ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
	if err != nil {
		log.Printf("SQLiteUPD->LDAP::Initialize() error: %v\n", err)
		sleepTime = 60
		return
	}

	//l.Debug = true
	defer l.Close()

	err = l.Bind(rconf.LDAP_URL[ldap_count][1], rconf.LDAP_URL[ldap_count][2])
	if err != nil {
		log.Printf("SQLiteUPD->LDAP::Bind() error: %v\n", err)
		sleepTime = 60
		return
	}

	for ckl1 = 0; ckl1 < len(rconf.WLB_LDAP_ATTR); ckl1++ {
		if rconf.WLB_LDAP_ATTR[ckl1][1] == "FullName" || rconf.WLB_LDAP_ATTR[ckl1][1] == "FirstName" || rconf.WLB_LDAP_ATTR[ckl1][1] == "LastName" || rconf.WLB_LDAP_ATTR[ckl1][1] == "DN" {
			ldap_Attr[ckl2] = rconf.WLB_LDAP_ATTR[ckl1][0]
			ckl2++
		}
	}

	//search := ldap.NewSearchRequest("ou=IA Quadra,ou=Quadra,o=Enterprise", 2, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=inetOrgPerson)", ldap_Attr, nil)
	search := ldap.NewSearchRequest(rconf.LDAP_URL[ldap_count][3], 2, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=inetOrgPerson)", ldap_Attr, nil)

	//log.Printf("Search: %v\n%v\n%v\n%v\n", search, rconf.LDAP_URL[ldap_count][3], ldap.NeverDerefAliases, ldap_Attr)

	sr, err := l.Search(search)
	if err != nil {
		log.Printf("SQLiteUPD->LDAP::Search() error: %v\n", err)
		sleepTime = 60
		return
	}

	log.Printf("SQLite Update ***** search: %s // found: %d\n", search.Filter, len(sr.Entries))

	if len(sr.Entries) > 0 {
		ckl2 = 0
		for _, entry := range sr.Entries {
			for _, attr := range entry.Attributes {
				for ckl1 = 0; ckl1 < len(rconf.WLB_LDAP_ATTR); ckl1++ {
					if rconf.WLB_LDAP_ATTR[ckl1][0] == attr.Name {
						userdb[rconf.WLB_LDAP_ATTR[ckl1][1]] = fmt.Sprintf("%s", strings.ToLower(strings.Join(attr.Values, ",")))
					}
				}
			}
			stmt, err := db.Prepare("INSERT INTO ldap_cache (FullName, LastName, FirstName, DN) values (?,?,?,?)")
			if err != nil {
				log.Printf("SQLiteUPD::Prepare() error: %v\n", err)
				sleepTime = 60
				return
			}
			_, err = stmt.Exec(userdb["FullName"], userdb["LastName"], userdb["FirstName"], userdb["DN"])
			//if err != nil {
			//	log.Printf("SQLite::Exec() error: %v\n", err)
			//}
			if ckl2 > 0 && ckl2 == int(ckl2/100)*100 {
				log.Printf("SQLite Update ***** %7d elements passed...\n", ckl2)
			}
			ckl2++
		}
		log.Printf("SQLite Update ***** %7d elements passed...\n", ckl2)
	}

	_, err = db.Exec(query_update)
	if err != nil {
		log.Printf("SQLiteUPD::Query() error: %v\n", err)
		sleepTime = 60
		return
	}

	_, err = db.Exec(query_clean)
	if err != nil {
		log.Printf("SQLite::Query() error: %v\n", err)
		sleepTime = 60
		return
	}

	log.Printf("SQLite Update ***** Completed!\n")
}
