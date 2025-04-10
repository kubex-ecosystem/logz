package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const BinaryName = "logz"

func installLogz() {
	err := installBinary(BinaryName)
	if err != nil {
		fmt.Printf("ErrorCtx during installation: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Installation complete!")
}

func installBinary(binaryName string) error {
	// Obter o diretório atual e o caminho do binário
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	targetDir := getInstallDirectory()
	targetPath := filepath.Join(targetDir, binaryName)
	if fileExists(targetPath) {
		fmt.Printf("✅ Binary already exists at: %s\n", targetPath)
	} else {
		// Copiar o binário
		fmt.Printf("Installing binary to: %s\n", targetPath)
		err := copyFile(binaryPath, targetPath)
		if err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}
		fmt.Printf("✅ Binary installed at: %s\n", targetPath)
	}
	if !isInPath(targetDir) {
		fmt.Printf("⚠️  Warning: %s is not in your PATH.\n", targetDir)
		err := addToPathInstruction(targetDir)
		if err != nil {
			return fmt.Errorf("failed to provide PATH instructions: %w", err)
		}
	}
	return nil
}

func getInstallDirectory() string {
	if os.Geteuid() == 0 {
		return "/usr/local/bin"
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "bin")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isInPath(dir string) bool {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")
	for _, p := range paths {
		if p == dir {
			return true
		}
	}
	return false
}

func addToPathInstruction(dir string) error {
	shell := filepath.Base(os.Getenv("SHELL"))
	var shellConfigFile string

	switch shell {
	case "bash":
		shellConfigFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
	case "zsh":
		shellConfigFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
	case "sh":
		shellConfigFile = filepath.Join(os.Getenv("HOME"), ".profile")
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	fmt.Printf("Add the following line to your %s file to include %s in your PATH:\n", shellConfigFile, dir)
	fmt.Printf("export PATH=%s:$PATH\n", dir)
	fmt.Printf("Then run: source %s\n", shellConfigFile)
	return nil
}

func copyFile(src, dest string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(input *os.File) {
		_ = input.Close()
	}(input)
	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(output *os.File) {
		_ = output.Close()
	}(output)
	_, err = io.Copy(output, input)
	if err != nil {
		return err
	}
	// Ajustar permissões do binário
	return os.Chmod(dest, 0755)
}
