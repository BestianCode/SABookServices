package SABFunctions

import (
	"fmt"
	"log"
	"os"
	"strings"

	// PostgreSQL
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABookServices/SABModules"

	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABDefine"
)

func AD_Analyze(conf *SABModules.Config_STR) int {

	var (
		/*			GlobalParentInsert = string("")

					GlobalParent   string
					GlobalParentID int
		*/
		ckl    int
		queryx string
		err    error
		x      string
		y      string
		z      string
		w      string
		v      string
		q      string
		dp     string

		/*
			lastx = int(0)
			lasty = int(0)


			ldap_table_check = int(1)

			ldap_que_check_tables = string("SELECT count(tablename) FROM pg_catalog.pg_tables where tablename like 'ldap%';")
		*/)

	log.Printf(".")
	log.Printf("..")
	log.Printf("...")
	log.Printf("Analyze AD-Logins database...")
	/*
		err = os.RemoveAll(conf.AD_ScriptDir)
		if err != nil {
			log.Printf("Directory %s remove error! %v\n", conf.AD_ScriptDir, err)
		}
	*/
	err = os.MkdirAll(conf.AD_ScriptDir, 0755)
	if err != nil {
		log.Printf("Directory %s create error! %v\n", conf.AD_ScriptDir, err)
	}

	/*
		for ckl = 0; ckl < len(conf.ROOT_DN); ckl++ {
			GlobalParentInsert = fmt.Sprintf("%sINSERT INTO ldap_entries VALUES (%d,'%s',%d,%d,%d,'%d','%d',%d); ", GlobalParentInsert, ckl+1, conf.ROOT_DN[ckl][0], 3, ckl, ckl+1, ckl+1, ckl, 0)
			GlobalParentInsert = fmt.Sprintf("%sINSERT INTO ldapx_institutes VALUES (%d,'%s','%d','%d',%d); ", GlobalParentInsert, ckl+1, conf.ROOT_DN[ckl][1], ckl+1, ckl, 0)
		}
		GlobalParentID = ckl
		GlobalParent = conf.ROOT_DN[ckl-1][0]
		//	log.Printf("%d %s\n\n%s\n", GlobalParentID, GlobalParent, GlobalParentInsert)
	*/
	db, err := sql.Open("postgres", conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return 10
	}

	defer db.Close()

	subParentCheck := ""
	for ckl = 0; ckl < len(conf.AD_LDAP_PARENT); ckl++ {
		if ckl > 0 {
			subParentCheck = fmt.Sprintf("%s or", subParentCheck)
		}
		subParentCheck = fmt.Sprintf("%s (xad.domain='%s' and cache.idorg='%s')", subParentCheck, conf.AD_LDAP_PARENT[ckl][0], conf.AD_LDAP_PARENT[ckl][1])
	}
	subParentCheck = fmt.Sprintf("(%s)", subParentCheck)

	log.Printf("\t\tCheck ad-logins for duplicates...")

	ofr, err := os.OpenFile(conf.AD_ScriptDir+"/AD-01-Duplicates.report", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error open report file: %v (%s)", err, conf.AD_ScriptDir+"/AD-01-Duplicates.report")
	}

	queryx = strings.Replace(SABDefine.PG_QUE_AD_GetDupInAD, "XYZDBADXYZ", SABDefine.PG_Table_AD, -1)
	//log.Printf("%s\n", queryx)
	rows, err := db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() AD Logins error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 11
	}

	for rows.Next() {
		rows.Scan(&x)
		log.Printf("Multiple AD USERS: %s\n", x)
		ofr.WriteString(fmt.Sprintf("%s\n", x))
	}

	ofr.Close()

	log.Printf("\t\tCheck ad-logins need for update...")

	ofr, err = os.OpenFile(conf.AD_ScriptDir+"/AD-02-Update.cmd", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error open report file: %v (%s)", err, conf.AD_ScriptDir+"/AD-02-Update.cmd")
	}

	queryx = strings.Replace(SABDefine.PG_QUE_AD_GetUpdateAD, "XYZDBADXYZ", SABDefine.PG_Table_AD, -1)
	queryx = strings.Replace(queryx, "XYZDBPersXYZ", SABDefine.PG_Table_MSSQL[2], -1)
	queryx = strings.Replace(queryx, "XYZSubParentCheckXYZ", subParentCheck, -1)
	//log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() AD Logins error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 12
	}

	for rows.Next() {
		rows.Scan(&x, &y)
		z = fmt.Sprintf("dsmod user \"%s\" -display \"%s\"\n", y, x)
		log.Printf("%s", z)
		//dsmod user "CN=Кашулин Андрей,OU=dit,OU=trg,DC=tgk-4,DC=ru" -display "Полное имя" -office "расположение офиса" -tel "телефон" -email "почта" -
		//mobile "сотовый" -iptel "еще телефон" -title "должность" -dept "подразделение" -company "организация"

		ofr.WriteString(z)
	}

	ofr.Close()

	log.Printf("\t\tCheck not connected ad-logins...")

	ofr, err = os.OpenFile(conf.AD_ScriptDir+"/AD-03-NotConnected.report", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error open report file: %v (%s)", err, conf.AD_ScriptDir+"/AD-03-NotConnected.report")
	}

	queryx = strings.Replace(SABDefine.PG_QUE_AD_GetNotConnected, "XYZDBADXYZ", SABDefine.PG_Table_AD, -1)
	queryx = strings.Replace(queryx, "XYZDBPersXYZ", SABDefine.PG_Table_MSSQL[2], -1)
	queryx = strings.Replace(queryx, "XYZSubParentCheckXYZ", subParentCheck, -1)
	//log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() AD Logins error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 12
	}

	for rows.Next() {
		rows.Scan(&x)
		log.Printf("Not connected AD USERS: %s\n", x)
		ofr.WriteString(fmt.Sprintf("dsrm -noprompt \"%s\"\n", x))
	}

	ofr.Close()

	log.Printf("\t\tCheck ad-logins need to update credentials...")

	ofr, err = os.OpenFile(conf.AD_ScriptDir+"/AD-04-UpdateCreds.cmd", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Error open report file: %v (%s)", err, conf.AD_ScriptDir+"/AD-04-UpdateCreds.cmd")
	}

	queryx = strings.Replace(SABDefine.PG_QUE_AD_SetCredentInfoToAD, "XYZDBADXYZ", SABDefine.PG_Table_AD, -1)
	//log.Printf("%s\n", queryx)
	rows, err = db.Query(queryx)
	if err != nil {
		log.Printf("PG::Query() AD Logins error: %v\n", err)
		log.Printf("%s\n", queryx)
		return 13
	}

	for rows.Next() {
		rows.Scan(&x, &y, &z, &w, &v, &q, &dp)
		log.Printf("AD USERS Upd Creds: %s\n", x)
		ofr.WriteString(fmt.Sprintf("dsmod user \"%s\" -email \"%s\" -title \"%s\" -mobile \"%s\" -tel \"%s\" -pager \"%s\" -dept \"%s\"\n", x, y, z, w, v, q, dp))
	}

	ofr.Close()

	log.Printf("\tComplete")

	return 94

}
