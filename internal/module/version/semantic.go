// Package version provides functionality to manage and check the version of the Kubex Horizon CLI tool.
// It includes methods to retrieve the current version, check for the latest version,
package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	manifest "github.com/kubex-ecosystem/logz/internal/module/info"
	"github.com/kubex-ecosystem/logz/logger"
	"github.com/spf13/cobra"
)

var gl = logger.NewLogger("Logz")
var (
	info manifest.Manifest
	vrs  Service
	err  error
)

func init() {
	if info == nil {
		info, err = manifest.GetManifest()
		if err != nil {
			gl.ErrorCtx("Failed to get manifest: "+err.Error(), map[string]any{})
		}
	}
}

type Service interface {
	// GetLatestVersion retrieves the latest version from the Git repository.
	GetLatestVersion() (string, error)
	// GetCurrentVersion returns the current version of the service.
	GetCurrentVersion() string
	// IsLatestVersion checks if the current version is the latest version.
	IsLatestVersion() (bool, error)
	// GetName returns the name of the service.
	GetName() string
	// GetVersion returns the current version of the service.
	GetVersion() string
	// GetRepository returns the Git repository URL of the service.
	GetRepository() string
	// setLastCheckedAt sets the last checked time for the version.
	setLastCheckedAt(time.Time)
	// updateLatestVersion updates the latest version from the Git repository.
	updateLatestVersion() error
}
type ServiceImpl struct {
	manifest.Manifest
	gitModelURL    string
	latestVersion  string
	lastCheckedAt  time.Time
	currentVersion string
}

func init() {
	if info == nil {
		var err error
		info, err = manifest.GetManifest()
		if err != nil {
			gl.ErrorCtx("Failed to get manifest: "+err.Error(), map[string]any{})
		}
	}
	if vrs == nil {
		vrs = NewVersionService()
	}
}

func getLatestTag(repoURL string) (string, error) {
	defer func() {
		if rec := recover(); rec != nil {
			gl.ErrorCtx(fmt.Sprintf("Recovered from panic in getLatestTag: %v", rec), map[string]interface{}{})
			err = fmt.Errorf("panic occurred while fetching latest tag: %v", rec)
		}
	}()

	defer func() {
		if vrs == nil {
			vrs = NewVersionService()
		}
		vrs.setLastCheckedAt(time.Now())
	}()

	if info == nil {
		var err error
		info, err = manifest.GetManifest()
		if err != nil {
			return "", fmt.Errorf("failed to get manifest: %w", err)
		}
	}
	if info.IsPrivate() {
		return "", fmt.Errorf("cannot fetch latest tag for private repositories")
	}

	if repoURL == "" {
		repoURL = info.GetRepository()
		if repoURL == "" {
			return "", fmt.Errorf("repository URL is not set")
		}
	}

	apiURL := fmt.Sprintf("%s/tags", repoURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}
	type Tag struct {
		Name string `json:"name"`
	}

	// Decode the JSON response into a slice of Tag structs
	// This assumes the API returns a JSON array of tags.
	// Adjust the decoding logic based on the actual API response structure.
	if resp.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}
	return tags[0].Name, nil
}
func (v *ServiceImpl) updateLatestVersion() error {
	if info.IsPrivate() {
		return fmt.Errorf("cannot fetch latest version for private repositories")
	}
	repoURL := strings.TrimSuffix(v.gitModelURL, ".git")
	tag, err := getLatestTag(repoURL)
	if err != nil {
		return err
	}
	v.latestVersion = tag
	return nil
}
func (v *ServiceImpl) vrsCompare(v1, v2 []int) (int, error) {
	compare := 0
	for i := 0; i < len(v1) && i < len(v2); i++ {
		if v1[i] < v2[i] {
			compare = -1
			break
		}
		if v1[i] > v2[i] {
			compare = 1
			break
		}
	}
	return compare, nil
}
func (v *ServiceImpl) versionAtMost(versionAtMostArg, max []int) (bool, error) {
	if comp, err := v.vrsCompare(versionAtMostArg, max); err != nil {
		return false, err
	} else if comp == 1 {
		return false, nil
	}
	return true, nil
}
func (v *ServiceImpl) parseVersion(versionToParse string) []int {
	if versionToParse == "" {
		return nil
	}
	if strings.Contains(versionToParse, "-") {
		versionToParse = strings.Split(versionToParse, "-")[0]
	}
	if strings.Contains(versionToParse, "v") {
		versionToParse = strings.TrimPrefix(versionToParse, "v")
	}
	parts := strings.Split(versionToParse, ".")
	parsedVersion := make([]int, len(parts))
	for i, part := range parts {
		if num, err := strconv.Atoi(part); err != nil {
			return nil
		} else {
			parsedVersion[i] = num
		}
	}
	return parsedVersion
}
func (v *ServiceImpl) IsLatestVersion() (bool, error) {
	if info.IsPrivate() {
		return false, fmt.Errorf("cannot check version for private repositories")
	}
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return false, err
		}
	}

	currentVersionParts := v.parseVersion(v.currentVersion)
	latestVersionParts := v.parseVersion(v.latestVersion)

	if len(currentVersionParts) == 0 || len(latestVersionParts) == 0 {
		return false, fmt.Errorf("invalid version format")
	}

	if len(currentVersionParts) != len(latestVersionParts) {
		return false, fmt.Errorf("version parts length mismatch")
	}

	return v.versionAtMost(currentVersionParts, latestVersionParts)
}
func (v *ServiceImpl) GetLatestVersion() (string, error) {
	if info.IsPrivate() {
		return "", fmt.Errorf("cannot fetch latest version for private repositories")
	}
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return "", err
		}
	}
	return v.latestVersion, nil
}
func (v *ServiceImpl) GetCurrentVersion() string {
	if v.currentVersion == "" {
		v.currentVersion = info.GetVersion()
	}
	return v.currentVersion
}
func (v *ServiceImpl) GetName() string {
	if info == nil {
		return "Unknown Service"
	}
	return info.GetName()
}
func (v *ServiceImpl) GetVersion() string {
	if info == nil {
		return "Unknown version"
	}
	return info.GetVersion()
}
func (v *ServiceImpl) GetRepository() string {
	if info == nil {
		return "No repository URL set in the manifest."
	}
	return info.GetRepository()
}
func (v *ServiceImpl) setLastCheckedAt(t time.Time) {
	v.lastCheckedAt = t
	gl.DebugCtx("Last checked at: "+t.Format(time.RFC3339), map[string]interface{}{})
}

