package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
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
	Flaws                []DetailedReportFlaw         `xml:"severity>category>cwe>staticflaws>flaw"`
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

type DetailedReportFlaw struct {
	XMLName                 xml.Name `xml:"flaw"`
	ID                      int      `xml:"issueid,attr"`
	CWE                     int      `xml:"cweid,attr"`
	AffectsPolicyCompliance bool     `xml:"affects_policy_compliance,attr"`
	Module                  string   `xml:"module,attr"`
	RemediationStatus       string   `xml:"remediation_status,attr"`
	MitigationStatus        string   `xml:"mitigation_status,attr"`
	SourceFile              string   `xml:"source_file,attr"`
	LineNumber              int      `xml:"line,attr"`
	ProcedureHash           string   `xml:"procedure_hash,attr"`
	PrototypeHash           string   `xml:"prototype_hash,attr"`
	StatementHash           string   `xml:"statement_hash,attr"`
}

func (api API) getDetailedReport(buildId int) DetailedReport {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/detailedreport.do?build_id=%d", buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	if strings.Contains(string(response[:]), "<error>No report available.</error>") {
		color.Red(fmt.Sprintf("Error: There was no detailed report for Build id %d. Has the scan finished?", buildId))
		os.Exit(1)
	}

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

func (report DetailedReport) getPolicyAffectingFlawCount() int {
	var count = 0

	for _, flaw := range report.Flaws {
		if flaw.AffectsPolicyCompliance {
			count++
		}
	}

	return count
}

func (flaw DetailedReportFlaw) isFlawOpen() bool {
	if flaw.RemediationStatus == "Fixed" {
		return false
	}

	if !(flaw.MitigationStatus == "none" || flaw.MitigationStatus == "rejected") {
		return false
	}

	return true
}

func (report DetailedReport) getOpenPolicyAffectingFlawCount() int {
	var count = 0

	for _, flaw := range report.Flaws {
		if flaw.isFlawOpen() && flaw.AffectsPolicyCompliance {
			count++
		}
	}

	return count
}

func (report DetailedReport) getOpenNonPolicyAffectingFlawCount() int {
	var count = 0

	for _, flaw := range report.Flaws {

		if flaw.isFlawOpen() && !flaw.AffectsPolicyCompliance {
			count++
		}
	}

	return count
}
