package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// ProjectResponse represents the result of a project creation operation
type ProjectResponse struct {
	Created bool   `json:"created"`
	Message string `json:"message"`
}

// ProjectResult represents the result of generic project operations
type ProjectResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var AppDataDir = "./app-data"

// GetProjectDir returns the base directory for lazy-backend projects
func GetProjectDir() (string, error) {
	const folderName = ".lazy-backend"

	basePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %v", err)
	}

	path := filepath.Join(basePath, folderName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

// CreateNewProject creates a new project directory with default config files
func CreateNewProject(name string) ProjectResponse {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return ProjectResponse{Created: false, Message: err.Error()}
	}

	newProjectPath := filepath.Join(projectRootDir, name)

	if _, err := os.Stat(newProjectPath); !os.IsNotExist(err) {
		return ProjectResponse{Created: false, Message: "project name already exists"}
	}

	err = os.MkdirAll(newProjectPath, 0755)
	if err != nil {
		return ProjectResponse{Created: false, Message: err.Error()}
	}

	// Create env.json
	envData := map[string]interface{}{
		"database": []interface{}{},
		"ssh":      []interface{}{},
		"apiKeys":  []interface{}{},
	}

	envBytes, _ := json.MarshalIndent(envData, "", "  ")
	err = os.WriteFile(filepath.Join(newProjectPath, "env.json"), envBytes, 0644)
	if err != nil {
		return ProjectResponse{Created: false, Message: err.Error()}
	}

	// Create schema.json
	err = os.WriteFile(filepath.Join(newProjectPath, "schema.json"), []byte("{}"), 0644)
	if err != nil {
		return ProjectResponse{Created: false, Message: err.Error()}
	}

	return ProjectResponse{Created: true, Message: "project was created"}
}

// DeleteProject deletes a project directory recursively
func DeleteProject(projectName string) ProjectResult {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return ProjectResult{Success: false, Message: err.Error()}
	}

	projectPath := filepath.Join(projectRootDir, projectName)

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return ProjectResult{Success: false, Message: "Failed to delete project: path does not exist"}
	}

	err = os.RemoveAll(projectPath)
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to delete project: %v", err)}
	}

	return ProjectResult{Success: true, Message: "Project deleted successfully"}
}

// RenameProject renames a project directory
func RenameProject(oldName string, newName string) ProjectResult {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return ProjectResult{Success: false, Message: err.Error()}
	}

	oldPath := filepath.Join(projectRootDir, oldName)
	newPath := filepath.Join(projectRootDir, newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return ProjectResult{Success: false, Message: "Failed to rename project: source path does not exist"}
	}

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		return ProjectResult{Success: false, Message: "Failed to rename project: target name already exists"}
	}

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to rename project: %v", err)}
	}

	return ProjectResult{Success: true, Message: "Project renamed successfully"}
}

// GetProjectEnv reads the env.json for a given project
func GetProjectEnv(projectName string) (map[string]interface{}, error) {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(projectRootDir, projectName, "env.json")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return make(map[string]interface{}), nil
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListProjects returns a list of project names found in the project root directory
func ListProjects() ([]string, error) {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return nil, err
	}
	fmt.Print("nnnjs")
	entries, err := os.ReadDir(projectRootDir)
	if err != nil {
		return nil, err
	}

	projects := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}
	return projects, nil
}

// SaveProjectEnv saves the env.json for a given project
func SaveProjectEnv(projectName string, envData map[string]interface{}) ProjectResult {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return ProjectResult{Success: false, Message: err.Error()}
	}

	envPath := filepath.Join(projectRootDir, projectName, "env.json")

	envBytes, err := json.MarshalIndent(envData, "", "  ")
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to marshal config: %v", err)}
	}

	err = os.WriteFile(envPath, envBytes, 0644)
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to write config file: %v", err)}
	}

	return ProjectResult{Success: true, Message: "Configuration saved successfully"}
}

// GetProjectSchema reads the schema.json for a given project
func GetProjectSchema(projectName string) (map[string]IntermediateSchema, error) {
	var schema map[string]IntermediateSchema
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return schema, err
	}

	schemaPath := filepath.Join(projectRootDir, projectName, "schema.json")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return schema, nil
	}

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return schema, err
	}

	err = json.Unmarshal(data, &schema)
	if err != nil {
		return schema, err
	}

	return schema, nil
}

// SaveProjectSchema saves the schema.json for a given project
func SaveProjectSchema(projectName string, schema map[string]IntermediateSchema) ProjectResult {
	projectRootDir, err := GetProjectDir()
	if err != nil {
		return ProjectResult{Success: false, Message: err.Error()}
	}

	schemaPath := filepath.Join(projectRootDir, projectName, "schema.json")

	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to marshal schema: %v", err)}
	}

	err = os.WriteFile(schemaPath, schemaBytes, 0644)
	if err != nil {
		return ProjectResult{Success: false, Message: fmt.Sprintf("Failed to write schema file: %v", err)}
	}

	return ProjectResult{Success: true, Message: "Schema saved successfully"}
}

func SaveRecent(project []map[string]string) {
	projectRootDir := "./app-data"
	if _, err := os.Stat(projectRootDir); os.IsNotExist(err) {
		os.MkdirAll(projectRootDir, 0755)
	}
	projectJson, _ := json.Marshal(project)
	os.WriteFile(path.Join(projectRootDir, "recent.json"), projectJson, 0644)
}

func GetRecent() []map[string]string {
	data, err := os.ReadFile(path.Join(AppDataDir, "recent.json"))
	if err != nil {
		return nil
	}
	var recent []map[string]string
	err = json.Unmarshal(data, &recent)
	if err != nil {
		return nil
	}
	return recent
}

func SaveApiSchema(data string, project string) error {
	dir, _ := GetProjectDir()
	err := os.WriteFile(path.Join(dir, project, "api.json"), []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetApiSchema(project string) ([]any, error) {
	dir, _ := GetProjectDir()

	var marshaledData = make([]any, 0)
	data, err := os.ReadFile(path.Join(dir, project, "api.json"))
	if err != nil {
		return marshaledData, err
	}

	err = json.Unmarshal(data, &marshaledData)
	return marshaledData, err

}
