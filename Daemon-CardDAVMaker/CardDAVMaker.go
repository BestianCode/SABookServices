package main

import (
	"fmt"
	//	"os/exec"
	"crypto/md5"
	"encoding/hex"
	"log"
	"strings"
	//"regexp"
	//"net"
	"time"

	"database/sql"

	// PostgreSQL
	//_ "github.com/lib/pq"

	// MySQL
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"

	// LDAP
	"github.com/go-ldap/ldap"

	//"github.com/BestianRU/SABookServices/SABModules"
)

func main() {

	const (
		pName = string("SABook CardDAVMaker")
		pVer  = string("4 2015.10.19.02.00")
	)
	/*
		type userInfo struct {
			uName   string
			uOrg    string
			uOU     string
			uPos    string
			uMail   []string
			uPhInt  []string
			uMobile []string
			uPhExt  []string
		}
	*/
	type userMap struct {
		uName string
		uPass string
		uDN   []string
	}

	var (
		userList = []userMap{
			userMap{"smirnov_oa", "123", []string{"ou=Upr IT,ou=Obosoblennoe podrazdelenie Quadra - IA,ou=IA Quadra,ou=Quadra,o=Enterprise", "ou=AUP,ou=TsG,ou=Quadra,o=Enterprise"}},
			userMap{"trifonov_av", "12345", []string{"ou=Sl IT,ou=AUP,ou=TsG,ou=Quadra,o=Enterprise", "ou=Upr IT,ou=Obosoblennoe podrazdelenie Quadra - IA,ou=IA Quadra,ou=Quadra,o=Enterprise"}},
			userMap{"bugrov_dg", "1234", []string{"ou=Upr IT,ou=Obosoblennoe podrazdelenie Quadra - IA,ou=IA Quadra,ou=Quadra,o=Enterprise"}},
			userMap{"ivanov_da", "12345678", []string{"ou=AUP,ou=TsG,ou=Quadra,o=Enterprise"}}}

		ldapServer = string("asterisk.tula.domino:389")
		ldapUser   = string("")
		ldapPass   = string("")
		ldapAttr   = []string{"displayName", "businessCategory", "mail", "telephoneNumber", "mobile", "pager"}
		ldapVCard  = []string{"FN", "ROLE", "EMAIL;WORK", "TEL;TYPE=VOICE;TYPE=PREF", "TEL;CELL", "TEL;WORK"}
		//ldapAttr   = []string{"displayName", "entryDN", "businessCategory", "mail", "telephoneNumber", "mobile", "pager"}
		//ldapVCard  = []string{"FN", "ORG", "ROLE", "EMAIL;WORK", "TEL;TYPE=VOICE;TYPE=PREF", "TEL;CELL", "TEL;WORK"}
		//ldapAttrForSum = string("entryDN")

		realm = string("SABookDAV")

		queryx string

		multiInsert = int(50)

		idxUsers = int(1)
		idxCards = int(1)

		mySQL_DN = string("tcp:mysql.domino:3306*cdav/cdav/dav69admin")
		//mySQL_DN     = string("cdav:dav69admin@tcp(mysql.domino:3306)/cdav")
		mySQL_InitDB = string(`
CREATE TABLE IF NOT EXISTS addressbooks (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  principaluri varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  description text COLLATE utf8_unicode_ci,
  ctag int(11) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (id),
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS calendarobjects (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  calendardata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  calendarid int(10) unsigned NOT NULL,
  lastmodified int(11) unsigned DEFAULT NULL,
  etag varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  size int(11) unsigned NOT NULL,
  componenttype varchar(8) COLLATE utf8_unicode_ci DEFAULT NULL,
  firstoccurence int(11) unsigned DEFAULT NULL,
  lastoccurence int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY calendarid (calendarid,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS calendars (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  principaluri varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  ctag int(10) unsigned NOT NULL DEFAULT '0',
  description text COLLATE utf8_unicode_ci,
  calendarorder int(10) unsigned NOT NULL DEFAULT '0',
  calendarcolor varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  timezone text COLLATE utf8_unicode_ci,
  components varchar(21) COLLATE utf8_unicode_ci DEFAULT NULL,
  transparent tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (id),
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS cards (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  addressbookid int(11) unsigned NOT NULL,
  carddata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  lastmodified int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS groupmembers (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  principal_id int(10) unsigned NOT NULL,
  member_id int(10) unsigned NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY principal_id (principal_id,member_id)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS locks (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  owner varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  timeout int(10) unsigned DEFAULT NULL,
  created int(11) DEFAULT NULL,
  token varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  scope tinyint(4) DEFAULT NULL,
  depth tinyint(4) DEFAULT NULL,
  uri varchar(1000) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  KEY token (token),
  KEY uri (uri(333))
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS principals (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  uri varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  email varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  vcardurl varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uri (uri)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;

CREATE TABLE IF NOT EXISTS users (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  username varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  digesta1 varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY username (username)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1;





CREATE TABLE IF NOT EXISTS z_cache_users (
  id int(10) unsigned NOT NULL,
  username varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  digesta1 varchar(32) COLLATE utf8_unicode_ci DEFAULT NULL,
  UNIQUE KEY username (username)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_principals (
  id  int(10) unsigned NOT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  email varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(80) COLLATE utf8_unicode_ci DEFAULT NULL,
  vcardurl varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  UNIQUE KEY uri (uri)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_cards (
  id int(11) unsigned NOT NULL,
  addressbookid int(11) unsigned NOT NULL,
  carddata mediumblob,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  lastmodified int(11) unsigned DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE IF NOT EXISTS z_cache_addressbooks (
  id int(11) unsigned NOT NULL,
  principaluri varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  displayname varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  uri varchar(200) COLLATE utf8_unicode_ci DEFAULT NULL,
  description text COLLATE utf8_unicode_ci,
  ctag int(11) unsigned NOT NULL DEFAULT '1',
  UNIQUE KEY principaluri (principaluri,uri)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

truncate table z_cache_users;
truncate table z_cache_principals;
truncate table z_cache_cards;
truncate table z_cache_addressbooks;
			`)
		mySQL_Update = []string{`
delete from users where username in 
	(select * from
		(select a.username from users as a where a.username not in
			(select b.username from z_cache_users as b where b.username=a.username and b.digesta1=a.digesta1)) as c);
`, `
delete from principals where uri in
	(select * from
		(select a.uri from principals as a where a.uri not in
			(select b.uri from z_cache_principals as b where b.uri=a.uri)) as c);
`, `
delete from addressbooks where principaluri in
	(select * from
		(select a.principaluri from addressbooks as a where a.principaluri not in
			(select b.principaluri from z_cache_addressbooks as b where b.principaluri=a.principaluri)) as c);
`, `
delete x,y from
	principals as x join addressbooks y on x.id=y.id
			where x.id not in
				(select * from (select a.id from principals as a,
					(select id, username from users) as subq
						where a.id=subq.id and subq.username=REPLACE(a.uri, 'principals/', '')) as c);
`, `
insert into users (username,digesta1)
	select a.username, a.digesta1 from z_cache_users as a
		where a.username not in (select b.username from users as b where b.username=a.username and
			b.digesta1=a.digesta1);
`, `
insert into principals (id,uri)
	select subq.id, a.uri from (select id,username from users) as subq, z_cache_principals as a
		where a.uri not in (select b.uri from principals as b where b.uri=a.uri) and
			subq.username=REPLACE(a.uri, 'principals/', '');
`, `
insert into addressbooks (id,principaluri,uri,ctag)
	select subq.id, a.principaluri,'default',1 from (select id,username from users) as subq, z_cache_addressbooks as a
		where a.principaluri not in
			(select b.principaluri from addressbooks as b where b.principaluri=a.principaluri) and
				subq.username=REPLACE(a.principaluri, 'principals/', '');
`, `
delete from cards where uri in
	(select * from
		(select a.uri from cards as a where a.uri not in
			(select b.uri from z_cache_cards as b where b.uri=a.uri and a.addressbookid=b.addressbookid)) as c) or
	addressbookid not in (select id from users);
`, `
insert into cards (addressbookid, carddata, uri)
	select a.addressbookid, a.carddata, a.uri from z_cache_cards as a
		where a.uri not in
			(select b.uri from cards as b where a.uri=b.uri and a.addressbookid=b.addressbookid);
			`}
	)

	//truncate table users;
	//truncate table principals;
	//truncate table cards;
	//truncate table addressbooks;

	l, err := ldap.Dial("tcp", ldapServer)
	if err != nil {
		log.Printf("LDAP::Initialize() error: %v\n", err)
		return
	}

	//l.Debug = true
	defer l.Close()

	err = l.Bind(ldapUser, ldapPass)
	if err != nil {
		log.Printf("LDAP::Bind() error: %v\n", err)
		return
	}

	db, err := sql.Open("mymysql", mySQL_DN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return
	}

	defer db.Close()

	log.Printf("\tInitialize DB...\n")
	rows, err := db.Query(mySQL_InitDB)
	if err != nil {
		log.Printf("01 MySQL::Query() error: %v\n", err)
		return
	}
	log.Printf("\t\tComplete!\n")

	time.Sleep(time.Duration(2) * time.Second)

	x := make(map[string]string, len(ldapAttr))

	password := ""
	multiCount := 0
	log.Printf("\tCreate cacheDB from LDAP...\n")
	for i := 0; i < len(userList); i++ {
		queryx = fmt.Sprintf("select id from users where username='%s';", userList[i].uName)
		//log.Printf("%s\n", queryx)
		rows, err = db.Query(queryx)
		if err != nil {
			log.Printf("02 MySQL::Query() error: %v\n", err)
			log.Printf("%s\n", queryx)
			return
		}
		userIDGet := 0
		rows.Next()
		rows.Scan(&userIDGet)
		if userIDGet > 0 {
			idxUsers = userIDGet
		} else {
			queryx = "select id from users order by id desc limit 1;"
			//log.Printf("%s\n", queryx)
			rows, err = db.Query(queryx)
			if err != nil {
				log.Printf("03 MySQL::Query() error: %v\n", err)
				log.Printf("%s\n", queryx)
				return
			}
			rows.Next()
			rows.Scan(&userIDGet)

			if userIDGet > 0 {
				userIDGet++
				idxUsers = userIDGet
			}
		}

		//fmt.Printf("%d\n", userIDGet)

		z := md5.New()
		z.Write([]byte(fmt.Sprintf("%s:%s:%s", userList[i].uName, realm, userList[i].uPass)))
		password = hex.EncodeToString(z.Sum(nil))

		queryx = fmt.Sprintf("INSERT INTO z_cache_users (id, username, digesta1)\n\tVALUES (%d, '%s', '%s');", idxUsers, userList[i].uName, password)
		queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_principals (id, uri, email, displayname, vcardurl)\n\tVALUES (%d, 'principals/%s', NULL, NULL, NULL);", queryx, idxUsers, userList[i].uName)
		queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_addressbooks (id, principaluri, uri, ctag)\n\tVALUES (%d, 'principals/%s', 'default', 1); select id from users order by id desc limit 1", queryx, idxUsers, userList[i].uName)
		//log.Printf("%s\n", queryx)
		_, err = db.Query(queryx)
		if err != nil {
			log.Printf("03 MySQL::Query() error: %v\n", err)
			log.Printf("%s\n", queryx)
			return
		}

		for j := 0; j < len(userList[i].uDN); j++ {

			log.Printf("\t\t\t%3d/%s - %s\n", idxUsers, userList[i].uName, userList[i].uDN[j])

			search := ldap.NewSearchRequest(userList[i].uDN[j], 2, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=inetOrgPerson)", ldapAttr, nil)

			sr, err := l.Search(search)
			if err != nil {
				log.Printf("LDAP::Search() error: %v\n", err)
				return
			}

			queryx = ""
			if len(sr.Entries) > 0 {
				for _, entry := range sr.Entries {
					for k := 0; k < len(ldapAttr); k++ {
						x[ldapVCard[k]] = ""
					}
					for _, attr := range entry.Attributes {
						for k := 0; k < len(ldapAttr); k++ {
							if attr.Name == ldapAttr[k] {
								x[ldapVCard[k]] = strings.Join(attr.Values, ",")
								x[ldapVCard[k]] = strings.Replace(x[ldapVCard[k]], ",", "\n"+ldapVCard[k]+":", -1)
							}
						}
					}
					y := fmt.Sprintf("BEGIN:VCARD\n")
					for k := 0; k < len(ldapAttr); k++ {
						if x[ldapVCard[k]] != "" {
							y = fmt.Sprintf("%s%s:%s\n", y, ldapVCard[k], x[ldapVCard[k]])
						}
					}
					z := md5.New()
					z.Write([]byte(y))
					uid := hex.EncodeToString(z.Sum(nil))
					y = fmt.Sprintf("%sUID:%s\n", y, uid)
					y = fmt.Sprintf("%sEND:VCARD\n", y)
					//fmt.Printf("%s\n\t%s.vcf\n\n", y, uid)

					queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_cards (id, addressbookid, carddata, uri, lastmodified)\n\tVALUES (%d, %d, '%s', '%s.vcf', NULL);", queryx, idxCards, idxUsers, y, uid)
					if multiCount > multiInsert {
						//log.Printf("%s\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("MySQL::Query() error: %v\n", err)
							log.Printf("%s\n", queryx)
							return
						}
						queryx = ""
						multiCount = 0
					}
					multiCount++
					idxCards++

				}
			}
			_, err = db.Query(queryx)
			if err != nil {
				log.Printf("MySQL::Query() error: %v\n", err)
				log.Printf("%s\n", queryx)
				return
			}
			queryx = ""
			multiCount = 0
		}
		idxUsers++
	}
	log.Printf("\t\tComplete!\n")

	log.Printf("\tUpdate tables...\n")
	for i := 0; i < len(mySQL_Update); i++ {
		log.Printf("\t\t\tstep %d of %d...\n", i, len(mySQL_Update))
		_, err = db.Query(mySQL_Update[i])
		if err != nil {
			log.Printf("MySQL::Query() error: %v\n", err)
			return
		}
	}
	log.Printf("\t\tComplete!\n")

	/*
		var (
			def_config_file = string("./AsteriskCIDUpdater.json")             // Default configuration file
			def_log_file    = string("/var/log/ABook/AsteriskCIDUpdater.log") // Default log file
			def_daemon_mode = string("NO")                                    // Default start in foreground

			sqlite_key   string
			sqlite_value string

			pg_name  string
			pg_phone string

			pg_array [100000][3]string
			sq_array [100000][3]string

			pg_array_len = int(0)
			sq_array_len = int(0)

			ckl1       = int(0)
			ckl2       = int(0)
			ckl_status = int(0)

			ast_cmd string

			rconf SABModules.Config_STR

			sql_mode = int(0)

			query string
		)

		fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

		rconf.LOG_File = def_log_file

		def_config_file, def_daemon_mode = SABModules.ParseCommandLine(def_config_file, def_daemon_mode)

		//	log.Printf("%s %s %s", def_config_file, def_daemon_mode, os.Args[0])

		SABModules.ReadConfigFile(def_config_file, &rconf)

		sqlite_select := fmt.Sprintf("SELECT key, value FROM astdb where key like '%%%s%%';", rconf.AST_CID_Group)
		pg_select := fmt.Sprintf("select x.cid_name, y.phone from ldapx_persons x, ldapx_phones y, (select a.phone, count(a.phone) as phone_count from ldapx_phones as a, ldapx_persons as b where a.pers_id=b.uid and b.contract=0 and a.pass=2 and b.lang=1 group by a.phone order by a.phone) as subq where x.uid=y.pers_id and y.pass=2 and x.lang=1 and subq.phone=y.phone and subq.phone_count<2 and y.phone like '%s%%' and x.contract=0 group by x.cid_name, y.phone order by y.phone;", rconf.AST_Num_Start)

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
			sql_mode = 1
		}

		if sql_mode == 0 {
			defer dbast.Close()
		}

		ast_gami := gami.NewAsterisk(&dbast, nil)
		ast_get := make(chan gami.Message, 10000)

		if sql_mode == 0 {
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

			err = rows.Scan(&sqlite_key, &sqlite_value)
			if err != nil {
				log.Printf("rows.Scan error: %v\n", err)
				return
			}

			sq_array[sq_array_len][0] = sqlite_key
			sq_array[sq_array_len][1] = sqlite_value
			sq_array[sq_array_len][2] = SABModules.PhoneMutation(sqlite_key)
			sq_array_len++

		}

		rows, err = dbpg.Query(pg_select)
		if err != nil {
			log.Printf("PG::Query() error: %v\n", err)
			return
		}

		pg_array_len = 0
		for rows.Next() {

			err = rows.Scan(&pg_name, &pg_phone)
			if err != nil {
				log.Printf("rows.Scan error: %v\n", err)
				return
			}

			pg_array[pg_array_len][0] = fmt.Sprintf("/%s/%s", rconf.AST_CID_Group, pg_phone)
			pg_array[pg_array_len][1] = pg_name
			pg_array[pg_array_len][2] = SABModules.PhoneMutation(pg_phone)
			pg_array_len++

		}

		for ckl1 = 0; ckl1 < sq_array_len; ckl1++ {
			ckl_status = 0
			for ckl2 = 0; ckl2 < pg_array_len; ckl2++ {
				if sq_array[ckl1][0] == pg_array[ckl2][0] && sq_array[ckl1][1] == pg_array[ckl2][1] {
					ckl_status = 1
					break
				}
			}
			if ckl_status == 0 {
				if sql_mode == 0 {
					ast_cmd = fmt.Sprintf("database del %s %s", rconf.AST_CID_Group, sq_array[ckl1][2])
					log.Printf("\t- %s\n", ast_cmd)

					ast_cb := func(m gami.Message) {
						ast_get <- m
					}

					err = ast_gami.Command(ast_cmd, &ast_cb)
					if err != nil {
						log.Printf("Asterisk ARI::Command() error: %v\n", err)
						return
					}

					for x1, x2 := range <-ast_get {
						if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
							log.Printf("\t\t\t%s\n", x2)
						}
					}
				} else {
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

		for ckl1 = 0; ckl1 < pg_array_len; ckl1++ {
			ckl_status = 0
			for ckl2 = 0; ckl2 < sq_array_len; ckl2++ {
				if pg_array[ckl1][0] == sq_array[ckl2][0] && pg_array[ckl1][1] == sq_array[ckl2][1] {
					ckl_status = 1
					break
				}
			}
			if ckl_status == 0 {

				if sql_mode == 0 {
					ast_cmd = fmt.Sprintf("database del %s %s", rconf.AST_CID_Group, pg_array[ckl1][2])
					log.Printf("\t- %s\n", ast_cmd)

					ast_cb := func(m gami.Message) {
						ast_get <- m
					}

					err = ast_gami.Command(ast_cmd, &ast_cb)
					if err != nil {
						log.Printf("Asterisk ARI::Command() error: %v\n", err)
						return
					}

					for x1, x2 := range <-ast_get {
						if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
							log.Printf("\t\t\t%s\n", x2)
						}
					}

					ast_cmd = fmt.Sprintf("database put %s %s \"%s\"", rconf.AST_CID_Group, pg_array[ckl1][2], pg_array[ckl1][1])
					log.Printf("\t+ %s\n", ast_cmd)

					ast_cb = func(m gami.Message) {
						ast_get <- m
					}

					err = ast_gami.Command(ast_cmd, &ast_cb)
					if err != nil {
						log.Printf("Asterisk ARI::Command() error: %v\n", err)
						return
					}

					for x1, x2 := range <-ast_get {
						if x1 == "ActionID" || x1 == "CmdData" || x1 == "Usage" {
							log.Printf("\t\t\t%s\n", x2)
						}
					}
				} else {
					query = fmt.Sprintf("delete from astdb where key='%s';", pg_array[ckl1][0])
					log.Printf("\t- %s\n", query)
					_, err := db.Exec(query)
					if err != nil {
						log.Printf("SQLite3::Query() DEL error: %v\n", err)
						return
					}

					query = fmt.Sprintf("insert into astdb (key,value) values ('%s','%s');", pg_array[ckl1][0], pg_array[ckl1][1])
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
	*/
}
