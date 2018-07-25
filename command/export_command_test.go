package command

import (
	"fmt"
	"os/user"
	"testing"
)

func TestIsSingleResource(t *testing.T) {
	tests := []struct {
		resourceToFetch string
		want            bool
	}{
		{"dashboards", false},
		{"dashboards/docker", true},
	}
	for _, test := range tests {
		t.Run(test.resourceToFetch, func(t *testing.T) {
			if got := isSingleResource(test.resourceToFetch); got != test.want {
				t.Errorf("isSingleResource() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetOutputFolder(t *testing.T) {
	user, _ := user.Current()

	tests := []struct {
		outputFolderFlag string
		resourceToFetch  string
		want             string
	}{
		{"", "dashboards", "dashboards/"},
		{"", "dashboards/docker", "dashboards/"},
		{"demo", "dashboards", "demo/dashboards/"},
		{"demo/", "dashboards", "demo/dashboards/"},
		{"demo/test1/test2", "dashboards", "demo/test1/test2/dashboards/"},
		{"demo", "dashboards/docker", "demo/dashboards/"},
		{"demo/", "dashboards/docker", "demo/dashboards/"},
		{".", "dashboards", "./dashboards/"},
		{"./demo", "dashboards", "./demo/dashboards/"},
		{"~", "dashboards", user.HomeDir + "/dashboards/"},
		{"~/demo", "dashboards", user.HomeDir + "/demo/dashboards/"},
		{"~/demo/test", "dashboards", user.HomeDir + "/demo/test/dashboards/"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s,%s", test.outputFolderFlag, test.resourceToFetch), func(t *testing.T) {
			if got := getOutputFolder(test.outputFolderFlag, test.resourceToFetch); got != test.want {
				t.Errorf("getOutputFolder() = %v, want %v", got, test.want)
			}
		})
	}
}
