package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"database/sql"
	// PostgreSQL
	_ "github.com/lib/pq"

	//"github.com/BestianRU/SABookServices/SABModules"
)

func modifyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		uid         string
		login       string
		action      string
		password_x1 string
		password_x2 string
		fullname    string
		role        string
		rolen       int
		uFullname   string
		xId         int
	)

	uid = r.FormValue("uid")
	action = r.FormValue("action")
	login = r.FormValue("login")
	password_x1 = r.FormValue("password_x1")
	password_x2 = r.FormValue("password_x2")
	fullname = r.FormValue("fullname")
	role = r.FormValue("role")

	remIPClient := getIPAddress(r)

	//mt.Printf("%s %s %s %s %s %s %s\n", uid, action, login, password_x1, password_x2, fullname, role)

	if len(uid) > 10 {

		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Fatalf("PG_INIT::Open() error: %v\n", err)
		}

		defer dbpg.Close()

		queryx := fmt.Sprintf("select fullname from ldapx_persons where uid='%s';", uid)
		rows, err := dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Change password: %v\n", err)
			return
		}
		rows.Next()
		rows.Scan(&uFullname)

		queryx = fmt.Sprintf("select id from aaa_logins where uid='%s';", uid)
		rows, err = dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Change password: %v\n", err)
			return
		}
		rows.Next()
		rows.Scan(&xId)

		switch action {
		case "change_password":
			if len(password_x1) > 4 && len(password_x2) > 4 && password_x1 == password_x2 && len(login) > 4 {
				queryx = fmt.Sprintf("update aaa_logins set password=md5('%s:%s:%s') where uid='%s';", login, rconf.SABRealm, password_x1, uid)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Get User Name for UID: %v\n", err)
					return
				}
				log.Printf("%s AAA Change password for: %s", remIPClient, uFullname)
				queryx = fmt.Sprintf("insert into aaa_dav_ntu (userid,updtime) select %d,%v where not exists (select userid from aaa_dav_ntu where userid=%d); update aaa_dav_ntu set updtime=%v where userid=%d;", xId, time.Now().Unix(), xId, time.Now().Unix(), xId)
				//fmt.Printf("%s\n", queryx)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Update NTU table: %v\n", err)
					return
				}
			}
		case "create_user":
			if len(password_x1) > 4 && len(password_x2) > 4 && password_x1 == password_x2 && len(login) > 4 && len(fullname) > 4 && len(role) > 0 {
				switch role {
				case "admin":
					rolen = roleAdmin
				case "user":
					rolen = roleUser
				default:
					rolen = 0
				}
				//queryx := fmt.Sprintf("insert into aaa_logins (id,login,fullname,password,role,uid) select id+1,'%s','%s',md5('%s:%s:%s'),%d,'%s' from aaa_logins order by id desc limit 1;", login, fullname, login, rconf.SABRealm, password_x1, rolen, uid)
				queryx = fmt.Sprintf("insert into aaa_logins (login,fullname,password,role,uid) values ('%s','%s',md5('%s:%s:%s'),%d,'%s');", login, fullname, login, rconf.SABRealm, password_x1, rolen, uid)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Create user: %v\n", err)
					return
				}
				log.Printf("%s AAA Create SAB account for: %s", remIPClient, uFullname)
			}
		case "change_role":
			if len(role) > 0 {
				switch role {
				case "admin":
					rolen = roleAdmin
				case "user":
					rolen = roleUser
				default:
					rolen = 0
				}
				queryx = fmt.Sprintf("update aaa_logins set role=%d where uid='%s';", rolen, uid)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Change role: %v\n", err)
					return
				}
				log.Printf("%s AAA Change role for: %s", remIPClient, uFullname)
			}
		case "delete_sab_login":
			queryx = fmt.Sprintf(`
delete from wb_auth_session where username in (select login from aaa_logins where uid='%s');
delete from aaa_dns where userid in (select id from aaa_logins where uid='%s');
delete from aaa_logins where uid='%s';
				`, uid, uid, uid)
			_, err = dbpg.Query(queryx)
			if err != nil {
				log.Printf("%s\n", queryx)
				log.Printf("PG::Query() Change role: %v\n", err)
				return
			}
			log.Printf("%s AAA Delete SAB account for: %s", remIPClient, uFullname)
			queryx = fmt.Sprintf("insert into aaa_dav_ntu (userid,updtime) select %d,%v where not exists (select userid from aaa_dav_ntu where userid=%d); update aaa_dav_ntu set updtime=%v where userid=%d;", xId, time.Now().Unix(), xId, time.Now().Unix(), xId)
			//fmt.Printf("%s\n", queryx)
			_, err = dbpg.Query(queryx)
			if err != nil {
				log.Printf("%s\n", queryx)
				log.Printf("PG::Query() Update NTU table: %v\n", err)
				return
			}
		}

	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)
}
