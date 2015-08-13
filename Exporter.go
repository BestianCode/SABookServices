package main

import (
	"log"
	"os"
	"syscall"
//	"time"
//	"fmt"
	"os/exec"

	"github.com/BestianRU/SABookServices/SABDefine"
	"github.com/BestianRU/SABookServices/SABFunctions"

)

func main() {

	var	(
		my_error		int

		st_LDAP_to_PG		=	int(0)
		st_Oracle_to_PG_m1	=	int(0)
		st_Oracle_to_PG_m2	=	int(0)
		st_Oracle_to_PG_m3	=	int(0)

		st_AsteriskCID_UP	=	int(0)

		st_LDAP_UP		=	int(0)

//		st_WorkTables		=	int(0)

		sleep_counter		=	int(0)
	)

	SABDefine.Conf.LOG_File = "/var/log/Exporter.log"

	SABDefine.Def_config_file=SABFunctions.ParseCommandLine(SABDefine.Def_config_file)

	SABFunctions.ReadConfigFile(SABDefine.Def_config_file)

	flog_file, err := os.OpenFile(SABDefine.Conf.LOG_File, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error open log file: %v", err)
	}
	flog_file.Close()

	ret, _, forkerr := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if forkerr != 0 {
		log.Fatalf("Fork error %v\n", forkerr)
	}

	if ret > 0 {
//		log.Printf("Daemon pid: %d", ret)
		os.Exit(0)
	}

	for {
//		log.Printf("Daemon pid: %d", ret)

//		SABFunctions.CheckState(&SABDefine.Conf, "1 Start")

		flog_file, err := os.OpenFile(SABDefine.Conf.LOG_File, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error open log file: %v", err)
		}

		log.SetOutput(flog_file)

		log.Printf("WakeUP!")

		if st_LDAP_to_PG == 0 {
			if my_error = SABFunctions.LDAP_to_PG(&SABDefine.Conf); my_error != 94 {
				log.Printf("LDAP_to_PG() my_error: %v\n", my_error)
			}else{
				st_LDAP_to_PG=1
			}
		}

		if st_Oracle_to_PG_m1 == 0 {
			if my_error = SABFunctions.Oracle_to_PG(1,&SABDefine.Conf); my_error != 94 {
				log.Printf("Oracle_to_PG() mode 1 my_error:%v\n", my_error)
			}else{
				st_Oracle_to_PG_m1=1
			}
		}

		if st_Oracle_to_PG_m2 == 0 {
			if my_error = SABFunctions.Oracle_to_PG(2,&SABDefine.Conf); my_error != 94 {
				log.Printf("Oracle_to_PG() mode 2 my_error:%v\n", my_error)
			}else{
				st_Oracle_to_PG_m2=1
			}
		}

		if st_Oracle_to_PG_m3 == 0 {
			if my_error = SABFunctions.Oracle_to_PG(3,&SABDefine.Conf); my_error != 94 {
				log.Printf("Oracle_to_PG() mode 3 my_error:%v\n", my_error)
			}else{
				st_Oracle_to_PG_m3=1
			}
		}

//		if st_LDAP_to_PG>0 && st_Oracle_to_PG_m1>0 && st_Oracle_to_PG_m2>0 && st_Oracle_to_PG_m3>0 && st_WorkTables==0 {
/*		if st_LDAP_to_PG>0 && st_Oracle_to_PG_m1>0 && st_Oracle_to_PG_m2>0 && st_Oracle_to_PG_m3>0 {
			if st_AsteriskCID_UP == 0 {
				if my_error = SABFunctions.MakeAsteriskCIDTable(&SABDefine.Conf); my_error != 94 {
					log.Printf("MakeAsteriskCIDTable error: %v\n", my_error)
				}else{
					st_AsteriskCID_UP=1
				}
			}
//			st_WorkTables=1
		}
*/
		if st_Oracle_to_PG_m1>0 && st_Oracle_to_PG_m2>0 && st_Oracle_to_PG_m3>0 {
			if st_AsteriskCID_UP == 0 {
				if my_error = SABFunctions.MakeAsteriskCIDTable(&SABDefine.Conf); my_error != 94 {
					log.Printf("MakeAsteriskCIDTable error: %v\n", my_error)
				}else{
					st_AsteriskCID_UP=1
				}
			}
			if st_LDAP_UP == 0 {
				if my_error = SABFunctions.GetORGs(&SABDefine.Conf); my_error != 94 {
					log.Printf("LDAP Update error: %v\n", my_error)
				}else{
					st_LDAP_UP=1
				}
			}
		}

		sleep_counter++

		log.Printf("----- %d ----- Sleep for %s sec...", sleep_counter, SABDefine.Sleep_Time)

		flog_file.Close()

		if sleep_counter>SABDefine.Sleep_cycles-1 {

			st_LDAP_to_PG=0
			st_Oracle_to_PG_m1=0
			st_Oracle_to_PG_m2=0
			st_Oracle_to_PG_m3=0

			st_AsteriskCID_UP=0
			st_LDAP_UP=0

//			st_WorkTables=0

			sleep_counter=0
		}

//		SABFunctions.CheckState(&SABDefine.Conf, "2 Before SLEEP")

//		time.Sleep(SABDefine.Sleep_Time*time.Second)

		if err = exec.Command("sleep", SABDefine.Sleep_Time).Run(); err != nil {
//			SABFunctions.CheckState(&SABDefine.Conf, "Sleep DIED")
			flog_file, err := os.OpenFile(SABDefine.Conf.LOG_File, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Error open log file: %v", err)
			}

			log.SetOutput(flog_file)

			log.Printf("!!!! Sleep error: %v", err)

			flog_file.Close()

			os.Exit(1)
		}

//		SABFunctions.CheckState(&SABDefine.Conf, "3 Cycle END")


	}
}

