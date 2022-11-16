package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
)

type PrescanModuleList struct {
	XMLName xml.Name        `xml:"prescanresults"`
	Modules []PrescanModule `xml:"module"`
}

type PrescanModule struct {
	XMLName        xml.Name             `xml:"module"`
	ID             int                  `xml:"id,attr"`
	Name           string               `xml:"name,attr"`
	Status         string               `xml:"status,attr"`
	Platform       string               `xml:"platform,attr"`
	Size           string               `xml:"size,attr"`
	MD5            string               `xml:"checksum,attr"`
	HasFatalErrors bool                 `xml:"has_fatal_errors,attr"`
	IsDependency   bool                 `xml:"is_dependency,attr"`
	Issues         []PrescanModuleIssue `xml:"issue"`
}

type PrescanModuleIssue struct {
	XMLName xml.Name `xml:"issue"`
	Details string   `xml:"details,attr"`
}

func (api API) getPrescanModuleList(appId, buildId int) PrescanModuleList {
	var url = fmt.Sprintf("https://analysiscenter.veracode.com/api/5.0/getprescanresults.do?app_id=%d&build_id=%d", appId, buildId)
	response := api.makeApiRequest(url, http.MethodGet)

	moduleList := PrescanModuleList{}
	xml.Unmarshal(response, &moduleList)

	// Sort modules by name for consistency
	sort.Slice(moduleList.Modules, func(i, j int) bool {
		return moduleList.Modules[i].Name < moduleList.Modules[j].Name
	})

	return moduleList
}

func (moduleList PrescanModuleList) getFromName(moduleName string) PrescanModule {
	for _, moduleFromlist := range moduleList.Modules {
		if moduleFromlist.Name == moduleName {
			return moduleFromlist
		}
	}

	return PrescanModule{}
}
