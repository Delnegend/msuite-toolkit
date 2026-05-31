package main

import (
	"reflect"
	"testing"

	get_user_apps "msuite-toolkit/pkg/endpoints/get-user-apps"
)

func TestExpandDestinationHosts(t *testing.T) {
	values, err := expandDestinationHosts([]string{"10.0.0.1", "10.0.0.2-4"})
	if err != nil {
		t.Fatalf("expandDestinationHosts returned error: %v", err)
	}

	expected := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4"}
	if !reflect.DeepEqual(values, expected) {
		t.Fatalf("unexpected hosts: got %v want %v", values, expected)
	}
}

func TestExpandDestinationPorts(t *testing.T) {
	values, err := expandDestinationPorts([]string{"442", "443-445"})
	if err != nil {
		t.Fatalf("expandDestinationPorts returned error: %v", err)
	}

	expected := []int{442, 443, 444, 445}
	if !reflect.DeepEqual(values, expected) {
		t.Fatalf("unexpected ports: got %v want %v", values, expected)
	}
}

func TestFilterAppsByDestination(t *testing.T) {
	matchingApp := get_user_apps.AuthorizedApp{}
	matchingApp.App.DestinationSetting.IPDef = "10.0.0.2"
	matchingApp.App.DestinationSetting.PortDef = "443"

	nonMatchingApp := get_user_apps.AuthorizedApp{}
	nonMatchingApp.App.DestinationSetting.IPDef = "10.0.0.9"
	nonMatchingApp.App.DestinationSetting.PortDef = "8080"

	appsMap := map[*get_user_apps.AuthorizedApp][]string{
		&matchingApp:    {"user@example.com"},
		&nonMatchingApp: {"other@example.com"},
	}

	filtered := filterAppsByDestination(appsMap, "10.0.0.1-2", "443-5000")
	if len(filtered) != 1 {
		t.Fatalf("unexpected filtered size: got %d want 1", len(filtered))
	}

	if _, ok := filtered[&matchingApp]; !ok {
		t.Fatalf("expected matching app to remain in filtered result")
	}

	if _, ok := filtered[&nonMatchingApp]; ok {
		t.Fatalf("expected non-matching app to be filtered out")
	}
}

func TestParseCommaSeparatedValues(t *testing.T) {
	values := parseCommaSeparatedValues("10.0.0.1, 10.0.0.2-4, ,10.0.0.9")
	expected := []string{"10.0.0.1", "10.0.0.2-4", "10.0.0.9"}
	if !reflect.DeepEqual(values, expected) {
		t.Fatalf("unexpected parsed values: got %v want %v", values, expected)
	}
}
