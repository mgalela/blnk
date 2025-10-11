package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type CoA struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Desc    string            `json:"desc"`
	Version string            `json:"version"`
	List    map[string]string `json:"list"`
}

func (c *CoA) GetID() string {
	return c.ID
}

func (c *CoA) GetName() string {
	return c.Name
}

func (c *CoA) GetDesc() string {
	return c.Desc
}

func (c *CoA) GetVersion() string {
	return c.Version
}

func (c *CoA) CheckCode(code string) string {
	return c.List[code]
}

func (c *CoA) CreateList(configDir string, files []os.DirEntry) error {

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(configDir, file.Name())

			// Create a temporary Viper instance for each file
			tempViper := viper.New()
			tempViper.SetConfigFile(filePath)
			tempViper.SetConfigType("json")

			if err := tempViper.ReadInConfig(); err != nil {
				fmt.Printf("Error reading config file %s: %v\n", filePath, err)
				continue
			}

			for key, value := range tempViper.AllSettings() {
				if value == nil {
					continue
				}
				switch key {
				case "id":
					strVal, ok := value.(string)
					if !ok || strVal == "" {
						continue
					}
					c.ID = strVal
				case "name":
					strVal, ok := value.(string)
					if !ok || strVal == "" {
						continue
					}
					c.Name = strVal
				case "desc":
					strVal, ok := value.(string)
					if !ok || strVal == "" {
						continue
					}
					c.Desc = strVal
				case "version":
					strVal, ok := value.(string)
					if !ok || strVal == "" {
						continue
					}
					c.Version = strVal
				case "list":
					// fmt.Printf("Adding code %T: %v\n", value, value)
					listMap, ok := value.([]interface{})
					if !ok {
						fmt.Printf("Invalid type for list in file %s\n", filePath)
						continue
					}
					for _, v := range listMap {
						// fmt.Printf("Adding listMap %T: %v\n", v, v)
						ltemp, ok := v.(map[string]interface{})
						if !ok {
							fmt.Println("Invalid type for ltemp")
							continue
						}
						var code, name, status string
						for k, v := range ltemp {
							strVal, ok := v.(string)
							if !ok || strVal == "" {
								continue
							}
							if k == "status" {
								status = strVal
							} else {
								code = strings.ReplaceAll(k, ".", "")
								name = strVal
							}
						}
						// fmt.Printf("Adding code %s: %s (%s)\n", code, name, status)
						if status == "enabled" && code != "" {
							if c.List == nil {
								c.List = make(map[string]string)
							}
							c.List[code] = name
						}
					}
					continue
				}
			}
		}
	}

	return nil
}

func main() {
	configDir := "./coa"
	var coa CoA

	files, err := os.ReadDir(configDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	if err := coa.CreateList(configDir, files); err != nil {
		fmt.Printf("Error loading CoA: %v\n", err)
		return
	}

	fmt.Printf("Loaded CoA: %s - %s v%s\n", coa.ID, coa.Name, coa.Version)
	for code, name := range coa.List {
		fmt.Printf("%s: %s\n", code, name)
	}
}
