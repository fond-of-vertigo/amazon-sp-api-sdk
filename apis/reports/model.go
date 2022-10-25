package reports

import (
	"fmt"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"net/url"
	"strings"
)

// ReportModel Detailed information about the report.
type ReportModel struct {
	// A list of marketplace identifiers for the report.
	MarketplaceIds []string `json:"marketplaceIds,omitempty"`
	// The identifier for the report. This identifier is unique only in combination with a seller ID.
	ReportId string `json:"reportId"`
	// The report type.
	ReportType string `json:"reportType"`
	// The start of a date and time range used for selecting the data to report.
	DataStartTime *apis.JsonTimeISO8601 `json:"dataStartTime,omitempty"`
	// The end of a date and time range used for selecting the data to report.
	DataEndTime *apis.JsonTimeISO8601 `json:"dataEndTime,omitempty"`
	// The identifier of the report schedule that created this report (if any). This identifier is unique only in combination with a seller ID.
	ReportScheduleId *string `json:"reportScheduleId,omitempty"`
	// The date and time when the report was created.
	CreatedTime apis.JsonTimeISO8601 `json:"createdTime"`
	// The processing status of the report.
	ProcessingStatus string `json:"processingStatus"`
	// The date and time when the report processing started, in ISO 8601 date time format.
	ProcessingStartTime *apis.JsonTimeISO8601 `json:"processingStartTime,omitempty"`
	// The date and time when the report processing completed, in ISO 8601 date time format.
	ProcessingEndTime *apis.JsonTimeISO8601 `json:"processingEndTime,omitempty"`
	// The identifier for the report document. Pass this into the getReportDocument operation to get the information you will need to retrieve the report document's contents.
	ReportDocumentId *string `json:"reportDocumentId,omitempty"`
}

type GetReportFilter struct {
	reportTypes        []string
	processingStatuses []string
	marketplaceIds     []string
	pageSize           int
	createdSince       apis.JsonTimeISO8601
	createdUntil       apis.JsonTimeISO8601
	nextToken          string
}

func (f *GetReportFilter) GetQuery() url.Values {
	q := url.Values{}
	q.Add("reportTypes", strings.Join(f.reportTypes, ","))
	q.Add("processingStatuses", strings.Join(f.processingStatuses, ","))
	q.Add("marketplaceIds", strings.Join(f.marketplaceIds, ","))
	q.Add("pageSize", fmt.Sprint(f.pageSize))
	q.Add("createdSince", f.createdSince.String())
	q.Add("createdUntil", f.createdUntil.String())
	q.Add("nextToken", f.nextToken)
	return q
}

// CreateReportSpecification Information required to create the report.
type CreateReportSpecification struct {
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// The report type.
	ReportType string `json:"reportType"`
	// The start of a date and time range, in ISO 8601 date time format, used for selecting the data to report. The default is now. The value must be prior to or equal to the current date and time. Not all report types make use of this.
	DataStartTime *apis.JsonTimeISO8601 `json:"dataStartTime,omitempty"`
	// The end of a date and time range, in ISO 8601 date time format, used for selecting the data to report. The default is now. The value must be prior to or equal to the current date and time. Not all report types make use of this.
	DataEndTime *apis.JsonTimeISO8601 `json:"dataEndTime,omitempty"`
	// A list of marketplace identifiers. The report document's contents will contain data for all of the specified marketplaces, unless the report type indicates otherwise.
	MarketplaceIds []string `json:"marketplaceIds"`
}

// CreateReportResponse Response schema.
type CreateReportResponse struct {
	// The identifier for the report. This identifier is unique only in combination with a seller ID.
	ReportId string `json:"reportId"`
}

// GetReportsResponse The response for the getReports operation.
type GetReportsResponse struct {
	// A list of reports.
	Reports []ReportModel `json:"reports"`
	// Returned when the number of results exceeds pageSize. To get the next page of results, call getReports with this token as the only parameter.
	NextToken *string `json:"nextToken,omitempty"`
}

// ReportDocument Information required for the report document.
type ReportDocument struct {
	// The identifier for the report document. This identifier is unique only in combination with a seller ID.
	ReportDocumentId string `json:"reportDocumentId"`
	// A presigned URL for the report document. This URL expires after 5 minutes.
	Url string `json:"url"`
	// If present, the report document contents have been compressed with the provided algorithm.
	CompressionAlgorithm *string `json:"compressionAlgorithm,omitempty"`
}

// ReportSchedule Detailed information about a report schedule.
type ReportSchedule struct {
	// The identifier for the report schedule. This identifier is unique only in combination with a seller ID.
	ReportScheduleId string `json:"reportScheduleId"`
	// The report type.
	ReportType string `json:"reportType"`
	// A list of marketplace identifiers. The report document's contents will contain data for all of the specified marketplaces, unless the report type indicates otherwise.
	MarketplaceIds []string `json:"marketplaceIds,omitempty"`
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// An ISO 8601 period value that indicates how often a report should be created.
	Period string `json:"period"`
	// The date and time when the schedule will create its next report, in ISO 8601 date time format.
	NextReportCreationTime *apis.JsonTimeISO8601 `json:"nextReportCreationTime,omitempty"`
}

// ReportScheduleList A list of report schedules.
type ReportScheduleList struct {
	ReportSchedules []ReportSchedule `json:"reportSchedules"`
}

// CreateReportScheduleResponse Response schema.
type CreateReportScheduleResponse struct {
	// The identifier for the report schedule. This identifier is unique only in combination with a seller ID.
	ReportScheduleId string `json:"reportScheduleId"`
}

// CreateReportScheduleSpecification struct for CreateReportScheduleSpecification
type CreateReportScheduleSpecification struct {
	// The report type.
	ReportType string `json:"reportType"`
	// A list of marketplace identifiers for the report schedule.
	MarketplaceIds []string `json:"marketplaceIds"`
	// Additional information passed to reports. This varies by report type.
	ReportOptions *map[string]string `json:"reportOptions,omitempty"`
	// One of a set of predefined ISO 8601 periods that specifies how often a report should be created.
	Period string `json:"period"`
	// The date and time when the schedule will create its next report, in ISO 8601 date time format.
	NextReportCreationTime *apis.JsonTimeISO8601 `json:"nextReportCreationTime,omitempty"`
}

// Error response returned when the request is unsuccessful.
type Error struct {
	// An error code that identifies the type of error that occurred.
	Code string `json:"code"`
	// A message that describes the error condition in a human-readable form.
	Message string `json:"message"`
	// Additional details that can help the caller understand or fix the issue.
	Details *string `json:"details,omitempty"`
}

// ErrorList A list of error responses returned when a request is unsuccessful.
type ErrorList struct {
	Errors []Error `json:"errors"`
}