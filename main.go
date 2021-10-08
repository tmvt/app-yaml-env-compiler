package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	fmt.Println("Ready to compile ...")

	filename, _ := filepath.Abs("app.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var mapResult map[interface{}]interface{}
	err = yaml.Unmarshal(yamlFile, &mapResult)
	if err != nil {
		panic(err)
	}

	for k, any := range mapResult {
		if k == "env_variables" {
			err := checkIsPointer(&any)
			if err != nil {
				panic(err)
			}
			valueOf := reflect.ValueOf(any)
			val := reflect.Indirect(valueOf)
			switch val.Type().Kind() {
			case reflect.Map:
				envMap := any.(map[interface{}]interface{})
				for in, iv := range envMap {

					envName := in.(string)
					envVal := iv.(string)

					// Check if var is supposed to be replaced
					if strings.HasPrefix(envVal, "$") {
						env := strings.Replace(strings.TrimSpace(envVal), "$", "", -1)
						// Do not replace if no new value has been set
						if nv := os.Getenv(env); nv != "" {
							envMap[envName] = nv
							fmt.Printf("Replaced %v!\n", envName)
						} else {
							fmt.Printf("No new value for %v!\n", envName)
						}
					}
				}
			default:
				panic(fmt.Sprintf("This is not supposed to happen, but if it does, good luck"))
			}
		}
	}

	fmt.Println(fmt.Sprintf("Compiled env variables: %v", mapResult["env_variables"]))

	out, err := yaml.Marshal(mapResult)
	// write the whole body at once
	err = ioutil.WriteFile("app.yaml", out, 0644)
	if err != nil {
		panic(err)
	}
}

func checkIsPointer(any interface{}) error {
	if reflect.ValueOf(any).Kind() != reflect.Ptr {
		return fmt.Errorf("you passed something that was not a pointer: %s", any)
	}
	return nil
}
