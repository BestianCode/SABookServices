package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"database/sql"

	// PostgreSQL
	_ "github.com/lib/pq"

	// MySQL
	_ "github.com/ziutek/mymysql/godrv"

	// LDAP
	"github.com/go-ldap/ldap"

	"github.com/BestianRU/SABookServices/SABModules"
)

func main() {

	const (
		pName = string("SABook CardDAVMaker")
		pVer  = string("4 2015.10.31.01.00")
	)

	type usIDPartList struct {
		id   int
		name string
	}

	var (
		ldap_Attr       []string
		ldap_VCard      []string
		queryx          string
		multiInsert     = int(50)
		idxUsers        = int(1)
		idxCards        = int(1)
		def_config_file = string("./CardDAVMaker.json") // Default configuration file
		def_log_file    = string("./CardDAVMaker.log")  // Default log file
		def_daemon_mode = string("NO")                  // Default start in foreground
		rconf           SABModules.Config_STR
		workMode        = string("FULL")
		i               int
		mySQL_InitDB    = string(`
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
		mySQL_Update_full1 = []string{`
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
			`}

		mySQL_Update_full2 = []string{`
delete from cards where uri in
	(select * from
		(select a.uri from cards as a where a.uri not in
			(select b.uri from z_cache_cards as b where b.uri=a.uri and a.addressbookid=b.addressbookid)) as c) or
	addressbookid not in (select id from users);
		`}

		mySQL_Update_part1 = []string{`
delete from users where id=XYZIDXYZ and username in 
	(select * from
		(select a.username from users as a where a.id=XYZIDXYZ and a.username not in
			(select b.username from z_cache_users as b where b.id=XYZIDXYZ and b.username=a.username and b.digesta1=a.digesta1)) as c);
`, `
delete from principals where id=XYZIDXYZ and uri in
	(select * from
		(select a.uri from principals as a where a.id=XYZIDXYZ and a.uri not in
			(select b.uri from z_cache_principals as b where a.id=XYZIDXYZ and b.uri=a.uri)) as c);
`, `
delete from addressbooks where id=XYZIDXYZ and principaluri in
	(select * from
		(select a.principaluri from addressbooks as a where a.id=XYZIDXYZ and a.principaluri not in
			(select b.principaluri from z_cache_addressbooks as b where b.id=XYZIDXYZ and b.principaluri=a.principaluri)) as c);
`, `
delete x,y from
	principals as x join addressbooks y on x.id=y.id
			where x.id=XYZIDXYZ and x.id not in
				(select * from (select a.id from principals as a,
					(select id, username from users) as subq
						where a.id=XYZIDXYZ and a.id=subq.id and subq.username=REPLACE(a.uri, 'principals/', '')) as c);
			`}
		mySQL_Update_part2 = []string{`
delete from cards where addressbookid=XYZIDXYZ and uri in
	(select * from
		(select a.uri from cards as a where a.addressbookid=XYZIDXYZ and a.uri not in
			(select b.uri from z_cache_cards as b where b.uri=a.uri and a.addressbookid=b.addressbookid)) as c) or
	addressbookid not in (select id from users);
		  `}

		mySQL_Update1 = []string{`
insert into users (id,username,digesta1)
	select a.id, a.username, a.digesta1 from z_cache_users as a
		where a.username not in (select b.username from users as b where b.username=a.username and
			b.digesta1=a.digesta1);
`, `
insert into principals (id,uri)
	select a.id, a.uri from (select id,username from users) as subq, z_cache_principals as a
		where a.uri not in (select b.uri from principals as b where b.uri=a.uri) and
			subq.username=REPLACE(a.uri, 'principals/', '');
`, `
insert into addressbooks (id,principaluri,uri,ctag)
	select subq.id, a.principaluri,'default',1 from (select id,username from users) as subq, z_cache_addressbooks as a
		where a.principaluri not in
			(select b.principaluri from addressbooks as b where b.principaluri=a.principaluri) and
				subq.username=REPLACE(a.principaluri, 'principals/', '');
			`}

		mySQL_Update2 = []string{`
insert into cards (addressbookid, carddata, uri)
	select a.addressbookid, a.carddata, a.uri from z_cache_cards as a
		where a.uri not in
			(select b.uri from cards as b where a.uri=b.uri and a.addressbookid=b.addressbookid);
		  			`}
	)

	fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

	rconf.LOG_File = def_log_file

	def_config_file, def_daemon_mode = SABModules.ParseCommandLine(def_config_file, def_daemon_mode)

	log.Printf("%s %s %s", def_config_file, def_daemon_mode, os.Args[0])

	SABModules.ReadConfigFile(def_config_file, &rconf)

	SABModules.Pid_Check(&rconf)

	if def_daemon_mode == "YES" {
		if err := exec.Command(os.Args[0], fmt.Sprintf("-daemon=GO -config=%s &", def_config_file)).Start(); err != nil {
			log.Fatalf("Fork daemon error: %v", err)
		} else {
			log.Printf("Forked!")
			os.Exit(0)
		}
	}

	SABModules.Log_ON(&rconf)
	SABModules.Log_OFF()

	SABModules.Pid_ON(&rconf)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		signalType := <-ch
		signal.Stop(ch)

		SABModules.Log_ON(&rconf)

		log.Printf(".")
		log.Printf("..")
		log.Printf("...")
		log.Printf("Exit command received. Exiting...")
		log.Println("Signal type: ", signalType)
		log.Printf("Bye...")
		log.Printf("...")
		log.Printf("..")
		log.Printf(".")

		SABModules.Log_OFF()

		SABModules.Pid_OFF(&rconf)

		os.Exit(0)
	}()

	ldap_Attr = make([]string, len(rconf.WLB_LDAP_ATTR))

	for i = 0; i < len(rconf.WLB_LDAP_ATTR); i++ {
		ldap_Attr[i] = rconf.WLB_LDAP_ATTR[i][0]
	}

	ldap_VCard = make([]string, len(rconf.WLB_LDAP_ATTR))

	for i = 0; i < len(rconf.WLB_LDAP_ATTR); i++ {
		ldap_VCard[i] = rconf.WLB_LDAP_ATTR[i][1]
	}

	for {

		SABModules.Log_ON(&rconf)

		log.Printf("--> WakeUP!")

		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Printf("PG::Open() error: %v\n", err)
			return
		}

		defer dbpg.Close()

		pgrows1, err := dbpg.Query("select count(userid) from aaa_dav_ntu")
		if err != nil {
			log.Printf("01 PG::Query() error: %v\n", err)
			return
		}

		pgrows1.Next()
		pgrows1.Scan(&i)

		if i > 0 {

			l, err := ldap.Dial("tcp", rconf.LDAP_URL[0][0])
			if err != nil {
				log.Printf("LDAP::Initialize() error: %v\n", err)
				return
			}

			//l.Debug = true
			defer l.Close()

			err = l.Bind(rconf.LDAP_URL[0][1], rconf.LDAP_URL[0][2])
			if err != nil {
				log.Printf("LDAP::Bind() error: %v\n", err)
				return
			}

			db, err := sql.Open("mymysql", rconf.MY_DSN)
			if err != nil {
				log.Printf("MySQL::Open() error: %v\n", err)
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

			time.Sleep(10 * time.Second)

			x := make(map[string]string, len(ldap_Attr))

			//password := ""
			multiCount := 0
			log.Printf("\tCreate cacheDB from LDAP...\n")

			time_now := time.Now().Unix()
			time_get := 0

			pgrows1, err := dbpg.Query("select updtime from aaa_dav_ntu where userid=0;")
			if err != nil {
				log.Printf("01 PG::Query() error: %v\n", err)
				return
			}

			pgrows1.Next()
			pgrows1.Scan(&time_get)

			if time_get > 0 {
				pgrows1, err = dbpg.Query("select x.id, x.login, x.password from aaa_logins as x where x.id in (select userid from aaa_dns where userid=x.id) order by login;")
				if err != nil {
					log.Printf("03 PG::Query() error: %v\n", err)
					return
				}
				workMode = "FULL"
			} else {
				pgrows1, err = dbpg.Query("select x.id, x.login, x.password from aaa_logins as x, aaa_dav_ntu as y where x.id=y.userid and x.id in (select userid from aaa_dns where userid=x.id) order by login;")
				if err != nil {
					log.Printf("04 PG::Query() error: %v\n", err)
					return
				}
				workMode = "PART"
			}

			usID := 0
			usName := ""
			usPass := ""
			usIDArray := make([]usIDPartList, 0)
			for pgrows1.Next() {
				//for i := 0; i < len(userList); i++ {

				pgrows1.Scan(&usID, &usName, &usPass)
				usIDArray = append(usIDArray, usIDPartList{id: usID, name: usName})
				queryx = fmt.Sprintf("select id from users where username='%s';", usName)
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

				//z := md5.New()
				//z.Write([]byte(fmt.Sprintf("%s:%s:%s", userList[i].uName, realm, userList[i].uPass)))
				//password = hex.EncodeToString(z.Sum(nil))

				queryx = fmt.Sprintf("INSERT INTO z_cache_users (id, username, digesta1)\n\tVALUES (%d, '%s', '%s');", usID, usName, usPass)
				queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_principals (id, uri, email, displayname, vcardurl)\n\tVALUES (%d, 'principals/%s', NULL, NULL, NULL);", queryx, usID, usName)
				queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_addressbooks (id, principaluri, uri, ctag)\n\tVALUES (%d, 'principals/%s', 'default', 1); select id from users order by id desc limit 1", queryx, usID, usName)
				//log.Printf("%s\n", queryx)
				_, err = db.Query(queryx)
				if err != nil {
					log.Printf("03 MySQL::Query() error: %v\n", err)
					log.Printf("%s\n", queryx)
					return
				}

				pgrows2, err := dbpg.Query(fmt.Sprintf("select dn from aaa_dns where userid=%d;", usID))
				if err != nil {
					log.Printf("02 PG::Query() error: %v\n", err)
					return
				}

				usDN := ""
				//for j := 0; j < len(usDN); j++ {
				for pgrows2.Next() {

					pgrows2.Scan(&usDN)
					//queryx = fmt.Sprintf("select id from users where username='%s';", usName)

					log.Printf("\t\t\t%3d/%s - %s\n", usID, usName, usDN)

					//log.Printf("%s|||%s|||%s\n", usDN, rconf.LDAP_URL[0][4], ldap_Attr)

					search := ldap.NewSearchRequest(usDN, 2, ldap.NeverDerefAliases, 0, 0, false, rconf.LDAP_URL[0][4], ldap_Attr, nil)

					sr, err := l.Search(search)
					if err != nil {
						log.Printf("LDAP::Search() error: %v\n", err)
						return
					}

					queryx = ""
					if len(sr.Entries) > 0 {
						for _, entry := range sr.Entries {
							for k := 0; k < len(ldap_Attr); k++ {
								x[ldap_VCard[k]] = ""
							}
							for _, attr := range entry.Attributes {
								for k := 0; k < len(ldap_Attr); k++ {
									if attr.Name == ldap_Attr[k] {
										x[ldap_VCard[k]] = strings.Join(attr.Values, ",")
										x[ldap_VCard[k]] = strings.Replace(x[ldap_VCard[k]], ",", "\n"+ldap_VCard[k]+":", -1)
									}
								}
							}
							y := fmt.Sprintf("BEGIN:VCARD\n")
							for k := 0; k < len(ldap_Attr); k++ {
								if x[ldap_VCard[k]] != "" {
									if ldap_VCard[k] == "FN" {
										fn_split := strings.Split(x[ldap_VCard[k]], " ")
										fn_nofam := strings.Replace(x[ldap_VCard[k]], fn_split[0], "", -1)
										fn_nofam = strings.Trim(fn_nofam, " ")
										y = fmt.Sprintf("%s%s:%s %s\n", y, ldap_VCard[k], fn_nofam, fn_split[0])
										//fmt.Printf("%s%s:%s %s\n", y, ldap_VCard[k], fn_nofam, fn_split[0])
									} else {
										y = fmt.Sprintf("%s%s:%s\n", y, ldap_VCard[k], x[ldap_VCard[k]])
									}
								}
							}
							z := md5.New()
							z.Write([]byte(y))
							uid := hex.EncodeToString(z.Sum(nil))
							y = fmt.Sprintf("%sUID:%s\n", y, uid)
							y = fmt.Sprintf("%sEND:VCARD\n", y)
							//fmt.Printf("%s\n\t%s.vcf\n\n", y, uid)

							queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_cards (id, addressbookid, carddata, uri, lastmodified)\n\tVALUES (%d, %d, '%s', '%s.vcf', NULL);", queryx, idxCards, usID, y, uid)
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

			if workMode == "PART" {
				log.Printf("\tUpdate tables in PartialUpdate mode...\n")
				for j := 0; j < len(usIDArray); j++ {
					log.Printf("\t\t\tUpdate %d/%s...\n", usIDArray[j].id, usIDArray[j].name)
					for i := 0; i < len(mySQL_Update_part1); i++ {
						log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update_part1))

						queryx = strings.Replace(mySQL_Update_part1[i], "XYZIDXYZ", fmt.Sprintf("%d", usIDArray[j].id), -1)
						//log.Printf("%s\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("MySQL::Query() error: %v\n", err)
							return
						}
						time.Sleep(2 * time.Second)
					}
					for i := 0; i < len(mySQL_Update1); i++ {
						log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update1))
						//log.Printf("%s\n", mySQL_Update1[i])
						_, err = db.Query(mySQL_Update1[i])
						if err != nil {
							log.Printf("%s\n", mySQL_Update1[i])
							log.Printf("MySQL::Query() error: %v\n", err)
							return
						}
						time.Sleep(2 * time.Second)
					}
					for i := 0; i < len(mySQL_Update_part2); i++ {
						log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update_part2))

						queryx = strings.Replace(mySQL_Update_part2[i], "XYZIDXYZ", fmt.Sprintf("%d", usIDArray[j].id), -1)
						//log.Printf("%s\n", queryx)
						_, err = db.Query(queryx)
						if err != nil {
							log.Printf("%s\n", queryx)
							log.Printf("MySQL::Query() error: %v\n", err)
							return
						}
						time.Sleep(2 * time.Second)
					}
					for i := 0; i < len(mySQL_Update2); i++ {
						log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update2))
						//log.Printf("%s\n", mySQL_Update2[i])
						_, err = db.Query(mySQL_Update2[i])
						if err != nil {
							log.Printf("%s\n", mySQL_Update2[i])
							log.Printf("MySQL::Query() error: %v\n", err)
							return
						}
						time.Sleep(2 * time.Second)
					}
					time.Sleep(2 * time.Second)
				}
			} else {
				log.Printf("\tUpdate tables...\n")
				for i := 0; i < len(mySQL_Update_full1); i++ {
					log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update_full1))
					_, err = db.Query(mySQL_Update_full1[i])
					if err != nil {
						log.Printf("%s\n", mySQL_Update_full1[i])
						log.Printf("MySQL::Query() error: %v\n", err)
						return
					}
					time.Sleep(2 * time.Second)
				}
				for i := 0; i < len(mySQL_Update1); i++ {
					log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update1))
					_, err = db.Query(mySQL_Update1[i])
					if err != nil {
						log.Printf("%s\n", mySQL_Update1[i])
						log.Printf("MySQL::Query() error: %v\n", err)
						return
					}
					time.Sleep(2 * time.Second)
				}
				for i := 0; i < len(mySQL_Update_full2); i++ {
					log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update_full2))
					_, err = db.Query(mySQL_Update_full2[i])
					if err != nil {
						log.Printf("%s\n", mySQL_Update_full2[i])
						log.Printf("MySQL::Query() error: %v\n", err)
						return
					}
					time.Sleep(2 * time.Second)
				}
				for i := 0; i < len(mySQL_Update2); i++ {
					log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update2))
					_, err = db.Query(mySQL_Update2[i])
					if err != nil {
						log.Printf("%s\n", mySQL_Update2[i])
						log.Printf("MySQL::Query() error: %v\n", err)
						return
					}
					time.Sleep(2 * time.Second)
				}
			}

			log.Printf("\t\tComplete!\n")

			log.Printf("\tClean NeedToUpdate table...\n")
			queryx = fmt.Sprintf("delete from aaa_dav_ntu where userid=0 or updtime<%d;", time_now)
			//log.Printf("%s\n", queryx)
			_, err = dbpg.Query(queryx)
			if err != nil {
				log.Printf("PG::Query() Clean NTU table error: %v\n", err)
				return
			}

			log.Printf("\tComplete!\n")

			l.Close()
			db.Close()
		}
		dbpg.Close()

		log.Printf("----- Sleep for %d sec...", rconf.Sleep_Time)

		SABModules.Log_OFF()

		time.Sleep(time.Duration(rconf.Sleep_Time) * time.Second)
	}
}
