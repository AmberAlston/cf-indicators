package indicator_test

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cf-indicators/pkg/indicator"
)

func TestValidDocument(t *testing.T) {
	t.Run("validation returns no errors if the document is valid", func(t *testing.T) {
		g := NewGomegaWithT(t)

		document := indicator.Document{
			Metrics: []indicator.Metric{
				{
					Title:       "Demo Latency",
					Description: "A test metric for testing",
					Name:        "latency",
					SourceID:    "demo",
					Origin:      "demo",
				},
			},
			Indicators: []indicator.Indicator{
				{
					Name:        "test_performance_indicator",
					Title:       "Test Performance Indicator",
					Description: "This is a valid markdown description.",
					PromQL:      "prom",
					Thresholds: []indicator.Threshold{
						{
							Level:    "warning",
							Dynamic:  true,
							Operator: indicator.GreaterThanOrEqualTo,
							Value:    50,
						},
					},
					MetricRefs: []indicator.MetricRef{{
						Name:     "latency",
						SourceID: "demo",
					}},
					Response:    "Panic!",
					Measurement: "Measurement Text",
				},
			},
			Documentation: indicator.Documentation{
				Title:       "Monitoring Test Product",
				Owner:       "Test Owner Team",
				Description: "Test description",
				Sections: []indicator.Section{{
					Title:       "Test Section",
					Description: "This section includes indicators and metrics",
					IndicatorRefs: []indicator.IndicatorRef{{
						Name: "test_performance_indicator",
					}},
					MetricRefs: []indicator.MetricRef{{
						Name:     "latency",
						SourceID: "demo",
					}},
				}},
			},
		}

		es := indicator.Validate(document)

		g.Expect(es).To(BeEmpty())
	})
}

func TestMetricValidation(t *testing.T) {

	t.Run("validation returns errors if any metric field is blank", func(t *testing.T) {
		g := NewGomegaWithT(t)

		document := indicator.Document{
			Metrics: []indicator.Metric{
				{
					Title:       " ",
					Description: " ",
					Name:        " ",
					SourceID:    " ",
					Origin:      " ",
				},
			},
		}

		es := indicator.Validate(document)

		g.Expect(es).To(ConsistOf(
			errors.New("metrics[0] title is required"),
			errors.New("metrics[0] description is required"),
			errors.New("metrics[0] name is required"),
			errors.New("metrics[0] source_id is required"),
			errors.New("metrics[0] origin is required"),
		))
	})
}

func TestIndicatorValidation(t *testing.T) {

	t.Run("validation returns errors if any indicator field is blank", func(t *testing.T) {
		g := NewGomegaWithT(t)

		document := indicator.Document{
			Indicators: []indicator.Indicator{
				{
					Name:        " ",
					Title:       " ",
					Description: " ",
					PromQL:      " ",
					Response:    " ",
					Measurement: " ",
					MetricRefs:  []indicator.MetricRef{},
				},
			},
		}

		es := indicator.Validate(document)

		g.Expect(es).To(ConsistOf(
			errors.New("indicators[0] name is required"),
			errors.New("indicators[0] title is required"),
			errors.New("indicators[0] description is required"),
			errors.New("indicators[0] promql is required"),
			errors.New("indicators[0] response is required"),
			errors.New("indicators[0] measurement is required"),
			errors.New("indicators[0] must reference at least 1 metric"),
		))
	})
}

func TestDocumentationValidation(t *testing.T) {

	t.Run("validation returns errors if metric or indicator is not found", func(t *testing.T) {
		g := NewGomegaWithT(t)

		document := indicator.Document{
			Documentation: indicator.Documentation{
				Sections: []indicator.Section{{
					IndicatorRefs: []indicator.IndicatorRef{{
						Name: "test_performance_indicator",
					}},
					MetricRefs: []indicator.MetricRef{{
						Name:     "latency",
						SourceID: "demo",
					}},
				}},
			},
		}

		es := indicator.Validate(document)

		g.Expect(es).To(ConsistOf(
			MatchError(ContainSubstring("references non-existent indicator")),
			MatchError(ContainSubstring("references non-existent metric")),
		))
	})
}
