package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"os/exec"

	"github.com/BestianRU/SABookServices/SABModules"
	"github.com/BestianRU/SABookServices/Daemon-Exporter/SABFunctions"
)

func main() {

	const (
		pName				=	string("SABook Exporter Daemon")
		pVer				=	string("3 2015.09.21.23.45")

		pg_MultiInsert		=	int(50)
	)

	var	(

		def_config_file		=	string ("./Exporter.json")			// Default configuration file
		def_log_file		=	string ("/var/log/ABook/Exporter.log")	// Default log file
		def_daemon_mode		=	string ("NO")						// Default start in foreground

		my_error				int

		st_MSSQL_to_PG		=	int(0)
		st_LDAP_to_PG		=	int(0)
		st_Oracle_to_PG		=	int(0)

		st_LDAP_UP			=	int(0)

		st_DB_clean			=	int(0)

		sleep_counter		=	int(0)

		rconf					SABModules.Config_STR
	)

	fmt.Printf("\n\t%s V%s\n\n", pName, pVer)

	rconf.LOG_File = def_log_file

	def_config_file, def_daemon_mode = SABModules.ParseCommandLine(def_config_file, def_daemon_mode)

//	log.Printf("%s %s %s", def_config_file, def_daemon_mode, os.Args[0])

	SABModules.ReadConfigFile(def_config_file, &rconf)

	SABModules.Pid_Check(&rconf)

	if def_daemon_mode=="YES" {
		if err :=  exec.Command(os.Args[0], fmt.Sprintf("-daemon=GO -config=%s &", def_config_file)).Start(); err != nil {
			log.Fatalf("Fork daemon error: %v", err)
		}else{
			log.Printf("Forked!")
			os.Exit(0)
		}
	}


	SABModules.Log_ON(&rconf)
	SABModules.Log_OFF()

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

	for {

		SABModules.Log_ON(&rconf)

		if sleep_counter==0 {
			log.Printf("-> %s V%s", pName, pVer)
			log.Printf("--> Go!")
		}else{
			log.Printf("--> WakeUP!")
		}

		if st_MSSQL_to_PG == 0 {
			if my_error = SABFunctions.MSSQL_to_PG(&rconf, pg_MultiInsert); my_error != 94 {
				log.Printf("MSSQL_to_PG() error:%v\n", my_error)
			}else{
				st_MSSQL_to_PG=1
			}
			log.Printf("-")
		}

		if st_LDAP_to_PG == 0 {
			if my_error = SABFunctions.LDAP_to_PG(&rconf, pg_MultiInsert); my_error != 94 {
				log.Printf("LDAP_to_PG() my_error: %v\n", my_error)
			}else{
				st_LDAP_to_PG=1
			}
			log.Printf("-")
		}

		if st_Oracle_to_PG == 0 {
			if my_error = SABFunctions.Oracle_to_PG(&rconf, pg_MultiInsert); my_error != 94 {
				log.Printf("Oracle_to_PG() mode 1 my_error:%v\n", my_error)
			}else{
				st_Oracle_to_PG=1
			}
			log.Printf("-")
		}

		if st_LDAP_to_PG>0 && st_Oracle_to_PG>0 && st_MSSQL_to_PG>0 {
			if st_DB_clean == 0 {
				if my_error = SABFunctions.RemoveNoChildrenCache(&rconf); my_error != 94 {
					log.Printf("RemoveNoChildrenCache error: %v\n", my_error)
				}else{
					st_DB_clean=1
				}
			}
			if st_DB_clean>0{
				if st_LDAP_UP == 0 {
					if my_error = SABFunctions.LDAP_Make(&rconf); my_error != 94 {
						log.Printf("LDAP_Make error: %v\n", my_error)
					}else{
						st_LDAP_UP=1
					}
				}
			}
		}

		sleep_counter++


		if sleep_counter>rconf.Sleep_cycles-1 {

			st_LDAP_to_PG=0
			st_Oracle_to_PG=0
			st_MSSQL_to_PG=0

			st_DB_clean=0

			st_LDAP_UP=0

			sleep_counter=0

		}

		log.Printf("----- Cycle %d of %d ----- Sleep for %d sec...", sleep_counter, rconf.Sleep_cycles, rconf.Sleep_Time)

		SABModules.Log_OFF()

		time.Sleep(time.Duration(rconf.Sleep_Time)*time.Second)

	}

}
