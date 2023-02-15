package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	fmt.Printf("Scan Compare v%s\nCopyright © Veracode, Inc. 2023. All Rights Reserved.\nThis is an unofficial Veracode product. It does not come with any support or warrenty.\n\n", AppVersion)
	vid := flag.String("vid", "", "Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds")
	vkey := flag.String("vkey", "", "Veracode API key - See https://docs.veracode.com/r/t_create_api_creds")
	profile := flag.String("profile", "default", "Veracode credential profile - See https://docs.veracode.com/r/c_httpie_tool")
	region := flag.String("region", "", "Veracode Region [global, us, eu]")
	scanA := flag.String("a", "", "Veracode Platform URL or build ID for scan \"A\"")
	scanB := flag.String("b", "", "Veracode Platform URL or build ID for scan \"B\"")

	flag.Parse()

	if !(*region == "" || *region == "global" || *region == "us" || *region == "eu") {
		color.HiRed("Error: Invalid region. Must be either \"global\", \"us\" or \"eu\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scanA) < 1 && len(*scanB) < 1 {
		color.HiRed("Error: No Veracode Platform URLs or build IDs specified for scans \"A\" and \"B\". Expected: \"scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp... -b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scanA) < 1 {
		color.HiRed("Error: No Veracode Platform URL or build ID specified for scan \"A\". Expected: \"scan_compare -a https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if len(*scanB) < 1 {
		color.HiRed("Error: No Veracode Platform URL or build ID specified for scan \"B\". Expected flag \"-b https://analysiscenter.veracode.com/auth/index.jsp...\"")
		print("\nUsage:\n")
		flag.PrintDefaults()
		return
	}

	if parseRegionFromUrl(*scanA) != parseRegionFromUrl(*scanB) {
		color.HiRed("Error: Cannot compare between different Veracode regions")
		os.Exit(1)
	}

	if *region != "" &&
		((strings.HasPrefix(*scanA, "https://") && parseRegionFromUrl(*scanA) != *region) ||
			(strings.HasPrefix(*scanB, "https://") && parseRegionFromUrl(*scanB) != *region)) {
		color.HiRed(fmt.Sprintf("Error: The region from the URL (%s) does not match that specified by the command line (%s)", parseRegionFromUrl(*scanA), *region))
		os.Exit(1)
	}

	var regionToUse string

	// Command line region takes precidence
	if *region == "" {
		regionToUse = parseRegionFromUrl(*scanA)
	} else {
		regionToUse = *region
	}

	notifyOfUpdates()

	var apiId, apiKey = getCredentials(*vid, *vkey, *profile)
	var api = API{apiId, apiKey, regionToUse}

	scanAAppId := parseAppIdFromPlatformUrl(*scanA)
	scanABuildId := parseBuildIdFromPlatformUrl(*scanA)
	scanBAppId := parseAppIdFromPlatformUrl(*scanB)
	scanBBuildId := parseBuildIdFromPlatformUrl(*scanB)

	if scanABuildId == scanBBuildId {
		color.HiRed("Error: These are both the same scan")
		os.Exit(1)
	}

	api.assertCredentialsWork()

	colorPrintf(fmt.Sprintf("Comparing scan %s against scan %s in the %s region\n",
		color.HiGreenString("\"A\" (Build id = %d)", scanABuildId),
		color.HiMagentaString("\"B\" (Build id = %d)", scanBBuildId),
		api.region))

	data := api.getData(scanAAppId, scanABuildId, scanBAppId, scanBBuildId)

	data.reportOnWarnings(*scanA, *scanB)
	data.assertPrescanModulesPresent()
	data.reportCommonalities()
	reportScanDetails(api.region, "A", data.ScanAReport, data.ScanBReport, data.ScanAPrescanFileList, data.ScanBPrescanFileList, data.ScanAPrescanModuleList, data.ScanBPrescanModuleList)
	reportScanDetails(api.region, "B", data.ScanBReport, data.ScanAReport, data.ScanBPrescanFileList, data.ScanAPrescanFileList, data.ScanBPrescanModuleList, data.ScanAPrescanModuleList)
	data.reportTopLevelModuleDifferences()
	data.reportNotSelectedModuleDifferences()
	data.reportDependencyModuleDifferences()
	reportDuplicateFiles("A", data.ScanAPrescanFileList)
	reportDuplicateFiles("B", data.ScanBPrescanFileList)
	data.reportModuleDifferences()
	data.reportFlawDifferences()
	data.reportSummary()
}

func (data Data) reportOnWarnings(scanAUrl, scanBUrl string) {
	var report strings.Builder

	if isPlatformURL(scanAUrl) && isPlatformURL(scanBUrl) {
		if parseAccountIdFromPlatformUrl(scanAUrl) != parseAccountIdFromPlatformUrl(scanBUrl) {
			report.WriteString("* These scans are from different accounts\n")
		} else if parseAppIdFromPlatformUrl(scanAUrl) != parseAppIdFromPlatformUrl(scanBUrl) {
			report.WriteString("* These scans are from different application profiles\n")
		}
	}

	if data.ScanAReport.StaticAnalysis.EngineVersion != data.ScanBReport.StaticAnalysis.EngineVersion {
		report.WriteString("* The scan engine versions are different. This means there has been one or more deployments of the Veracode scan engine between these scans. This can sometimes explain why new flaws might be reported (due to improved scan coverage), and others are no longer reported (due to a reduction of False Positives)\n")
	}

	if time.Since(data.ScanAReport.SubmittedDate).Hours() >= 30*24 && time.Since(data.ScanBReport.SubmittedDate).Hours() >= 30*24 {
		report.WriteString("* Both scans are older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	} else if time.Since(data.ScanAReport.SubmittedDate).Hours() >= 30*24 {
		report.WriteString("* Scan A is older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	} else if time.Since(data.ScanBReport.SubmittedDate).Hours() >= 30*24 {
		report.WriteString("* Scan B is older than 30 days. This means the files will have been deleted and Veracode support therefore require a newer scan to investigate any issues further.\n")
	}

	if report.Len() > 0 {
		printTitle("Warnings")
		color.HiYellow(report.String())
	}
}

func (data Data) reportCommonalities() {
	var report strings.Builder

	if data.ScanAReport.AppName == data.ScanBReport.AppName {
		report.WriteString(fmt.Sprintf("Application:        \"%s\"\n", data.ScanAReport.AppName))
	}

	if data.ScanAReport.SandboxId == data.ScanBReport.SandboxId && len(data.ScanAReport.SandboxName) > 0 {
		report.WriteString(fmt.Sprintf("Sandbox:            \"%s\"\n", data.ScanAReport.SandboxName))
	}

	if data.ScanAReport.StaticAnalysis.ScanName == data.ScanBReport.StaticAnalysis.ScanName {
		report.WriteString(fmt.Sprintf("Scan name:          \"%s\"\n", data.ScanAReport.StaticAnalysis.ScanName))
	}

	if len(data.ScanAPrescanFileList.Files) == len(data.ScanBPrescanFileList.Files) {
		report.WriteString(fmt.Sprintf("Files uploaded:     %d\n", len(data.ScanAPrescanFileList.Files)))
	}

	if len(data.ScanAPrescanModuleList.Modules) == len(data.ScanBPrescanModuleList.Modules) {
		report.WriteString(fmt.Sprintf("Total modules:      %d\n", len(data.ScanAPrescanModuleList.Modules)))
	}

	if len(data.ScanAReport.StaticAnalysis.Modules) == len(data.ScanBReport.StaticAnalysis.Modules) {
		report.WriteString(fmt.Sprintf("Modules selected:   %d\n", len(data.ScanAReport.StaticAnalysis.Modules)))
	}

	if data.ScanAReport.StaticAnalysis.EngineVersion == data.ScanBReport.StaticAnalysis.EngineVersion {
		report.WriteString(fmt.Sprintf("Engine version:     %s\n", data.ScanAReport.StaticAnalysis.EngineVersion))
	}

	if data.ScanAReport.TotalFlaws == data.ScanBReport.TotalFlaws && data.ScanAReport.UnmitigatedFlaws == data.ScanBReport.UnmitigatedFlaws && data.ScanAReport.getPolicyAffectingFlawCount() == data.ScanBReport.getPolicyAffectingFlawCount() && data.ScanAReport.getOpenPolicyAffectingFlawCount() == data.ScanBReport.getOpenPolicyAffectingFlawCount() && data.ScanAReport.getOpenNonPolicyAffectingFlawCount() == data.ScanBReport.getOpenNonPolicyAffectingFlawCount() {
		flawsFormatted := fmt.Sprintf("Flaws:              %d total, %d mitigated, %d policy affecting, %d open affecting policy, %d open not affecting policy\n", data.ScanAReport.TotalFlaws, data.ScanAReport.TotalFlaws-data.ScanAReport.UnmitigatedFlaws, data.ScanAReport.getPolicyAffectingFlawCount(), data.ScanAReport.getOpenPolicyAffectingFlawCount(), data.ScanAReport.getOpenNonPolicyAffectingFlawCount())

		if data.ScanAReport.TotalFlaws == 0 {
			report.WriteString(color.HiYellowString(flawsFormatted))
		} else {
			report.WriteString(flawsFormatted)
		}
	}

	if report.Len() > 0 {
		printTitle("In common with both scans")
		colorPrintf(report.String())
	}
}

func reportScanDetails(region, side string, thisDetailedReport, otherDetailedReport DetailedReport, thisPrescanFileList, otherPrescanFileList PrescanFileList, thisPrescanModuleList, otherPrescanModuleList PrescanModuleList) {
	colorPrintf(getFormattedSideStringWithMessage(side, fmt.Sprintf("\nScan %s", side)))
	fmt.Println("\n======")

	if thisDetailedReport.AccountId != otherDetailedReport.AccountId {
		fmt.Printf("Account ID:         %d\n", thisDetailedReport.AccountId)
	}

	if thisDetailedReport.AppName != otherDetailedReport.AppName {
		fmt.Printf("Application:        \"%s\"\n", thisDetailedReport.AppName)
	}

	if thisDetailedReport.SandboxId != otherDetailedReport.SandboxId && len(thisDetailedReport.SandboxName) > 0 {
		fmt.Printf("Sandbox:            \"%s\"\n", thisDetailedReport.SandboxName)
	}

	if thisDetailedReport.StaticAnalysis.ScanName != otherDetailedReport.StaticAnalysis.ScanName {
		fmt.Printf("Scan name:          \"%s\"\n", thisDetailedReport.StaticAnalysis.ScanName)
	}

	fmt.Printf("Review Modules URL: %s\n", thisDetailedReport.getReviewModulesUrl(region))
	fmt.Printf("Triage Flaws URL:   %s\n", thisDetailedReport.getTriageFlawsUrl(region))

	if len(thisPrescanFileList.Files) != len(otherPrescanFileList.Files) {
		fmt.Printf("Files uploaded:     %d\n", len(thisPrescanFileList.Files))
	}

	if len(thisPrescanModuleList.Modules) != len(otherPrescanModuleList.Modules) {
		fmt.Printf("Total modules:      %d\n", len(thisPrescanModuleList.Modules))
	}

	if len(thisDetailedReport.StaticAnalysis.Modules) != len(otherDetailedReport.StaticAnalysis.Modules) {
		fmt.Printf("Modules selected:   %d\n", len(thisDetailedReport.StaticAnalysis.Modules))
	}

	if thisDetailedReport.StaticAnalysis.EngineVersion != otherDetailedReport.StaticAnalysis.EngineVersion {
		fmt.Printf("Engine version:     %s\n", thisDetailedReport.StaticAnalysis.EngineVersion)
	}

	fmt.Printf("Submitted:          %s (%s ago)\n", thisDetailedReport.SubmittedDate, formatDuration(time.Since(thisDetailedReport.SubmittedDate)))
	fmt.Printf("Published:          %s (%s ago)\n", thisDetailedReport.PublishedDate, formatDuration(time.Since(thisDetailedReport.PublishedDate)))
	fmt.Printf("Duration:           %s\n", thisDetailedReport.Duration)

	if !(thisDetailedReport.TotalFlaws == otherDetailedReport.TotalFlaws && thisDetailedReport.UnmitigatedFlaws == otherDetailedReport.UnmitigatedFlaws && thisDetailedReport.getPolicyAffectingFlawCount() == otherDetailedReport.getPolicyAffectingFlawCount() && thisDetailedReport.getOpenNonPolicyAffectingFlawCount() == otherDetailedReport.getOpenNonPolicyAffectingFlawCount()) {
		flawsFormatted := fmt.Sprintf("Flaws:              %d total, %d mitigated, %d policy affecting, %d open affecting policy, %d open not affecting policy\n", thisDetailedReport.TotalFlaws, thisDetailedReport.TotalFlaws-thisDetailedReport.UnmitigatedFlaws, thisDetailedReport.getPolicyAffectingFlawCount(), thisDetailedReport.getOpenPolicyAffectingFlawCount(), thisDetailedReport.getOpenNonPolicyAffectingFlawCount())

		if thisDetailedReport.TotalFlaws == 0 {
			color.HiYellow(flawsFormatted)
		} else {
			fmt.Print(flawsFormatted)
		}
	}
}

func (data Data) assertPrescanModulesPresent() {
	if len(data.ScanAPrescanModuleList.Modules) == 0 && len(data.ScanBPrescanModuleList.Modules) == 0 {
		color.HiRed("Error: Could not retrieve pre-scan modules for either scan")
		os.Exit(1)
	}

	if len(data.ScanAPrescanModuleList.Modules) == 0 {
		color.HiRed("Error: Could not retrieve pre-scan modules for scan A")
		os.Exit(1)
	}

	if len(data.ScanBPrescanModuleList.Modules) == 0 {
		color.HiRed("Error: Could not retrieve pre-scan modules for scan B")
		os.Exit(1)
	}
}

func (data Data) reportSummary() {
	var report strings.Builder

	if data.ScanAReport.SubmittedDate.Before(data.ScanBReport.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", getFormattedSideString("B"), formatDuration(data.ScanBReport.SubmittedDate.Sub(data.ScanAReport.SubmittedDate)), getFormattedSideString("A")))
	} else if data.ScanAReport.SubmittedDate.After(data.ScanBReport.SubmittedDate) {
		report.WriteString(fmt.Sprintf("%s was submitted %s after %s\n", getFormattedSideString("A"), formatDuration(data.ScanAReport.SubmittedDate.Sub(data.ScanBReport.SubmittedDate)), getFormattedSideString("B")))
	}

	if data.ScanAReport.Duration > data.ScanBReport.Duration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", getFormattedSideString("A"), formatDuration(data.ScanAReport.Duration-data.ScanBReport.Duration)))
	} else if data.ScanAReport.Duration < data.ScanBReport.Duration {
		report.WriteString(fmt.Sprintf("%s took longer by %s\n", getFormattedSideString("B"), formatDuration(data.ScanBReport.Duration-data.ScanAReport.Duration)))
	}

	if report.Len() > 0 {
		printTitle("Summary")
		colorPrintf(report.String())
	}
}
