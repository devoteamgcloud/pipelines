// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"regexp"

	"github.com/golang/glog"
	apiv1beta1 "github.com/kubeflow/pipelines/backend/api/v1beta1/go_client"
	apiv2beta1 "github.com/kubeflow/pipelines/backend/api/v2beta1/go_client"
	"github.com/kubeflow/pipelines/backend/src/common/util"
	"google.golang.org/grpc/codes"
)

const (
	// This regex expresses the following constraints:
	// * Allows lowercase/uppercase letters
	// * Allows "_", "-" and numbers in the middle
	// * Additionally, numbers are also allowed at the end
	// * At most 64 characters
	metricNamePattern = "^[a-zA-Z]([-_a-zA-Z0-9]{0,62}[a-zA-Z0-9])?$"
)

// ValidateRunMetricV1 validates RunMetric fields from request.
func ValidateRunMetricV1(metric *apiv1beta1.RunMetric) error {
	matched, err := regexp.MatchString(metricNamePattern, metric.GetName())
	if err != nil {
		// This should never happen.
		return util.NewInternalServerError(
			err, "failed to compile pattern '%s'", metricNamePattern)
	}
	if !matched {
		return util.NewInvalidInputError(
			"metric.name '%s' doesn't match with the pattern '%s'", metric.GetName(), metricNamePattern)
	}
	if metric.GetNodeId() == "" {
		return util.NewInvalidInputError("metric.node_id must not be empty")
	}
	if len(metric.GetNodeId()) > 128 {
		return util.NewInvalidInputError(
			"metric.node_id '%s' cannot be longer than 128 characters", metric.GetNodeId())
	}
	return nil
}

// NewReportRunMetricResult turns error into a ReportRunMetricResult.
func NewReportRunMetricResultV1(
	metricName string, nodeID string, err error) *apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult {
	result := &apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult{
		MetricName:   metricName,
		MetricNodeId: nodeID,
	}
	if err == nil {
		result.Status = apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_OK
		return result
	}
	userError, ok := err.(*util.UserError)
	if !ok {
		result.Status = apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR
		return result
	}
	switch userError.ExternalStatusCode() {
	case codes.AlreadyExists:
		result.Status = apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_DUPLICATE_REPORTING
	case codes.InvalidArgument:
		result.Status = apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_INVALID_ARGUMENT
	default:
		result.Status = apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR
	}
	result.Message = userError.ExternalMessage()
	if result.Status == apiv1beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR {
		glog.Errorf("Internal error '%v' when reporting metric '%s/%s'", err, nodeID, metricName)
	}
	return result
}

// ValidateRunMetric validates RunMetric fields from request.
func ValidateRunMetric(metric *apiv2beta1.RunMetric) error {
	matched, err := regexp.MatchString(metricNamePattern, metric.GetDisplayName())
	if err != nil {
		// This should never happen.
		return util.NewInternalServerError(
			err, "failed to compile pattern '%s'", metricNamePattern)
	}
	if !matched {
		return util.NewInvalidInputError(
			"metric.name '%s' doesn't match with the pattern '%s'", metric.GetDisplayName(), metricNamePattern)
	}
	if metric.GetNodeId() == "" {
		return util.NewInvalidInputError("metric.node_id must not be empty")
	}
	if len(metric.GetNodeId()) > 128 {
		return util.NewInvalidInputError(
			"metric.node_id '%s' cannot be longer than 128 characters", metric.GetNodeId())
	}
	return nil
}

// NewReportRunMetricResult turns error into a ReportRunMetricResult.
func NewReportRunMetricResult(metricName string, nodeID string, err error) *apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult {
	result := &apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult{
		MetricName:   metricName,
		MetricNodeId: nodeID,
	}
	if err == nil {
		result.Status = apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_OK
		return result
	}
	userError, ok := err.(*util.UserError)
	if !ok {
		result.Status = apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR
		return result
	}
	switch userError.ExternalStatusCode() {
	case codes.AlreadyExists:
		result.Status = apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_DUPLICATE_REPORTING
	case codes.InvalidArgument:
		result.Status = apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_INVALID_ARGUMENT
	default:
		result.Status = apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR
	}
	result.Message = userError.ExternalMessage()
	if result.Status == apiv2beta1.ReportRunMetricsResponse_ReportRunMetricResult_INTERNAL_ERROR {
		glog.Errorf("Internal error '%v' when reporting metric '%s/%s'", err, nodeID, metricName)
	}
	return result
}
