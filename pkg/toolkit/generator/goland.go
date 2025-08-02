package generator

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	log "go-micro.dev/v5/logger"
	"os"
	"path/filepath"
)

type Cicd struct {
	AppType          string   `json:"app_type"`
	CmdBinDir        string   `json:"cmd_bin_dir"`
	ServiceName      string   `json:"service_name"`
	Port             int      `json:"port"`
	ConfigRemoteKeys []string `json:"config_remote_keys"`
}

type moduleConfig struct {
	Name string `xml:"name,attr"`
}

type parameter struct {
	Value string `xml:"value,attr"`
}

type method struct {
	Value string `xml:"v,attr"`
}

type Configuration struct {
	Default          bool         `xml:"default,attr"`
	Name             string       `xml:"name,attr"`
	Type             string       `xml:"type,attr"`
	FactoryName      string       `xml:"factoryName,attr"`
	Module           moduleConfig `xml:"module"`
	WorkingDirectory parameter    `xml:"working_directory"`
	GoParameters     parameter    `xml:"go_parameters"`
	Parameters       parameter    `xml:"parameters"`
	Kind             parameter    `xml:"kind"`
	FilePath         parameter    `xml:"filePath"`
	Package          parameter    `xml:"package"`
	Directory        parameter    `xml:"directory"`
	Method           method       `xml:"method"`
}

type component struct {
	Name          string        `xml:"name,attr"`
	Configuration Configuration `xml:"configuration"`
}

type IdeConfig struct {
	Component component `xml:"component"`
}

func New(cicd Cicd, namePackage string) component {
	return component{
		Name: "ProjectRunConfigurationManager",
		Configuration: Configuration{
			Default:     false,
			Name:        fmt.Sprintf("%v-%v", cicd.ServiceName, cicd.AppType),
			Type:        "GoApplicationRunConfiguration",
			FactoryName: "Go Application",
			Module: moduleConfig{
				Name: namePackage,
			},
			WorkingDirectory: parameter{
				Value: fmt.Sprintf("$PROJECT_DIR$/%v", cicd.CmdBinDir),
			},
			//GoParameters: parameter{
			//	Value: "-i",
			//},
			//Parameters: parameter{
			//	Value: strings.Join(AdditionArgs, " "),
			//},
			Kind: parameter{
				Value: "FILE",
			},
			FilePath: parameter{
				Value: fmt.Sprintf("$PROJECT_DIR$/%v/main.go", cicd.CmdBinDir),
			},
			Package: parameter{
				Value: namePackage,
			},
			Directory: parameter{
				Value: "$PROJECT_DIR$/",
			},
			Method: method{
				Value: "2",
			},
		},
	}
}

func ExportGolandConfig(args []string) {
	var currentDir string
	if len(args) > 0 {
		currentDir = args[0]
	} else {
		currentDir, _ = os.Getwd()
	}

	pathProject := fmt.Sprintf("%v/.idea", currentDir)

	if _, err := os.Stat(pathProject); os.IsNotExist(err) {
		log.Warn("Directory a has never been opened by goland")
		return
	}

	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "cicd.json" {
			RunCommand(currentDir, filepath.Dir(path))
		}
		return nil
	})

	if err != nil {
		log.Warn(fmt.Sprintf("Error %v", err))
	}
}

func RunCommand(currentDir, dir string) {
	println(dir)

	file, err := os.ReadFile(dir + "/cicd.json")

	cicd := Cicd{}

	err = json.Unmarshal(file, &cicd)
	path := currentDir + "/.idea/runConfigurations"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	nameFile := fmt.Sprintf("%v/%v-%v.xml", path, cicd.ServiceName, cicd.AppType)

	writer, err := os.OpenFile(nameFile, os.O_RDWR|os.O_CREATE, os.ModePerm)

	output, err := xml.MarshalIndent(New(cicd, filepath.Base(currentDir)), "  ", "    ")
	if err != nil {
		println(fmt.Sprintf("Error 1 %v", err))
	}

	_, err = writer.Write(output)

	if err != nil {
		println(fmt.Sprintf("Error 2 %v", err))
	}

}
