package main

import (
	"sync"
)

type Data struct {
	ScanAReport            DetailedReport
	ScanBReport            DetailedReport
	ScanAPrescanFileList   PrescanFileList
	ScanBPrescanFileList   PrescanFileList
	ScanAPrescanModuleList PrescanModuleList
	ScanBPrescanModuleList PrescanModuleList
}

func (api API) getData(scanAAppId, scanABuildId, scanBAppId, scanBBuildId int) Data {
	var data = Data{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		data.ScanAReport = api.getDetailedReport(scanABuildId)
	}()

	go func() {
		defer wg.Done()
		data.ScanBReport = api.getDetailedReport(scanBBuildId)
	}()

	wg.Wait()
	wg.Add(4)

	// We can't rely on the passed-in app IDs as they may not be present if not using a URL, so get the app ID from the detailed report

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

	data.ScanAReport.SubmittedDate = parseVeracodeDate(data.ScanAReport.StaticAnalysis.SubmittedDate).Local()
	data.ScanAReport.PublishedDate = parseVeracodeDate(data.ScanAReport.StaticAnalysis.PublishedDate).Local()
	data.ScanBReport.SubmittedDate = parseVeracodeDate(data.ScanBReport.StaticAnalysis.SubmittedDate).Local()
	data.ScanBReport.PublishedDate = parseVeracodeDate(data.ScanBReport.StaticAnalysis.PublishedDate).Local()
	data.ScanAReport.Duration = data.ScanAReport.PublishedDate.Sub(data.ScanAReport.SubmittedDate)
	data.ScanBReport.Duration = data.ScanBReport.PublishedDate.Sub(data.ScanBReport.SubmittedDate)

	return data
}
