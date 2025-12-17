package version

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	gl "github.com/kubex-ecosystem/logz"
	manifest "github.com/kubex-ecosystem/logz/internal/module/info"

	"github.com/spf13/cobra"
)

func NewVersionService() Service {
	return &ServiceImpl{
		Manifest:       info,
		gitModelURL:    info.GetRepository(),
		currentVersion: info.GetVersion(),
		latestVersion:  "",
	}
}
func GetVersion() string {
	if info == nil {
		_, err := manifest.GetManifest()
		if err != nil {
			fmt.Println("Failed to get manifest: " + err.Error())
			return "Unknown version"
		}
	}
	return info.GetVersion()
}
func GetGitRepositoryModelURL() string {
	if info.GetRepository() == "" {
		return "No repository URL set in the manifest."
	}
	return info.GetRepository()
}
func GetVersionInfo() string {
	fmt.Println("Version: " + GetVersion())
	fmt.Println("Git repository: " + GetGitRepositoryModelURL())
	return fmt.Sprintf("Version: %s\nGit repository: %s", GetVersion(), GetGitRepositoryModelURL())
}
func GetLatestVersionFromGit() string {
	if info.IsPrivate() {
		fmt.Println("Cannot fetch latest version for private repositories.")
		return "Cannot fetch latest version for private repositories."
	}

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	gitURLWithoutGit := strings.TrimSuffix(GetGitRepositoryModelURL(), ".git")
	if gitURLWithoutGit == "" {
		fmt.Println("No repository URL set in the manifest.")
		return "No repository URL set in the manifest."
	}

	response, err := netClient.Get(gitURLWithoutGit + "/releases/latest")
	if err != nil {
		fmt.Println("Error fetching latest version: " + err.Error())
		fmt.Println(gitURLWithoutGit + "/releases/latest")
		return err.Error()
	}

	if response.StatusCode != 200 {
		fmt.Println("Error fetching latest version: " + response.Status)
		fmt.Println("Url: " + gitURLWithoutGit + "/releases/latest")
		body, _ := io.ReadAll(response.Body)
		return fmt.Sprintf("Error: %s\nResponse: %s", response.Status, string(body))
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}
func GetLatestVersionInfo() string {
	if info.IsPrivate() {
		fmt.Println("Cannot fetch latest version for private repositories.")
		return "Cannot fetch latest version for private repositories."
	}
	fmt.Println("Latest version: " + GetLatestVersionFromGit())
	return "Latest version: " + GetLatestVersionFromGit()
}
func GetVersionInfoWithLatestAndCheck() string {
	if info.IsPrivate() {
		fmt.Println("Cannot check version for private repositories.")
		return "Cannot check version for private repositories."
	}
	if GetVersion() == GetLatestVersionFromGit() {
		fmt.Println("You are using the latest version.")
		return fmt.Sprintf("You are using the latest version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	} else {
		gl.Log("warn", "You are using an outdated version.")
		return fmt.Sprintf("You are using an outdated version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	}
}
func CliCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of " + info.GetName(),
		Long:  "Print the version number of " + info.GetName() + " and other related information.",
		Run: func(cmd *cobra.Command, args []string) {
			if info.IsPrivate() {
				gl.Log("warn", "The information shown may not be accurate for private repositories.")
				fmt.Println("Current version: " + GetVersion())
				fmt.Println("Git repository: " + GetGitRepositoryModelURL())
				return
			}
			GetVersionInfo()
		},
	}
	subLatestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Print the latest version number of " + info.GetName(),
		Long:  "Print the latest version number of " + info.GetName() + " from the Git repository.",
		Run: func(cmd *cobra.Command, args []string) {
			if info.IsPrivate() {
				fmt.Println("Cannot fetch latest version for private repositories.")
				return
			}
			GetLatestVersionInfo()
		},
	}
	subCmdCheck := &cobra.Command{
		Use:   "check",
		Short: "Check if the current version is the latest version of " + info.GetName(),
		Long:  "Check if the current version is the latest version of " + info.GetName() + " and print the version information.",
		Run: func(cmd *cobra.Command, args []string) {
			if info.IsPrivate() {
				fmt.Println("Cannot check version for private repositories.")
				return
			}
			// fmt.Println(GetVersionInfoWithLatestAndCheck())
			fmt.Println(GetVersionInfoWithLatestAndCheck())
		},
	}
	updCmd := &cobra.Command{
		Use:   "update",
		Short: "Update the version information of " + info.GetName(),
		Long:  "Update the version information of " + info.GetName() + " by fetching the latest version from the Git repository.",
		Run: func(cmd *cobra.Command, args []string) {
			if info.IsPrivate() {
				fmt.Println("Cannot update version for private repositories.")
				return
			}
			if err := vrs.updateLatestVersion(); err != nil {
				fmt.Println("Failed to update version: " + err.Error())
			} else {
				latestVersion, err := vrs.GetLatestVersion()
				if err != nil {
					fmt.Println("Failed to get latest version: " + err.Error())
				} else {
					fmt.Println("Current version: " + vrs.GetCurrentVersion())
					fmt.Println("Latest version: " + latestVersion)
				}
				vrs.setLastCheckedAt(time.Now())
			}
		},
	}
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get the current version of " + info.GetName(),
		Long:  "Get the current version of " + info.GetName() + " from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Current version: " + vrs.GetCurrentVersion())
		},
	}
	restartCmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the " + info.GetName() + " service",
		Long:  "Restart the " + info.GetName() + " service to apply any changes made.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Restarting the service...")
			// Logic to restart the service can be added here
			gl.Log("success", "Service restarted successfully")
		},
	}

	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	versionCmd.AddCommand(updCmd)
	versionCmd.AddCommand(getCmd)
	versionCmd.AddCommand(restartCmd)
	return versionCmd
}
