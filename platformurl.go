package main

import (
	"fmt"
	"strconv"
	"strings"
)

func platformUrlInvalid(url string) int {
	panic(fmt.Sprintf("%s is not a valid or supported Veracode Platform URL", url))
}

func isPlatformURL(url string) bool {
	return strings.HasPrefix(url, "https://analysiscenter.veracode.com/auth/index.jsp")
}

func parseAccountIdFromPlatformUrl(urlOrAccountId string) int {
	accountId, err := strconv.Atoi(urlOrAccountId)

	if err == nil {
		return accountId
	}

	if !isPlatformURL(urlOrAccountId) {
		platformUrlInvalid(urlOrAccountId)
	}

	var urlFragment = strings.Split(urlOrAccountId, "#")[1]

	if strings.HasPrefix(urlFragment, "ReviewResultsStaticFlaws") || strings.HasPrefix(urlFragment, "AnalyzeAppModuleList") || strings.HasPrefix(urlFragment, "StaticOverview") || strings.HasPrefix(urlFragment, "AnalyzeAppSourceFiles") || strings.HasPrefix(urlFragment, "ViewReportsResultSummary") || strings.HasPrefix(urlFragment, "ViewReportsDetailedReport") {
		accountId, err := strconv.Atoi(strings.Split(urlFragment, ":")[1])

		if err != nil {
			platformUrlInvalid(urlOrAccountId)
		}

		return accountId

	}

	platformUrlInvalid(urlOrAccountId)
	return -1
}

func parseAppIdFromPlatformUrl(urlOrAppId string) int {
	appId, err := strconv.Atoi(urlOrAppId)

	if err == nil {
		return appId
	}

	if !isPlatformURL(urlOrAppId) {
		platformUrlInvalid(urlOrAppId)
	}

	var urlFragment = strings.Split(urlOrAppId, "#")[1]

	if strings.HasPrefix(urlFragment, "ReviewResultsStaticFlaws") || strings.HasPrefix(urlFragment, "AnalyzeAppModuleList") || strings.HasPrefix(urlFragment, "StaticOverview") || strings.HasPrefix(urlFragment, "AnalyzeAppSourceFiles") || strings.HasPrefix(urlFragment, "ViewReportsResultSummary") || strings.HasPrefix(urlFragment, "ViewReportsDetailedReport") {
		appId, err := strconv.Atoi(strings.Split(urlFragment, ":")[2])

		if err != nil {
			platformUrlInvalid(urlOrAppId)
		}

		return appId

	}

	platformUrlInvalid(urlOrAppId)
	return -1
}

func parseBuildIdFromPlatformUrl(urlOrBuildId string) int {
	buildId, err := strconv.Atoi(urlOrBuildId)

	if err == nil {
		return buildId
	}

	if !isPlatformURL(urlOrBuildId) {
		platformUrlInvalid(urlOrBuildId)
	}

	var urlFragment = strings.Split(urlOrBuildId, "#")[1]

	if strings.HasPrefix(urlFragment, "ReviewResultsStaticFlaws") || strings.HasPrefix(urlFragment, "AnalyzeAppModuleList") || strings.HasPrefix(urlFragment, "StaticOverview") || strings.HasPrefix(urlFragment, "AnalyzeAppSourceFiles") || strings.HasPrefix(urlFragment, "ViewReportsResultSummary") || strings.HasPrefix(urlFragment, "ViewReportsDetailedReport") {
		buildId, err := strconv.Atoi(strings.Split(urlFragment, ":")[3])

		if err != nil {
			platformUrlInvalid(urlOrBuildId)
		}

		return buildId

	}

	platformUrlInvalid(urlOrBuildId)
	return -1
}
