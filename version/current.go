package version

import (
	_ "embed"
	"fmt"
	"github.com/faelmori/logz/logger"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"strings"
	"time"
)

const moduleAlias = "Logz"
const moduleName = "logz"
const gitModelUrl = "https://github.com/faelmori/" + moduleName + ".git"
const currentVersionFallback = "v1.1.8" // First version with the version file

var (
	l          = logger.NewLogger("Logz")
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of " + moduleAlias,
		Long:  "Print the version number of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetVersionInfo()
		},
	}
	subLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Print the latest version number of " + moduleAlias,
		Long:  "Print the latest version number of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetLatestVersionInfo()
		},
	}
	subCmdCheck = &cobra.Command{
		Use:   "check",
		Short: "Check if the current version is the latest version of " + moduleAlias,
		Long:  "Check if the current version is the latest version of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetVersionInfoWithLatestAndCheck()
		},
	}
)

//go:embed CLI_VERSION
var currentVersion string

func GetVersion() string {
	if currentVersion == "" {
		return currentVersionFallback
	}
	return currentVersion
}

func GetGitModelUrl() string {
	return gitModelUrl
}

func GetVersionInfo() string {
	l.Info("Version: "+GetVersion(), map[string]interface{}{})
	l.Info("Git repository: "+GetGitModelUrl(), map[string]interface{}{})
	return fmt.Sprintf("Version: %s\nGit repository: %s", GetVersion(), GetGitModelUrl())
}

func GetLatestVersionFromGit() string {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	gitUrlWithoutGit := strings.TrimSuffix(gitModelUrl, ".git")

	response, err := netClient.Get(gitUrlWithoutGit + "/releases/latest")
	if err != nil {
		l.Error("Error fetching latest version: "+err.Error(), map[string]interface{}{})
		l.Error("Url: "+gitUrlWithoutGit+"/releases/latest", map[string]interface{}{})
		return err.Error()
	}

	if response.StatusCode != 200 {
		l.Error("Error fetching latest version: "+response.Status, map[string]interface{}{})
		l.Error("Url: "+gitUrlWithoutGit+"/releases/latest", map[string]interface{}{})
		body, _ := io.ReadAll(response.Body)
		return fmt.Sprintf("Error: %s\nResponse: %s", response.Status, string(body))
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}

func GetLatestVersionInfo() string {
	l.Info("Latest version: "+GetLatestVersionFromGit(), map[string]interface{}{})
	return "Latest version: " + GetLatestVersionFromGit()
}

func GetVersionInfoWithLatestAndCheck() string {
	if GetVersion() == GetLatestVersionFromGit() {
		l.Info("You are using the latest version.", map[string]interface{}{})
		return fmt.Sprintf("You are using the latest version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	} else {
		l.Warn("You are using an outdated version.", map[string]interface{}{})
		return fmt.Sprintf("You are using an outdated version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	}
}

func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	return versionCmd
}
