package io

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DependencySchema struct {
	Version       int64    `json:"version"`
	LangExtension string   `json:"lang-extension"`
	Description   string   `json:"description"`
	Publisher     string   `json:"publisher"`
	Install       []string `json:"install"`
	Modules       []Module `json:"modules"`
}

type Module struct {
	Path    string    `json:"path"`
	Exports []Exports `json:"exports"`
}

type Exports struct {
	Description *string        `json:"description,omitempty"`
	Name        string         `json:"name"`
	Params      map[string]any `json:"params"`
	ReturnType  string         `json:"returnType"`
}

type Dependence struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	ExportsCount int     `json:"exportsCount"`
	Version      float64 `json:"version"`
	Publisher    string  `json:"publisher"`

	//absolute path to the module folder
	Location string `json:"location"`
}

// LoadLocalDependencies returns error and path to the dependencies if validation is completed
func LoadLocalDependencies(ctx context.Context) (Dependence, error) {
	folder, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: "select lazy builder dependency folder",
	})
	if err != nil {
		return Dependence{}, err
	}

	manifestPath := path.Join(folder, "manifest.json")
	exports, err := GetDependencyExports(manifestPath)
	if err != nil {
		return Dependence{}, err
	}

	// Read the manifest file to get dependency metadata
	bytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return Dependence{}, err
	}

	var schema DependencySchema
	err = json.Unmarshal(bytes, &schema)
	if err != nil {
		return Dependence{}, err
	}

	// Create dependence info
	dependence := Dependence{
		Version:      float64(schema.Version),
		Name:         path.Base(folder),
		Description:  schema.Description,
		ExportsCount: len(exports),
		Location:     folder,
		Publisher:    schema.Publisher,
	}
	InstallDependencyModule(dependence)
	return dependence, nil
}

// GetDependencyExports reads and validates a dependency manifest file
func GetDependencyExports(path string) ([]Exports, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schema DependencySchema
	err = json.Unmarshal(bytes, &schema)
	if err != nil {
		return nil, err
	}

	// Validate the schema
	if schema.Version <= 0 {
		return nil, ErrInvalidVersion
	}
	if schema.LangExtension == "" {
		return nil, ErrMissingLangExtension
	}

	// Collect all exports from all modules
	allExports := make([]Exports, 0)
	for _, module := range schema.Modules {
		allExports = append(allExports, module.Exports...)
	}

	return allExports, nil
}

// Custom errors
var (
	ErrInvalidVersion       = &DependencyError{Message: "invalid or missing version in dependency manifest"}
	ErrMissingLangExtension = &DependencyError{Message: "missing language extension in dependency manifest"}
	ErrInvalidManifest      = &DependencyError{Message: "invalid manifest file format"}
)

type DependencyError struct {
	Message string
}

func (e *DependencyError) Error() string {
	return e.Message
}

func InstallDependencyModule(dependency Dependence) error {
	dPath := path.Join(APP_FOLDER, "installed-dependencies")
	if _, err := os.Stat(dPath); os.IsNotExist(err) {
		err := os.MkdirAll(dPath, 0755)
		if err != nil {
			return err
		}
	}

	// Read existing dependency-list.json
	dependencyListPath := path.Join(dPath, "dependency-list.json")
	var dependencyList map[string]Dependence

	if _, err := os.Stat(dependencyListPath); err == nil {
		data, err := os.ReadFile(dependencyListPath)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			json.Unmarshal(data, &dependencyList)
		}
	}

	if dependencyList == nil {
		dependencyList = make(map[string]Dependence)
	}

	// Add/update dependency
	name := dependency.Name
	dependencyList[fmt.Sprint(name, '@', dependency.Version)] = dependency

	// Marshal and write
	marshalled, err := json.MarshalIndent(dependencyList, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dependencyListPath, marshalled, 0644)
	if err != nil {
		return err
	}

	return nil
}

// resolveProjectPath accepts either a project name or a real project directory path.
// It resolves the value into a directory that can safely be used for dependency metadata.
func resolveProjectPath(projectPath string) (string, error) {
	projectPath = strings.TrimSpace(projectPath)
	if projectPath == "" {
		return "", fmt.Errorf("project path is required")
	}

	if info, err := os.Stat(projectPath); err == nil {
		if info.IsDir() {
			return filepath.Clean(projectPath), nil
		}
		return "", fmt.Errorf("project path is not a directory")
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat project path: %w", err)
	}

	projectRootDir, err := GetProjectDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(projectRootDir, projectPath), nil
}

