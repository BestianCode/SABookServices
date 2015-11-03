package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"

	// PostgreSQL
	_ "github.com/lib/pq"

	// MySQL
	_ "github.com/ziutek/mymysql/godrv"

	// LDAP
	"github.com/go-ldap/ldap"

	"github.com/BestianRU/SABModules/SBMConnect"
	"github.com/BestianRU/SABModules/SBMSystem"
)

func checkNTUWishes(conf SBMSystem.ReadJSONConfig, logRedirect SBMSystem.LogFile) int {
	var (
		i  int
		pg SBMConnect.PgSQL
	)

	pg.Init(conf)
	defer pg.Close()

	pgrows, err := pg.D.Query("select count(userid) from aaa_dav_ntu")
	if err != nil {
		log.Printf("PG::Query() checkNTUWishes error: %v\n", err)
		return -1
	}

	pgrows.Next()
	pgrows.Scan(&i)

	return i
}

func goNTUWork(conf SBMSystem.ReadJSONConfig, logRedirect SBMSystem.LogFile) {

	type usIDPartList struct {
		id   int
		name string
	}

	var (
		ldap_Attr   []string
		ldap_VCard  []string
		queryx      string
		multiInsert = int(50)
		idxUsers    = int(1)
		idxCards    = int(1)
		workMode    = string("FULL")
		i           int
		pg          SBMConnect.PgSQL
	)

	logRedirect.Log("--> WakeUP!")

	ldap_Attr = make([]string, len(conf.Conf.WLB_LDAP_ATTR))

	for i = 0; i < len(conf.Conf.WLB_LDAP_ATTR); i++ {
		ldap_Attr[i] = conf.Conf.WLB_LDAP_ATTR[i][0]
	}

	ldap_VCard = make([]string, len(conf.Conf.WLB_LDAP_ATTR))

	for i = 0; i < len(conf.Conf.WLB_LDAP_ATTR); i++ {
		ldap_VCard[i] = conf.Conf.WLB_LDAP_ATTR[i][1]
	}

	pg.Init(conf)
	defer pg.Close()

	l, err := ldap.Dial("tcp", conf.Conf.LDAP_URL[0][0])
	if err != nil {
		log.Printf("LDAP::Initialize() error: %v\n", err)
		return
	}

	//l.Debug = true
	defer l.Close()

	err = l.Bind(conf.Conf.LDAP_URL[0][1], conf.Conf.LDAP_URL[0][2])
	if err != nil {
		log.Printf("LDAP::Bind() error: %v\n", err)
		return
	}

	db, err := sql.Open("mymysql", conf.Conf.MY_DSN)
	if err != nil {
		log.Printf("MySQL::Open() error: %v\n", err)
		return
	}

	defer db.Close()

	logRedirect.Log("\tInitialize DB...")
	rows, err := db.Query(mySQL_InitDB)
	if err != nil {
		log.Printf("01 MySQL::Query() error: %v\n", err)
		return
	}
	logRedirect.Log("\t\tComplete!")

	time.Sleep(10 * time.Second)

	x := make(map[string]string, len(ldap_Attr))

	//password := ""
	multiCount := 0
	logRedirect.Log("\tCreate cacheDB from LDAP...")

	time_now := time.Now().Unix()
	time_get := 0

	pgrows1, err := pg.D.Query("select updtime from aaa_dav_ntu where userid=0;")
	if err != nil {
		log.Printf("01 PG::Query() error: %v\n", err)
		return
	}

	pgrows1.Next()
	pgrows1.Scan(&time_get)

	if time_get > 0 {
		pgrows1, err = pg.D.Query("select x.id, x.login, x.password from aaa_logins as x where x.id in (select userid from aaa_dns where userid=x.id) order by login;")
		if err != nil {
			log.Printf("03 PG::Query() error: %v\n", err)
			return
		}
		workMode = "FULL"
	} else {
		pgrows1, err = pg.D.Query("select x.id, x.login, x.password from aaa_logins as x, aaa_dav_ntu as y where x.id=y.userid and x.id in (select userid from aaa_dns where userid=x.id) order by login;")
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

		pgrows2, err := pg.D.Query(fmt.Sprintf("select dn from aaa_dns where userid=%d;", usID))
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

			//log.Printf("%s|||%s|||%s\n", usDN, conf.Conf.LDAP_URL[0][4], ldap_Attr)

			search := ldap.NewSearchRequest(usDN, 2, ldap.NeverDerefAliases, 0, 0, false, conf.Conf.LDAP_URL[0][4], ldap_Attr, nil)

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

	logRedirect.Log("\t\tComplete!")

	if workMode == "PART" {
		logRedirect.Log("\tUpdate tables in PartialUpdate mode...")
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
		logRedirect.Log("\tUpdate tables...")
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

	logRedirect.Log("\t\tComplete!")

	logRedirect.Log("\tClean NeedToUpdate table...")
	queryx = fmt.Sprintf("delete from aaa_dav_ntu where userid=0 or updtime<%d;", time_now)
	//log.Printf("%s\n", queryx)
	_, err = pg.D.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() Clean NTU table error: %v\n", err)
		return
	}

	logRedirect.Log("\tComplete!")
	logRedirect.Bye()
}

func main() {

	const (
		pName = string("SABook CardDAVMaker")
		pVer  = string("5 2015.11.03.21.00")
	)

	var (
		jsonConfig  SBMSystem.ReadJSONConfig
		logRedirect SBMSystem.LogFile
		pid         SBMSystem.PidFile
		sleepWatch  = int(0)
	)

	fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

	jsonConfig.Init()

	pid.ON(jsonConfig)
	pid.OFF(jsonConfig)

	logRedirect.ON(jsonConfig)
	logRedirect.OFF()

	SBMSystem.Fork(jsonConfig)
	SBMSystem.Signal(jsonConfig, pid)

	pid.ON(jsonConfig)
	defer pid.OFF(jsonConfig)

	logRedirect.ON(jsonConfig)
	logRedirect.Hello(pName, pVer)
	logRedirect.OFF()

	for {
		logRedirect.ON(jsonConfig)

		if checkNTUWishes(jsonConfig, logRedirect) > 0 {
			logRedirect.Hello(pName, pVer)
			goNTUWork(jsonConfig, logRedirect)
			sleepWatch = 0
		}

		if sleepWatch > 3600 {
			logRedirect.Log("<-- I'm alive ... :)")
			sleepWatch = 0
		}

		logRedirect.OFF()
		time.Sleep(time.Duration(jsonConfig.Conf.Sleep_Time) * time.Second)
		sleepWatch += jsonConfig.Conf.Sleep_Time
	}
}
