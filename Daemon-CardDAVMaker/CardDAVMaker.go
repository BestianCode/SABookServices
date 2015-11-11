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

func checkNTUWishes(conf SBMSystem.ReadJSONConfig, rLog SBMSystem.LogFile) int {
	var pg SBMConnect.PgSQL
	if pg.Init(conf, "") != 0 {
		return -1
	}
	defer pg.Close()
	return pg.QSimple("select count(userid) from aaa_dav_ntu")
}

func goNTUWork(conf SBMSystem.ReadJSONConfig, rLog SBMSystem.LogFile) {

	type usIDPartList struct {
		id   int
		name string
	}

	var (
		ldap_Attr   []string
		ldap_VCard  []string
		queryx      string
		multiInsert = int(64)
		idxUsers    = int(1)
		idxCards    = int(1)
		workMode    = string("FULL")
		i           int
		pg          SBMConnect.PgSQL
		my          SBMConnect.MySQL
		ld          SBMConnect.LDAP
		err         error
		pgrows1     *sql.Rows
	)

	rLog.Log("--> WakeUP!")

	for i = 0; i < len(conf.Conf.WLB_LDAP_ATTR); i++ {
		ldap_Attr = append(ldap_Attr, conf.Conf.WLB_LDAP_ATTR[i][0])
		ldap_VCard = append(ldap_VCard, conf.Conf.WLB_LDAP_ATTR[i][1])
	}

	if pg.Init(conf, "") < 0 {
		log.Println("Init PgSQL Error!")
		return
	}
	defer pg.Close()

	if my.Init(conf, mySQL_InitDB) < 0 {
		log.Println("Init MySQL Error!")
		return
	}
	defer my.Close()

	if ld.Init(conf) < 0 {
		log.Println("Init LDAP Error!")
		return
	}
	defer ld.Close()

	time.Sleep(10 * time.Second)

	x := make(map[string]string, len(ldap_Attr))

	multiCount := 0
	rLog.Log("\tCreate cacheDB from LDAP...")

	time_now := time.Now().Unix()
	time_get := pg.QSimple("select updtime from aaa_dav_ntu where userid=0;")

	switch {
	case time_get < 0:
		log.Println("PgSQL: Error reading NTU table")
		return
	case time_get > 0:
		pgrows1, err = pg.D.Query("select x.id, x.login, x.password from aaa_logins as x where x.id in (select userid from aaa_dns where userid=x.id) order by login;")
		if err != nil {
			log.Printf("PG::Query() 02 error: %v\n", err)
			return
		}
		workMode = "FULL"
	case time_get == 0:
		pgrows1, err = pg.D.Query("select x.id, x.login, x.password from aaa_logins as x, aaa_dav_ntu as y where x.id=y.userid and x.id in (select userid from aaa_dns where userid=x.id) order by login;")
		if err != nil {
			log.Printf("PG::Query() 03 error: %v\n", err)
			return
		}
		workMode = "PART"
	}

	usID := 0
	usName := ""
	usPass := ""
	userIDGet := 0
	usIDArray := make([]usIDPartList, 0)
	for pgrows1.Next() {
		pgrows1.Scan(&usID, &usName, &usPass)
		usIDArray = append(usIDArray, usIDPartList{id: usID, name: usName})

		userIDGet = my.QSimple("select id from users where username='", usName, "';")
		switch {
		case userIDGet < 0:
			log.Println("Error get user ID")
			return
		case userIDGet > 0:
			idxUsers = userIDGet
		case userIDGet == 0:
			userIDGet = my.QSimple("select id from users order by id desc limit 1;")
			if userIDGet > 0 {
				userIDGet++
				idxUsers = userIDGet
			} else {
				log.Println("Error get max user ID")
				return
			}
		}

		queryx = fmt.Sprintf("INSERT INTO z_cache_users (id, username, digesta1)\n\tVALUES (%d, '%s', '%s');", usID, usName, usPass)
		queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_principals (id, uri, email, displayname, vcardurl)\n\tVALUES (%d, 'principals/%s', NULL, NULL, NULL);", queryx, usID, usName)
		queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_addressbooks (id, principaluri, uri, ctag)\n\tVALUES (%d, 'principals/%s', 'default', 1); select id from users order by id desc limit 1", queryx, usID, usName)
		//log.Printf("%s\n", queryx)
		_, err = my.D.Query(queryx)
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
		for pgrows2.Next() {
			pgrows2.Scan(&usDN)
			log.Printf("\t\t\t%3d/%s - %s\n", usID, usName, usDN)

			search := ldap.NewSearchRequest(usDN, 2, ldap.NeverDerefAliases, 0, 0, false, conf.Conf.LDAP_URL[0][4], ldap_Attr, nil)

			sr, err := ld.D.Search(search)
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
					queryx = fmt.Sprintf("%s\nINSERT INTO z_cache_cards (id, addressbookid, carddata, uri, lastmodified)\n\tVALUES (%d, %d, '%s', '%s.vcf', NULL);", queryx, idxCards, usID, y, uid)
					if multiCount > multiInsert {
						_, err = my.D.Query(queryx)
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
			_, err = my.D.Query(queryx)
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

	rLog.Log("\t\tComplete!")

	if workMode == "PART" {
		rLog.Log("\tUpdate tables in PartialUpdate mode...")
		for j := 0; j < len(usIDArray); j++ {
			log.Printf("\t\t\tUpdate %d/%s...\n", usIDArray[j].id, usIDArray[j].name)
			for i := 0; i < len(mySQL_Update_part1); i++ {
				log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update_part1))
				if my.QSimple(strings.Replace(mySQL_Update_part1[i], "XYZIDXYZ", fmt.Sprintf("%d", usIDArray[j].id), -1)) < 0 {
					log.Println("MySQL: Update error")
					return
				}
				time.Sleep(2 * time.Second)
			}
			for i := 0; i < len(mySQL_Update1); i++ {
				log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update1))

				if my.QSimple(mySQL_Update1[i]) < 0 {
					log.Println("MySQL: Update error")
					return
				}
				time.Sleep(2 * time.Second)
			}
			for i := 0; i < len(mySQL_Update_part2); i++ {
				log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update_part2))
				if my.QSimple(strings.Replace(mySQL_Update_part2[i], "XYZIDXYZ", fmt.Sprintf("%d", usIDArray[j].id), -1)) < 0 {
					log.Println("MySQL: Update error")
					return
				}
				time.Sleep(2 * time.Second)
			}
			for i := 0; i < len(mySQL_Update2); i++ {
				log.Printf("\t\t\tstep %d (%d of %d)...\n", j+1, i+1, len(mySQL_Update2))
				if my.QSimple(mySQL_Update2[i]) < 0 {
					log.Println("MySQL: Update error")
					return
				}
				time.Sleep(2 * time.Second)
			}
			time.Sleep(2 * time.Second)
		}
	} else {
		rLog.Log("\tUpdate tables...")
		for i := 0; i < len(mySQL_Update_full1); i++ {
			log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update_full1))
			if my.QSimple(mySQL_Update_full1[i]) < 0 {
				log.Println("MySQL: Update error")
				return
			}
			time.Sleep(2 * time.Second)
		}
		for i := 0; i < len(mySQL_Update1); i++ {
			log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update1))
			if my.QSimple(mySQL_Update1[i]) < 0 {
				log.Println("MySQL: Update error")
				return
			}
			time.Sleep(2 * time.Second)
		}
		for i := 0; i < len(mySQL_Update_full2); i++ {
			log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update_full2))
			if my.QSimple(mySQL_Update_full2[i]) < 0 {
				log.Println("MySQL: Update error")
				return
			}
			time.Sleep(2 * time.Second)
		}
		for i := 0; i < len(mySQL_Update2); i++ {
			log.Printf("\t\t\tstep %d of %d...\n", i+1, len(mySQL_Update2))
			if my.QSimple(mySQL_Update2[i]) < 0 {
				log.Println("MySQL: Update error")
				return
			}
			time.Sleep(2 * time.Second)
		}
	}

	rLog.Log("\t\tComplete!")

	rLog.Log("\tClean NeedToUpdate table...")

	if pg.QSimple("delete from aaa_dav_ntu where userid=0 or updtime<", time_now, ";") < 0 {
		log.Println("PgSQL: Clean NTU table error")
		return
	}

	rLog.Log("\tComplete!")
	rLog.Bye()
}

