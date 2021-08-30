package yaml

/*var configInstance interface{}

var configFile = "./config/application.yaml"

func GetMapConfig()(conf map[string]interface{}, err error){
	configFile, err := readFile()
	conf = make(map[string]interface{})
	err = yaml.Unmarshal(configFile, conf)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetYamlConfig(config interface{})(configMap map[interface{}]interface{}, err error){
	configYaml, err := readFile()
	if err != nil {
		return
	}
	configMap = make(map[interface{}]interface{})
	err = yaml.Unmarshal(configYaml, config)
	fmt.Printf("%+v", config)
	return
	if err != nil {
		log.Println(err)
		panic(err.Error())
	} else {
		// fmt.Printf("读取配置文件: %+v\n", conf)
		parseConfig(config, configMap)
		// fmt.Printf("更新配置参数: %+v\n", conf)

		fmt.Printf("%+v", config)
	}
	// configInstance = conf
	return
}

func readFile() (config []byte, err error){
	config, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("yamlFile.Get err %v ", err)
	}
	return
}*/

func Get(key string) interface{} {
	return getRaw(key).Interface()
}