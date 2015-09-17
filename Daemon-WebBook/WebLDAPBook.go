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

//LDAP
	"github.com/BestianRU/SABookServices/ForeignModules/ldap.v1"

	"github.com/BestianRU/SABookServices/SABModules"
//	"github.com/kabukky/httpscerts"
//	"github.com/gavruk/go-blog-example/models"
)

const (
	pName				=	string("SABook Web Address Book")
	pVer				=	string("1 alpha 2015.09.17.23.45")
)

var	(
	def_config_file		=	string ("./WebLDAPBook.json")				// Default configuration file
	def_log_file		=	string ("/var/log/ABook/WebLDAPBook.log")	// Default log file
	def_daemon_mode		=	string ("NO")								// Default start in foreground

	pVersion				string

	rconf					SABModules.Config_STR

	ldap_count		=	int(0)
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	type tList struct {
		URL 				string
		URLName				string
		Dn 					string
		Name 				string
		BusinessCategory	string
		TelephoneNumber		string
		Mobile				string
		Pager				string
		Mail 				string
	}

	var (
		ftype				string
		fdn					string
		foname				string
		fname				string
		fbusinessCategory	string
		ftelephoneNumber	string
		fmobile				string
		fpager				string
		fmail				string
		fPath				string
		fURL				string
		fURL_Name			string

		dn 					string
		ldap_Search		string

		ldapSearchMode	=	int(1)

		ckl1, ckl2			int

		ldap_Attr			[]string

	)

	ldap_Attr = make ([]string, len(rconf.WLB_LDAP_ATTR))

	for ckl1:=0;ckl1<len(rconf.WLB_LDAP_ATTR);ckl1++ {
			ldap_Attr[ckl1]=rconf.WLB_LDAP_ATTR[ckl1][0]
	}

	SABModules.Log_ON(&rconf)

	get_dn := r.FormValue("dn")
	get_cn := r.FormValue("cn")
//	log.Printf("DN: %s --- CN: %s", get_dn, get_cn)

	if get_dn == "" {
		dn=rconf.LDAP_URL[ldap_count][3]
	}else{
		dn=get_dn
	}

	if get_cn == "" {
		ldap_Search=rconf.LDAP_URL[ldap_count][4]
	}else{
//		ldap_Search=fmt.Sprintf("(&(objectClass=*)(cn=*%s*))",unidecode.Unidecode(get_cn))
//		ldap_Search=fmt.Sprintf("(&(objectClass=*)((displayName=*%s*)))",get_cn)
		ldap_Search=fmt.Sprintf("(|(displayName=*%s*)(cn=*%s*))", get_cn, get_cn)
//		ldap_Search=fmt.Sprintf("(cn=*%s*)",get_cn)
		ldapSearchMode=2
	}

	remIPClient:=strings.Split(r.RemoteAddr,":")[0]

	log.Printf("->")
	log.Printf("--> %s", pVersion)
	log.Printf("->")
	log.Println(remIPClient+" --> https://"+r.Host+r.RequestURI)
	log.Printf("%s ++> DN: %s / CN: %s / Mode: %d", remIPClient, dn, ldap_Search, ldapSearchMode)

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
			fmt.Fprintf(w, err.Error())
			log.Printf("LDAP::Initialize() error: %v\n", err)
			continue
		}

		defer l.Close()
	//	l.Debug = true

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


	search := ldap.NewSearchRequest(dn, ldapSearchMode, ldap.NeverDerefAliases, 0, 0, false, ldap_Search, ldap_Attr, nil)

//	log.Printf("Search: %v\n%v\n%v\n%v\n%v\n%v\n", search, dn, ldapSearchMode, ldap.NeverDerefAliases, ldap_Search, ldap_Attr)

	sr, err := l.Search(search)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Printf("LDAP::Search() error: %v\n", err)
		return
	}

