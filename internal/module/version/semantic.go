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

	gl "github.com/kubex-ecosystem/logz"
)

var info manifest.Manifest
var vrs Service

func GetVersionService() Service {
	if vrs == nil {
		vrs = NewVersionService()
	}
	return vrs
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
			fmt.Println("Failed to get manifest: " + err.Error())
		}
	}
	if vrs == nil {
		vrs = NewVersionService()
	}
}

func getLatestTag(repoURL string) (string, error) {
	defer func() {
		if rec := recover(); rec != nil {
			gl.Errorf("panic recovered in getLatestTag: %v", rec)
			return
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
	fmt.Println("Last checked at: " + t.Format(time.RFC3339))
}
