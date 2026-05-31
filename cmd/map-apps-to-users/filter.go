package main

import (
	"fmt"
	get_user_apps "msuite-toolkit/pkg/endpoints/get-user-apps"
	"net"
	"strconv"
	"strings"
)

func filterAppsByDestination(appsMap map[*get_user_apps.AuthorizedApp][]string, destinationHosts string, destinationPorts string) map[*get_user_apps.AuthorizedApp][]string {
	if strings.TrimSpace(destinationHosts) == "" && strings.TrimSpace(destinationPorts) == "" {
		return appsMap
	}

	allowedHosts, err := expandDestinationHosts(parseCommaSeparatedValues(destinationHosts))
	if err != nil {
		return map[*get_user_apps.AuthorizedApp][]string{}
	}

	allowedPorts, err := expandDestinationPorts(parseCommaSeparatedValues(destinationPorts))
	if err != nil {
		return map[*get_user_apps.AuthorizedApp][]string{}
	}

	filtered := make(map[*get_user_apps.AuthorizedApp][]string)
	for app, users := range appsMap {
		appHosts, err := expandDestinationHosts([]string{app.App.DestinationSetting.IPDef})
		if err != nil {
			continue
		}

		appPorts, err := expandDestinationPorts([]string{app.App.DestinationSetting.PortDef})
		if err != nil {
			continue
		}

		if matchesDestinationFilter(appHosts, allowedHosts) && matchesPortFilter(appPorts, allowedPorts) {
			filtered[app] = users
		}
	}

	return filtered
}

func parseCommaSeparatedValues(value string) []string {
	var values []string
	for item := range strings.SplitSeq(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		values = append(values, item)
	}

	return values
}

func matchesDestinationFilter(appValues []string, filterValues []string) bool {
	if len(filterValues) == 0 {
		return true
	}

	for _, appValue := range appValues {
		for _, filterValue := range filterValues {
			if appValue == filterValue {
				return true
			}
		}
	}

	return false
}

func matchesPortFilter(appValues []int, filterValues []int) bool {
	if len(filterValues) == 0 {
		return true
	}

	for _, appValue := range appValues {
		for _, filterValue := range filterValues {
			if appValue == filterValue {
				return true
			}
		}
	}

	return false
}

func expandDestinationHosts(values []string) ([]string, error) {
	var expanded []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if strings.Contains(value, "-") {
			hosts, err := expandDestinationHostRange(value)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, hosts...)
			continue
		}

		if ip := net.ParseIP(value); ip == nil {
			return nil, fmt.Errorf("invalid destination host: %s", value)
		}
		expanded = append(expanded, value)
	}

	return expanded, nil
}

func expandDestinationHostRange(value string) ([]string, error) {
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid destination host range: %s", value)
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return nil, fmt.Errorf("invalid destination host range start: %s", parts[0])
	}
	start4 := startIP.To4()
	if start4 == nil {
		return nil, fmt.Errorf("only IPv4 destination host ranges are supported: %s", value)
	}

	startParts := strings.Split(strings.TrimSpace(parts[0]), ".")
	if len(startParts) != 4 {
		return nil, fmt.Errorf("invalid destination host range start: %s", parts[0])
	}

	startOctet, err := strconv.Atoi(startParts[3])
	if err != nil {
		return nil, fmt.Errorf("invalid destination host range start octet: %w", err)
	}

	endPart := strings.TrimSpace(parts[1])
	endOctet, err := strconv.Atoi(endPart)
	if err != nil {
		return nil, fmt.Errorf("invalid destination host range end octet: %w", err)
	}

	if endOctet < startOctet {
		return nil, fmt.Errorf("destination host range end must be >= start: %s", value)
	}

	prefix := strings.Join(startParts[:3], ".")
	expanded := make([]string, 0, endOctet-startOctet+1)
	for octet := startOctet; octet <= endOctet; octet++ {
		expanded = append(expanded, fmt.Sprintf("%s.%d", prefix, octet))
	}

	return expanded, nil
}

func expandDestinationPorts(values []string) ([]int, error) {
	var expanded []int
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if strings.Contains(value, "-") {
			ports, err := expandDestinationPortRange(value)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, ports...)
			continue
		}

		port, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid destination port: %s", value)
		}
		expanded = append(expanded, port)
	}

	return expanded, nil
}

func expandDestinationPortRange(value string) ([]int, error) {
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid destination port range: %s", value)
	}

	startPort, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid destination port range start: %w", err)
	}
	endPort, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid destination port range end: %w", err)
	}

	if endPort < startPort {
		return nil, fmt.Errorf("destination port range end must be >= start: %s", value)
	}

	expanded := make([]int, 0, endPort-startPort+1)
	for port := startPort; port <= endPort; port++ {
		expanded = append(expanded, port)
	}

	return expanded, nil
}
