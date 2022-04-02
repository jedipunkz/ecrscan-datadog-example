// Submit metrics returns "Payload accepted" response

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/jedipunkz/ecrscan/pkg/myecr"
)

func main() {
	e := myecr.Ecr{}

	e.Repositories = [][]string{
		{"scantest", "latest"},
	}
	e.Resion = "ap-northeast-1"

	finding, _, err := e.ListFindings()
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range e.Repositories {
		for k, v := range finding.FindingSeverityCounts {
			body := datadog.MetricsPayload{
				Series: []datadog.Series{
					{
						Metric: "ecrscan.image." + r[i],
						Type:   datadog.PtrString("gauge"),
						Points: [][]*float64{
							{
								datadog.PtrFloat64(float64(time.Now().Unix())),
								datadog.PtrFloat64(float64(*v)),
							},
						},
						Tags: &[]string{
							"imagename:" + k,
						},
					},
				},
			}
			ctx := datadog.NewDefaultContext(context.Background())
			configuration := datadog.NewConfiguration()
			apiClient := datadog.NewAPIClient(configuration)
			resp, r, err := apiClient.MetricsApi.SubmitMetrics(ctx, body, *datadog.NewSubmitMetricsOptionalParameters())

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
				fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			}

			responseContent, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", responseContent)
		}
	}
}
