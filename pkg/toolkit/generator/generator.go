package generator

import (
	"bytes"
	"fmt"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/toolkit/templates"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GeneratorData struct {
	Name         string
	Handler      string
	ServiceName  string
	QueueName    string
	TableName    string
	Port         int
	HandlerLower string
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

	logger.DefaultLogger.Infof("API %s generated successfully\nExploring at %s", name, basePath)
	return nil
}

// Replace with:
func GenerateService(name string) error {
	basePath := fmt.Sprintf("cmd/services/%s", name)
	handler := strings.Title(strings.TrimSuffix(name, "-service"))
	serviceName := strings.TrimSuffix(name, "-service")

	data := GeneratorData{
		Name:         name,
		Handler:      handler,
		ServiceName:  serviceName,
		TableName:    strings.Replace(name, "-service", "_events", 1),
		Port:         40000,
		HandlerLower: strings.ToLower(handler),
	}

	// Generate proto files
	protoFiles := map[string]string{
		fmt.Sprintf("proto/models/%s.proto", serviceName):   templates.ModelProtoTemplate,
		fmt.Sprintf("proto/services/%s.proto", serviceName): templates.ServiceProtoTemplate,
	}

	// Generate service files
	files := map[string]string{
		fmt.Sprintf("%s/main.go", basePath):                  templates.ServiceMainTemplate,
		fmt.Sprintf("%s/handlers/handler.go", basePath):      templates.ServiceHandlerTemplate,
		fmt.Sprintf("%s/handlers/handler_test.go", basePath): templates.ServiceHandlerTestTemplate,
		fmt.Sprintf("%s/tables/table.go", basePath):          templates.ServiceTableTemplate,
		fmt.Sprintf("%s/config/cicd.json", basePath):         templates.ServiceCICDTemplate,
	}

	// Generate all files
	for path, tmpl := range protoFiles {
		if err := generateFile(tmpl, data, path); err != nil {
			return err
		}
	}

	for path, tmpl := range files {
		if err := generateFile(tmpl, data, path); err != nil {
			return err
		}
	}

	logger.DefaultLogger.Infof("Service %s generated successfully\nExploring at %s", name, basePath)
	logger.DefaultLogger.Infof("Please add common function GetGRPCServiceClient() at internal/grpc.go")
	logger.DefaultLogger.Infof("And function define service name at pkg/common/service_name.go")
	logger.DefaultLogger.Infof("\ngcli gen proto proto/models/%s.proto\ngcli gen proto proto/services/%s.proto", name, name)
	logger.DefaultLogger.Info("Run these command to generate proto file")
	return nil
}

func GenerateConsumer(name string) error {
	basePath := fmt.Sprintf("cmd/consumers/%s", name)
	handler := strings.Title(strings.TrimSuffix(name, "-consumer"))
	serviceName := strings.TrimSuffix(name, "-consumer")

	data := GeneratorData{
		Name:        name,
		Handler:     handler,
		ServiceName: serviceName,
	}

	files := map[string]string{
		fmt.Sprintf("%s/main.go", basePath):             templates.ConsumerMainTemplate,
		fmt.Sprintf("%s/handlers/handler.go", basePath): templates.ConsumerHandlerTemplate,
		fmt.Sprintf("%s/config/cicd.json", basePath):    templates.ConsumerCICDTemplate,
		fmt.Sprintf("%s/config/config.toml", basePath):  templates.ConsumerSampleConfig,
	}

	for path, tmpl := range files {
		if err := generateFile(tmpl, data, path); err != nil {
			return err
		}
	}

	logger.DefaultLogger.Infof("Consumer %s generated successfully\nExploring at %s", name, basePath)
	logger.DefaultLogger.Infof("Create key %s at consul with content from file %s to start consumer",
		fmt.Sprintf("consumers/%s.toml", name),
		fmt.Sprintf("%s/config/config.toml", basePath))
	return nil
}
