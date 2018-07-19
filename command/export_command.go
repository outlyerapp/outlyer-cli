package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	outlyer "github.com/outlyer/outlyer-cli"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// NewExportCommand creates a Command for exporting Outlyer resources to disk.
func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export all|[resource]|[resource/name]",
		Short: "Export resources from the specified account. The available resources are: alerts, checks, dashboards and plugins",
		Example: `
Export the entire account resources (alerts, checks, dashboards and plugins) to the current folder:
$ outlyer export all --account=<your_account>

Export the entire account resources (alerts, checks, dashboards and plugins) to a specific folder:
$ outlyer export all --account=<your_account> --folder=<your_folder>

Export the account's alerts and dashboards to the current folder:
$ outlyer export alerts dashboards --account=<your_account>

Export the account's alerts and only two single dashboards to a specific folder:
$ outlyer export alerts dashboards/docker dashboards/kafka --account=<your_account> --folder=<your_folder>
`,
		Run: exportCommand,
	}

	cmd.PersistentFlags().StringP("account", "a", "", "User account to use")
	cmd.PersistentFlags().StringP("folder", "f", "", "Folder to save resources")
	return cmd
}

// exportCommand fetches the resources from a specific user account
// and persists them locally.
func exportCommand(cmd *cobra.Command, args []string) {
	account := cmd.PersistentFlags().Lookup("account").Value.String()
	if account == "" {
		ExitWithError(ExitBadArgs, fmt.Errorf("Account is required"))
	}

	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("Resource is required"))
	}

	baseOutputFolder := getBaseOutputFolder(cmd.PersistentFlags().Lookup("folder").Value.String())

	// Avoids fetching specific resources if the arguments contain "all"
	for _, resourceFlagValue := range args {
		if resourceFlagValue == "all" {
			fmt.Print("Exporting alerts... ")
			export("alerts", baseOutputFolder, account)
			fmt.Println("Done!")

			fmt.Print("Exporting checks... ")
			export("checks", baseOutputFolder, account)
			fmt.Println("Done!")

			fmt.Print("Exporting dashboards... ")
			export("dashboards", baseOutputFolder, account)
			fmt.Println("Done!")

			fmt.Print("Exporting plugins... ")
			export("plugins", baseOutputFolder, account)
			fmt.Println("Done!")

			ExitWithSuccess("Your account was successfully exported")
		}
	}

	// There is no "all" argument, so fetches all listed resources
	for _, resourceToFetch := range args {
		fmt.Printf("Exporting %s... ", resourceToFetch)
		export(resourceToFetch, baseOutputFolder, account)
		fmt.Println("Done!")
	}
	ExitWithSuccess("Resources successfully exported")
}

// export queries the Outlyer API for the provided resource(s) based on the user account set
// and then persists the response to yaml file(s) locally
func export(resourceToFetch, baseOutputFolder, account string) {
	resp, err := outlyer.Get("/accounts/" + account + "/" + resourceToFetch + "?view=export")
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Could not fetch resource %s from account %s%s", resourceToFetch, account, err))
	}

	// Builds the correct output folder based on the resource to export
	var baseResource string
	slashIndex := strings.Index(resourceToFetch, "/")
	if slashIndex != -1 {
		baseResource = resourceToFetch[:slashIndex]
	} else {
		baseResource = resourceToFetch
	}
	outputFolder := baseOutputFolder + baseResource + "/"

	if strings.Contains(resourceToFetch, "/") {
		exportSingleResourceToDisk(outputFolder, resp)
	} else {
		exportMultipleResourcesToDisk(outputFolder, resp)
	}
}

// exportSingleResourceToDisk converts the []byte response to a single yaml file
// and saves it into the given outputFolder
func exportSingleResourceToDisk(outputFolder string, resp []byte) {
	var resource map[string]interface{}
	yaml.Unmarshal(resp, &resource)

	os.MkdirAll(outputFolder, 0755)

	resourceFileName := outputFolder + strings.Replace(resource["name"].(string), ".py", "", 1) + ".yaml"
	resourceInBytes, _ := yaml.Marshal(&resource)

	err := ioutil.WriteFile(resourceFileName, resourceInBytes, 0644)
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Could not write resource %s to disk\n%s", resourceFileName, err))
	}
}

// exportMultipleResourcesToDisk converts the []byte response to multiple yaml files
// and saves them into the given outputFolder
func exportMultipleResourcesToDisk(outputFolder string, resp []byte) {
	var resources []map[string]interface{}
	yaml.Unmarshal(resp, &resources)

	os.MkdirAll(outputFolder, 0755)

	for _, resource := range resources {
		resourceFileName := outputFolder + strings.Replace(resource["name"].(string), ".py", "", 1) + ".yaml"
		resourceInBytes, _ := yaml.Marshal(&resource)

		err := ioutil.WriteFile(resourceFileName, resourceInBytes, 0644)
		if err != nil {
			ExitWithError(ExitError, fmt.Errorf("Could not write resource %s to disk\n%s", resourceFileName, err))
		}
	}
}

// getBaseOutputFolder is a helper function to build the correct output base folder name to export resources
func getBaseOutputFolder(outputFolder string) string {
	if outputFolder == "" {
		outputFolder = "./"
	}
	if outputFolder[:1] == "~" {
		outputFolder = strings.Replace(outputFolder, "~", os.Getenv("HOME"), 1)
	}
	if outputFolder[len(outputFolder)-1:] != "/" {
		outputFolder = outputFolder + "/"
	}
	return outputFolder
}
