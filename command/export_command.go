package command

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

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

	cmd.PersistentFlags().StringP("account", "a", "", "(Required) User account to use")
	cmd.PersistentFlags().StringP("folder", "f", "", "(Optional) Folder to export resources. If not provided, exports to the current folder")
	return cmd
}

// exportCommand validates the user input and calls export for each resource
// provided by the user
func exportCommand(cmd *cobra.Command, args []string) {
	account := cmd.PersistentFlags().Lookup("account").Value.String()
	if account == "" {
		ExitWithError(ExitBadArgs, fmt.Errorf("Account is required"))
	}
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("Resource is required"))
	}

	outputFolderFlag := cmd.PersistentFlags().Lookup("folder").Value.String()

	// Creates WaitGroup to wait for goroutines to finish exporting resources concurrently
	var wg sync.WaitGroup

	// Avoids fetching specific resources if the arguments contain "all"
	for _, resourceToFetch := range args {
		if resourceToFetch == "all" {
			wg.Add(4)
			go export("alerts", account, getOutputFolder(outputFolderFlag, "alerts"), &wg)
			go export("checks", account, getOutputFolder(outputFolderFlag, "checks"), &wg)
			go export("dashboards", account, getOutputFolder(outputFolderFlag, "dashboards"), &wg)
			go export("plugins", account, getOutputFolder(outputFolderFlag, "plugins"), &wg)
			wg.Wait()
			ExitWithSuccess("Done! Your account was successfully exported")
		}
	}

	// There is no "all" argument, so fetches all listed resources
	for _, resourceToFetch := range args {
		wg.Add(1)
		go export(resourceToFetch, account, getOutputFolder(outputFolderFlag, resourceToFetch), &wg)
	}
	wg.Wait()
	ExitWithSuccess("Done! Resources successfully exported")
}

// export queries the resources for the given user account and persists them locally
func export(resourceToFetch, account, outputFolder string, wg *sync.WaitGroup) {
	fmt.Printf("Exporting %s...\n", resourceToFetch)

	resp, err := outlyer.Get("/accounts/" + account + "/" + resourceToFetch + "?view=export")
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Could not fetch resource %s from account %s%s", resourceToFetch, account, err))
	}

	var resources []map[string]interface{}

	if isSingleResource(resourceToFetch) {
		var singleResource map[string]interface{}
		yaml.Unmarshal(resp, &singleResource)
		resources = make([]map[string]interface{}, 1)
		resources[0] = singleResource
	} else {
		yaml.Unmarshal(resp, &resources)
	}

	convertCheckFields(resources, outputFolder)
	os.MkdirAll(outputFolder, 0755)

	for _, resource := range resources {
		var resourceInBytes []byte
		var resourceFileName string
		resourceName := resource["name"].(string)

		if strings.Contains(resourceToFetch, "plugins") {
			resourceInBytes, err = base64.StdEncoding.DecodeString(resource["content"].(string))
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Could not decode plugin %s\n%s", resourceName, err))
			}
			resourceFileName = outputFolder + resourceName
		} else {
			resourceInBytes, err = yaml.Marshal(&resource)
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Error marshalling resource %s\n%s", resourceName, err))
			}
			resourceFileName = outputFolder + resourceName + ".yaml"
		}

		err := ioutil.WriteFile(resourceFileName, resourceInBytes, 0644)
		if err != nil {
			ExitWithError(ExitError, fmt.Errorf("Could not write resource %s to disk\n%s", resourceFileName, err))
		}
	}
	wg.Done()
}

// convertCheckFields converts the format and variables field names into handler and env
// this function can be removed as soon as https://github.com/outlyerapp/public_api/issues/349
// is merged.
func convertCheckFields(checks []map[string]interface{}, outputFolder string) {
	if strings.Contains(outputFolder, "checks") {
		for _, check := range checks {
			handler := check["format"]
			env := check["variables"]

			check["handler"] = handler
			check["env"] = env

			delete(check, "format")
			delete(check, "variables")
		}
	}
}

// getOutputFolder is a helper function to build the correct output folder to export the given resource
func getOutputFolder(outputFolderFlag, resourceToFetch string) string {
	if outputFolderFlag == "" { // Replace empty by the current directory
		outputFolderFlag = "."
	}
	if outputFolderFlag[:1] == "~" { // Replace ~ by the user's full home path
		outputFolderFlag = strings.Replace(outputFolderFlag, "~", os.Getenv("HOME"), 1)
	}
	if outputFolderFlag[len(outputFolderFlag)-1:] != "/" { // Appends a / to the end of the folder path
		outputFolderFlag = outputFolderFlag + "/"
	}

	var baseResource string
	slashIndex := strings.Index(resourceToFetch, "/")
	if slashIndex != -1 { // The user specified a single resource like 'dashboards/docker', so ignore '/docker'
		baseResource = resourceToFetch[:slashIndex] + "/"
	} else {
		baseResource = resourceToFetch + "/"
	}
	outputFolder := outputFolderFlag + baseResource
	return outputFolder
}

// isSingleResource checks whether the user provided a single resource like dashboards/docker
func isSingleResource(resourceToFetch string) bool {
	return strings.Contains(resourceToFetch, "/")
}
