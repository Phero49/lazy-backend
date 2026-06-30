package main

import (
	"context"
	"fmt"
	"lazy-backend/io"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CreateNewProject creates a new project
func (a *App) CreateNewProject(name string) io.ProjectResponse {
	return io.CreateNewProject(name)
}

// ListProjects returns a list of projects
func (a *App) ListProjects() ([]string, error) {
	return io.ListProjects()
}

// DeleteProject deletes a project
func (a *App) DeleteProject(name string) io.ProjectResult {
	return io.DeleteProject(name)
}

// RenameProject renames a project
func (a *App) RenameProject(oldName, newName string) io.ProjectResult {
	return io.RenameProject(oldName, newName)
}

// GetProjectEnv gets the environment config for a project
func (a *App) GetProjectEnv(name string) (map[string]interface{}, error) {
	return io.GetProjectEnv(name)
}

// SaveProjectEnv saves the environment config for a project
func (a *App) SaveProjectEnv(name string, envData map[string]interface{}) io.ProjectResult {
	return io.SaveProjectEnv(name, envData)
}

// GetDatabaseList returns a list of databases for a given host
func (a *App) GetDatabaseList(host string, port int, user, password string) ([]string, error) {
	//parser := io.MariaDbParser{}
	return make([]string, 0), nil
}

// GetTableList returns a list of tables for a given database
func (a *App) GetTableList(host string, port int, user, password, database string) ([]string, error) {
	//parser := io.MariaDbParser{}
	return make([]string, 0), nil

}

// GetTableSchema returns the intermediate schema for a specific table
func (a *App) GetTableSchema(host string, port int, user, password, database, tableName string) ([]io.Field, error) {
	//parser := io.MariaDbParser{}
	return make([]io.Field, 0), nil
}

// ImportDatabaseSchema introspects the entire database and returns the intermediate schema
func (a *App) ImportDatabaseSchema(host string, port int, user, password, database string) (io.IntermediateSchema, error) {
	//parser := io.MariaDbParser{}
	var p io.IntermediateSchema
	return p, nil
}

// GetProjectSchema gets the schema for a project
func (a *App) GetProjectSchema(name string) (map[string]io.IntermediateSchema, error) {
	return io.GetProjectSchema(name)
}

// SaveProjectSchema saves the schema for a project
func (a *App) SaveProjectSchema(name string, schema map[string]io.IntermediateSchema) io.ProjectResult {
	return io.SaveProjectSchema(name, schema)
}

// SaveProjectSchema saves the schema for a project
func (a *App) SaveApiSchema(data string, project string) error {
	return io.SaveApiSchema(data, project)
}

// SaveProjectSchema saves the schema for a project
func (a *App) GetApiSchema(project string) ([]any, error) {
	return io.GetApiSchema(project)
}
