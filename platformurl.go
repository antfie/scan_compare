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

func isParseableURL(urlFragment string) bool {
	var supportedPages = []string{
		"ReviewResultsStaticFlaws",
		"ReviewResultsAllFlaws",
		"AnalyzeAppModuleList",
		"StaticOverview",
		"AnalyzeAppSourceFiles",
		"ViewReportsResultSummary",
		"ViewReportsDetailedReport"}

	for _, page := range supportedPages {
		if strings.HasPrefix(urlFragment, page) {
			return true
		}
	}
	return false
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

	if isParseableURL(urlFragment) {
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

	if isParseableURL(urlFragment) {
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

	if isParseableURL(urlFragment) {
		buildId, err := strconv.Atoi(strings.Split(urlFragment, ":")[3])

		if err != nil {
			platformUrlInvalid(urlOrBuildId)
		}

		return buildId

	}

	platformUrlInvalid(urlOrBuildId)
	return -1
}
