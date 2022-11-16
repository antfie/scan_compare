package main

import (
	"flag"
	"fmt"
	"strings"
)

type API struct {
	id  string
	key string
}

func main() {
	print("Scan Compare v1.0\nCopyright © Veracode, Inc. 2022. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warrenty.\n\n")
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	scanA := flag.String("a", "", "Veracode Platform URL for scan 'A'")
	scanB := flag.String("b", "", "Veracode Platform URL for scan 'B'")

	flag.Parse()

	var apiId, apiKey = getCredentials(*vid, *vkey)
	var api = API{apiId, apiKey}

	if len(*scanA) < 1 {
		panic("No Veracode Platform URL specified for scan 'A'. Expected flag '-a https://analysiscenter.veracode.com/auth/index.jsp...'. Try -h for help.")
	}

	if len(*scanB) < 1 {
		panic("No Veracode Platform URL specified for scan 'B'. Expected flag '-b https://analysiscenter.veracode.com/auth/index.jsp...'. Try -h for help.")
	}

	scanAAppId := parseAppIdFromPlatformUrl(*scanA)
	scanABuildId := parseBuildIdFromPlatformUrl(*scanA)
	scanBAppId := parseAppIdFromPlatformUrl(*scanB)
	scanBBuildId := parseBuildIdFromPlatformUrl(*scanB)

	if scanABuildId == scanBBuildId {
		panic("These are the same scans")
	}

	fmt.Printf("Comparing scan 'A' (Build id = %d) against scan 'B' (Build id = %d)\n", scanABuildId, scanBBuildId)
	data := api.getData(scanAAppId, scanABuildId, scanBAppId, scanBBuildId)

	data.reportOnWarnings(*scanA, *scanB)
	data.reportCommonalities()
	data.reportScanADetails()
	data.reportScanBDetails()
	data.reportTopLevelModuleDifferences()
	// data.reportNotSelectedModuleDifferences()
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
		report.WriteString("The scan engine versions are different. This means there has been one or more deployments to the Veracode scan engine between these scans\n")
	}

	if report.Len() > 0 {
		fmt.Println("Warnings")
		fmt.Println("========")
		fmt.Println(report.String())
	}
}

func (data Data) reportCommonalities() {
	var report strings.Builder

	if data.ScanAReport.AppName == data.ScanBReport.AppName {
		report.WriteString(fmt.Sprintf("Application: %s\n", data.ScanAReport.AppName))
	}

	if data.ScanAReport.SandboxId == data.ScanBReport.SandboxId {
		report.WriteString(fmt.Sprintf("Sandbox: %s\n", data.ScanAReport.SandboxName))
	}

	if data.ScanAReport.TotalFlaws == data.ScanBReport.TotalFlaws && data.ScanAReport.UnmitigatedFlaws == data.ScanBReport.UnmitigatedFlaws {
		report.WriteString(fmt.Sprintf("Flaws: %d total, %d not mitigated\n", data.ScanAReport.TotalFlaws, data.ScanAReport.UnmitigatedFlaws))
	}

	if len(data.ScanAPrescanFileList.Files) == len(data.ScanBPrescanFileList.Files) {
		report.WriteString(fmt.Sprintf("Files uploaded: %d\n", len(data.ScanAPrescanFileList.Files)))
	}

	if len(data.ScanAReport.StaticAnalysis.Modules) == len(data.ScanBReport.StaticAnalysis.Modules) {
		report.WriteString(fmt.Sprintf("Top-level modules selected for analysis: %d", len(data.ScanAReport.StaticAnalysis.Modules)))
	}

	if report.Len() > 0 {
		fmt.Println("\nIn common with both scans")
		fmt.Println("=========================")
		fmt.Println(report.String())
	}
}

func (data Data) reportScanADetails() {
	fmt.Println("\nScan 'A'")
	fmt.Println("========")

	if data.ScanAReport.StaticAnalysis.EngineVersion != data.ScanBReport.StaticAnalysis.EngineVersion {
		fmt.Printf("Engine version: %s\n", data.ScanAReport.StaticAnalysis.EngineVersion)
	}

	if data.ScanAReport.AppName != data.ScanBReport.AppName {
		fmt.Printf("Application: '%s'\n", data.ScanAReport.AppName)
	}

	if data.ScanAReport.SandboxId != data.ScanBReport.SandboxId {
		fmt.Printf("Sandbox: '%s'\n", data.ScanBReport.SandboxName)
	}

	fmt.Printf("Scan Name: '%s'\n", data.ScanAReport.StaticAnalysis.ScanName)
	fmt.Printf("Review Modules URL: %s\n", data.ScanAReport.getReviewModulesUrl())

	if len(data.ScanAPrescanFileList.Files) != len(data.ScanBPrescanFileList.Files) {
		fmt.Printf("Files uploaded: %d\n", len(data.ScanAPrescanFileList.Files))
	}

	if len(data.ScanAReport.StaticAnalysis.Modules) != len(data.ScanBReport.StaticAnalysis.Modules) {
		fmt.Printf("Top-level modules selected for analysis: %d\n", len(data.ScanAReport.StaticAnalysis.Modules))
	}

	fmt.Printf("Submitted: %s\n", parseVeracodeDate(data.ScanAReport.StaticAnalysis.SubmittedDate).Local())
	fmt.Printf("Duration: %s\n", parseVeracodeDate(data.ScanAReport.StaticAnalysis.PublishedDate).Sub(parseVeracodeDate(data.ScanAReport.StaticAnalysis.SubmittedDate)))

	if !(data.ScanAReport.TotalFlaws == data.ScanBReport.TotalFlaws && data.ScanAReport.UnmitigatedFlaws == data.ScanBReport.UnmitigatedFlaws) {
		fmt.Printf("Flaws: %d total, %d not mitigated\n", data.ScanAReport.TotalFlaws, data.ScanAReport.UnmitigatedFlaws)
	}
}

