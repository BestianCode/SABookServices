package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	//LDAP
	"github.com/go-ldap/ldap"

	"database/sql"
	// PostgreSQL
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"
	//	"github.com/kabukky/httpscerts"
	//	"github.com/gavruk/go-blog-example/models"
)

type tList struct {
	URL         string
	URLName     string
	ORGName     string
	USERName    string
	FullName    string
	PhoneInt    string
	PhoneExt    string
	Mobile      string
	Mail        string
	Position    string
	ADLogin     string
	ADDomain    string
	AdminMode   string
	UID         string
	AAALogin    string
	AAAPassword string
	AAAFullName string
	AAARole     string
	NewSABLogin string
	DavDN       string
}

const (
	pName     = string("Web Address Book")
	pVer      = string("4 alpha 2015.10.30.21.00")
	userLimit = 20
	COOKIE_ID = "SABookSessionID"
	roleAdmin = 100
	roleUser  = 1
)

var (
	def_config_file = string("./WebLDAPBook.json")             // Default configuration file
	def_log_file    = string("/var/log/ABook/WebLDAPBook.log") // Default log file
	def_daemon_mode = string("NO")                             // Default start in foreground
	pVersion        string
	rconf           SABModules.Config_STR
	ldap_count      = int(0)
	sleepTime       = 60
	//dbpg            *sql.DB
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

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	_, err = dbpg.Exec(query_create)
	if err != nil {
		log.Fatalf("PG_INIT::Exec() create tables error: %v\n", err)
	}

}

func initLDAPConnector() string {
	var (
		ckl = int(0)
		err error
		l   *ldap.Conn
	)

	for {
		if ckl > 9 {
			log.Printf("LDAP Init SRV ***** Error connect to all LDAP servers...")
			return "error"
		}

		ldap_count++
		if ldap_count > len(rconf.LDAP_URL)-1 {
			ldap_count = 0
		}

		log.Printf("LDAP Init SRV ***** Trying connect to server %d of %d: %s", ldap_count+1, len(rconf.LDAP_URL), rconf.LDAP_URL[ldap_count][0])
		l, err = ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
		if err != nil {
			continue
		}

		defer l.Close()

		break

		ckl++
	}
	return rconf.LDAP_URL[ldap_count][0]
}

// ====================================================================================================

