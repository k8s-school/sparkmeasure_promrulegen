package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func toHelpText(metric string) string {
	parts := strings.Split(metric, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, " ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	re := regexp.MustCompile(`# HELP (\S+) Attribute exposed for management sparkmeasure\.metrics:name=(\S+)\.(\S+)\.(\S+)\.(\S+),type=\S+,attribute=Value`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		if len(matches) == 6 {
			metricName := matches[1]
			namespace := matches[2]
			appID := matches[3]
			jmxMetric := matches[4]
			jmxType := matches[5]

			promType := "gauge"
			if strings.Contains(metricName, "total") || strings.Contains(jmxType, "counter") {
				promType = "counter"
			}

			help := toHelpText(strings.TrimPrefix(metricName, "sparkmeasure_"))

			fmt.Printf("- pattern: 'sparkmeasure\\.metrics<name=%s\\.%s\\.%s\\.%s,\\s*type=\\S+><>Value'\n", namespace, appID, jmxMetric, jmxType)
			fmt.Printf("  name: %s\n", metricName)
			fmt.Printf("  type: %s\n", strings.ToLower(promType))
			fmt.Printf("  help: \"%s\"\n", help)
			fmt.Printf("  labels:\n")
			fmt.Printf("    app_namespace: \"%s\"\n", namespace)
			fmt.Printf("    app_id: \"%s\"\n\n", appID)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
