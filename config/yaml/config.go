package yaml

import (
	"github.com/duanchi/min/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var configInstance interface{}

var configFile = "./config/application.yaml"

func GetConfig(config interface{}) (err error){
	configYaml, err := readFile()
	if err != nil {
		return
	}

	pattern, _ := regexp.Compile(`\${.+?}`)
	configYaml = pattern.ReplaceAllFunc(configYaml, func(b []byte) []byte {
		s := string(b)
		value := strings.SplitN(s[2:len(s) - 1], ":", 2)
		if len(value) > 1 {
			return []byte(util.Getenv(value[0], value[1]))
		} else {
			return []byte(util.Getenv(value[0], ""))
		}
	})

	err = yaml.Unmarshal(configYaml, config)
	return
}

func readFile() (config []byte, err error){
	config, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	return
}

func SetConfigFile (config string) {
	configFile = config
}