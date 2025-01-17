package generator

import (
	"bytes"
	"fmt"
	"gitlab.com/a7923/athena-go/pkg/logger"
	"gitlab.com/a7923/athena-go/pkg/toolkit/templates"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GeneratorData struct {
	Name        string
	Handler     string
	ServiceName string
	QueueName   string
	TableName   string
	Port        int
}

func generateFile(tmpl string, data GeneratorData, outputPath string) error {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func GenerateAPI(name string) error {
	basePath := fmt.Sprintf("cmd/apis/%s", name)
	handler := strings.Title(strings.TrimSuffix(name, "-api"))
	serviceName := strings.TrimSuffix(name, "-api")

	data := GeneratorData{
		Name:        name,
		Handler:     handler,
		ServiceName: serviceName,
		QueueName:   strings.Replace(name, "-api", "_events", 1),
		Port:        30000, // You might want to make this configurable
	}

	files := map[string]string{
		fmt.Sprintf("%s/main.go", basePath):             templates.APIMainTemplate,
		fmt.Sprintf("%s/app/webserver.go", basePath):    templates.APIWebServerTemplate,
		fmt.Sprintf("%s/handlers/handler.go", basePath): templates.APIHandlerTemplate,
		fmt.Sprintf("%s/config/cicd.json", basePath):    templates.APICICDTemplate,
	}

	for path, tmpl := range files {
		if err := generateFile(tmpl, data, path); err != nil {
			return err
		}
	}

	// Create empty config.toml
	configPath := fmt.Sprintf("%s/config/config.toml", basePath)
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(configPath, []byte("# Add your configuration here"), 0644); err != nil {
		return err
	}

	logger.AthenaLogger.Infof("API %s generated successfully\nExploring at %s", name, basePath)
	return nil
}

func GenerateConsumer(name string) error {
	basePath := fmt.Sprintf("cmd/consumers/%s", name)
	handler := strings.Title(strings.TrimSuffix(name, "-consumer"))

	data := GeneratorData{
		Name:        name,
		Handler:     handler,
		ServiceName: strings.Replace(name, "-consumer", "", 1),
	}

	files := map[string]string{
		fmt.Sprintf("%s/main.go", basePath):             templates.ConsumerMainTemplate,
		fmt.Sprintf("%s/handlers/handler.go", basePath): templates.ConsumerHandlerTemplate,
	}

	for path, tmpl := range files {
		if err := generateFile(tmpl, data, path); err != nil {
			return err
		}
	}

	return nil
}
