package devserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	projectIO "lazy-backend/io"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var PORT = 8003
var server *http.Server // Keep a reference so we can shut it down
var route = http.NewServeMux()

var RunningPlugins = make([]struct {
	Name string
	Port string
}, 0)

func GetDevServerPort() int {

	return PORT
}
func portIsAvailable(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	ln.Close()
	return nil
}

func StartDevServer() {
	Register(route)
	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Dev server running on port %d", PORT)
	})

	getPort := func() (string, error) {
		portStr := fmt.Sprintf(":%d", PORT)
		if err := portIsAvailable(portStr); err == nil {
			return portStr, nil
		}

		for i := 1; i <= 100; i++ {
			tryPort := PORT + i
			tryPortStr := fmt.Sprintf(":%d", tryPort)
			if err := portIsAvailable(tryPortStr); err == nil {
				PORT = tryPort
				fmt.Printf("Port %d was taken, using %d instead\n", PORT-i, tryPort)
				return tryPortStr, nil
			}
		}

		return "", fmt.Errorf("failed to find available port after 100 attempts (starting from %d)", PORT)
	}

	finalPort, err := getPort()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Starting dev server on http://localhost%s\n", finalPort)

	// Create the server with proper configuration
	server = &http.Server{
		Addr:         finalPort,
		Handler:      route,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Handle graceful shutdown on Ctrl+C
	go handleShutdownSignal()

	// Start the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
	}
}

// shutdownServer gracefully stops the server
func ShutdownServer() {
	fmt.Println("\nShutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	} else {
		fmt.Println(" dev Server stopped successfully")
	}
}

// handleShutdownSignal listens for Ctrl+C or system signals
func handleShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Wait for signal
	ShutdownServer()
}

const registrationTimeout = 30 * time.Second

// PluginRegistration is the JSON body plugins POST to /register-plugin
type PluginRegistration struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

// RunningPlugin tracks a started plugin's process and metadata
type RunningPlugin struct {
	Name   string
	Port   string
	Cmd    *exec.Cmd
	Cancel context.CancelFunc
}

var (
	runningPlugins = make(map[string]*RunningPlugin)
	pluginsMu      sync.RWMutex
	wailsCtx       context.Context
)

// SetWailsContext must be called once during app startup so events can be emitted.
func SetWailsContext(ctx context.Context) {
	wailsCtx = ctx
}

// Register adds the /register-plugin endpoint to the existing route mux.
// Call this once during init or before StartDevServer.
func Register(route *http.ServeMux) {
	route.HandleFunc("/register-plugin", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reg PluginRegistration
		if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if reg.Name == "" || reg.Port == "" {
			http.Error(w, "name and port are required", http.StatusBadRequest)
			return
		}

		pluginsMu.Lock()
		if rp, exists := runningPlugins[reg.Name]; exists {
			rp.Port = reg.Port
			p, err := projectIO.GetPlugin(reg.Name)
			if err != nil {
				fmt.Println("error getting plugin", p)
			}
			fmt.Println(p.Manifest, "plugin-----", reg.Name)
			runtime.EventsEmit(wailsCtx, "pluginRegistered", map[string]any{
				"port":   reg.Port,
				"plugin": p.Manifest,
			})
		}
		pluginsMu.Unlock()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "registered")
		fmt.Printf("[devserver] Plugin '%s' registered on port %s\n", reg.Name, reg.Port)
	})
}

// StartProjectDependencies launches all active plugins for a project,
// waits for each to register within the timeout, and emits Wails events.
func StartProjectDependencies(project string) {
	activePlugins, err := projectIO.GetActivePlugins(project)
	if err != nil {
		emitPluginError("", fmt.Sprintf("failed to load active plugins: %v", err))
		return
	}

	if len(activePlugins) == 0 {
		fmt.Printf("[devserver] No active plugins for project '%s'\n", project)
		return
	}

	var wg sync.WaitGroup
	for _, plugin := range activePlugins {
		wg.Add(1)
		go func(p projectIO.PluginItem) {
			defer wg.Done()
			startPlugin(p)
		}(plugin)
	}
	wg.Wait()
}

func startPlugin(plugin projectIO.PluginItem) {
	if plugin.Path == "" {
		emitPluginError(plugin.Name, "plugin has no path")
		return
	}

	startCmd := plugin.Run

	if startCmd == "" {
		emitPluginError(plugin.Name, "plugin has no start command in manifest")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Try to parse the command
	parts := strings.Fields(startCmd)
	var cmd *exec.Cmd

	if len(parts) > 0 {
		// Check if the executable exists in PATH
		if path, err := exec.LookPath(parts[0]); err == nil {
			// Direct execution (faster)
			if len(parts) > 1 {
				cmd = exec.CommandContext(ctx, path, parts[1:]...)
			} else {
				cmd = exec.CommandContext(ctx, path)
			}
		} else {
			// Fallback to shell
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/sh"
			}
			cmd = exec.CommandContext(ctx, shell, "-c", startCmd)
		}
	} else {
		// Empty command, fallback to shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.CommandContext(ctx, shell, "-c", startCmd)
	}

	fmt.Println(plugin.Path)
	cmd.Dir = plugin.Path

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("DEV_SERVER_PORT=%d", GetDevServerPort()),
		fmt.Sprintf("PLUGIN_NAME=%s", plugin.Name),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		emitPluginError(plugin.Name, fmt.Sprintf("failed to start: %v", err))
		return
	}

	pluginsMu.Lock()
	runningPlugins[plugin.Name] = &RunningPlugin{
		Name:   plugin.Name,
		Cmd:    cmd,
		Cancel: cancel,
	}
	pluginsMu.Unlock()

	fmt.Printf("[devserver] Started '%s' (pid %d), waiting for registration...\n", plugin.Name, cmd.Process.Pid)

	if !waitForRegistration(plugin.Name, registrationTimeout) {
		cancel()
		pluginsMu.Lock()
		delete(runningPlugins, plugin.Name)
		pluginsMu.Unlock()
		emitPluginError(plugin.Name, fmt.Sprintf("failed to register within %s", registrationTimeout))
		return
	}

	emitPluginSuccess(plugin.Name)
}

// waitForRegistration polls until the plugin has a non-empty Port or timeout expires.
func waitForRegistration(name string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		pluginsMu.RLock()
		rp, exists := runningPlugins[name]
		port := ""
		if exists {
			port = rp.Port
		}
		pluginsMu.RUnlock()

		if port != "" {
			return true
		}
		time.Sleep(250 * time.Millisecond)
	}
	return false
}

// StopAllPlugins kills all running plugin processes.
// Call this from your existing ShutdownServer or onShutdown handler.
func StopAllPlugins() {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	for name, rp := range runningPlugins {
		rp.Cancel()
		delete(runningPlugins, name)
		fmt.Printf("[devserver] Stopped plugin '%s'\n\n", name)
	}
}

func emitPluginSuccess(pluginName string) {
	if wailsCtx == nil {
		return
	}

	runtime.EventsEmit(wailsCtx, "plugin:started", pluginName)
	fmt.Printf("[devserver] ✅ Plugin '%s' ready\n", pluginName)
}

func emitPluginError(pluginName string, reason string) {
	if wailsCtx == nil {
		return
	}
	payload := map[string]string{
		"name":   pluginName,
		"reason": reason,
	}
	runtime.EventsEmit(wailsCtx, "plugin:error", payload)
	fmt.Printf("[devserver] ❌ Plugin '%s': %s\n", pluginName, reason)
}