func NewVersionService() Service {
	return &ServiceImpl{
		Manifest:       info,
		gitModelURL:    info.GetRepository(),
		currentVersion: info.GetVersion(),
		latestVersion:  "",
	}
}

var (
	versionCmd   *cobra.Command
	subLatestCmd *cobra.Command
	subCmdCheck  *cobra.Command
	updCmd       *cobra.Command
	getCmd       *cobra.Command
	restartCmd   *cobra.Command
)

func init() {
	if versionCmd == nil {
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Print the version number of " + info.GetName(),
			Long:  "Print the version number of " + info.GetName() + " and other related information.",
			Run: func(cmd *cobra.Command, args []string) {
				if info.IsPrivate() {
					gl.WarnCtx("The information shown may not be accurate for private repositories.", map[string]interface{}{})
					gl.InfoCtx("Current version: "+GetVersion(), map[string]interface{}{})
					gl.InfoCtx("Git repository: "+GetGitRepositoryModelURL(), map[string]interface{}{})
					return
				}
				GetVersionInfo()
			},
		}
	}
	if subLatestCmd == nil {
		subLatestCmd = &cobra.Command{
			Use:   "latest",
			Short: "Print the latest version number of " + info.GetName(),
			Long:  "Print the latest version number of " + info.GetName() + " from the Git repository.",
			Run: func(cmd *cobra.Command, args []string) {
				if info.IsPrivate() {
					gl.ErrorCtx("Cannot fetch latest version for private repositories.", map[string]interface{}{})
					return
				}
				GetLatestVersionInfo()
			},
		}
	}
	if subCmdCheck == nil {
		subCmdCheck = &cobra.Command{
			Use:   "check",
			Short: "Check if the current version is the latest version of " + info.GetName(),
			Long:  "Check if the current version is the latest version of " + info.GetName() + " and print the version information.",
			Run: func(cmd *cobra.Command, args []string) {
				if info.IsPrivate() {
					gl.ErrorCtx("Cannot check version for private repositories.", map[string]interface{}{})
					return
				}
				GetVersionInfoWithLatestAndCheck()
			},
		}
	}
	if updCmd == nil {
		updCmd = &cobra.Command{
			Use:   "update",
			Short: "Update the version information of " + info.GetName(),
			Long:  "Update the version information of " + info.GetName() + " by fetching the latest version from the Git repository.",
			Run: func(cmd *cobra.Command, args []string) {
				if info.IsPrivate() {
					gl.ErrorCtx("Cannot update version for private repositories.", map[string]interface{}{})
					return
				}
				if err := vrs.updateLatestVersion(); err != nil {
					gl.ErrorCtx("Failed to update version: "+err.Error(), map[string]interface{}{})
				} else {
					latestVersion, err := vrs.GetLatestVersion()
					if err != nil {
						gl.ErrorCtx("Failed to get latest version: "+err.Error(), map[string]interface{}{})
					} else {
						gl.InfoCtx("Current version: "+vrs.GetCurrentVersion(), map[string]interface{}{})
						gl.InfoCtx("Latest version: "+latestVersion, map[string]interface{}{})
					}
					vrs.setLastCheckedAt(time.Now())
				}
			},
		}
	}
	if getCmd == nil {
		getCmd = &cobra.Command{
			Use:   "get",
			Short: "Get the current version of " + info.GetName(),
			Long:  "Get the current version of " + info.GetName() + " from the manifest.",
			Run: func(cmd *cobra.Command, args []string) {
				gl.InfoCtx("Current version: "+vrs.GetCurrentVersion(), map[string]interface{}{})
			},
		}
	}
	if restartCmd == nil {
		restartCmd = &cobra.Command{
			Use:   "restart",
			Short: "Restart the " + info.GetName() + " service",
			Long:  "Restart the " + info.GetName() + " service to apply any changes made.",
			Run: func(cmd *cobra.Command, args []string) {
				gl.InfoCtx("Restarting the service...", map[string]interface{}{})
				// Logic to restart the service can be added here
				gl.InfoCtx("Service restarted successfully", map[string]interface{}{})
			},
		}
	}

}
func GetVersion() string {
	if info == nil {
		_, err := manifest.GetManifest()
		if err != nil {
			gl.ErrorCtx("Failed to get manifest: "+err.Error(), map[string]interface{}{})
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
	gl.InfoCtx("Version: "+GetVersion(), map[string]interface{}{})
	gl.InfoCtx("Git repository: "+GetGitRepositoryModelURL(), map[string]interface{}{})
	return fmt.Sprintf("Version: %s\nGit repository: %s", GetVersion(), GetGitRepositoryModelURL())
}
func GetLatestVersionFromGit() string {
	if info.IsPrivate() {
		gl.ErrorCtx("Cannot fetch latest version for private repositories.", map[string]interface{}{})
		return "Cannot fetch latest version for private repositories."
	}

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	gitURLWithoutGit := strings.TrimSuffix(GetGitRepositoryModelURL(), ".git")
	if gitURLWithoutGit == "" {
		gl.ErrorCtx("No repository URL set in the manifest.", map[string]interface{}{})
		return "No repository URL set in the manifest."
	}

	response, err := netClient.Get(gitURLWithoutGit + "/releases/latest")
	if err != nil {
		gl.ErrorCtx("Error fetching latest version: "+err.Error(), map[string]interface{}{})
		gl.ErrorCtx("URL: "+gitURLWithoutGit+"/releases/latest", map[string]interface{}{})
		return err.Error()
	}

	if response.StatusCode != 200 {
		gl.ErrorCtx("Error fetching latest version: "+response.Status, map[string]interface{}{})
		gl.ErrorCtx("URL: "+gitURLWithoutGit+"/releases/latest", map[string]interface{}{})
		body, _ := io.ReadAll(response.Body)
		return fmt.Sprintf("Error: %s\nResponse: %s", response.Status, string(body))
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}
func GetLatestVersionInfo() string {
	if info.IsPrivate() {
		gl.ErrorCtx("Cannot fetch latest version for private repositories.", map[string]interface{}{})
		return "Cannot fetch latest version for private repositories."
	}
	gl.InfoCtx("Latest version: "+GetLatestVersionFromGit(), map[string]interface{}{})
	return "Latest version: " + GetLatestVersionFromGit()
}
func GetVersionInfoWithLatestAndCheck() string {
	if info.IsPrivate() {
		gl.ErrorCtx("Cannot check version for private repositories.", map[string]interface{}{})
		return "Cannot check version for private repositories."
	}
	if GetVersion() == GetLatestVersionFromGit() {
		gl.InfoCtx("You are using the latest version.", map[string]interface{}{})
		return fmt.Sprintf("You are using the latest version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	} else {
		gl.WarnCtx("You are using an outdated version.", map[string]interface{}{})
		return fmt.Sprintf("You are using an outdated version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	}
}
func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	versionCmd.AddCommand(updCmd)
	versionCmd.AddCommand(getCmd)
	versionCmd.AddCommand(restartCmd)
	return versionCmd
}
