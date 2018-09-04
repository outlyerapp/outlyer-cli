package command

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"sync"

	"github.com/outlyerapp/outlyer-cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// NewExportCommand creates a Command for exporting Outlyer resources to disk.
func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export .|[resource]|[resource/name]",
		Short: "Export resources from the specified account. The available resources are: alerts, checks, dashboards and plugins",
		Example: `
Export the entire account resources (alerts, checks, dashboards and plugins) to the current folder:
$ outlyer export . --account=<your_account>

Export the entire account resources (alerts, checks, dashboards and plugins) to a specific folder:
$ outlyer export . --account=<your_account> --folder=<your_folder>

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

	// Revert this code block when the view export is implemented for all endpoints
	for _, resourceToFetch := range args {
		if resourceToFetch == "." {
			args = args[:0]
			args = append(args, Alerts)
			args = append(args, Checks)
			args = append(args, Dashboards)
			args = append(args, Plugins)
			break
		}
	}

	// Remove this code block when the view export is implemented for all endpoints
	var resourceNames []string
	for _, resourceToFetch := range args {
		var resources []map[string]interface{}
		if resourceToFetch == Alerts {
			resp, err := api.Get("/accounts/" + account + "/" + resourceToFetch)
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Could not fetch %s from account %s\n%s", resourceToFetch, account, err))
			}

			yaml.Unmarshal(resp, &resources)

			for _, resource := range resources {
				resourceNames = append(resourceNames, "alerts/"+resource["name"].(string))
			}
		}
		if resourceToFetch == Checks {
			resp, err := api.Get("/accounts/" + account + "/" + resourceToFetch)
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Could not fetch %s from account %s\n%s", resourceToFetch, account, err))
			}

			yaml.Unmarshal(resp, &resources)

			for _, resource := range resources {
				resourceNames = append(resourceNames, "checks/"+resource["name"].(string))
			}
		}
		if resourceToFetch == Dashboards {
			resp, err := api.Get("/accounts/" + account + "/" + resourceToFetch)
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Could not fetch %s from account %s\n%s", resourceToFetch, account, err))
			}

			yaml.Unmarshal(resp, &resources)

			for _, resource := range resources {
				resourceNames = append(resourceNames, "dashboards/"+resource["name"].(string))
			}
		}
		if resourceToFetch == Plugins {
			resp, err := api.Get("/accounts/" + account + "/" + resourceToFetch)
			if err != nil {
				ExitWithError(ExitError, fmt.Errorf("Could not fetch %s from account %s\n%s", resourceToFetch, account, err))
			}

			yaml.Unmarshal(resp, &resources)

			for _, resource := range resources {
				resourceNames = append(resourceNames, "plugins/"+resource["name"].(string))
			}
		}
	}
	args = remove(args, Alerts)
	args = remove(args, Dashboards)
	args = remove(args, Checks)
	args = remove(args, Plugins)
	args = append(args, resourceNames...)
	args = removeDuplicates(args)

	// There is no "." argument, so fetches all listed resources
	for _, resourceToFetch := range args {
		wg.Add(1)
		go export(resourceToFetch, account, getOutputFolder(outputFolderFlag, resourceToFetch), &wg)
	}
	wg.Wait()

	fmt.Println("Resources successfully exported:")
	for _, resource := range args {
		fmt.Printf("- %s\n", resource)
	}
}

// export queries the resources for the given user account and persists them locally
func export(resourceToFetch, account, outputFolder string, wg *sync.WaitGroup) {
	resp, err := api.Get("/accounts/" + account + "/" + resourceToFetch + "?view=export")
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Could not fetch %s from account %s\n%s", resourceToFetch, account, err))
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

// getOutputFolder is a helper function to build the correct output folder to export the given resource
func getOutputFolder(outputFolderFlag, resourceToFetch string) string {
	if outputFolderFlag != "" {
		if outputFolderFlag[:1] == "~" { // Replace ~ by the user's full home path
			user, err := user.Current()
			if err != nil {
				ExitWithError(ExitError, err)
			}
			outputFolderFlag = strings.Replace(outputFolderFlag, "~", user.HomeDir, 1)
		}
		if outputFolderFlag[len(outputFolderFlag)-1:] != "/" { // Appends a / to the end of the folder path if it not exists
			outputFolderFlag = outputFolderFlag + "/"
		}
	}

	var baseResource string
	slashIndex := strings.Index(resourceToFetch, "/")
	if slashIndex != -1 { // The user specified a single resource like 'dashboards/docker', so ignore '/docker'
		baseResource = resourceToFetch[:slashIndex+1]
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

// remove removes a string from a []string
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
