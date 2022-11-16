package main

import (
	"sync"
	"time"
)

type Data struct {
	ScanAReport            SummaryReport
	ScanBReport            SummaryReport
	ScanAPrescanFileList   PrescanFileList
	ScanBPrescanFileList   PrescanFileList
	ScanAPrescanModuleList PrescanModuleList
	ScanBPrescanModuleList PrescanModuleList
	ScanASubmittedDate     time.Time
	ScanBSubmittedDate     time.Time
	ScanADuration          time.Duration
	ScanBDuration          time.Duration
}

func (api API) getData(scanAAppId, scanABuildId, scanBAppId, scanBBuildId int) Data {
	var data = Data{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		data.ScanAReport = api.getSummaryReport(scanABuildId)
	}()

	go func() {
		defer wg.Done()
		data.ScanBReport = api.getSummaryReport(scanBBuildId)
	}()

	wg.Wait()
	wg.Add(4)

	// We can't rely on the passed-in app IDs as they may not be present if not using a URL, so get the app ID from the summary report

	go func() {
		defer wg.Done()
		data.ScanAPrescanFileList = api.getPrescanFileList(data.ScanAReport.AppId, scanABuildId)
	}()

	go func() {
		defer wg.Done()
		data.ScanBPrescanFileList = api.getPrescanFileList(data.ScanBReport.AppId, scanBBuildId)
	}()

	go func() {
		defer wg.Done()
		data.ScanAPrescanModuleList = api.getPrescanModuleList(data.ScanAReport.AppId, scanABuildId)
	}()

	go func() {
		defer wg.Done()
		data.ScanBPrescanModuleList = api.getPrescanModuleList(data.ScanBReport.AppId, scanBBuildId)
	}()

	wg.Wait()

	data.ScanASubmittedDate = parseVeracodeDate(data.ScanAReport.StaticAnalysis.SubmittedDate).Local()
	data.ScanBSubmittedDate = parseVeracodeDate(data.ScanBReport.StaticAnalysis.SubmittedDate).Local()
	data.ScanADuration = parseVeracodeDate(data.ScanAReport.StaticAnalysis.PublishedDate).Local().Sub(data.ScanASubmittedDate)
	data.ScanBDuration = parseVeracodeDate(data.ScanBReport.StaticAnalysis.PublishedDate).Local().Sub(data.ScanBSubmittedDate)

	return data
}
