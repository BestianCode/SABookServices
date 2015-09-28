package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"log"
	"strings"
	"html/template"
	"net/http"
	"strconv"

//LDAP
	"github.com/go-ldap/ldap"

// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"


	"github.com/BestianRU/SABookServices/SABModules"
//	"github.com/kabukky/httpscerts"
//	"github.com/gavruk/go-blog-example/models"
)

type tList struct {
	URL 				string
	URLName				string
	ORGName 			string
	USERName 			string
	FullName 			string
	PhoneInt			string
	PhoneExt			string
	Mobile				string
	Mail 				string
	Position			string
	ADLogin				string
}


const (
	pName				=	string("Web Address Book")
	pVer				=	string("1 alpha 2015.09.29.00.00")
	userLimit			=	20
)

var	(
	def_config_file		=	string ("./WebLDAPBook.json")				// Default configuration file
	def_log_file		=	string ("/var/log/ABook/WebLDAPBook.log")	// Default log file
	def_daemon_mode		=	string ("NO")								// Default start in foreground
	pVersion				string
	rconf					SABModules.Config_STR
	ldap_count		=	int(0)
)

func getMore(remIPClient string, fField map[string]string, fType string, l *ldap.Conn, dnList map[string]tList){
	var		(
		fPath				string
		fURL				string
		fURLName			string
		ckl1, ckl2, ckl3	int
		ldap_Attr			[]string
	)

	if fField["DN"]!="" && (fField["USERName"]!="" || fField["ORGName"]!=""){
		fPath=fField["DN"]
		fPath=strings.Replace(strings.ToLower(fPath), ","+strings.ToLower(rconf.LDAP_URL[ldap_count][3]), "", -1)
		fPath_Split:=strings.Split(fPath, ",")
		fURLName=""
		for ckl1=0;ckl1<len(fPath_Split)-1;ckl1++ {
			fPath_Strip:=""
			for ckl2=ckl1+1;ckl2<len(fPath_Split);ckl2++ {
				fPath_Strip=fmt.Sprintf("%s%s,", fPath_Strip, fPath_Split[ckl2])
			}
			if fType=="User" {
				fPath_Strip=fmt.Sprintf("%s%s", fPath_Strip, rconf.LDAP_URL[ldap_count][3])
				if ckl1==0 {
					fURL=fPath_Strip
				}
//						log.Printf("X1: %s", fPath_Strip)
				subsearch := ldap.NewSearchRequest(fPath_Strip, 0, ldap.NeverDerefAliases, 0, 0, false, rconf.LDAP_URL[ldap_count][4], ldap_Attr, nil)
				subsr, err := l.Search(subsearch)
				if err != nil {
//								fmt.Fprintf(w, err.Error())
					log.Printf("LDAP::Search() error: %v\n", err)
				}

//						log.Printf("Y1: %s / %s / %d\n", fPath_Strip, rconf.LDAP_URL[ldap_count][4], len(subsr.Entries))
				if len(subsr.Entries)>0 {
					for _, subentry := range subsr.Entries {
						for _, subattr := range subentry.Attributes {
							for ckl3=0;ckl3<len(rconf.WLB_LDAP_ATTR);ckl3++ {
								if subattr.Name == rconf.WLB_LDAP_ATTR[ckl3][0] {
									if rconf.WLB_LDAP_ATTR[ckl3][1] == "ORGName" {
										if ckl1==0 {
											fURLName=fmt.Sprintf("%s", strings.Join(subattr.Values, ","))
										}else{
											fURLName=fmt.Sprintf("%s / %s", strings.Join(subattr.Values, ","), fURLName)
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

		fField["DN"]=strings.Replace(strings.ToLower(fField["DN"]), "/", ",", -1)
		fmt.Sprintf("/Go%s?dn=%s", fType, fField["DN"])
		fField["DN"]=fmt.Sprintf("/Go%s?dn=%s", fType, fField["DN"])
		fURL=fmt.Sprintf("/Go%s?dn=%s", fType, fURL)
		log.Printf("%s <-- %s", remIPClient, fField["DN"])
		dnList[fField["DN"]]=tList{URL: fURL, URLName: fURLName, ORGName: fField["ORGName"], USERName: fField["USERName"], FullName: fField["FullName"], Position: fField["Position"], PhoneInt: fField["PhoneInt"], Mobile: fField["Mobile"], PhoneExt: fField["PhoneExt"], Mail: fField["Mail"], ADLogin: fField["ADLogin"]}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var (

		xSearchPplMode	=	int(0)
		xSearch 			string
		xMessage			string

		dn 					string
		dn_back				string
		dn_back_tmp			[]string

		go_home_button		string

		ldap_Search			string

		ldapSearchMode	=	int(1)

		ckl1				int

		ldap_Attr			[]string

		xGetDN				[1000]string
		xGetCkl				int
	)

	ldap_Attr = make ([]string, len(rconf.WLB_LDAP_ATTR))

	for ckl1:=0;ckl1<len(rconf.WLB_LDAP_ATTR);ckl1++ {
			ldap_Attr[ckl1]=rconf.WLB_LDAP_ATTR[ckl1][0]
	}

	SABModules.Log_ON(&rconf)

	get_dn		:= r.FormValue("dn")
	get_cn		:= r.FormValue("cn")
	get_fn		:= r.FormValue("FirstName")
	get_ln		:= r.FormValue("LastName")
	remIPClient	:=strings.Split(r.RemoteAddr,":")[0]
//	log.Printf("DN: %s --- CN: %s", get_dn, get_cn)

	if get_dn == ""{
		dn=rconf.LDAP_URL[ldap_count][3]
	}else{
		dn=get_dn
	}

	if len(dn)<len(rconf.LDAP_URL[ldap_count][3]) {
		dn=rconf.LDAP_URL[ldap_count][3]
	}

	log.Printf("->")
	log.Printf("--> %s", pVersion)
	log.Printf("->")
	ucurl, _ :=strconv.Unquote(r.RequestURI)
	log.Println(remIPClient+" --> http://"+r.Host+ucurl)
	log.Printf("%s ++> DN: %s / CN: %s / Mode: %d / Def.DN: %s", remIPClient, dn, ldap_Search, ldapSearchMode, rconf.LDAP_URL[ldap_count][3])

	if get_cn == "" && get_ln == "" && get_fn == "" {
		ldap_Search=rconf.LDAP_URL[ldap_count][4]
	}else{
		if strings.ToLower(rconf.WLB_SQL_PreFetch)=="yes" {
			log.Printf("%s ++> SQL Search: %s\n", remIPClient, get_cn)
			pgdb, err := sql.Open("postgres", rconf.PG_DSN)
			if err != nil {
				log.Printf("PG::Open() error: %v\n", err)
				return
			}
			defer pgdb.Close()

//			queryx := fmt.Sprintf("select x.dn from ldap_entries as x, ldapx_persons as y where y.uid=x.uid and lower(y.fullname) like lower('%%%s%%') and lower(x.dn) like lower('%%%s');", get_cn, dn)
//			queryx := fmt.Sprintf("select x.dn from ldap_entries as x, ldapx_persons as y where y.uid=x.uid and lower(y.fullname) like lower('%%%s%%');", get_cn)
			queryx := "select x.dn from ldap_entries as x, ldapx_persons as y where y.uid=x.uid"
			if len(get_cn)>1 {
				queryx = fmt.Sprintf("%s and lower(y.fullname) like lower('%%%s%%')", queryx, get_cn)
			}
			if len(get_ln)>1 {
				queryx = fmt.Sprintf("%s and lower(y.surname) like lower('%s%%')", queryx, get_ln)
			}
			if len(get_fn)>1 {
				queryx = fmt.Sprintf("%s and lower(y.name) like lower('%s%%')", queryx, get_fn)
			}
			queryx = fmt.Sprintf("%s;", queryx)
			log.Printf("SQL: %s\n", queryx)
			rows, err := pgdb.Query(queryx)
			if err != nil {
				fmt.Printf("SQL Error: %s\n", queryx)
				log.Printf("PG::Query() Check LDAP tables error: %v\n", err)
				return
			}
			xGetCkl=0
			for rows.Next() {
				rows.Scan(&xGetDN[xGetCkl])
//				fmt.Println(xGetDN[xGetCkl], dn)
				if strings.Contains(strings.ToLower(xGetDN[xGetCkl]), strings.ToLower(dn)) {
					log.Printf("%s <-- SQL Found: %s\n", remIPClient, xGetDN[xGetCkl])
					xGetCkl++
					if xGetCkl > userLimit {
						xMessage = fmt.Sprintf("Количество персон по вашему запросу привысило %d человек. Пожалуйста, задайте критерии более конкретно!", userLimit)
						break
					}
				}
			}
			xSearchPplMode=1
		}else{
	//		ldap_Search=fmt.Sprintf("(&(objectClass=*)(cn=*%s*))",unidecode.Unidecode(get_cn))
	//		ldap_Search=fmt.Sprintf("(&(objectClass=*)((displayName=*%s*)))",get_cn)
			ldap_Search=fmt.Sprintf("(&(objectClass=inetOrgPerson)(displayName=*%s*))", get_cn)
	//		ldap_Search=fmt.Sprintf("(|(displayName=*%s*))", get_cn)
	//		ldap_Search=fmt.Sprintf("(%s)", search_str)
	//		ldap_Search=fmt.Sprintf("(displayName=*%s*)", get_cn)
	//		ldap_Search=fmt.Sprintf("(cn=%s)",get_cn)
	//		ldap_Search=fmt.Sprintf("(&(objectClass=)(cn=*%s*))",unidecode.Unidecode(get_cn))
		}
		ldapSearchMode=2
	}

	if strings.ToLower(dn) != strings.ToLower(rconf.LDAP_URL[ldap_count][3]) || xSearchPplMode==1 {
		go_home_button="+"
	}
	if ldapSearchMode !=2 {
		xSearch = "+"
	}


	if strings.ToLower(dn) != strings.ToLower(rconf.LDAP_URL[ldap_count][3]) {
		if ldapSearchMode == 1 && xSearchPplMode==0 {
			dn_back_tmp = strings.Split(dn, ",")
			for ckl1=1;ckl1<len(dn_back_tmp);ckl1++ {
				if ckl1 == 1 {
					dn_back = dn_back_tmp[ckl1]
				}else{
					dn_back += fmt.Sprintf(",%s", dn_back_tmp[ckl1])
				}
			}
		}
	}

	log.Printf("%s ... Initialize connector...", remIPClient)

	l, err := ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
	if err == nil {
		l.Close()
	}

	ckl1=0

	for {
		if ckl1>9 {
			fmt.Fprintf(w, "Error connect to all LDAP servers...")
			log.Printf("Error connect to all LDAP servers...")
			return
		}

		ldap_count++
		if ldap_count>len(rconf.LDAP_URL)-1 {
			ldap_count=0
		}

		log.Printf("%s ... Trying to connect server %d of %d: %s", remIPClient, ldap_count+1, len(rconf.LDAP_URL), rconf.LDAP_URL[ldap_count][0])
		l, err = ldap.Dial("tcp", rconf.LDAP_URL[ldap_count][0])
		if err != nil {
//			fmt.Fprintf(w, err.Error())
//			log.Printf("LDAP::Initialize() error: %v\n", err)
			continue
		}

		defer l.Close()
//		l.Debug = true

		break

		ckl1++
	}

	log.Printf("%s =!= Connected to server %d of %d: %s", remIPClient, ldap_count+1, len(rconf.LDAP_URL), rconf.LDAP_URL[ldap_count][0])

	err = l.Bind(rconf.LDAP_URL[ldap_count][1],rconf.LDAP_URL[ldap_count][2])
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

	t.ExecuteTemplate(w, "header", template.FuncMap{"Pagetitle":"PhoneBook"})

	t, err = template.ParseFiles("templates/search.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "search", template.FuncMap{"GoHome":go_home_button, "PrevDN":dn_back, "DN":dn, "xSearch":xSearch, "xMessage":xMessage})

	t, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	if xSearchPplMode==0 {

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

		if len(sr.Entries)>0 {
			dnList := make (map[string]tList, len(sr.Entries))
			for _, entry := range sr.Entries {
				fType		:= ""
				fField 		:= make	(map[string]string, len(rconf.WLB_LDAP_ATTR))
				for _, attr := range entry.Attributes {
					for ckl1:=0;ckl1<len(rconf.WLB_LDAP_ATTR);ckl1++ {
						if attr.Name == rconf.WLB_LDAP_ATTR[ckl1][0] {
							fField[rconf.WLB_LDAP_ATTR[ckl1][1]]=fmt.Sprintf("%s", strings.Join(attr.Values, ","))
	//						fmt.Printf("Name: %s==%s --> %s = %s\n", attr.Name, rconf.WLB_LDAP_ATTR[ckl1][0], rconf.WLB_LDAP_ATTR[ckl1][1], fField[rconf.WLB_LDAP_ATTR[ckl1][1]])
							if rconf.WLB_LDAP_ATTR[ckl1][1] == "ORGName" {
								fType="Org"
							}
							if rconf.WLB_LDAP_ATTR[ckl1][1] == "USERName" {
								fType="User"
							}
						}
					}
				}
				getMore(remIPClient, fField, fType, l, dnList)
			}
			t.ExecuteTemplate(w, "index", dnList)
		}
	}else{
		dnList := make (map[string]tList, xGetCkl)
		for ckl1=0;ckl1<xGetCkl;ckl1++{
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
			fType		:= "User"
			fField 		:= make	(map[string]string, len(rconf.WLB_LDAP_ATTR))
			fField["DN"] = xGetDN[ckl1]
			if len(sr.Entries)>0 {
				for _, entry := range sr.Entries {
					for _, attr := range entry.Attributes {
						for ckl2:=0;ckl2<len(rconf.WLB_LDAP_ATTR);ckl2++ {
							if attr.Name == rconf.WLB_LDAP_ATTR[ckl2][0] {
								fField[rconf.WLB_LDAP_ATTR[ckl2][1]]=fmt.Sprintf("%s", strings.Join(attr.Values, ","))
//								fmt.Printf("Name: %s==%s --> %s = %s\n", attr.Name, rconf.WLB_LDAP_ATTR[ckl1][0], rconf.WLB_LDAP_ATTR[ckl1][1], fField[rconf.WLB_LDAP_ATTR[ckl1][1]])
							}
						}

					}
				}
			}
			getMore(remIPClient, fField, fType, l, dnList)
		}
		t.ExecuteTemplate(w, "index", dnList)
	}

	t, err = template.ParseFiles("templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "footer", template.FuncMap{"WebBookVersion":pVersion, "xMailBT":rconf.WLB_MailBT})

	SABModules.Log_OFF()
}

func main() {

	pVersion=fmt.Sprintf("%s V%s", pName, pVer)

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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/GoOrg", indexHandler)
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

	SABModules.Log_OFF()

	http.ListenAndServe(rconf.WLB_Listen_IP+":"+fmt.Sprintf("%d",rconf.WLB_Listen_PORT), nil)
}

