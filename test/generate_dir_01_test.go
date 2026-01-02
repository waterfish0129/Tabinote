package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var stSeparator = string(filepath.Separator)
var stWorkDir, _ = os.Getwd()
var stRootDir = stWorkDir[:strings.LastIndex(stWorkDir, stSeparator)]

const stJsonFileName = "dir.json"

var iJsonData map[string]any

func loadJson() {
	gnJsonBytes, _ := os.ReadFile(stWorkDir + stSeparator + stJsonFileName)
	err := json.Unmarshal(gnJsonBytes, &iJsonData)
	if err != nil {
		panic("Load Json Data Error:" + err.Error())
	}
}

func parseMap(mapData map[string]any, stParentDir string) {
	for _, value := range mapData {
		switch value.(type) {
		case string:
			{
				path, _ := value.(string)
				if path == "" {
					continue
				}
				if stParentDir != "" {
					path = stParentDir + stSeparator + path
				}
				stParentDir = path

				createDir(path)
			}
		case []any:
			{
				parseArray(value.([]any), stParentDir)
			}
		}
	}
}
func parseArray(slJsonData []any, stParentDir string) {
	for _, value := range slJsonData {
		mapV, _ := value.(map[string]any)
		parseMap(mapV, stParentDir)
	}
}
func createDir(path string) {
	if path == "" {
		return
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic("Create Dir Error:" + err.Error())
	}
}
func TestGenerateDir01(T *testing.T) {
	loadJson()
	parseMap(iJsonData, stRootDir)
}
