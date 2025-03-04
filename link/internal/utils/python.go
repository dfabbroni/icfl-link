package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type PythonEnv struct {
	VenvPath     string
	BinPath      string
	Python       string
	Pip          string
	mu           sync.Mutex
	SuperLinkCmd *exec.Cmd
	FlwrExecCmd *exec.Cmd
}

var (
	sharedEnv *PythonEnv
	once      sync.Once
)

// InitializeExperimentEnvironment sets up the base directory and Python environment
func InitializeExperimentEnvironment(baseDir string) error {
	// Ensure we're on Linux
	if runtime.GOOS != "linux" {
		return fmt.Errorf("this feature is only supported on Linux systems")
	}

	absPath, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %v", err)
	}

	// Create and initialize the shared venv
	_, err = GetSharedPythonEnv(absPath)
	return err
}

// GetSharedPythonEnv returns a singleton PythonEnv instance
func GetSharedPythonEnv(baseDir string) (*PythonEnv, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("this feature is only supported on Linux systems")
	}

	var initErr error
	once.Do(func() {
		venvPath := filepath.Join(baseDir, "flower")

		sharedEnv = &PythonEnv{
			VenvPath: venvPath,
			BinPath:  filepath.Join(venvPath, "bin"),
			Python:   filepath.Join(venvPath, "bin", "python"),
			Pip:      filepath.Join(venvPath, "bin", "pip"),
		}

		// Create virtual environment if it doesn't exist
		if _, err := os.Stat(venvPath); os.IsNotExist(err) {
			cmd := exec.Command("python3", "-m", "venv", venvPath)
			if err := cmd.Run(); err != nil {
				initErr = fmt.Errorf("failed to create shared virtual environment: %v", err)
				return
			}
		}

		if err := sharedEnv.InstallFlwr(); err != nil {
			initErr = fmt.Errorf("failed to install Flower: %v", err)
			return
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return sharedEnv, nil
}

// InstallExperimentDependencies installs packages in the shared environment using the experiments pyproject.toml
func (env *PythonEnv) InstallExperimentDependencies(experimentDir string) error {
	env.mu.Lock()
	defer env.mu.Unlock()

	// Check if pyproject.toml exists
	if _, err := os.Stat(filepath.Join(experimentDir, "pyproject.toml")); os.IsNotExist(err) {
		return fmt.Errorf("pyproject.toml not found in %s", experimentDir)
	}

	cmd := exec.Command(env.Pip, "install", "-e", ".")
	cmd.Dir = experimentDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", env.VenvPath),
		fmt.Sprintf("PATH=%s%c%s", env.BinPath, os.PathListSeparator, os.Getenv("PATH")),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip install failed: %v\nOutput: %s", err, output)
	}

	return nil
}

func (env *PythonEnv) InstallFlwr() error {
	cmd := exec.Command(env.Pip, "install", "flwr==1.15.0")
	return cmd.Run()
}

// RunFlwr executes the experiment using the flwr command
func (env *PythonEnv) RunFlwr(experimentDir, experimentID, experimentName string) error {
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	logFileName := fmt.Sprintf("flwr_%s.log", timestamp)
	logLocation := filepath.Join("uploads", experimentID, "logs")
	logFile := filepath.Join(logLocation, logFileName)

	if err := os.MkdirAll(logLocation, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	flwrLogFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open flwr log file: %v", err)
	}
	defer flwrLogFile.Close()

	// Run the experiment using flwr
	cmd := exec.Command(filepath.Join(env.BinPath, "flwr"), "run", ".", experimentName, "--stream")
	cmd.Dir = experimentDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", env.VenvPath),
		fmt.Sprintf("PATH=%s%c%s", env.BinPath, os.PathListSeparator, os.Getenv("PATH")),
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = flwrLogFile
	cmd.Stderr = flwrLogFile

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start flwr: %v", err)
	}
	env.FlwrExecCmd = cmd
	log.Printf("Started flwr with PID: %d", cmd.Process.Pid)

	return nil
}

// InitializeSuperLink starts the SuperLink process with SSL and authentication
func (env *PythonEnv) InitializeSuperLink() error {
	timestamp := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	logFileName := fmt.Sprintf("superlink_%s.log", timestamp)
	logFile := filepath.Join("logs", logFileName)

	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	superLinkLogFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open SuperLink log file: %v", err)
	}
	defer superLinkLogFile.Close()

	// Start SuperLink with SSL and authentication
	superLinkCmd := exec.Command(filepath.Join(env.BinPath, "flower-superlink"),
		"--ssl-ca-certfile", "authentication/certificates/ca.crt",
		"--ssl-certfile", "authentication/certificates/server.pem",
		"--ssl-keyfile", "authentication/certificates/server.key",
		"--auth-list-public-keys", "authentication/keys/client_public_keys.csv")
	superLinkCmd.Env = append(os.Environ(),
		fmt.Sprintf("VIRTUAL_ENV=%s", env.VenvPath),
		fmt.Sprintf("PATH=%s%c%s", env.BinPath, os.PathListSeparator, os.Getenv("PATH")),
	)
	superLinkCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	superLinkCmd.Stdout = superLinkLogFile
	superLinkCmd.Stderr = superLinkLogFile

	if err := superLinkCmd.Start(); err != nil {
		return fmt.Errorf("failed to start SuperLink: %v", err)
	}
	env.SuperLinkCmd = superLinkCmd
	log.Printf("Started SuperLink with PID: %d", superLinkCmd.Process.Pid)

	return nil
}

func (env *PythonEnv) CleanupSuperLink() error {
	if env.SuperLinkCmd != nil && env.SuperLinkCmd.Process != nil {
		if pgid, err := syscall.Getpgid(env.SuperLinkCmd.Process.Pid); err == nil {
			if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
				return fmt.Errorf("failed to stop SuperLink: %v", err)
			}
		}
		env.SuperLinkCmd = nil
	}
	return nil
}

func (env *PythonEnv) CleanupFlwr() error {
	if env.FlwrExecCmd != nil && env.FlwrExecCmd.Process != nil {
		if pgid, err := syscall.Getpgid(env.FlwrExecCmd.Process.Pid); err == nil {
			if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
				return fmt.Errorf("failed to stop SuperExec: %v", err)
			}
		}
		env.FlwrExecCmd = nil
	}
	return nil
}
