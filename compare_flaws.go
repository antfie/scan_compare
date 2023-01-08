package main

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/fatih/color"
)

func (data Data) reportFlawDifferences() {
	data.reportFlawStateDifferences()
	data.reportFlawMitigationDifferences()
	data.reportFlawLineNumberChanges()

	// Disable this for now. Results are not stable. Unsure if memory leak or API returningindederministic result
	// data.reportMatchedFlawMovements()

	data.reportPolicyAffectingFlawDifferences()
	data.reportNonPolicyAffectingFlawDifferences()
	data.reportClosedFlawDifferences()
}

func (data Data) reportFlawStateDifferences() {
	var report strings.Builder

	compareFlawStates(&report, data.ScanAReport, data.ScanBReport)

	if report.Len() > 0 {
		color.HiCyan("\nFlaw State Differences")
		fmt.Print("======================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportFlawMitigationDifferences() {
	var report strings.Builder

	compareFlawMitigations(&report, data.ScanAReport, data.ScanBReport)

	if report.Len() > 0 {
		color.HiCyan("\nFlaw Mitigation Differences")
		fmt.Print("===========================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportFlawLineNumberChanges() {
	var report strings.Builder

	compareFlawLineNumberChanges(&report, data.ScanAReport, data.ScanBReport)

	if report.Len() > 0 {
		color.HiCyan("\nFlaw Line Number Differences")
		fmt.Print("============================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportMatchedFlawMovements() {
	var report strings.Builder

	compareFlawMatchMovements(&report, data.ScanAReport, data.ScanBReport)

	if report.Len() > 0 {
		color.HiCyan("\nMatched Flaw Movements")
		fmt.Print("======================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportPolicyAffectingFlawDifferences() {
	var report strings.Builder

	compareFlaws(&report, "A", data.ScanAReport, data.ScanBReport, true, false)
	compareFlaws(&report, "B", data.ScanBReport, data.ScanAReport, true, false)

	if report.Len() > 0 {
		color.HiCyan("\nPolicy Affecting Open Flaw Differences")
		fmt.Print("======================================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportNonPolicyAffectingFlawDifferences() {
	var report strings.Builder

	compareFlaws(&report, "A", data.ScanAReport, data.ScanBReport, false, false)
	compareFlaws(&report, "B", data.ScanBReport, data.ScanAReport, false, false)

	if report.Len() > 0 {
		color.HiCyan("\nNon Policy Affecting Open Flaw Differences")
		fmt.Print("==========================================\n")
		colorPrintf(report.String())
	}
}

func (data Data) reportClosedFlawDifferences() {
	var report strings.Builder

	compareFlaws(&report, "A", data.ScanAReport, data.ScanBReport, false, true)
	compareFlaws(&report, "B", data.ScanBReport, data.ScanAReport, false, true)

	if report.Len() > 0 {
		color.HiCyan("\nClosed Flaw Differences")
		fmt.Print("=======================\n")
		colorPrintf(report.String())
	}
}

func getSortedCwes(report DetailedReport) []int {
	var cwes []int
	for _, thisSideFlaw := range report.Flaws {
		if !isInIntArray(thisSideFlaw.CWE, cwes) {
			cwes = append(cwes, thisSideFlaw.CWE)
		}
	}

	sort.Ints(cwes[:])
	return cwes
}

func compareFlaws(report *strings.Builder, side string, thisSideReport, otherSideReport DetailedReport, policyAffecting bool, onlyClosed bool) {
	for _, cwe := range getSortedCwes(thisSideReport) {
		var flawsOnlyInThisScan []int

		for _, thisSideFlaw := range thisSideReport.Flaws {
			if thisSideFlaw.CWE != cwe {
				continue
			}

			if onlyClosed {
				if thisSideFlaw.isFlawOpen() {
					continue
				}
			} else {
				if policyAffecting && !(thisSideFlaw.isFlawOpen() && thisSideFlaw.AffectsPolicyCompliance) {
					continue
				}

				if !policyAffecting && !(thisSideFlaw.isFlawOpen() && !thisSideFlaw.AffectsPolicyCompliance) {
					continue
				}
			}

			if !otherSideReport.isFlawInReport(thisSideFlaw.ID) {
				flawsOnlyInThisScan = append(flawsOnlyInThisScan, thisSideFlaw.ID)
			}
		}

		if len(flawsOnlyInThisScan) > 0 {
			report.WriteString(fmt.Sprintf("%s: %dx CWE-%d = %s\n",
				getFormattedOnlyInSideString(side),
				len(flawsOnlyInThisScan),
				cwe,
				getSortedIntArrayAsFormattedString(flawsOnlyInThisScan)))
		}
	}
}

func compareFlawStates(report *strings.Builder, thisSideReport, otherSideReport DetailedReport) {
	stateChanges := make(map[string][]int)

	for _, thisSideFlaw := range thisSideReport.Flaws {
		for _, otherSideFlaw := range otherSideReport.Flaws {
			if thisSideFlaw.ID != otherSideFlaw.ID {
				continue
			}

			if thisSideFlaw.RemediationStatus == otherSideFlaw.RemediationStatus {
				continue
			}

			var stateChange = fmt.Sprintf("%s %-9s => %s %-9s: CWE-%d",
				getFormattedSideString("A"),
				thisSideFlaw.RemediationStatus,
				getFormattedSideString("B"),
				otherSideFlaw.RemediationStatus,
				thisSideFlaw.CWE)

			stateChanges[stateChange] = append(stateChanges[stateChange], thisSideFlaw.ID)

		}
	}

	sortedKeys := make([]string, 0, len(stateChanges))
	for k := range stateChanges {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		var flawIds = stateChanges[key]

		var formattedS = strings.Replace(key, "CWE", fmt.Sprintf("%dx CWE", len(flawIds)), 1)
		report.WriteString(fmt.Sprintf("%s = %s\n", formattedS, getSortedIntArrayAsFormattedString(flawIds)))
	}
}

func compareFlawMitigations(report *strings.Builder, thisSideReport, otherSideReport DetailedReport) {
	for _, thisSideFlaw := range thisSideReport.Flaws {
		for _, otherSideFlaw := range otherSideReport.Flaws {
			if thisSideFlaw.ID != otherSideFlaw.ID {
				continue
			}

			if thisSideFlaw.MitigationStatus != otherSideFlaw.MitigationStatus {
				report.WriteString(fmt.Sprintf("%d (CWE-%d): %s: %s, %s: %s\n",
					thisSideFlaw.ID,
					thisSideFlaw.CWE,
					getFormattedSideString("A"),
					cases.Title(language.English).String(thisSideFlaw.MitigationStatus),
					getFormattedSideString("B"),
					cases.Title(language.English).String(otherSideFlaw.MitigationStatus)))
			}
		}
	}
}

func compareFlawLineNumberChanges(report *strings.Builder, thisSideReport, otherSideReport DetailedReport) {
	for _, thisSideFlaw := range thisSideReport.Flaws {
		for _, otherSideFlaw := range otherSideReport.Flaws {
			if thisSideFlaw.ID != otherSideFlaw.ID {
				continue
			}

			if thisSideFlaw.LineNumber != otherSideFlaw.LineNumber {
				report.WriteString(fmt.Sprintf("%d (CWE-%d): %s: %d, %s: %d\n",
					thisSideFlaw.ID,
					thisSideFlaw.CWE,
					getFormattedSideString("A"),
					thisSideFlaw.LineNumber,
					getFormattedSideString("B"),
					otherSideFlaw.LineNumber))
			}
		}
	}
}

func compareFlawMatchMovements(report *strings.Builder, thisSideReport, otherSideReport DetailedReport) {
	for _, thisSideFlaw := range thisSideReport.Flaws {
		for _, otherSideFlaw := range otherSideReport.Flaws {
			if thisSideFlaw.ID != otherSideFlaw.ID {
				continue
			}

			// If we are missing hashes move on
			if len(thisSideFlaw.ProcedureHash) < 5 || len(thisSideFlaw.PrototypeHash) < 5 || len(thisSideFlaw.StatementHash) < 5 || len(otherSideFlaw.ProcedureHash) < 5 || len(otherSideFlaw.PrototypeHash) < 5 || len(otherSideFlaw.StatementHash) < 5 {
				continue
			}

			if thisSideFlaw.ProcedureHash != otherSideFlaw.ProcedureHash || thisSideFlaw.PrototypeHash != otherSideFlaw.PrototypeHash || thisSideFlaw.StatementHash != otherSideFlaw.StatementHash {
				report.WriteString(fmt.Sprintf("%d (CWE-%d) - %s: Procedure = %s, Prototype = %s, Statement = %s, %s: Procedure = %s, Prototype = %s, Statement = %s\n",
					thisSideFlaw.ID,
					thisSideFlaw.CWE,
					getFormattedSideString("A"),
					thisSideFlaw.ProcedureHash,
					thisSideFlaw.PrototypeHash,
					thisSideFlaw.StatementHash,
					getFormattedSideString("B"),
					otherSideFlaw.ProcedureHash,
					otherSideFlaw.PrototypeHash,
					otherSideFlaw.StatementHash))
			}
		}
	}
}
