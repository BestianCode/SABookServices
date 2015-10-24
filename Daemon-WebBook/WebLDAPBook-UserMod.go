package main

import (
	//"crypto/md5"
	//"crypto/rand"
	//"encoding/hex"
	"fmt"
	//"html/template"
	"log"
	"net/http"
	//"strconv"
	//"strings"
	//"syscall"
	//"os"
	//"time"

	//LDAP
	//"github.com/go-ldap/ldap"

	//"database/sql"
	// PostgreSQL
	//_ "github.com/lib/pq"
	// SQLite
	//_ "github.com/mattn/go-sqlite3"

	//"github.com/BestianRU/SABookServices/SABModules"
	//	"github.com/kabukky/httpscerts"
	//	"github.com/gavruk/go-blog-example/models"
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
	)

	uid = r.FormValue("uid")
	action = r.FormValue("action")
	login = r.FormValue("login")
	password_x1 = r.FormValue("password_x1")
	password_x2 = r.FormValue("password_x2")
	fullname = r.FormValue("fullname")
	role = r.FormValue("role")

	//mt.Printf("%s %s %s %s %s %s %s\n", uid, action, login, password_x1, password_x2, fullname, role)

	if len(uid) > 10 {
		switch action {
		case "change_password":
			if len(password_x1) > 4 && len(password_x2) > 4 && password_x1 == password_x2 && len(login) > 4 {
				queryx := fmt.Sprintf("update aaa_logins set password=md5('%s:%s:%s') where uid='%s';", login, rconf.SABRealm, password_x1, uid)
				_, err := dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Change password: %v\n", err)
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
				queryx := fmt.Sprintf("insert into aaa_logins (id,login,fullname,password,role,uid) select id+1,'%s','%s',md5('%s:%s:%s'),%d,'%s' from aaa_logins order by id desc limit 1;", login, fullname, login, rconf.SABRealm, password_x1, rolen, uid)
				_, err := dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Create user: %v\n", err)
					return
				}
			}
		}

	}

	http.Redirect(w, r, r.Referer(), http.StatusMovedPermanently)

	/*
		RedirectDN := r.FormValue("RPR")

		if len(RedirectDN) < 1 {
			RedirectDN = "/"
		} else {
			RedirectDN = strings.Replace(RedirectDN, "'", "", -1)
		}

		if r.FormValue("go") == "unlogin" {
			RemoveUserSession(r)
			http.Redirect(w, r, RedirectDN, http.StatusMovedPermanently)
		}

		remIPClient := getIPAddress(r)

		SABModules.Log_ON(&rconf)
		defer SABModules.Log_OFF()

		if len(username) < 2 || len(password) < 2 {

			log.Printf("%s AAA Get login form...", remIPClient)

			t, err := template.ParseFiles("templates/header.html")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				log.Println(err.Error())
				return
			}

			t.ExecuteTemplate(w, "header", template.FuncMap{"Pagetitle": rconf.WLB_HTML_Title, "FRColor": "#FF0000", "BGColor": "#FFEEEE"})

			t, err = template.ParseFiles("templates/search.html")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				log.Println(err.Error())
				return
			}

			t.ExecuteTemplate(w, "search", template.FuncMap{"GoHome": "Yes", "LineColor": "#FFDDDD"})

			t, err = template.ParseFiles("templates/login.html")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				log.Println(err.Error())
				return
			}

			t.ExecuteTemplate(w, "login", template.FuncMap{"RedirectDN": RedirectDN})

			t, err = template.ParseFiles("templates/footer.html")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				log.Println(err.Error())
				return
			}

			t.ExecuteTemplate(w, "footer", template.FuncMap{"WebBookVersion": pVersion, "xMailBT": rconf.WLB_MailBT, "LineColor": "#FFDDDD"})

		} else {
			queryx = fmt.Sprintf("select distinct login from aaa_logins where login='%s' and password=md5('%s:%s:%s') limit 1;", username, username, rconf.SABRealm, password)
			//fmt.Printf("%s\n", queryx)
			rows, err := dbpg.Query(queryx)
			if err != nil {
				log.Printf("PG::Query() Check login and password: %v\n", err)
				return
			}

			rows.Next()
			rows.Scan(&get_login)

			if get_login == username {
				userID := StoreUserSession(username, w)
				if userID == "error" {
					return
				}

				log.Printf("%s AAA Login enter with username %s (%s)\n", remIPClient, username, userID)

				http.Redirect(w, r, RedirectDN, http.StatusMovedPermanently)

			} else {
				log.Printf("%s AAA Login ERROR with username %s\n", remIPClient, username)
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			}
		}
	*/
}
