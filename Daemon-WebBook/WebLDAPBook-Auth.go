package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"strconv"
	"strings"
	//"syscall"
	//"os"
	"time"

	//LDAP
	//"github.com/go-ldap/ldap"

	"database/sql"
	// PostgreSQL
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"
	//	"github.com/kabukky/httpscerts"
	//	"github.com/gavruk/go-blog-example/models"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		username  string
		password  string
		queryx    string
		get_login string
	)

	username = r.FormValue("username")
	password = r.FormValue("password")

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
		/*		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
				if err != nil {
					log.Printf("PG::Open() error: %v\n", err)
					return
				}
				defer dbpg.Close()
		*/

		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Fatalf("PG_INIT::Open() error: %v\n", err)
		}

		defer dbpg.Close()

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

}

func GenerateMD5(user string) string {
	z := md5.New()
	z.Write([]byte(fmt.Sprintf("%s", user)))
	return hex.EncodeToString(z.Sum(nil))
}

func GenerateSessionId(user string) string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", GenerateMD5(user), b)
}

func StoreUserSession(user string, w http.ResponseWriter) string {

	var err error

	userID := GenerateSessionId(user)
	user_cookie := &http.Cookie{
		Name:    COOKIE_ID,
		Value:   userID,
		Expires: time.Now().Add(time.Duration(rconf.WLB_SessTimeOut) * time.Minute),
	}

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	_, err = dbpg.Query(fmt.Sprintf("insert into wb_auth_session (username, sessname, exptime) values ('%s','%s',%v);", user, userID, time.Now().Unix()))
	if err != nil {
		log.Printf("PG::Query() Insert session error: %v\n", err)
		return "error"
	}
	//fmt.Printf("%v / %v\n", user_cookie, time.Now().Unix())
	http.SetCookie(w, user_cookie)
	return userID
}

func UpdateUserSession(sess string, w http.ResponseWriter) string {

	var err error

	user_cookie := &http.Cookie{
		Name:    COOKIE_ID,
		Value:   sess,
		Expires: time.Now().Add(time.Duration(rconf.WLB_SessTimeOut) * time.Minute),
	}

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	_, err = dbpg.Query(fmt.Sprintf("update wb_auth_session set exptime='%v' where sessname='%s';", time.Now().Unix(), sess))
	if err != nil {
		log.Printf("PG::Query() Update session error: %v\n", err)
		return "error"
	}
	//fmt.Printf("%v / %v\n", user_cookie, time.Now().Unix())
	http.SetCookie(w, user_cookie)
	return "Ok"
}

func CheckUserSession(r *http.Request, w http.ResponseWriter) (string, int) {
	var (
		get_user = string("")
		get_time = int(0)
		get_role = int(0)
	)

	cookie, _ := r.Cookie(COOKIE_ID)
	if cookie != nil {

		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Fatalf("PG_INIT::Open() error: %v\n", err)
		}

		defer dbpg.Close()

		rows, err := dbpg.Query(fmt.Sprintf("select username,exptime from wb_auth_session where sessname='%s';", cookie.Value))
		if err != nil {
			log.Printf("PG::Query() Select session error: %v\n", err)
			return "error", 101
		}

		rows.Next()
		rows.Scan(&get_user, &get_time)

		if len(get_user) > 0 && get_time > 0 {
			now_time := time.Now().Unix()

			//fmt.Printf("%v / %s / %d / %v / %v\n", cookie.Value, get_user, get_time, int64(get_time+rconf.WLB_SessTimeOut*60), int64(now_time+10))

			if int64(get_time+rconf.WLB_SessTimeOut*60) < int64(now_time+10) {
				_, err = dbpg.Query(fmt.Sprintf("delete from wb_auth_session where sessname='%s';", cookie.Value))
				if err != nil {
					log.Printf("PG::Query() Delete session error: %v\n", err)
					return "error", 102
				}
			} else {
				rows, err = dbpg.Query(fmt.Sprintf("select role from aaa_logins where login='%s';", get_user))
				if err != nil {
					log.Printf("PG::Query() Select role error: %v\n", err)
					return "error", 103
				}
				rows.Next()
				rows.Scan(&get_role)
				if get_role > 0 {
					if UpdateUserSession(cookie.Value, w) != "error" {
						return get_user, get_role
					}
				}
			}
		}

	}

	return "guest", 0
}

func RemoveUserSession(r *http.Request) {

	cookie, _ := r.Cookie(COOKIE_ID)
	if cookie != nil {
		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Fatalf("PG_INIT::Open() error: %v\n", err)
		}

		defer dbpg.Close()

		_, err = dbpg.Query(fmt.Sprintf("delete from wb_auth_session where sessname='%s';", cookie.Value))
		if err != nil {
			log.Printf("PG::Query() Delete session error: %v\n", err)
		}

	}
}
