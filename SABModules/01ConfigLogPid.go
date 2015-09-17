package SABModules

import (
	"log"
	"os"
	"fmt"
	"flag"
	"strconv"
	"syscall"

	"encoding/json"
)

type	Config_STR	struct {

	PG_DSN				string

	AST_SQLite_DB		string
	AST_CID_Group		string
	AST_Num_Start		string

	AST_ARI_Host		string
	AST_ARI_Port		int
	AST_ARI_User		string
	AST_ARI_Pass		string

	Oracle_SRV			[][]string

	MSSQL_DSN			[][]string

	LDAP_URL			[][]string

	ROOT_OU				string

	ROOT_DN				[][]string

	Sleep_Time			int
	Sleep_cycles		int

	LOG_File			string
	PID_File			string

	TRANS_NAMES			[][]string

	BlackList_OU		[]string

	WLB_Listen_IP		string
	WLB_Listen_PORT		int
	WLB_LDAP_ATTR		[][]string
}

var Flog_File *os.File

func ParseCommandLine(config_file string, daemon_mode string) (string, string) {
	ConfigPtr := flag.String("config", config_file, "Path to Configuration file")
	DaemonPtr := flag.String("daemon", daemon_mode, "Fork as system daemon (YES or NO)")
	flag.Parse()
	config_file=*ConfigPtr
	daemon_mode=*DaemonPtr
	log.Printf("path to Configuration file: %s", config_file)
//	log.Printf("Daemon mode: %s", daemon_mode)
	return config_file, daemon_mode
}


func ReadConfigFile (config_file string, conf *Config_STR){
	conf_file, err := os.Open(config_file)
	if err != nil {
		log.Fatalf("Error open Configuration file %s %v\n", config_file, err)
	}

	conf_decoder := json.NewDecoder(conf_file)
	err = conf_decoder.Decode(&conf)
	if err != nil {
		log.Fatalf("Error read SABDefine.Configuration file %s %v\n", config_file, err)
	}

	conf_file.Close()
}

func Log_ON(conf *Config_STR){
	Flog_File, err := os.OpenFile(conf.LOG_File, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error open log file: %v (%s)", err, conf.LOG_File)
	}

	log.SetOutput(Flog_File)
}

func Log_OFF(){
	Flog_File.Close()
}

func Pid_Check(conf *Config_STR){

	var	(
		fpid_File	*os.File
	)

	fpid_File, err := os.OpenFile(conf.PID_File, os.O_RDONLY, 0666)
	if err != nil {
		return
	}

	defer fpid_File.Close()

	pid_read := make([]byte, 10)

	pid_bytes, err := fpid_File.Read(pid_read)
	if err != nil {
		log.Printf("> Remove old pid file")
		os.Remove(conf.PID_File)
		return
	}

	if pid_bytes>0 {
		pid_read_int, err := strconv.Atoi(fmt.Sprintf("%s", pid_read[0:pid_bytes]))
		if err != nil {
			log.Printf("> Remove old pid file")
			os.Remove(conf.PID_File)
			return
		}

		pid_proc, err := os.FindProcess(pid_read_int)
		if err != nil {
			log.Printf("> Remove old pid file")
			os.Remove(conf.PID_File)
			return
		}

		err = pid_proc.Signal(syscall.Signal(0))
		if err != nil {
			log.Printf("> Remove old pid file")
			os.Remove(conf.PID_File)
			return
		}

		log.Printf("> Another copy of the program with PID %d is running! Exiting...", pid_read_int)
		os.Exit(1)

	}else{
		log.Printf("> Remove old pid file")
		os.Remove(conf.PID_File)
		return
	}
}

func Pid_ON(conf *Config_STR){

	var	fpid_File *os.File

	fpid_File, err := os.OpenFile(conf.PID_File, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error open pid file: %v (%s)", err, conf.PID_File)
	}

	fpid_File.WriteString(fmt.Sprintf("%d", os.Getpid()))

	fpid_File.Close()

}

func Pid_OFF(conf *Config_STR){
	os.Remove(conf.PID_File)
}