func getMore(remIPClient string, fField map[string]string, fType string, l *ldap.Conn, dnList map[string]tList, setAdminMode string) {
	var (
		fPath            string
		fURL             string
		fURLName         string
		ckl1, ckl2, ckl3 int
		ldap_Attr        []string
		aaa_login        = string("")
		aaa_password     = string("")
		aaa_fullname     = string("")
		aaa_role         = string("")
		newSABLogin      string
		get_davdn        = string("")
		err              error
	)

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	if fField["DN"] != "" && (fField["USERName"] != "" || fField["ORGName"] != "") {
		fPath = fField["DN"]
		fPath = strings.Replace(strings.ToLower(fPath), ","+strings.ToLower(rconf.LDAP_URL[ldap_count][3]), "", -1)
		fPath_Split := strings.Split(fPath, ",")
		fURLName = ""
		for ckl1 = 0; ckl1 < len(fPath_Split)-1; ckl1++ {
			fPath_Strip := ""
			for ckl2 = ckl1 + 1; ckl2 < len(fPath_Split); ckl2++ {
				fPath_Strip = fmt.Sprintf("%s%s,", fPath_Strip, fPath_Split[ckl2])
			}
			if fType == "User" {
				fPath_Strip = fmt.Sprintf("%s%s", fPath_Strip, rconf.LDAP_URL[ldap_count][3])
				if ckl1 == 0 {
					fURL = fPath_Strip
				}
				//						log.Printf("X1: %s", fPath_Strip)
				subsearch := ldap.NewSearchRequest(fPath_Strip, 0, ldap.NeverDerefAliases, 0, 0, false, rconf.LDAP_URL[ldap_count][4], ldap_Attr, nil)
				subsr, err := l.Search(subsearch)
				if err != nil {
					//								fmt.Fprintf(w, err.Error())
					log.Printf("LDAP::Search() error: %v\n", err)
				}

				//						log.Printf("Y1: %s / %s / %d\n", fPath_Strip, rconf.LDAP_URL[ldap_count][4], len(subsr.Entries))
				if len(subsr.Entries) > 0 {
					for _, subentry := range subsr.Entries {
						for _, subattr := range subentry.Attributes {
							for ckl3 = 0; ckl3 < len(rconf.WLB_LDAP_ATTR); ckl3++ {
								if subattr.Name == rconf.WLB_LDAP_ATTR[ckl3][0] {
									if rconf.WLB_LDAP_ATTR[ckl3][1] == "ORGName" {
										if ckl1 == 0 {
											fURLName = fmt.Sprintf("%s", strings.Join(subattr.Values, ","))
										} else {
											fURLName = fmt.Sprintf("%s / %s", strings.Join(subattr.Values, ","), fURLName)
										}
										//												log.Printf("Z1: %s", fURLName)
									}
								}
							}
						}
					}
				}
			}
		}

		fField["DN"] = strings.Replace(strings.ToLower(fField["DN"]), "/", ",", -1)
		fmt.Sprintf("/Go%s?dn=%s", fType, fField["DN"])
		fField["DN"] = fmt.Sprintf("/Go%s?dn=%s", fType, fField["DN"])
		fURL = fmt.Sprintf("/Go%s?dn=%s", fType, fURL)
		log.Printf("%s <-- %s", remIPClient, fField["DN"])
		davDN := "LIST:\n"
		if setAdminMode == "Yes" {
			queryx := fmt.Sprintf("select x.dn from aaa_dns as x, aaa_logins as y where y.uid='%s' and y.id=x.userid;", fField["UID"])
			//fmt.Printf("%s\n", queryx)
			rows, err := dbpg.Query(queryx)
			if err != nil {
				log.Printf("PG::Query() Select info from aaa_logins: %v\n", err)
				return
			}

			for rows.Next() {
				rows.Scan(&get_davdn)
				davDN = fmt.Sprintf("%s%s\n", davDN, get_davdn)
			}

			queryx = fmt.Sprintf("select login,password,fullname,role from aaa_logins where uid='%s';", fField["UID"])
			//fmt.Printf("%s\n", queryx)
			rows, err = dbpg.Query(queryx)
			if err != nil {
				log.Printf("PG::Query() Select info from aaa_logins: %v\n", err)
				return
			}

			rows.Next()
			rows.Scan(&aaa_login, &aaa_password, &aaa_fullname, &aaa_role)
			if len(aaa_password) > 0 {
				aaa_password = "Ok"
			} else {
				aaa_password = "NONE!"
			}
			xt, _ := strconv.Atoi(aaa_role)
			switch xt {
			case roleAdmin:
				aaa_role = "Administrator"
			case roleUser:
				aaa_role = "User"
			default:
				aaa_role = "Guest"
			}
			if len(fField["ADLogin"]) < 2 || len(fField["ADDomain"]) < 2 {
				newSABLogin = fField["Mail"]
			} else {
				newSABLogin = fmt.Sprintf("%s@%s", fField["ADLogin"], fField["ADDomain"])
			}
		}

		dnList[fField["DN"]] = tList{URL: fURL, URLName: fURLName, ORGName: fField["ORGName"], USERName: fField["USERName"], FullName: fField["FullName"], Position: fField["Position"], PhoneInt: fField["PhoneInt"], Mobile: fField["Mobile"], PhoneExt: fField["PhoneExt"], Mail: fField["Mail"], ADLogin: fField["ADLogin"], ADDomain: fField["ADDomain"], AdminMode: setAdminMode, UID: fField["UID"], AAALogin: aaa_login, AAAPassword: aaa_password, AAAFullName: aaa_fullname, AAARole: aaa_role, NewSABLogin: newSABLogin, DavDN: davDN}
		//fmt.Printf("%v\n", dnList)
	}
}

func getIPAddress(r *http.Request) string {
	return fmt.Sprintf("%s (%v)", strings.Split(r.RemoteAddr, ":")[0], strings.Trim(strings.Trim(strings.Replace(r.Header.Get("X-FORWARDED-FOR"), "127.0.0.1", "", -1), " "), ","))
}

