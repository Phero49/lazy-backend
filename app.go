package main

import (
	"context"
	"fmt"
	devserver "lazy-backend/dev-server"
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
	devserver.SetWailsContext(ctx)

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

func (a *App) LoadLocalDependencies() (io.Dependence, error) {
	return io.LoadLocalDependencies(a.ctx)
}

// InstallDependencyModule exposes dependency installation to the Wails app layer.
func (a *App) InstallDependencyModule(dependency io.Dependence) error {
	return io.InstallDependencyModule(dependency)
}

// GetAllInstalledDependencies exposes the globally installed dependency catalog to the UI.
func (a *App) GetAllInstalledDependencies() ([]io.Dependence, error) {
	return io.GetAllInstalledDependencies()
}

// GetProjectUsedDependencies exposes project-scoped dependency resolution with the required project path argument.
func (a *App) GetProjectUsedDependencies(projectPath string) ([]io.Dependence, error) {
	return io.GetProjectUsedDependencies(projectPath)
}

// RemoveProjectDependency keeps the explicit project-path API available on the app binding.
func (a *App) RemoveProjectDependency(projectPath string, name string) error {
	return io.RemoveProjectDependency(projectPath, name)
}

// RemoveDependency is retained as a compatibility alias for existing frontend code.
func (a *App) RemoveDependency(projectPath string, name string) error {
	return a.RemoveProjectDependency(projectPath, name)
}

// AddProjectDependency exposes the dependency attachment workflow to the app binding.
func (a *App) AddProjectDependency(projectPath string, dependency io.Dependence) error {
	return io.AddProjectDependency(projectPath, dependency)
}

// AddProjectDependency exposes the dependency attachment workflow to the app binding.
func (a *App) GetDependencyManifest(dependencyName string) (*io.DependencySchema, error) {
	return io.GetDependencyManifest(dependencyName)
}

// LoadLocalPlugin opens a native directory picker and registers the selected
// folder as a local plugin. Bound to frontend as window.go.main.App.LoadLocalPlugin()
func (a *App) LoadLocalPlugin() error {
	return io.LoadLocalPlugin(a.ctx)
}

// InstallPlugin registers a non-local plugin at the given path.
// Bound to frontend as window.go.main.App.InstallPlugin(path)
func (a *App) InstallPlugin(pluginPath string) error {
	return io.InstallPlugin(pluginPath, false)
}

// DeletePlugin removes a plugin by its exact name@version identifier.
// Non-local plugins also have their folder deleted from disk.
// Bound to frontend as window.go.main.App.DeletePlugin(name)
func (a *App) DeletePlugin(pluginName string) error {
	return io.DeletePlugin(pluginName)
}

// GetPluginsList returns all registered plugin entries.
// Bound to frontend as window.go.main.App.GetPluginsList()
func (a *App) GetPluginsList() ([]io.PluginItem, error) {
	return io.GetPluginsList()
}

// GetPlugin retrieves full manifest details for a single plugin by name@version.
// Bound to frontend as window.go.main.App.GetPlugin(name)
func (a *App) GetPlugin(pluginName string) (*io.PluginDetail, error) {
	return io.GetPlugin(pluginName)
}

// --- Project Plugin Management ---

// AddPluginToProject activates an installed plugin for the given project.
func (a *App) AddPluginToProject(projectName string, pluginName string) error {
	return io.AddPluginToProject(projectName, pluginName)
}

// RemovePluginFromProject deactivates a plugin for the given project.
// The plugin remains installed globally.
func (a *App) RemovePluginFromProject(projectName string, pluginName string) error {
	return io.RemovePluginFromProject(projectName, pluginName)
}

// GetActivePlugins returns full details for all plugins activated in the given project.
func (a *App) GetActivePlugins(projectName string) ([]io.PluginItem, error) {
	return io.GetActivePlugins(projectName)
}

// --- Dev Server / Plugin Runtime ---

// StartProjectDependencies launches all active plugins for the given project.
// Emits "plugin:started" on success and "plugin:error" on failure per plugin.
func (a *App) StartProjectDependencies(projectName string) error {
	devserver.StartProjectDependencies(projectName)
	return nil
}

// StopAllPlugins kills all running plugin processes.
// Typically called from onShutdown, but exposed here for manual control.
func (a *App) StopAllPlugins() {
	devserver.StopAllPlugins()
} 
  