package config

import (
	"path/filepath"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var conf *Configuration

func InitConfig(configPath string) {
	if conf == nil {
		filename, _ := filepath.Abs(configPath)
		log.Println("trying to read file ", filename)
		yamlFile, err := ioutil.ReadFile(filename)
		var confi Configuration
		err = yaml.Unmarshal(yamlFile, &confi)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		conf = &confi
	}
}

func Get() *Configuration {
	return conf
}

type Configuration struct {
	BrokersList []string  `yaml:"brokers_list"`
	ProducerTopic		string `yaml:"producer_topic"`
	ConsumerTopics		[]string `yaml:"consumer_topics"`
	ConsumerGroupId     string `yaml:"consumer_group_id"`
}
