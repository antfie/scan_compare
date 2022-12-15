package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func (data Data) reportTopLevelModuleDifferences() {
	var report strings.Builder

	compareTopLevelSelectedModules(&report, "A", data.ScanAReport.StaticAnalysis.Modules, data.ScanBReport.StaticAnalysis.Modules, data.ScanAPrescanFileList, data.ScanAPrescanModuleList)
	compareTopLevelSelectedModules(&report, "B", data.ScanBReport.StaticAnalysis.Modules, data.ScanAReport.StaticAnalysis.Modules, data.ScanBPrescanFileList, data.ScanBPrescanModuleList)

	if report.Len() > 0 {
		color.HiCyan("\nDifferences of Top-Level Modules Selected As An Entry Point For Scanning")
		fmt.Println("========================================================================")
		colorPrintf(report.String())
	}
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

func getFormattedModuleMD5(input string) string {
	if len(strings.TrimSpace(input)) > 0 {
		return fmt.Sprintf(", MD5 = %s", input)
	}

	return ""
}

func compareTopLevelSelectedModules(report *strings.Builder, side string, modulesInThisSideReport, modulesInTheOtherSideReport []DetailedReportModule, thisSidePrescanFileList PrescanFileList, thisSidePrescanModuleList PrescanModuleList) {
	for _, moduleFoundInThisSide := range modulesInThisSideReport {
		if !moduleFoundInThisSide.isModuleNameInDetailedReportModuleArray(modulesInTheOtherSideReport) {
			prescanModule := thisSidePrescanModuleList.getFromName(moduleFoundInThisSide.Name)
			var formattedSupportIssues = ""

			if len(prescanModule.Issues) > 0 {
				formattedSupportIssues = fmt.Sprintf(", %s", color.HiYellowString("Support issues = %d", len(prescanModule.Issues)))
			}

			var formattedMissingSupportedFiles = ""

			missingSupportedFileCount := getMissingSupportedFileCountFromPreScanModuleStatus(prescanModule)

			if missingSupportedFileCount > 1 {
				formattedMissingSupportedFiles = fmt.Sprintf(", %s", color.HiYellowString("Missing Supporting Files = %d", missingSupportedFileCount))
			}

			var formattedIsDependency = ""

			if prescanModule.IsDependency {
				formattedIsDependency = fmt.Sprintf(", %s", color.HiYellowString("Module is Dependency"))
			}

			report.WriteString(fmt.Sprintf("%s: \"%s\" - Size = %s%s%s%s%s, Platform = %s / %s / %s\n",
				getFormattedOnlyInSideString(side),
				moduleFoundInThisSide.Name,
				prescanModule.Size,
				formattedSupportIssues,
				formattedMissingSupportedFiles,
				formattedIsDependency,
				getFormattedModuleMD5(thisSidePrescanFileList.getFromName(moduleFoundInThisSide.Name).MD5),
				moduleFoundInThisSide.Architecture,
				moduleFoundInThisSide.Os,
				moduleFoundInThisSide.Compiler))
		}
	}
}

func (data Data) reportNotSelectedModuleDifferences() {
	var report strings.Builder

	compareTopLevelNotSelectedModules(&report, "A", data.ScanAPrescanModuleList, data.ScanBPrescanModuleList, data.ScanAReport.StaticAnalysis.Modules, false)
	compareTopLevelNotSelectedModules(&report, "B", data.ScanBPrescanModuleList, data.ScanAPrescanModuleList, data.ScanBReport.StaticAnalysis.Modules, false)

	if report.Len() > 0 {
		if strings.Contains(report.String(), "files extracted from") {
			color.HiCyan("\nDifferences of Top-Level Modules Which May or May Not Have Been Selected")
			fmt.Println("========================================================================")
		} else {
			color.HiCyan("\nDifferences of Top-Level Modules Not Selected As An Entry Point (And Not Scanned) - Unselected Potential First Party Components")
			fmt.Println("===============================================================================================================================")
		}

		colorPrintf(report.String())
	}
}

func (data Data) reportDependencyModuleDifferences() {
	var report strings.Builder

	compareTopLevelNotSelectedModules(&report, "A", data.ScanAPrescanModuleList, data.ScanBPrescanModuleList, data.ScanAReport.StaticAnalysis.Modules, true)
	compareTopLevelNotSelectedModules(&report, "B", data.ScanBPrescanModuleList, data.ScanAPrescanModuleList, data.ScanBReport.StaticAnalysis.Modules, true)

	if report.Len() > 0 {
		color.HiCyan("\nDifferences of Dependency Modules Not Selected As An Entry Point")
		fmt.Println("================================================================")
		colorPrintf(report.String())
	}
}

func isModuleNotSelectedTopLevel(prescanModuleFoundInThisSide PrescanModule, thisSideReportModuleList []DetailedReportModule, onlyDependencies bool) bool {
	if prescanModuleFoundInThisSide.IsDependency != onlyDependencies {
		return false
	}

	for _, DetailedReportModule := range thisSideReportModuleList {
		if prescanModuleFoundInThisSide.Name == DetailedReportModule.Name {
			return false
		}
	}

	return true
}

func compareTopLevelNotSelectedModules(report *strings.Builder, side string, prescanModulesInThisSide, prescanModulesInTheOtherSide PrescanModuleList, thisSideReportModuleList []DetailedReportModule, onlyDependencies bool) {
	for _, prescanModuleFoundInThisSide := range prescanModulesInThisSide.Modules {
		if !isModuleNotSelectedTopLevel(prescanModuleFoundInThisSide, thisSideReportModuleList, onlyDependencies) {
			continue
		}

		if prescanModulesInTheOtherSide.getFromName(prescanModuleFoundInThisSide.Name).Name != prescanModuleFoundInThisSide.Name {
			var formattedSupportIssues = ""

			if len(prescanModuleFoundInThisSide.Issues) > 0 {
				formattedSupportIssues = fmt.Sprintf(", %s", color.HiYellowString("Support issues = %d", len(prescanModuleFoundInThisSide.Issues)))
			}

			var formattedFatalError = ""

			if prescanModuleFoundInThisSide.HasFatalErrors {
				formattedFatalError = fmt.Sprintf(", %s", color.HiRedString(fmt.Sprintf("Unscannable%s", prescanModuleFoundInThisSide.getFatalReason())))
			}

			var formattedMissingSupportedFiles = ""

			missingSupportedFileCount := getMissingSupportedFileCountFromPreScanModuleStatus(prescanModuleFoundInThisSide)

			if missingSupportedFileCount > 1 {
				formattedMissingSupportedFiles = fmt.Sprintf(", %s", color.HiYellowString("Missing Supporting Files = %d", missingSupportedFileCount))
			}

			report.WriteString(fmt.Sprintf("%s: \"%s\" - Size = %s%s%s%s%s, Platform = %s\n",
				getFormattedOnlyInSideString(side),
				prescanModuleFoundInThisSide.Name,
				prescanModuleFoundInThisSide.Size,
				formattedSupportIssues,
				formattedFatalError,
				formattedMissingSupportedFiles,
				getFormattedModuleMD5(prescanModuleFoundInThisSide.MD5),
				prescanModuleFoundInThisSide.Platform))
		}
	}
}

func reportDuplicateFiles(side string, prescanFileList PrescanFileList) {
	var report strings.Builder
	var processedFiles []string

	for _, thisFile := range prescanFileList.Files {
		if isStringInStringArray(thisFile.Name, processedFiles) {
			continue
		}

		md5s := []string{thisFile.MD5}
		var count = 0

		for _, otherFile := range prescanFileList.Files {
			if thisFile.Name == otherFile.Name {
				count++
				if !isStringInStringArray(otherFile.MD5, md5s) {
					md5s = append(md5s, otherFile.MD5)
				}
			}
		}

		if len(md5s) > 1 {
			if count == len(md5s) {
				report.WriteString(fmt.Sprintf("\"%s\": %d occurances each with different MD5 hashes\n", thisFile.Name, count))
			} else {
				report.WriteString(fmt.Sprintf("\"%s\": %d occurances with %d different MD5 hashes\n", thisFile.Name, count, len(md5s)))
			}
		}

		processedFiles = append(processedFiles, thisFile.Name)
	}

	if report.Len() > 0 {
		colorPrintf(getFormattedSideStringWithMessage(side, fmt.Sprintf("\nDuplicate Files Within Scan %s\n", side)))
		fmt.Print("=============================\n")
		color.HiYellow(report.String())
	}
}

func getNonDuplicatedFileNames(fileList PrescanFileList) []string {
	var duplicateFiles []string
	var processedFiles []string

	for _, file := range fileList.Files {
		if isStringInStringArray(file.Name, processedFiles) && !isStringInStringArray(file.Name, duplicateFiles) {
			duplicateFiles = append(duplicateFiles, file.Name)
		}

		processedFiles = append(processedFiles, file.Name)
	}

	var nonDuplicatedFiles []string

	for _, file := range fileList.Files {
		if !isStringInStringArray(file.Name, duplicateFiles) {
			nonDuplicatedFiles = append(nonDuplicatedFiles, file.Name)
		}
	}

	return nonDuplicatedFiles
}

func (data Data) reportModuleDifferences() {
	var report strings.Builder

	var scanANonDuplicatedFiles = getNonDuplicatedFileNames(data.ScanAPrescanFileList)
	var scanBNonDuplicatedFiles = getNonDuplicatedFileNames(data.ScanBPrescanFileList)

	for _, thisFile := range data.ScanAPrescanFileList.Files {
		if !isStringInStringArray(thisFile.Name, scanANonDuplicatedFiles) {
			continue
		}

		if !isStringInStringArray(thisFile.Name, scanBNonDuplicatedFiles) {
			continue
		}

		for _, otherFile := range data.ScanBPrescanFileList.Files {
			if thisFile.Name == otherFile.Name {
				if thisFile.MD5 != otherFile.MD5 {
					report.WriteString(
						fmt.Sprintf("\"%s\" %s: MD5 = %s, %s: MD5 = %s \n",
							thisFile.Name,
							getFormattedSideString("A"),
							thisFile.MD5,
							getFormattedSideString("B"),
							otherFile.MD5))
				}
			}
		}
	}

	if report.Len() > 0 {
		color.HiCyan("\nModule Differences (Ignoring any duplicates)")
		fmt.Print("============================================\n")
		colorPrintf(report.String())
	}
}