func AddProjectDependency(projectPath string, dependency Dependence) error {
	resolvedProjectPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return err
	}

	filePath := filepath.Join(resolvedProjectPath, "dependency-list.json")
	bytes, err := os.ReadFile(filePath)
	exist := !os.IsNotExist(err)

	dependencies := make(map[string]Dependence)

	if exist {
		if err := json.Unmarshal(bytes, &dependencies); err != nil {
			return fmt.Errorf("failed to unmarshal dependencies: %w", err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read dependency file: %w", err)
	}

	dependencies[dependency.Name] = dependency

	data, err := json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dependencies: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write dependency file: %w", err)
	}

	return nil
}

func RemoveProjectDependency(projectPath string, name string) error {
	resolvedProjectPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return err
	}

	filePath := filepath.Join(resolvedProjectPath, "dependency-list.json")
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("dependency file does not exist")
		}
		return fmt.Errorf("failed to read dependency file: %w", err)
	}

	dependencies := make(map[string]Dependence)
	if err := json.Unmarshal(bytes, &dependencies); err != nil {
		return fmt.Errorf("failed to unmarshal dependencies: %w", err)
	}

	// Check if dependency exists before deleting
	if _, exists := dependencies[name]; !exists {
		return fmt.Errorf("dependency '%s' not found", name)
	}

	// Remove dependency
	delete(dependencies, name)

	// Marshal back to JSON
	data, err := json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal dependencies: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write dependency file: %w", err)
	}

	return nil
}

func GetAllInstalledDependencies() ([]Dependence, error) {
	dPath := path.Join(APP_FOLDER, "installed-dependencies")
	dependencyListPath := path.Join(dPath, "dependency-list.json")

	if _, err := os.Stat(dependencyListPath); err != nil {
		if os.IsNotExist(err) {
			return []Dependence{}, nil
		}
		return nil, fmt.Errorf("failed to stat installed dependency list: %w", err)
	}

	data, err := os.ReadFile(dependencyListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read installed dependency list: %w", err)
	}
	if len(data) == 0 {
		return []Dependence{}, nil
	}

	dependencyList := make(map[string]Dependence)
	if err := json.Unmarshal(data, &dependencyList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installed dependency list: %w", err)
	}

	dependencies := make([]Dependence, 0, len(dependencyList))
	for _, dependency := range dependencyList {
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

// GetProjectUsedDependencies reads the dependency list for the provided project path.
// Fixed: the function now accepts the project directory explicitly instead of always reading from the global app folder.
func GetProjectUsedDependencies(projectPath string) ([]Dependence, error) {
	resolvedProjectPath, err := resolveProjectPath(projectPath)
	if err != nil {
		return nil, err
	}

	projectDependencyListPath := filepath.Join(resolvedProjectPath, "dependency-list.json")

	if _, err := os.Stat(projectDependencyListPath); err != nil {
		if os.IsNotExist(err) {
			return []Dependence{}, nil
		}
		return nil, fmt.Errorf("failed to stat project dependency list: %w", err)
	}

	data, err := os.ReadFile(projectDependencyListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project dependency list: %w", err)
	}
	if len(data) == 0 {
		return []Dependence{}, nil
	}

	dependencyList := make(map[string]Dependence)
	if err := json.Unmarshal(data, &dependencyList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project dependency list: %w", err)
	}

	dependencies := make([]Dependence, 0, len(dependencyList))
	for _, dependency := range dependencyList {
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func GetDependencyManifest(name string) (*DependencySchema, error) {
	dependencyListPath := path.Join(APP_FOLDER, "dependency-list.json")
	if _, err := os.Stat(dependencyListPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to stat installed dependency list: %w", err)
	}

	data, err := os.ReadFile(dependencyListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read installed dependency list: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("failed to dependency the manifest is empty or corrupt")
	}

	var dependencyList *DependencySchema
	if err := json.Unmarshal(data, &dependencyList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installed dependency list: %w", err)
	}

	return dependencyList, nil
}