func (data Data) reportScanBDetails() {
	fmt.Println("\nScan 'B'")
	fmt.Println("========")

	if data.ScanAReport.StaticAnalysis.EngineVersion != data.ScanBReport.StaticAnalysis.EngineVersion {
		fmt.Printf("Engine version: '%s'\n", data.ScanBReport.StaticAnalysis.EngineVersion)
	}

	if data.ScanAReport.AppName != data.ScanBReport.AppName {
		fmt.Printf("Application: '%s'\n", data.ScanBReport.AppName)
	}

	if data.ScanAReport.SandboxId != data.ScanBReport.SandboxId {
		fmt.Printf("Sandbox: '%s'\n", data.ScanBReport.SandboxName)
	}

	fmt.Printf("Scan Name: '%s'\n", data.ScanBReport.StaticAnalysis.ScanName)
	fmt.Printf("Review Modules URL: %s\n", data.ScanBReport.getReviewModulesUrl())

	if len(data.ScanAPrescanFileList.Files) != len(data.ScanBPrescanFileList.Files) {
		fmt.Printf("Files uploaded: %d\n", len(data.ScanBPrescanFileList.Files))
	}

	if len(data.ScanAReport.StaticAnalysis.Modules) != len(data.ScanBReport.StaticAnalysis.Modules) {
		fmt.Printf("Top-level modules selected for analysis: %d\n", len(data.ScanBReport.StaticAnalysis.Modules))
	}

	fmt.Printf("Submitted: %s\n", parseVeracodeDate(data.ScanBReport.StaticAnalysis.SubmittedDate).Local())
	fmt.Printf("Duration: %s\n", parseVeracodeDate(data.ScanBReport.StaticAnalysis.PublishedDate).Sub(parseVeracodeDate(data.ScanBReport.StaticAnalysis.SubmittedDate)))

	if !(data.ScanAReport.TotalFlaws == data.ScanBReport.TotalFlaws && data.ScanAReport.UnmitigatedFlaws == data.ScanBReport.UnmitigatedFlaws) {
		fmt.Printf("Flaws: %d total, %d not mitigated\n", data.ScanBReport.TotalFlaws, data.ScanBReport.UnmitigatedFlaws)
	}
}

func (data Data) reportTopLevelModuleDifferences() {
	var report strings.Builder

	compareTopLevelSelectedModules(&report, "A", data.ScanAReport.StaticAnalysis.Modules, data.ScanBReport.StaticAnalysis.Modules, data.ScanAPrescanFileList, data.ScanAPrescanModuleList)
	compareTopLevelSelectedModules(&report, "B", data.ScanBReport.StaticAnalysis.Modules, data.ScanAReport.StaticAnalysis.Modules, data.ScanBPrescanFileList, data.ScanBPrescanModuleList)

	if report.Len() > 0 {
		fmt.Println("\nDifferences of Top-Level Modules Selected As An Entry Point")
		fmt.Println("===========================================================")
		fmt.Println(report.String())
	}
}

func compareTopLevelSelectedModules(report *strings.Builder, side string, modulesInThisSide, modulesInTheOtherSide []SummaryReportModule, thisSidePrescanFileList PrescanFileList, thisSidePrescanModuleList PrescanModuleList) {
	for _, moduleFoundInThisSide := range modulesInThisSide {
		if !isModuleNameInArray(moduleFoundInThisSide, modulesInTheOtherSide) {
			prescanModule := thisSidePrescanModuleList.getFromName(moduleFoundInThisSide.Name)
			report.WriteString(fmt.Sprintf("Only in %s: '%s' - Size = %s, Issues = %d, MD5 = %s, Compiler = %s, OS = %s, Architecture = %s\n",
				side,
				moduleFoundInThisSide.Name,
				prescanModule.Size,
				len(prescanModule.Issues),
				thisSidePrescanFileList.getFromName(moduleFoundInThisSide.Name).MD5,
				moduleFoundInThisSide.Compiler,
				moduleFoundInThisSide.Os,
				moduleFoundInThisSide.Architecture))
		}
	}
}

func isModuleNameInArray(module SummaryReportModule, modules []SummaryReportModule) bool {
	for _, moduleInList := range modules {
		if module.Name == moduleInList.Name {
			return true
		}
	}

	return false
}

func (data Data) reportNotSelectedModuleDifferences() {
	var report strings.Builder

	compareTopLevelSelectedModules(&report, "A", data.ScanAReport.StaticAnalysis.Modules, data.ScanBReport.StaticAnalysis.Modules, data.ScanAPrescanFileList, data.ScanAPrescanModuleList)
	compareTopLevelSelectedModules(&report, "B", data.ScanBReport.StaticAnalysis.Modules, data.ScanAReport.StaticAnalysis.Modules, data.ScanBPrescanFileList, data.ScanBPrescanModuleList)

	if report.Len() > 0 {
		fmt.Println("\nDifferences of Top-Level Modules Not Selected As An Entry Point")
		fmt.Println("===================================================================")
		fmt.Println(report.String())
	}
}
