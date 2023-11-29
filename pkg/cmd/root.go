/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/spf13/cobra"
)

type RootCmdValues struct {
	// The file to which logs are being directed, if any. "" if not.
	logFileName string

	// We don't trace anything until this flag is true.
	// This means that any errors which occur in the cobra framework are not
	// followed by stack traces all the time.
	isCapturingLogs bool

	// The path to GALASA_HOME. Over-rides the environment variable.
	CmdParamGalasaHomePath string
}

var rootCmdValues *RootCmdValues

type RootCommand struct {
	values       *RootCmdValues
	cobraCommand *cobra.Command
}

// -------------------------------------------------------------------------------
// Constructor
// -------------------------------------------------------------------------------
func NewRootCommand(factory Factory) (*RootCommand, error) {
	cmd := new(RootCommand)

	err := cmd.init(factory)

	return cmd, err
}

// -------------------------------------------------------------------------------
// Public methods
// -------------------------------------------------------------------------------
func (cmd *RootCommand) init(factory Factory) error {

	var err error = nil

	// Flags parsed by this command put values into this instance of the structure.
	cmd.values = &RootCmdValues{
		isCapturingLogs: false,
	}

	cmd.cobraCommand, err = newRootCobraCommand(factory, cmd.values)

	// if err == nil {
	// 	err = createRootCmdChildren(factory, cmd.cobraCommand, rootCmdValues)
	// }

	return err
}

func (cmd *RootCommand) GetName() string {
	return COMMAND_NAME_ROOT
}

func (cmd *RootCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

//-------------------------------------------------------------------------------
// Private methods
//-------------------------------------------------------------------------------

// TODO: make a method on the RootCommand
func CreateRootCmd(factory Factory) (*cobra.Command, error) {
	// Flags parsed by this command put values into this instance of the structure.
	rootCmdValues = &RootCmdValues{
		isCapturingLogs: false,
	}

	cobraCommand, err := newRootCobraCommand(factory, rootCmdValues)
	if err == nil {
		err = createRootCmdChildren(factory, cobraCommand, rootCmdValues)
	}

	if err == nil {
		sanitiseCommandHelpDescriptions(cobraCommand)
	}

	return cobraCommand, err
}

// TODO: Turn into an object-method.
func newRootCobraCommand(factory Factory, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	version, err := embedded.GetGalasaCtlVersion()
	var rootCmd *cobra.Command
	if err == nil {

		rootCmd = &cobra.Command{
			Use:     "galasactl",
			Short:   "CLI for Galasa",
			Long:    `A tool for controlling Galasa resources using the command-line.`,
			Version: version,
		}

		rootCmd.SetErr(factory.GetStdErrConsole())
		rootCmd.SetOut(factory.GetStdOutConsole())

		var galasaCtlVersion string
		galasaCtlVersion, err = embedded.GetGalasaCtlVersion()
		if err == nil {

			rootCmd.Version = galasaCtlVersion

			rootCmd.PersistentFlags().StringVarP(&(rootCmdValues.logFileName), "log", "l", "",
				"File to which log information will be sent. Any folder referred to must exist. "+
					"An existing file will be overwritten. "+
					"Specify \"-\" to log to stderr. "+
					"Defaults to not logging.")

			rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

			rootCmd.PersistentFlags().StringVarP(&(rootCmdValues.CmdParamGalasaHomePath), "galasahome", "", "",
				"Path to a folder where Galasa will read and write files and configuration settings. "+
					"The default is '${HOME}/.galasa'. "+
					"This overrides the GALASA_HOME environment variable which may be set instead.",
			)

		}
	}
	return rootCmd, err
}

func sanitiseCommandHelpDescriptions(rootCmd *cobra.Command) {
	setHelpFlagForAllCommands(rootCmd, func(cobra *cobra.Command) {
		alias := cobra.NameAndAliases()
		//if the command has an alias,
		//the format would be cobra.Name, cobra.Aliases
		//otherwise it is just cobra.Name
		nameAndAliases := strings.Split(alias, ", ")
		if len(nameAndAliases) > 1 {
			alias = nameAndAliases[1]
		}

		cobra.Flags().BoolP("help", "h", false, "Displays the options for the "+alias+" command.")
	})
}

func createRootCmdChildren(factory Factory, rootCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createLocalCmd(factory, rootCmd, rootCmdValues)
	if err == nil {
		_, err = createProjectCmd(factory, rootCmd, rootCmdValues)
	}
	if err == nil {
		_, err = createPropertiesCmd(factory, rootCmd, rootCmdValues)
	}
	if err == nil {
		_, err = createRunsCmd(factory, rootCmd, rootCmdValues)
	}
	if err == nil {
		_, err = createAuthCmd(factory, rootCmd, rootCmdValues)
	}
	return err
}

// The main entry point into the cmd package.
func Execute(factory Factory, args []string) error {
	var err error

	finalWordHandler := factory.GetFinalWordHandler()

	var commands CommandCollection
	commands, err = NewCommandCollection(factory)

	if err == nil {

		// Catch execution if a panic happens.
		defer func() {
			err := recover()

			// Display the error and exit.
			finalWordHandler.FinalWord(err)
		}()

		// Execute the command
		err = commands.Execute(args)
	}
	finalWordHandler.FinalWord(err)
	return err
}

func setHelpFlagForAllCommands(command *cobra.Command, setHelpFlag func(*cobra.Command)) {
	setHelpFlag(command)

	//for all the commands eg properties get, set etc
	for _, cobraCommand := range command.Commands() {
		setHelpFlagForAllCommands(cobraCommand, setHelpFlag)
	}
}
