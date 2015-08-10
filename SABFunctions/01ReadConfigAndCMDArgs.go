package SABFunctions

import (
	"log"
	"os"
	"flag"

	"encoding/json"

	"github.com/BestianRU/SABookServices/SABDefine"

)

func ParseCommandLine(config_file string) string {
	ConfigPtr := flag.String("config", config_file, "Path to Configuration file")
	flag.Parse()
	config_file=*ConfigPtr
	log.Printf("path to Configuration file: %s", config_file)
	return config_file
}


func ReadConfigFile (config_file string){
	conf_file, err := os.Open(config_file)
	if err != nil {
		log.Fatalf("Error open Configuration file %s %v\n", config_file, err)
	}

	conf_decoder := json.NewDecoder(conf_file)
	err = conf_decoder.Decode(&SABDefine.Conf)
	if err != nil {
		log.Fatalf("Error read SABDefine.Configuration file %s %v\n", config_file, err)
	}

	conf_file.Close()
}