func main() {
	var (
		jsonConfig SBMSystem.ReadJSONConfig
		rLog       SBMSystem.LogFile
		pid        SBMSystem.PidFile
		sleepWatch = int(0)
	)

	const (
		pName = string("SABook CardDAVMaker")
		pVer  = string("5 2015.11.05.21.00")
	)

	fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

	jsonConfig.Init()

	rLog.ON(jsonConfig)
	pid.ON(jsonConfig)
	pid.OFF(jsonConfig)
	rLog.OFF()

	SBMSystem.Fork(jsonConfig)
	SBMSystem.Signal(jsonConfig, pid)

	rLog.ON(jsonConfig)
	pid.ON(jsonConfig)
	defer pid.OFF(jsonConfig)
	rLog.Hello(pName, pVer)
	rLog.OFF()

	for {
		rLog.ON(jsonConfig)
		jsonConfig.Update()

		if checkNTUWishes(jsonConfig, rLog) > 0 {
			rLog.Hello(pName, pVer)
			goNTUWork(jsonConfig, rLog)
			sleepWatch = 0
		}

		if sleepWatch > 3600 {
			rLog.Log("<-- I'm alive ... :)")
			sleepWatch = 0
		}

		rLog.OFF()
		time.Sleep(time.Duration(jsonConfig.Conf.Sleep_Time) * time.Second)
		sleepWatch += jsonConfig.Conf.Sleep_Time
	}
}
