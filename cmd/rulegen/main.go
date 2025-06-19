package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var logger = log.Default()

func toHelpText(metric string) string {
	parts := strings.Split(metric, "_")
	caser := cases.Title(language.Und)
	for i := range parts {
		parts[i] = caser.String(parts[i])
	}
	return strings.Join(parts, " ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	re := regexp.MustCompile(`# HELP (\S+) Attribute exposed for management sparkmeasure\.metrics:name=(\S+)\.(\S+)\.(\S+)\.(\S+),type=gauges,attribute=Value`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		logger.Println("Processing line: ", line)
		logger.Println("Matches found: ", matches)
		logger.Println("Number of matches: ", len(matches))

		if len(matches) == 6 {
			metricName := matches[1]
			namespace := matches[2]
			appID := matches[3]
			jmxMetric := matches[5]

			logger.Printf("- name: %s, namespace: %s, app_id: %s, jmx_metric: %s\n", metricName, namespace, appID, jmxMetric)

			promType := jmxMetric

			// error if promType is not one of the expected types: gauge or counter
			if promType != "gauge" && promType != "counter" {
				fmt.Fprintf(os.Stderr, "Error: unexpected metric type '%s' for metric '%s'\n", promType, metricName)
				os.Exit(1)
			}

			help := toHelpText(strings.TrimPrefix(metricName, "sparkmeasure_"))

			f, err := os.OpenFile("output.yaml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			writer := bufio.NewWriter(f)
			fmt.Fprintf(writer, "- pattern: 'sparkmeasure\\.metrics<name=%s\\.%s\\.%s\\.%s,\\s*type=gauges><>Value'\n", namespace, appID, metricName, jmxMetric)
			fmt.Fprintf(writer, "  name: %s\n", metricName)
			fmt.Fprintf(writer, "  type: %s\n", strings.ToLower(promType))
			fmt.Fprintf(writer, "  help: \"%s\"\n", help)
			fmt.Fprintf(writer, "  labels:\n")
			fmt.Fprintf(writer, "    app_namespace: \"%s\"\n", namespace)
			fmt.Fprintf(writer, "    app_id: \"%s\"\n\n", appID)
			writer.Flush()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
