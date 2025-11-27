// Package module provides internal types and functions for the GoBE application.
package module

import (
	"fmt"

	"github.com/kubex-ecosystem/logz/cmd/cli"
	"github.com/kubex-ecosystem/logz/internal/module/info"
	"github.com/kubex-ecosystem/logz/internal/module/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type glgr interface {
	Log(level string, parts ...any)
}

var gl glgr

func SetLogger(logger glgr) {
	gl = logger
}

type LogZ struct {
	parentCmdName string
	hideBanner    bool
	certPath      string
	keyPath       string
	configPath    string
}

func (m *LogZ) Alias() string {
	return ""
}
func (m *LogZ) ShortDescription() string {
	return "LogZ: A Kubex Ecosystem Logging Tool"
}
func (m *LogZ) LongDescription() string {
	return `LogZ: A Kubex Ecosystem Logging Tool for managing logs and events, providing insights into system performance and security.`
}
func (m *LogZ) Usage() string {
	return "logz [command] [args]"
}
func (m *LogZ) Examples() []string {
	return []string{"logz -l info -m file -o /path/to/logfile.log",
		"logz -l debug -m json -o /path/to/logfile.json",
		"logz -l warn -m stdout",
		"logz -l error -m stderr",
		"logz -l fatal -m file -o /path/to/logfile.log",
	}
}
func (m *LogZ) Active() bool {
	return true
}
func (m *LogZ) Module() string {
	return "logz"
}
func (m *LogZ) Execute() error {
	dbChanData := make(chan interface{})
	defer close(dbChanData)

	if spyderErr := m.Command().Execute(); spyderErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", spyderErr)
		return spyderErr
	} else {
		return nil
	}
}
func (m *LogZ) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: m.Module(),
		//Aliases:     []string{m.Alias(), "w", "wb", "webServer", "http"},
		Example: m.concatenateExamples(),
		Annotations: m.GetDescriptions(
			[]string{
				m.LongDescription(),
				m.ShortDescription(),
			}, m.hideBanner,
		),
		Version: version.GetVersion(),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(version.CliCommand())
	cmd.AddCommand(cli.LoggerCmd())
	cmd.AddCommand(cli.LogzCmd())
	// cmd.AddCommand(cli.MetricsCmd())

	setUsageDefinition(cmd)
	for _, c := range cmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return cmd
}

func (m *LogZ) GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	return info.GetDescriptions(descriptionArg, (m.hideBanner || hideBanner))
}
func (m *LogZ) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *LogZ) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}