func getLDAPdnList(l *ldap.Conn, dn string, xCount int, xCountMax int, w http.ResponseWriter, r *http.Request, dnList []string) {
	var (
		queryx   string
		o        string
		checkSet string
	)

	_, userperm := CheckUserSession(r, w)

	switch userperm {
	case roleAdmin:
	default:
		return
	}

	xCount++

	if xCount > xCountMax {
		return
	}

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	t, err := template.ParseFiles("templates/tree-01.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}
	t.ExecuteTemplate(w, "tree-01", nil)

	queryx = fmt.Sprintf("select uid,name from ldapx_institutes where idparent='%s' order by name;", dn)
	rows, err := dbpg.Query(queryx)
	if err != nil {
		log.Printf("%s\n", queryx)
		log.Printf("PG::Query() Get Children for DN: %v\n", err)
		return
	}

	for rows.Next() {
		rows.Scan(&dn, &o)
		randomCheckID := make([]byte, 16)
		rand.Read(randomCheckID)
		//fmt.Printf("\tList DN: %x %s (%s)\n", randomCheckID, dn, o)
		checkSet = "unchecked"
		for i := 0; i < len(dnList); i++ {
			if dn == dnList[i] {
				checkSet = "checked"
				//fmt.Printf("%s / %s\n", dn, dnList[i])
			}
			//fmt.Printf("%s / %s\n", dn, dnList[i])
		}
		t, err = template.ParseFiles("templates/tree-02.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
			return
		}
		t.ExecuteTemplate(w, "tree-02", template.FuncMap{"RANDOM": fmt.Sprintf("%x", randomCheckID), "DNName": dn, "ORG": o, "CHECKX": checkSet})
		//fmt.Printf("\tList DN: %x %s (%s) %s\n", randomCheckID, dn, o, checkORnot)
		getLDAPdnList(l, dn, xCount, xCountMax, w, r, dnList)
		t, err = template.ParseFiles("templates/tree-08.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
			return
		}
		t.ExecuteTemplate(w, "tree-08", nil)
		/*}
		}*/
	}
	t, err = template.ParseFiles("templates/tree-09.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}
	t.ExecuteTemplate(w, "tree-09", nil)
}

func davDNHandler(w http.ResponseWriter, r *http.Request) {
	var (
		uid       string
		uAction   string
		uFullname string
		queryx    string
		err       error
		l         *ldap.Conn
		xFRColor  = string("#FFFFFF")
		xBGColor  = string("#FFFFFF")
	)

	_, userperm := CheckUserSession(r, w)

	switch userperm {
	case roleAdmin:
		//xFRColor = "#FF0000"
		//xBGColor = "#FFFFFF"
	default:
		return
	}

	SABModules.Log_ON(&rconf)
	defer SABModules.Log_OFF()

	remIPClient := getIPAddress(r)

	uid = r.FormValue("uid")
	uAction = r.FormValue("action")

	dbpg, err := sql.Open("postgres", rconf.PG_DSN)
	if err != nil {
		log.Fatalf("PG_INIT::Open() error: %v\n", err)
	}

	defer dbpg.Close()

	queryx = fmt.Sprintf("select fullname from ldapx_persons where uid='%s';", uid)
	rows, err := dbpg.Query(queryx)
	if err != nil {
		log.Printf("%s\n", queryx)
		log.Printf("PG::Query() Get User Name for UID: %v\n", err)
		return
	}
	rows.Next()
	rows.Scan(&uFullname)

	if uAction == "SaveDN" && len(uFullname) > 0 {

		queryx = fmt.Sprintf("select id from aaa_logins where uid='%s';", uid)
		rows, err := dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Get User ID for UID: %v\n", err)
			return
		}

		xId := 0
		rows.Next()
		rows.Scan(&xId)

		queryx = fmt.Sprintf("delete from aaa_dns where userid=%d;", xId)
		_, err = dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Delete DNs for UID: %v\n", err)
			return
		}
		time.Sleep(time.Second)
		r.ParseForm()
		for parName := range r.Form {
			if strings.Contains(parName, "SaveDN") {

				queryx = fmt.Sprintf("insert into aaa_dns (userid,dn) select %d, dn from ldap_entries where uid='%s';", xId, r.FormValue(parName))
				//fmt.Printf("%s\n", queryx)
				_, err = dbpg.Query(queryx)
				if err != nil {
					log.Printf("%s\n", queryx)
					log.Printf("PG::Query() Insert DNs for UID: %v\n", err)
					return
				}
				//fmt.Fprintf(w, "%v\n", r.FormValue(parName))
			}
		}

		queryx = fmt.Sprintf("insert into aaa_dav_ntu (userid,updtime) select %d,%v where not exists (select userid from aaa_dav_ntu where userid=%d); update aaa_dav_ntu set updtime=%v where userid=%d;", xId, time.Now().Unix(), xId, time.Now().Unix(), xId)
		//fmt.Printf("%s\n", queryx)
		_, err = dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Update NTU table: %v\n", err)
			return
		}

		log.Printf("%s --> Set DavDN List for %s", remIPClient, uFullname)

		time.Sleep(time.Second)
		fmt.Fprintf(w, "<script type=\"text/javascript\">window.close();</script>")

	} else {

		log.Printf("%s <-- Get DavDN List for %s", remIPClient, uFullname)

		if initLDAPConnector() == "error" {
			return
		}

		l, err = ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Printf("LDAP::Initialize() error: %v\n", err)
			return
		}

		//		l.Debug = true
		defer l.Close()

		log.Printf("%s =!= Connected to server %d of %d: %s", remIPClient, ldap_count+1, len(rconf.LDAP_URL), rconf.LDAP_URL[ldap_count][0])

		err = l.Bind(rconf.LDAP_URL[ldap_count][1], rconf.LDAP_URL[ldap_count][2])
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Printf("LDAP::Bind() error: %v\n", err)
			return
		}

		t, err := template.ParseFiles("templates/header.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
			return
		}
		t.ExecuteTemplate(w, "header", template.FuncMap{"Pagetitle": rconf.WLB_HTML_Title, "FRColor": xFRColor, "BGColor": xBGColor, "TREEOn": "Yes"})

		t, err = template.ParseFiles("templates/tree-00.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
			return
		}
		t.ExecuteTemplate(w, "tree-00", template.FuncMap{"UID": uid})

		queryx = fmt.Sprintf("select distinct uid from ldap_entries where lower(dn)=lower('%s') limit 1;", rconf.LDAP_URL[ldap_count][3])
		rows, err := dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Get UID for DN: %v\n", err)
			return
		}
		rows.Next()
		uidDN := ""
		rows.Scan(&uidDN)

		queryx = fmt.Sprintf("select z.uid from aaa_dns as x, aaa_logins as y, ldap_entries as z where x.userid=y.id and x.dn=z.dn and y.uid='%s';", uid)
		//fmt.Printf("%s\n", queryx)
		rows, err = dbpg.Query(queryx)
		if err != nil {
			log.Printf("%s\n", queryx)
			log.Printf("PG::Query() Get UID for DN: %v\n", err)
			return
		}

		dnList := make([]string, 0)

		x := ""
		for rows.Next() {
			rows.Scan(&x)
			dnList = append(dnList, x)
			//fmt.Printf("%s\n", x)
		}

		//fmt.Printf("%v\n", dnList)

		getLDAPdnList(l, uidDN, 0, rconf.WLB_DavDNTreeDepLev, w, r, dnList)

		t, err = template.ParseFiles("templates/tree-10.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Println(err.Error())
			return
		}
		t.ExecuteTemplate(w, "tree-10", nil)

	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var (
		xSearchPplMode = int(0)
		xSearch        string
		xMessage       string

		dn          string
		dn_back     string
		dn_back_tmp []string

		go_home_button string

		ldap_Search string

		ldapSearchMode = int(1)

		ckl1 int

		ldap_Attr []string

		xGetDN  [1000]string
		xGetCkl int

		l   *ldap.Conn
		err error

		xFRColor     = string("#FFFFFF")
		xBGColor     = string("#FFFFFF")
		LUserName    = string("")
		setAdminMode = string("")
	)

	username, userperm := CheckUserSession(r, w)

	//fmt.Printf("%s / %d\n", username, userperm)

	switch userperm {
	case roleAdmin:
		xFRColor = "#FF0000"
		xBGColor = "#FFFFFF"
		LUserName = username
		setAdminMode = "Yes"
	case roleUser:
		xFRColor = "#0000FF"
		xBGColor = "#FFFFFF"
		LUserName = username
	default:
		xFRColor = "#FFFFFF"
		xBGColor = "#FFFFFF"
		LUserName = ""
	}

	ldap_Attr = make([]string, len(rconf.WLB_LDAP_ATTR))

	for ckl1 := 0; ckl1 < len(rconf.WLB_LDAP_ATTR); ckl1++ {
		ldap_Attr[ckl1] = rconf.WLB_LDAP_ATTR[ckl1][0]
	}

	SABModules.Log_ON(&rconf)
	defer SABModules.Log_OFF()

	get_dn := r.FormValue("dn")
	get_cn := r.FormValue("cn")
	get_fn := r.FormValue("FirstName")
	get_ln := r.FormValue("LastName")
	searchMode := r.FormValue("SearchMode")

	remIPClient := getIPAddress(r)
	//	log.Printf("DN: %s --- CN: %s", get_dn, get_cn)

	if get_dn == "" {
		dn = rconf.LDAP_URL[ldap_count][3]
	} else {
		dn = get_dn
	}

	if len(dn) < len(rconf.LDAP_URL[ldap_count][3]) {
		dn = rconf.LDAP_URL[ldap_count][3]
	}

	log.Printf("->")
	log.Printf("--> %s", pVersion)
	log.Printf("->")
	ucurl, _ := strconv.Unquote(r.RequestURI)
	log.Println(remIPClient + " --> http://" + r.Host + ucurl)
	log.Printf("%s ++> DN: %s / CN: %s / Mode: %d / Def.DN: %s", remIPClient, dn, ldap_Search, ldapSearchMode, rconf.LDAP_URL[ldap_count][3])

	if get_cn == "" && get_ln == "" && get_fn == "" {
		ldap_Search = rconf.LDAP_URL[ldap_count][4]
	} else {
		log.Printf("%s ++> SQL Search: %s/%s/%s\n", remIPClient, get_cn, get_fn, get_ln)
		dbpg, err := sql.Open("postgres", rconf.PG_DSN)
		if err != nil {
			log.Fatalf("PG_INIT::Open() error: %v\n", err)
		}

		defer dbpg.Close()

		queryx := "select x.dn from ldap_entries as x, ldapx_persons as y where x.uid=y.uid"
		if len(get_cn) > 2 {
			queryx = fmt.Sprintf("%s and lower(fullname) like lower('%%%s%%')", queryx, strings.ToLower(get_cn))
		}
		if len(get_ln) > 2 {
			queryx = fmt.Sprintf("%s and lower(surname) like lower('%s%%')", queryx, strings.ToLower(get_ln))
		}
		if len(get_fn) > 2 {
			queryx = fmt.Sprintf("%s and lower(name) like lower('%s%%')", queryx, strings.ToLower(get_fn))
		}
		if len(get_cn) <= 2 && len(get_ln) <= 2 && len(get_fn) <= 2 {
			queryx = fmt.Sprintf("%s and 2=3;", queryx)
		} else {
			queryx = fmt.Sprintf("%s;", queryx)
		}
		//		log.Printf("Search QUERY: %s\n", queryx)
		rows, err := dbpg.Query(queryx)
		if err != nil {
			fmt.Printf("SQL Error: %s\n", queryx)
			log.Printf("PG::Query() Check LDAP tables error: %v\n", err)
			return
		}
		xGetCkl = 0
		for rows.Next() {
			rows.Scan(&xGetDN[xGetCkl])
			//fmt.Println("XXX:", xGetDN[xGetCkl], dn)
			if strings.Contains(strings.ToLower(xGetDN[xGetCkl]), strings.ToLower(dn)) || searchMode == "Full" {
				log.Printf("%s <-- SQL Found: %s\n", remIPClient, xGetDN[xGetCkl])
				xGetCkl++
				if xGetCkl > userLimit {
					xMessage = fmt.Sprintf("Количество персон по вашему запросу превысило %d! Пожалуйста, задайте критерии более конкретно!", userLimit)
					break
				}
			}
		}
		xSearchPplMode = 1
		ldapSearchMode = 2
	}

	if strings.ToLower(dn) != strings.ToLower(rconf.LDAP_URL[ldap_count][3]) || xSearchPplMode == 1 {
		go_home_button = "+"
	}
	if ldapSearchMode != 2 {
		xSearch = "+"
	}

	if strings.ToLower(dn) != strings.ToLower(rconf.LDAP_URL[ldap_count][3]) {
		if ldapSearchMode == 1 && xSearchPplMode == 0 {
			dn_back_tmp = strings.Split(dn, ",")
			for ckl1 = 1; ckl1 < len(dn_back_tmp); ckl1++ {
				if ckl1 == 1 {
					dn_back = dn_back_tmp[ckl1]
				} else {
					dn_back += fmt.Sprintf(",%s", dn_back_tmp[ckl1])
				}
			}
		}
	}

	//	log.Printf("%s ... Initialize LDAP connector...", remIPClient)

	if initLDAPConnector() == "error" {
		return
	}

	l, err = ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Printf("LDAP::Initialize() error: %v\n", err)
		return
	}

	//		l.Debug = true
	defer l.Close()

	log.Printf("%s =!= Connected to server %d of %d: %s", remIPClient, ldap_count+1, len(rconf.LDAP_URL), rconf.LDAP_URL[ldap_count][0])

	err = l.Bind(rconf.LDAP_URL[ldap_count][1], rconf.LDAP_URL[ldap_count][2])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Printf("LDAP::Bind() error: %v\n", err)
		return
	}

	t, err := template.ParseFiles("templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "header", template.FuncMap{"Pagetitle": rconf.WLB_HTML_Title, "FRColor": xFRColor, "BGColor": xBGColor})

	t, err = template.ParseFiles("templates/search.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "search", template.FuncMap{"GoHome": go_home_button, "PrevDN": dn_back, "DN": dn, "xSearch": xSearch, "xMessage": xMessage, "LineColor": "#EEEEEE", "LUserName": LUserName, "LoginShow": "Yes", "RedirectDN": r.RequestURI})

	t, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	if xSearchPplMode == 0 {

		search := ldap.NewSearchRequest(dn, ldapSearchMode, ldap.NeverDerefAliases, 0, 0, false, ldap_Search, ldap_Attr, nil)

		//	log.Printf("Search: %v\n%v\n%v\n%v\n%v\n%v\n", search, dn, ldapSearchMode, ldap.NeverDerefAliases, ldap_Search, ldap_Attr)

		sr, err := l.Search(search)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			log.Printf("LDAP::Search() error: %v\n", err)
			return
		}

		//	fmt.Printf("\n\nSearch: %v", search)

		log.Printf("%s ++> search: %s // found: %d\n", remIPClient, search.Filter, len(sr.Entries))

		if len(sr.Entries) > 0 {
			dnList := make(map[string]tList, len(sr.Entries))
			for _, entry := range sr.Entries {
				fType := ""
				fField := make(map[string]string, len(rconf.WLB_LDAP_ATTR))
				for _, attr := range entry.Attributes {
					for ckl1 := 0; ckl1 < len(rconf.WLB_LDAP_ATTR); ckl1++ {
						if attr.Name == rconf.WLB_LDAP_ATTR[ckl1][0] {
							fField[rconf.WLB_LDAP_ATTR[ckl1][1]] = fmt.Sprintf("%s", strings.Join(attr.Values, ","))
							//						fmt.Printf("Name: %s==%s --> %s = %s\n", attr.Name, rconf.WLB_LDAP_ATTR[ckl1][0], rconf.WLB_LDAP_ATTR[ckl1][1], fField[rconf.WLB_LDAP_ATTR[ckl1][1]])
							if rconf.WLB_LDAP_ATTR[ckl1][1] == "ORGName" {
								fType = "Org"
							}
							if rconf.WLB_LDAP_ATTR[ckl1][1] == "USERName" {
								fType = "User"
							}
						}
					}
				}
				getMore(remIPClient, fField, fType, l, dnList, setAdminMode)
			}
			t.ExecuteTemplate(w, "index", dnList)
		}
	} else {
		dnList := make(map[string]tList, xGetCkl)
		for ckl1 = 0; ckl1 < xGetCkl; ckl1++ {
			//			fmt.Printf("GET: %s / %d\n", xGetDN[ckl1], ckl1)
			search := ldap.NewSearchRequest(xGetDN[ckl1], 0, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=inetOrgPerson)", ldap_Attr, nil)
			//			fmt.Printf("GET: %v\n", search)
			sr, err := l.Search(search)
			if err != nil {
				fmt.Printf(err.Error())
				//				fmt.Fprintf(w, err.Error())
				log.Printf("LDAP::Search() error: %v %v\n", search, err)
				continue
			}
			fType := "User"
			fField := make(map[string]string, len(rconf.WLB_LDAP_ATTR))
			fField["DN"] = xGetDN[ckl1]
			if len(sr.Entries) > 0 {
				for _, entry := range sr.Entries {
					for _, attr := range entry.Attributes {
						for ckl2 := 0; ckl2 < len(rconf.WLB_LDAP_ATTR); ckl2++ {
							if attr.Name == rconf.WLB_LDAP_ATTR[ckl2][0] {
								fField[rconf.WLB_LDAP_ATTR[ckl2][1]] = fmt.Sprintf("%s", strings.Join(attr.Values, ","))
								//								fmt.Printf("Name: %s==%s --> %s = %s\n", attr.Name, rconf.WLB_LDAP_ATTR[ckl1][0], rconf.WLB_LDAP_ATTR[ckl1][1], fField[rconf.WLB_LDAP_ATTR[ckl1][1]])
							}
						}

					}
				}
			}
			getMore(remIPClient, fField, fType, l, dnList, setAdminMode)
		}
		t.ExecuteTemplate(w, "index", dnList)
	}

	t, err = template.ParseFiles("templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "footer", template.FuncMap{"WebBookVersion": pVersion, "xMailBT": rconf.WLB_MailBT, "LineColor": "#EEEEEE"})
}

// ====================================================================================================

func main() {

	pVersion = fmt.Sprintf("%s V%s", pName, pVer)

	fmt.Printf("\n\t%s\n\n", pVersion)

	rconf.LOG_File = def_log_file

	def_config_file, def_daemon_mode = SABModules.ParseCommandLine(def_config_file, def_daemon_mode)

	SABModules.ReadConfigFile(def_config_file, &rconf)

	SABModules.Pid_Check(&rconf)
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

	/*	err := httpscerts.Check("cert.pem", "key.pem")
		if err != nil {
			err = httpscerts.Generate("cert.pem", "key.pem", rconf.WLB_Listen_IP)
			if err != nil {
				log.Println("Error: Couldn't create https certs.")
				os.Exit(1)
			}
		}*/
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/modify", modifyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/GoOrg", indexHandler)
	http.HandleFunc("/DavDN", davDNHandler)
	//	fmt.Printf("1 %v\n", rconf)
	//	fmt.Printf("2 %s / %s\n", rconf.WLB_Listen_IP, fmt.Sprintf("%s",rconf.WLB_Listen_PORT))
	//	fmt.Printf("3\n")
	//	fmt.Printf("4\n")

	//	go http.ListenAndServeTLS(rconf.WLB_Listen_IP+":443", "cert.pem", "key.pem", nil)
	//	http.ListenAndServe(rconf.WLB_Listen_IP+":80", http.HandlerFunc(redirectToHttps))
	//	fmt.Printf("5 %s / %s\n", rconf.WLB_Listen_IP, fmt.Sprintf("%s",rconf.WLB_Listen_PORT))

	SABModules.Log_ON(&rconf)

	log.Printf("->")
	log.Printf("--> %s", pVersion)
	log.Printf("---> I'm Ready...")
	log.Printf(" _")

	pgInit()

	SABModules.Log_OFF()

	go func() {
		for {
			SABModules.Log_ON(&rconf)
			log.Printf("Session cleaner ***** Remove old sessions...")

			dbpg, err := sql.Open("postgres", rconf.PG_DSN)
			if err != nil {
				log.Fatalf("PG_INIT::Open() error: %v\n", err)
			}

			queryx := fmt.Sprintf("delete from wb_auth_session where exptime<%v;", int64(time.Now().Unix()-int64(rconf.WLB_SessTimeOut*60+30)))
			_, err = dbpg.Query(queryx)
			if err != nil {
				log.Printf("%s\n", queryx)
				log.Printf("PG::Query() Get User Name for UID: %v\n", err)
				return
			}

			dbpg.Close()

			log.Printf("Session cleaner ***** Sleep for %d seconds...", rconf.Sleep_Time)
			SABModules.Log_OFF()
			time.Sleep(time.Duration(rconf.Sleep_Time) * time.Second)

		}
	}()

	http.ListenAndServe(rconf.WLB_Listen_IP+":"+fmt.Sprintf("%d", rconf.WLB_Listen_PORT), nil)
}
