package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

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

func (data Data) reportCloseDFlawDifferences() {
	var report strings.Builder

	compareFlaws(&report, "A", data.ScanAReport, data.ScanBReport, false, true)
	compareFlaws(&report, "B", data.ScanBReport, data.ScanAReport, false, true)

	if report.Len() > 0 {
		color.HiCyan("\nClosed Flaw Differences")
		fmt.Print("=======================\n")
		colorPrintf(report.String())
	}
}

func compareFlaws(report *strings.Builder, side string, thisSideReport, otherSideReport DetailedReport, policyAffecting bool, onlyClosed bool) {
	var cwes []int
	for _, thisSideFlaw := range thisSideReport.Flaws {
		if !isInIntArray(thisSideFlaw.CWE, cwes) {
			cwes = append(cwes, thisSideFlaw.CWE)
		}
	}

	sort.Ints(cwes[:])

	for _, cwe := range cwes {
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
				getSortedIntArrayAsSFormattedString(flawsOnlyInThisScan)))
		}
	}
}
