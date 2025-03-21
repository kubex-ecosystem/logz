package logger

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CheckLogSize checks and manages the size of the logs
func CheckLogSize(config Config) error {
	logDir := config.Output()
	files, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Println("Error reading the log directory:", err)
		return err
	}

	// Fetch maximum log sizes from the configuration
	maxLogSize := config.GetInt("maxLogSize", 20*1024*1024)      // Default 20 MB
	moduleLogSize := config.GetInt("moduleLogSize", 5*1024*1024) // Default 5 MB

	var totalSize int64
	filesToRotate := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			fileInfo, err := file.Info()
			if err != nil {
				fmt.Println("Error getting file information:", err)
				continue
			}
			totalSize += fileInfo.Size()
			if fileInfo.Size() > int64(moduleLogSize) {
				filesToRotate = append(filesToRotate, filepath.Join(logDir, file.Name()))
			}
		}
	}

	// Rotation based on total size
	if totalSize > int64(maxLogSize) {
		fmt.Println("Archiving logs due to excessive size...")
		if err := ArchiveLogs(filesToRotate); err != nil {
			fmt.Println("Error archiving logs:", err)
			return err
		}
	}

	// Individual rotation of large files
	if len(filesToRotate) > 0 {
		fmt.Println("Rotating large log files...")
		if err := RotateLogFiles(filesToRotate); err != nil {
			fmt.Println("Error rotating log files:", err)
			return err
		}
	}

	return nil
}

// RotateLogFiles compresses and recreates the log files
func RotateLogFiles(files []string) error {
	for _, logFile := range files {
		if err := RotateLogFile(logFile); err != nil {
			fmt.Println("Error rotating the log file:", err)
			continue
		}
		fmt.Println("Log file rotated successfully:", logFile)
	}
	return nil
}

// RotateLogFile compresses a single log file
func RotateLogFile(logFilePath string) error {
	archivePath := fmt.Sprintf("%s.tar.gz", logFilePath)
	if err := CreateTarGz(archivePath, []string{logFilePath}); err != nil {
		return err
	}

	if err := os.Remove(logFilePath); err != nil {
		return fmt.Errorf("error removing the log file: %v", err)
	}

	if err := os.WriteFile(logFilePath, []byte{}, 0644); err != nil {
		return fmt.Errorf("error recreating the log file: %v", err)
	}
	return nil
}

// CreateTarGz creates a tar.gz file from the logs
func CreateTarGz(archivePath string, files []string) error {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error creating the tar.gz file: %v", err)
	}
	defer archiveFile.Close()

	gw := gzip.NewWriter(archiveFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		if err := addFileToTar(tw, file); err != nil {
			return err
		}
	}
	fmt.Println("Logs archived successfully:", archivePath)
	return nil
}

// addFileToTar adds a file to the tar archive
func addFileToTar(tw *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening the file: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file information: %v", err)
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("error creating the tar header: %v", err)
	}
	header.Name = filepath.Base(filePath)

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("error writing the tar header: %v", err)
	}

	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("error copying the file content to the tar: %v", err)
	}

	return nil
}

// ArchiveLogs archives old logs into a zip file
func ArchiveLogs(files []string) error {
	logDir := GetLogPath()
	if len(files) == 0 {
		err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".log") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error listing the log files: %v", err)
		}
	}
	tempDir := os.TempDir()
	archiveName := fmt.Sprintf("logs_archive_%s.zip", time.Now().Format("20060102_150405"))
	archivePath := filepath.Join(tempDir, archiveName)

	zipFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error creating the zip file: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}

	fmt.Println("Logs archived successfully:", archivePath)
	return nil
}

// addFileToZip adds a file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening the file for zip: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file information: %v", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("error creating the zip header: %v", err)
	}
	header.Name = filepath.Base(filePath)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("error creating the file in the zip: %v", err)
	}

	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("error copying the file content to the zip: %v", err)
	}

	return nil
}

func GetLogDirectorySize(directory string) (int64, error) {
	if directory == "" {
		directory = filepath.Dir(GetLogPath())
	}
	var totalSize int64

	// Traverse the specified directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing the path %s: %v", path, err)
		}

		// Only files are considered in the total size
		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("error calculating the directory size: %v", err)
	}

	return totalSize, nil
}
