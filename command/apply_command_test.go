package command

import (
	"testing"
)

func TestGetType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"alerts/docker.yaml", "alerts"},
		{"/dir1/alerts/docker.yaml", "alerts"},
		{"/dir1/dir2/alerts/docker.yaml", "alerts"},
		{"checks/docker.yaml", "checks"},
		{"/dir1/checks/docker.yaml", "checks"},
		{"/dir1/dir2/checks/docker.yaml", "checks"},
		{"dashboards/docker.yaml", "dashboards"},
		{"dir1/dashboards/docker.yaml", "dashboards"},
		{"/dir1/dir2/dashboards/docker.yaml", "dashboards"},
		{"plugins/docker.py", "plugins"},
		{"dir1/plugins/docker.py", "plugins"},
		{"/dir1/dir2/plugins/docker.py", "plugins"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			r := &resource{
				path: test.path,
			}
			if got := r.getType(); got != test.want {
				t.Errorf("resource.getType() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetTypeAndName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"alerts/docker.yaml", "alerts/docker"},
		{"/dir1/alerts/docker.yaml", "alerts/docker"},
		{"/dir1/dir2/alerts/docker.yaml", "alerts/docker"},
		{"checks/docker.yaml", "checks/docker"},
		{"/dir1/checks/docker.yaml", "checks/docker"},
		{"/dir1/dir2/checks/docker.yaml", "checks/docker"},
		{"dashboards/docker.yaml", "dashboards/docker"},
		{"dir1/dashboards/docker.yaml", "dashboards/docker"},
		{"/dir1/dir2/dashboards/docker.yaml", "dashboards/docker"},
		{"plugins/docker.py", "plugins/docker"},
		{"dir1/plugins/docker.py", "plugins/docker"},
		{"/dir1/dir2/plugins/docker.py", "plugins/docker"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			r := &resource{
				path: test.path,
			}
			if got := r.getTypeAndName(); got != test.want {
				t.Errorf("resource.getTypeAndName() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetTypeAndNameWithExtension(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"alerts/docker.yaml", "alerts/docker.yaml"},
		{"/dir1/alerts/docker.yaml", "alerts/docker.yaml"},
		{"/dir1/dir2/alerts/docker.yaml", "alerts/docker.yaml"},
		{"checks/docker.yaml", "checks/docker.yaml"},
		{"/dir1/checks/docker.yaml", "checks/docker.yaml"},
		{"/dir1/dir2/checks/docker.yaml", "checks/docker.yaml"},
		{"dashboards/docker.yaml", "dashboards/docker.yaml"},
		{"dir1/dashboards/docker.yaml", "dashboards/docker.yaml"},
		{"/dir1/dir2/dashboards/docker.yaml", "dashboards/docker.yaml"},
		{"plugins/docker.py", "plugins/docker.py"},
		{"dir1/plugins/docker.py", "plugins/docker.py"},
		{"/dir1/dir2/plugins/docker.py", "plugins/docker.py"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			r := &resource{
				path: test.path,
			}
			if got := r.getTypeAndNameWithExtension(); got != test.want {
				t.Errorf("resource.getTypeAndNameWithExtension() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetNameWithExtension(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"alerts/docker.yaml", "docker.yaml"},
		{"/dir1/alerts/docker.yaml", "docker.yaml"},
		{"/dir1/dir2/alerts/docker.yaml", "docker.yaml"},
		{"checks/docker.yaml", "docker.yaml"},
		{"/dir1/checks/docker.yaml", "docker.yaml"},
		{"/dir1/dir2/checks/docker.yaml", "docker.yaml"},
		{"dashboards/docker.yaml", "docker.yaml"},
		{"dir1/dashboards/docker.yaml", "docker.yaml"},
		{"/dir1/dir2/dashboards/docker.yaml", "docker.yaml"},
		{"plugins/docker.py", "docker.py"},
		{"dir1/plugins/docker.py", "docker.py"},
		{"/dir1/dir2/plugins/docker.py", "docker.py"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			r := &resource{
				path: test.path,
			}
			if got := r.getNameWithExtension(); got != test.want {
				t.Errorf("resource.getNameWithExtension() = %v, want %v", got, test.want)
			}
		})
	}
}
