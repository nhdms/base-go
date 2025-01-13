package generator

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/nhdms/base-go/pkg/common"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	appImagePrefix = "agbiz/go"
)

func listCICDPaths(rootDir string) ([]string, error) {
	cicdPaths := []string{}
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "cicd.json" {
			cicdPaths = append(cicdPaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return cicdPaths, nil
}

func GenerateCICD(path string) {
	// check folder exists
	folderToDeploy := `overlays/dev`
	if _, err := os.Stat(folderToDeploy); os.IsNotExist(err) {
		log.Fatalf("You need to cd to the deployment repository")
	}

	var err error
	allCICDs := []string{
		path,
	}

	if !strings.HasSuffix(path, "cicd.json") {
		allCICDs, err = listCICDPaths(path)
		if err != nil {
			log.Fatalf("Error listing cicd paths: %v", err)
			return
		}
	}

	for _, cicdPath := range allCICDs {
		// read cicd file
		serviceConfig := &common.ServiceConfig{}
		jsonBytes, err := os.ReadFile(cicdPath)
		if err != nil {
			log.Fatalf("Error reading cicd file: %v", err)
			return
		}

		_ = json.Unmarshal(jsonBytes, serviceConfig)
		err = serviceConfig.Validate()
		if err != nil {
			log.Fatalf("Error validating cicd file: %v", err)
			return
		}

		generateFromConfig(serviceConfig)

	}
}

type CICD struct {
	FilePath string
	Content  string
}

const (
	serviceDetailTemplate = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../../../base

namePrefix: {{.serviceName}}-{{.appType}}-

commonLabels:
  app: {{.serviceName}}
  service-type: {{.appType}}
  environment: dev

configMapGenerator:
- name: app-config
  behavior: merge
  literals:
  - SERVICE_NAME={{.serviceName}}
  - APP_TYPE={{.appType}}
  - CONFIG_REMOTE_KEYS={{.remoteKeys}}

images:
- name: placeholder
  newName: {{.imageName}}
  newTag: '1'
`
)

func generateFromConfig(cfg *common.ServiceConfig) {
	resp := []CICD{}

	//apps/service/user-service.yaml
	resp = append(resp, CICD{
		FilePath: fmt.Sprintf(`apps/%s/%s.yaml`, cfg.AppType, cfg.ServiceName),
		Content: fmt.Sprintf(`service_name: %s
app_type: %s
image:
  repository: registry.your-domain.com/img
`, cfg.ServiceName, cfg.AppType),
	})

	tmpl := template.Must(template.New("service").Parse(serviceDetailTemplate))
	vars := map[string]string{
		"serviceName": cfg.ServiceName,
		"appType":     cfg.AppType,
		"remoteKeys":  strings.Join(cfg.ConfigRemoteKeys, ","),
		"imageName":   fmt.Sprintf(`%s-%s-%s`, appImagePrefix, cfg.AppType, cfg.ServiceName),
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, vars); err != nil {
		log.Fatalln("Failed to execute template", err.Error())
	}

	resp = append(resp, CICD{
		FilePath: fmt.Sprintf(`overlays/dev/%s/%s/kustomization.yaml`, cfg.AppType, cfg.ServiceName),
		Content:  result.String(),
	})

	if len(resp) == 0 {
		log.Println("No CICD files generated")
		return
	}

	for _, r := range resp {
		_ = os.MkdirAll(filepath.Dir(r.FilePath), os.ModePerm)
		err := os.WriteFile(r.FilePath, []byte(r.Content), 0600)
		if err != nil {
			log.Fatalf("Error writing file %s: %v", r.FilePath, err)
		}

		log.Println("Generated file", r.FilePath)
	}

	log.Printf("All CICD files generated, please register service %s/%s at `overlays/dev/kustomization.yaml\n`", cfg.AppType, cfg.ServiceName)
	log.Println("Please run \nkubectl kustomize overlays/dev\n to preview before applying")
	return
}