//	log.Printf("\n\nXXX2: %v", search)

	log.Printf("%s ++> search: %s // found: %d\n", remIPClient, search.Filter, len(sr.Entries))

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

	t.ExecuteTemplate(w, "search", template.FuncMap{"DN":dn})

	t, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	if len(sr.Entries)>0 {
		dnList := make (map[string]tList, len(sr.Entries))
		for _, entry := range sr.Entries {
			fdn=""
			foname=""
			fname=""
			ftype=""
			fbusinessCategory=""
			ftelephoneNumber=""
			fmobile=""
			fpager=""
			fmail=""
			for _, attr := range entry.Attributes {
				if attr.Name == "o" {
					x  := strings.Join(attr.Values, ",")
					foname=fmt.Sprintf("%s", x)
					ftype="Org"
				}
				if attr.Name == "displayName" {
					x  := strings.Join(attr.Values, ",")
					fname=fmt.Sprintf("%s", x)
					ftype="User"
				}
				if attr.Name == "entryDN" {
					x  := strings.Join(attr.Values, ",")
					fdn=fmt.Sprintf("%s", x)
				}
				if attr.Name == "businessCategory" {
					x  := strings.Join(attr.Values, ",")
					fbusinessCategory=fmt.Sprintf("%s", x)
				}
				if attr.Name == "telephoneNumber" {
					x  := strings.Join(attr.Values, ",")
					ftelephoneNumber=fmt.Sprintf("%s", x)
				}
				if attr.Name == "mobile" {
					x  := strings.Join(attr.Values, ",")
					fmobile=fmt.Sprintf("%s", x)
				}
				if attr.Name == "pager" {
					x  := strings.Join(attr.Values, ",")
					fpager=fmt.Sprintf("%s", x)
				}
				if attr.Name == "mail" {
					x  := strings.Join(attr.Values, ",")
					fmail=fmt.Sprintf("%s", x)
				}
			}
			if fdn!="" && (fname!="" || foname!=""){
				fPath=fdn
				fPath=strings.Replace(strings.ToLower(fPath), ","+strings.ToLower(rconf.LDAP_URL[ldap_count][3]), "", -1)
				fPath_Split:=strings.Split(fPath, ",")
				if ftype=="User" {
//					log.Printf("%s", fPath)
				}
				fURL_Name=""
				for ckl1=0;ckl1<len(fPath_Split)-1;ckl1++ {
					fPath_Strip:=""
					for ckl2=ckl1+1;ckl2<len(fPath_Split);ckl2++ {
						fPath_Strip=fmt.Sprintf("%s%s,", fPath_Strip, fPath_Split[ckl2])
					}
					if ftype=="User" {
						fPath_Strip=fmt.Sprintf("%s%s", fPath_Strip, rconf.LDAP_URL[ldap_count][3])
						if ckl1==0 {
							fURL=fPath_Strip
						}
//						log.Printf("%s", fPath_Strip)


						subsearch := ldap.NewSearchRequest(fPath_Strip, 0, ldap.NeverDerefAliases, 0, 0, false, rconf.LDAP_URL[ldap_count][4], ldap_Attr, nil)
						subsr, err := l.Search(subsearch)
						if err != nil {
							fmt.Fprintf(w, err.Error())
							log.Printf("LDAP::Search() error: %v\n", err)
						}

//						log.Printf("\t\t\t%s / %s / %d\n", fPath_Strip, rconf.LDAP_URL[ldap_count][4], len(subsr.Entries))

						if len(subsr.Entries)>0 {
							for _, subentry := range subsr.Entries {
								for _, subattr := range subentry.Attributes {
									if subattr.Name == "o" {
										if ckl1==0 {
											fURL_Name=fmt.Sprintf("%s", strings.Join(subattr.Values, ","))
										}else{
											fURL_Name=fmt.Sprintf("%s / %s", strings.Join(subattr.Values, ","), fURL_Name)
										}
//										log.Printf("%s", fURL_Name)
									}
								}
							}

						}



					}
				}

				fdn=fmt.Sprintf("/Go%s?dn=%s", ftype, fdn)
				fURL=fmt.Sprintf("/Go%s?dn=%s", ftype, fURL)
				log.Printf("%s <-- %s", remIPClient, fdn)
				dnList[fdn]=tList{URL: fURL, URLName: fURL_Name, Dn: foname, Name: fname, BusinessCategory: fbusinessCategory, TelephoneNumber: ftelephoneNumber, Mobile: fmobile, Pager: fpager, Mail: fmail}
			}
		}

		t.ExecuteTemplate(w, "index", dnList)
	}

	t, err = template.ParseFiles("templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		log.Println(err.Error())
		return
	}

	t.ExecuteTemplate(w, "footer", template.FuncMap{"WebBookVersion":pVersion})

	SABModules.Log_OFF()
}

/*
func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "redirect ro: https://%s/%s", rconf.WLB_Listen_IP, r.RequestURI)
	log.Printf("redirect ro: https://%s/%s", rconf.WLB_Listen_IP, r.RequestURI)
	time.Sleep(time.Duration(2)*time.Second)
	http.Redirect(w, r, "https://"+rconf.WLB_Listen_IP+r.RequestURI, http.StatusMovedPermanently)
}
*/

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

