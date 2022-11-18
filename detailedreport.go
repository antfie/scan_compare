package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type DetailedReport struct {
	XMLName              xml.Name                     `xml:"detailedreport"`
	AccountId            int                          `xml:"account_id,attr"`
	AppId                int                          `xml:"app_id,attr"`
	AppName              string                       `xml:"app_name,attr"`
	SandboxId            int                          `xml:"sandbox_id,attr"`
	SandboxName          string                       `xml:"sandbox_name,attr"`
	BuildId              int                          `xml:"build_id,attr"`
	AnalysisId           int                          `xml:"analysis_id,attr"`
	StaticAnalysisUnitId int                          `xml:"static_analysis_unit_id,attr"`
	TotalFlaws           int                          `xml:"total_flaws,attr"`
	UnmitigatedFlaws     int                          `xml:"flaws_not_mitigated,attr"`
	StaticAnalysis       DetailedReportStaticAnalysis `xml:"static-analysis"`
	SubmittedDate        time.Time
	Duration             time.Duration
}

type DetailedReportStaticAnalysis struct {
	XMLName           xml.Name               `xml:"static-analysis"`
	EngineVersion     string                 `xml:"engine_version,attr"`
	SubmittedDate     string                 `xml:"submitted_date,attr"`
	PublishedDate     string                 `xml:"published_date,attr"`
	ScanName          string                 `xml:"version,attr"`
	Score             int                    `xml:"score,attr"`
	AnalysisSizeBytes string                 `xml:"analysis_size_bytes,attr"`
	Modules           []DetailedReportModule `xml:"modules>module"`
}

type DetailedReportModule struct {
	XMLName      xml.Name `xml:"module"`
	Name         string   `xml:"name,attr"`
	Compiler     string   `xml:"compiler,attr"`
	Os           string   `xml:"os,attr"`
	Architecture string   `xml:"architecture,attr"`
}

func (api API) getDetailedReport(buildId int) DetailedReport {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/detailedreport.do?build_id=%d", buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	report := DetailedReport{}
	xml.Unmarshal(response, &report)

	// Sort modules by name for consistency
	sort.Slice(report.StaticAnalysis.Modules, func(i, j int) bool {
		return report.StaticAnalysis.Modules[i].Name < report.StaticAnalysis.Modules[j].Name
	})

	return report
}

func (report DetailedReport) getReviewModulesUrl() string {
	return fmt.Sprintf("https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:%d:%d:%d:%d:%d::::%d",
		report.AccountId,
		report.AppId,
		report.BuildId,
		report.AnalysisId,
		report.StaticAnalysisUnitId,
		report.SandboxId)
}
