package command

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/outlyer/outlyer-cli/api"
	"github.com/spf13/cobra"
)

type plugin struct {
	Content  string `yaml:"content"`
	Encoding string `yaml:"encoding"`
	Name     string `yaml:"name"`
}

type resource struct {
	path   string
	bytes  []byte
	status string
	err    error
}

func (r *resource) getType() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getTypeAndName() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)/[^.]+`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getTypeAndNameWithExtension() string {
	regex := regexp.MustCompile(`(alerts|checks|dashboards|plugins)/.+`)
	res := regex.FindStringSubmatch(r.path)
	return res[0]
}

func (r *resource) getNameWithExtension() string {
	regex := regexp.MustCompile(`(.*)/(.+)$`)
	res := regex.FindStringSubmatch(r.path)
	return res[2]
}

// NewApplyCommand creates a Command for applying resources to the user's Outlyer account
func NewApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply .|[folder]|[file]",
		Short: "Apply resources to the specified account. The available resources are: alerts, checks, dashboards and plugins",
		Run:   applyCommand,
	}

	cmd.PersistentFlags().StringP("account", "a", "", "(Required) User account to use")
	return cmd
}

func applyCommand(cmd *cobra.Command, args []string) {
	account := cmd.PersistentFlags().Lookup("account").Value.String()
	if account == "" {
		ExitWithError(ExitBadArgs, fmt.Errorf("Account is required"))
	}

	if len(args) < 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("Resource is required"))
	}

	paths := getPaths(args)
	resources := getResources(paths)

	fmt.Printf("\nResources to apply...\n\n")
	for _, resource := range resources {
		fmt.Printf("\t- %s\n", resource.path)
	}
	fmt.Printf("\nAre you sure you want to apply to account '%s'? [y/n] ", account)

	reader := bufio.NewReader(os.Stdin)
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.Replace(confirmation, "\n", "", -1) // removes return character on *unix and darwin
	confirmation = strings.Replace(confirmation, "\r", "", -1) // removes return character on windows

	if confirmation == "y" || confirmation == "Y" {
		// Creates WaitGroup to wait for goroutines to finish applying resources concurrently
		var wg sync.WaitGroup

		for i := 0; i < len(resources); i++ {
			wg.Add(1)
			go apply(account, &resources[i], &wg)
		}

		wg.Wait()
		fmt.Println("")
		fmt.Printf(getColumnPattern(), "ACCOUNT", "RESOURCE", "STATUS", "REASON")
		for _, resource := range resources {
			if resource.err == nil {
				fmt.Printf(getColumnPattern(), account, resource.getTypeAndNameWithExtension(), resource.status, "")
			} else {
				fmt.Printf(getColumnPattern(), account, resource.getTypeAndNameWithExtension(), "FAIL", resource.err)
			}
		}
		fmt.Println("")
	} else {
		fmt.Println("Skipping apply. 0 resources applied.")
	}
}

func apply(account string, resource *resource, wg *sync.WaitGroup) {
	var resp *api.Response
	var err error
	if resource.getType() == "plugins" {
		resp, err = api.Patch("/accounts/"+account+"/"+resource.getTypeAndNameWithExtension(), resource.bytes)
	} else {
		resp, err = api.Patch("/accounts/"+account+"/"+resource.getTypeAndName(), resource.bytes)
	}
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("could not process request"))
	}

	resource.status = "OK [UPDATED]"
	resource.err = resp.ErrorDetail

	if resp.Code == 404 {
		resp, err = api.Post("/accounts/"+account+"/"+resource.getType(), resource.bytes)
		if err != nil {
			ExitWithError(ExitError, fmt.Errorf("could not process request"))
		}
		resource.status = "OK [CREATED]"
		resource.err = resp.ErrorDetail
	}
	wg.Done()
}

func getPaths(args []string) []string {
	var paths []string

	for _, arg := range args {
		if !fileOrDirExists(arg) {
			ExitWithError(ExitError, fmt.Errorf("%s: no such file or directory", arg))
		}

		fileInfo, _ := os.Stat(arg)
		if fileInfo.IsDir() {
			arg = appendSlashTo(arg)
			dirWithResourceName, _ := regexp.Compile("(alerts|checks|dashboards|plugins)/")
			if dirWithResourceName.MatchString(arg) { // Is the dir a resource name?
				files, _ := ioutil.ReadDir(arg)
				for _, file := range files { // Then add all resources from it
					paths = append(paths, arg+file.Name())
				}
			} else {
				// The dir name is not a valid resource name (like dir1), but does it have any subdir containing resources?
				// Covers the case "apply dir1/ --account=my-account", where dir1 has subdirs "alerts", "checks", etc
				files, _ := ioutil.ReadDir(arg)
				for _, file := range files {
					if file.IsDir() {
						regex, _ := regexp.Compile("(alerts|checks|dashboards|plugins)")
						if regex.MatchString(file.Name()) {
							resources, _ := ioutil.ReadDir(arg + file.Name())
							for _, resource := range resources {
								paths = append(paths, arg+file.Name()+"/"+resource.Name())
							}
						}
					}
				}
			}
		} else {
			validResourcePath, _ := regexp.Compile("(.*)(alerts|checks|dashboards|plugins)/[^.]+...[^.]")
			if validResourcePath.MatchString(arg) {
				paths = append(paths, arg)
			}
		}
	}

	paths = removeDuplicates(paths)

	if len(paths) == 0 {
		ExitWithError(ExitError, fmt.Errorf("could not find any resources to apply"))
	}

	return paths
}

func getResources(paths []string) []resource {
	resources := make([]resource, len(paths))
	for i, path := range paths {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			ExitWithError(ExitError, err)
		}

		res := resource{path, bytes, "FAIL", nil}
		if res.getType() == "checks" {
			res.bytes = convertCheckFieldsToSend(bytes)
		} else if res.getType() == "plugins" {
			res.bytes = bytes
			res = convertPlugin(res)
		}
		resources[i] = res
	}
	return resources
}

func convertPlugin(res resource) resource {
	pluginBase64 := base64.StdEncoding.EncodeToString(res.bytes)
	plugin := &plugin{Content: pluginBase64, Name: res.getNameWithExtension(), Encoding: "base64"}
	pluginInBytes, _ := yaml.Marshal(&plugin)
	res.bytes = pluginInBytes
	return res
}

func convertCheckFieldsToSend(bytes []byte) []byte {
	check := make(map[string]interface{})
	yaml.Unmarshal(bytes, &check)

	format := check["handler"]
	variables := check["env"]

	check["format"] = format
	check["variables"] = variables

	delete(check, "handler")
	delete(check, "env")

	checkInBytes, err := yaml.Marshal(check)
	if err != nil {
		ExitWithError(ExitError, fmt.Errorf("Error marshalling resource %s\n%s", check["name"], err))
	}
	return checkInBytes
}

func getColumnPattern() string {
	return "%-20s\t%-40s\t%-20s\t%-25v\n"
}
