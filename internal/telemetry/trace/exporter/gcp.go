package exporter

import (
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

func NewGCP() *texporter.Exporter {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.New: %v", err)
	}
	return exporter
}
