package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type API struct {
	id  string
	key string
}

var AppVersion string = "DEV"

func colorPrintf(format string) {
	color.New().Printf(format)
}

func main() {
	fmt.Printf("Scan Compare %s\nCopyright © Veracode, Inc. 2022. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warrenty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	scanA := flag.String("a", "", "Veracode Platform URL for scan \"A\"")
	scanB := flag.String("b", "", "Veracode Platform URL for scan \"B\"")

	flag.Parse()

	if len(*scanA) < 1 && len(*scanB) < 1 {
		color.Red("Error: No Veracode Platform URLs specified for scans \"A\" and \"B\". Expected: \"scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp... -b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scanA) < 1 {
		color.Red("Error: No Veracode Platform URL specified for scan \"A\". Expected: \"scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scanB) < 1 {
		color.Red("Error: No Veracode Platform URL specified for scan \"B\". Expected flag \"-b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	var apiId, apiKey = getCredentials(*vid, *vkey)
	var api = API{apiId, apiKey}

	scanAAppId := parseAppIdFromPlatformUrl(*scanA)
	scanABuildId := parseBuildIdFromPlatformUrl(*scanA)
	scanBAppId := parseAppIdFromPlatformUrl(*scanB)
	scanBBuildId := parseBuildIdFromPlatformUrl(*scanB)

	if scanABuildId == scanBBuildId {
		panic("These are the same scans")
	}

	colorPrintf(fmt.Sprintf("Comparing scan %s against scan %s\n",
		color.GreenString("\"A\" (Build id = %d)", scanABuildId),
		color.MagentaString("\"B\" (Build id = %d)", scanBBuildId)))

	data := api.getData(scanAAppId, scanABuildId, scanBAppId, scanBBuildId)

	data.reportOnWarnings(*scanA, *scanB)
	data.reportCommonalities()
	reportScanDetails("A", data.ScanAReport, data.ScanBReport, data.ScanAPrescanFileList, data.ScanBPrescanFileList)
	reportScanDetails("B", data.ScanBReport, data.ScanAReport, data.ScanBPrescanFileList, data.ScanAPrescanFileList)
	data.reportTopLevelModuleDifferences()
	data.reportNotSelectedModuleDifferences()
	data.reportDependencyModuleDifferences()
	data.reportSummary()
}

func (data Data) reportOnWarnings(scanAUrl, scanBUrl string) {
	var report strings.Builder

	if isPlatformURL(scanAUrl) && isPlatformURL(scanBUrl) {
		if parseAccountIdFromPlatformUrl(scanAUrl) != parseAccountIdFromPlatformUrl(scanBUrl) {
			report.WriteString("These scans are from different accounts\n")
		} else if parseAppIdFromPlatformUrl(scanAUrl) != parseAppIdFromPlatformUrl(scanBUrl) {
			report.WriteString("These scans are from different application profiles\n")
		}
	}

	if data.ScanAReport.StaticAnalysis.EngineVersion != data.ScanBReport.StaticAnalysis.EngineVersion {
		report.WriteString("The scan engine versions are different. This means there has been one or more deployments of the Veracode scan engine between these scans. This can sometimes explain why new flaws might be reported (due to improved scan coverage), and others are no longer reported (due to a reduction of Flase Positives)\n")
	}

	if report.Len() > 0 {
		color.Cyan("\nWarnings")
		fmt.Println("========")
		color.Yellow(report.String())
	}
}

func (data Data) reportCommonalities() {
	var report strings.Builder

	if data.ScanAReport.AppName == data.ScanBReport.AppName {
		report.WriteString(fmt.Sprintf("Application: \"%s\"\n", data.ScanAReport.AppName))
	}

	if data.ScanAReport.SandboxId == data.ScanBReport.SandboxId && len(data.ScanAReport.SandboxName) > 0 {
		report.WriteString(fmt.Sprintf("Sandbox: \"%s\"\n", data.ScanAReport.SandboxName))
	}

	if data.ScanAReport.StaticAnalysis.ScanName == data.ScanBReport.StaticAnalysis.ScanName {
		report.WriteString(fmt.Sprintf("Scan name: \"%s\"\n", data.ScanAReport.StaticAnalysis.ScanName))
	}

	if data.ScanAReport.TotalFlaws == data.ScanBReport.TotalFlaws && data.ScanAReport.UnmitigatedFlaws == data.ScanBReport.UnmitigatedFlaws {
		report.WriteString(fmt.Sprintf("Flaws: %d total, %d not mitigated\n", data.ScanAReport.TotalFlaws, data.ScanAReport.UnmitigatedFlaws))
	}

	if len(data.ScanAPrescanFileList.Files) == len(data.ScanBPrescanFileList.Files) {
		report.WriteString(fmt.Sprintf("Files uploaded: %d\n", len(data.ScanAPrescanFileList.Files)))
	}

	if len(data.ScanAReport.StaticAnalysis.Modules) == len(data.ScanBReport.StaticAnalysis.Modules) {
		report.WriteString(fmt.Sprintf("Top-level modules selected for analysis: %d\n", len(data.ScanAReport.StaticAnalysis.Modules)))
	}

	if data.ScanAReport.StaticAnalysis.EngineVersion == data.ScanBReport.StaticAnalysis.EngineVersion {
		report.WriteString(fmt.Sprintf("Engine version: %s\n", data.ScanAReport.StaticAnalysis.EngineVersion))
	}

	if report.Len() > 0 {
		color.Cyan("\nIn common with both scans")
		fmt.Println("=========================")
		colorPrintf(report.String())
	}
}

func reportScanDetails(side string, thisSumaryReport, otherSummaryReport SummaryReport, thisPrescanFileList, otherPrescanFileList PrescanFileList) {
	color.Magenta(fmt.Sprintf("\nScan %s", side))
	fmt.Println("======")

	if thisSumaryReport.AppName != otherSummaryReport.AppName {
		fmt.Printf("Application: \"%s\"\n", thisSumaryReport.AppName)
	}

	if thisSumaryReport.SandboxId != otherSummaryReport.SandboxId && len(thisSumaryReport.SandboxName) > 0 {
		fmt.Printf("Sandbox: \"%s\"\n", thisSumaryReport.SandboxName)
	}

	if thisSumaryReport.StaticAnalysis.ScanName != otherSummaryReport.StaticAnalysis.ScanName {
		fmt.Printf("Scan name: \"%s\"\n", thisSumaryReport.StaticAnalysis.ScanName)
	}

	fmt.Printf("Review Modules URL: %s\n", thisSumaryReport.getReviewModulesUrl())

	if thisSumaryReport.StaticAnalysis.EngineVersion != otherSummaryReport.StaticAnalysis.EngineVersion {
		fmt.Printf("Engine version: \"%s\"\n", thisSumaryReport.StaticAnalysis.EngineVersion)
	}

	if len(thisPrescanFileList.Files) != len(otherPrescanFileList.Files) {
		fmt.Printf("Files uploaded: %d\n", len(thisPrescanFileList.Files))
	}

	if len(thisSumaryReport.StaticAnalysis.Modules) != len(otherSummaryReport.StaticAnalysis.Modules) {
		fmt.Printf("Top-level modules selected for analysis: %d\n", len(thisSumaryReport.StaticAnalysis.Modules))
	}

	fmt.Printf("Submitted: %s\n", thisSumaryReport.SubmittedDate)
	fmt.Printf("Duration: %s\n", thisSumaryReport.Duration)

	if thisSumaryReport.StaticAnalysis.EngineVersion != otherSummaryReport.StaticAnalysis.EngineVersion {
		fmt.Printf("Engine version: %s\n", thisSumaryReport.StaticAnalysis.EngineVersion)
	}

	if !(thisSumaryReport.TotalFlaws == otherSummaryReport.TotalFlaws && thisSumaryReport.UnmitigatedFlaws == otherSummaryReport.UnmitigatedFlaws) {
		fmt.Printf("Flaws: %d total, %d mitigated\n", thisSumaryReport.TotalFlaws, thisSumaryReport.TotalFlaws-thisSumaryReport.UnmitigatedFlaws)
	}
}

func (data Data) reportTopLevelModuleDifferences() {
	var report strings.Builder

	compareTopLevelSelectedModules(&report, "A", data.ScanAReport.StaticAnalysis.Modules, data.ScanBReport.StaticAnalysis.Modules, data.ScanAPrescanFileList, data.ScanAPrescanModuleList)
	compareTopLevelSelectedModules(&report, "B", data.ScanBReport.StaticAnalysis.Modules, data.ScanAReport.StaticAnalysis.Modules, data.ScanBPrescanFileList, data.ScanBPrescanModuleList)

	if report.Len() > 0 {
		color.Cyan("\nDifferences of Top-Level Modules Selected As An Entry Point For Scanning")
		fmt.Println("========================================================================")
		colorPrintf(report.String())
	}
}

func getFormattedOnlyInSideString(side string) string {
	if side == "A" {
		return color.GreenString("Only in A")
	}

	return color.MagentaString("Only in B")
}

func getFormattedSideString(side string) string {
	if side == "A" {
		return color.GreenString("A")
	}

	return color.MagentaString("B")
}

func getMissingSupportedFileCountFromPreScanModuleStatus(module PrescanModule) int {
	for _, issue := range strings.Split(module.Status, ",") {
		if strings.HasPrefix(issue, "Missing Supporting Files") {
			trimmedPrefix := strings.Replace(issue, "Missing Supporting Files - ", "", 1)
			count, err := strconv.Atoi(strings.Split(trimmedPrefix, " ")[0])

			if err != nil {
				return 0
			}

			return count
		}
	}

	return 0
}

func getFatalReason(module PrescanModule) string {
	for _, issue := range strings.Split(module.Status, ",") {
		if strings.HasPrefix(issue, "(Fatal)") {
			return strings.Replace(issue, "(Fatal)", ": ", 1)
		}
	}

	return ""
}

func compareTopLevelSelectedModules(report *strings.Builder, side string, modulesInThisSideReport, modulesInTheOtherSideReport []SummaryReportModule, thisSidePrescanFileList PrescanFileList, thisSidePrescanModuleList PrescanModuleList) {
	for _, moduleFoundInThisSide := range modulesInThisSideReport {
		if !isModuleNameInSummaryReportModuleArray(moduleFoundInThisSide, modulesInTheOtherSideReport) {
			prescanModule := thisSidePrescanModuleList.getFromName(moduleFoundInThisSide.Name)
			var formattedSupportIssues = ""

			if len(prescanModule.Issues) > 0 {
				formattedSupportIssues = fmt.Sprintf(", %s", color.YellowString("Support issues = %d", len(prescanModule.Issues)))
			}

			var formattedMissingSupportedFiles = ""

			missingSupportedFileCount := getMissingSupportedFileCountFromPreScanModuleStatus(prescanModule)

			if missingSupportedFileCount > 1 {
				formattedMissingSupportedFiles = fmt.Sprintf(", %s", color.YellowString("Missing Supporting Files = %d", missingSupportedFileCount))
			}

			report.WriteString(fmt.Sprintf("%s: \"%s\" - Size = %s%s%s, MD5 = %s, Platform = %s / %s / %s\n",
				getFormattedOnlyInSideString(side),
				moduleFoundInThisSide.Name,
				prescanModule.Size,
				formattedSupportIssues,
				formattedMissingSupportedFiles,
				thisSidePrescanFileList.getFromName(moduleFoundInThisSide.Name).MD5,
				moduleFoundInThisSide.Architecture,
				moduleFoundInThisSide.Os,
				moduleFoundInThisSide.Compiler))
		}
	}
}

func isModuleNameInSummaryReportModuleArray(module SummaryReportModule, modules []SummaryReportModule) bool {
	for _, moduleInList := range modules {
		if module.Name == moduleInList.Name {
			return true
		}
	}

	return false
}

func (data Data) reportNotSelectedModuleDifferences() {
	var report strings.Builder

	compareTopLevelNotSelectedModules(&report, "A", data.ScanAPrescanModuleList, data.ScanBPrescanModuleList, data.ScanAReport.StaticAnalysis.Modules, false)
	compareTopLevelNotSelectedModules(&report, "B", data.ScanBPrescanModuleList, data.ScanAPrescanModuleList, data.ScanBReport.StaticAnalysis.Modules, false)

	if report.Len() > 0 {
		color.Cyan("\nDifferences of Top-Level Modules Not Selected As An Entry Point (And Not Scanned) - Unselected Potential First Party Components")
		fmt.Println("===============================================================================================================================")
		colorPrintf(report.String())
	}
}

func (data Data) reportDependencyModuleDifferences() {
	var report strings.Builder

	compareTopLevelNotSelectedModules(&report, "A", data.ScanAPrescanModuleList, data.ScanBPrescanModuleList, data.ScanAReport.StaticAnalysis.Modules, true)
	compareTopLevelNotSelectedModules(&report, "B", data.ScanBPrescanModuleList, data.ScanAPrescanModuleList, data.ScanBReport.StaticAnalysis.Modules, true)

	if report.Len() > 0 {
		color.Cyan("\nDifferences of Dependency Modules Not Selected As An Entry Point")
		fmt.Println("================================================================")
		colorPrintf(report.String())
	}
}

func isModuleNotSelectedTopLevel(prescanModuleFoundInThisSide PrescanModule, thisSideReportModuleList []SummaryReportModule, onlyDependencies bool) bool {
	if prescanModuleFoundInThisSide.IsDependency != onlyDependencies {
		return false
	}

	for _, summaryReportModule := range thisSideReportModuleList {
		if prescanModuleFoundInThisSide.Name == summaryReportModule.Name {
			return false
		}
	}

	return true
}

func compareTopLevelNotSelectedModules(report *strings.Builder, side string, prescanModulesInThisSide, prescanModulesInTheOtherSide PrescanModuleList, thisSideReportModuleList []SummaryReportModule, onlyDependencies bool) {
	for _, prescanModuleFoundInThisSide := range prescanModulesInThisSide.Modules {
		if !isModuleNotSelectedTopLevel(prescanModuleFoundInThisSide, thisSideReportModuleList, onlyDependencies) {
			continue
		}

		if prescanModulesInTheOtherSide.getFromName(prescanModuleFoundInThisSide.Name).Name != prescanModuleFoundInThisSide.Name {
			var formattedSupportIssues = ""

			if len(prescanModuleFoundInThisSide.Issues) > 0 {
				formattedSupportIssues = fmt.Sprintf(", %s", color.YellowString("Support issues = %d", len(prescanModuleFoundInThisSide.Issues)))
			}

			var formattedFatalError = ""

			if prescanModuleFoundInThisSide.HasFatalErrors {
				formattedFatalError = fmt.Sprintf(", %s", color.RedString(fmt.Sprintf("Unscannable%s", getFatalReason(prescanModuleFoundInThisSide))))
			}

			var formattedMissingSupportedFiles = ""

			missingSupportedFileCount := getMissingSupportedFileCountFromPreScanModuleStatus(prescanModuleFoundInThisSide)

			if missingSupportedFileCount > 1 {
				formattedMissingSupportedFiles = fmt.Sprintf(", %s", color.YellowString("Missing Supporting Files = %d", missingSupportedFileCount))
			}

			report.WriteString(fmt.Sprintf("%s: \"%s\" - Size = %s%s%s%s, MD5 = %s, Platform = %s\n",
				getFormattedOnlyInSideString(side),
				prescanModuleFoundInThisSide.Name,
				prescanModuleFoundInThisSide.Size,
				formattedSupportIssues,
				formattedFatalError,
				formattedMissingSupportedFiles,
				prescanModuleFoundInThisSide.MD5,
				prescanModuleFoundInThisSide.Platform))
		}
	}
}

func (data Data) reportSummary() {
	var report strings.Builder

	if data.ScanAReport.SubmittedDate.Before(data.ScanBReport.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", getFormattedSideString("B"), data.ScanBReport.SubmittedDate.Sub(data.ScanAReport.SubmittedDate), getFormattedSideString("A")))
	} else if data.ScanAReport.SubmittedDate.After(data.ScanBReport.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", getFormattedSideString("A"), data.ScanAReport.SubmittedDate.Sub(data.ScanBReport.SubmittedDate), getFormattedSideString("B")))
	}

	if data.ScanAReport.Duration > data.ScanBReport.Duration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", getFormattedSideString("A"), data.ScanAReport.Duration-data.ScanBReport.Duration))
	} else if data.ScanAReport.Duration < data.ScanBReport.Duration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", getFormattedSideString("B"), data.ScanBReport.Duration-data.ScanAReport.Duration))
	}

	if report.Len() > 0 {
		color.Cyan("\nSummary")
		fmt.Print("========\n")
		colorPrintf(report.String())
	}
}
