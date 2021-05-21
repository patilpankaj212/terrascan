package filters

import (
	"github.com/accurics/terrascan/pkg/policy"
)

type RegoMetadataPreLoadFilter struct {
	scanRules  []string
	skipRules  []string
	categories []string
	severity   string
}

func NewRegoMetadataPreLoadFilter(scanRules, skipRules, categories []string, severity string) *RegoMetadataPreLoadFilter {
	return &RegoMetadataPreLoadFilter{
		scanRules:  scanRules,
		skipRules:  skipRules,
		categories: categories,
		severity:   severity,
	}
}

func (r RegoMetadataPreLoadFilter) IsFiltered(regoMetadata *policy.RegoMetadata) bool {
	if len(r.skipRules) > 0 {
		refIDsSpec := RerefenceIDsFilterSpecification{r.skipRules}
		return refIDsSpec.IsSatisfied(regoMetadata)
	}

	return false
}

func (r RegoMetadataPreLoadFilter) IsAllowed(regoMetadata *policy.RegoMetadata) bool {
	isSeverityAllowed, isCategoryAllowed, isScanRuleAllowed := true, true, true

	if len(r.severity) > 0 {
		sevSpec := SeverityFilterSpecification{r.severity}
		isSeverityAllowed = sevSpec.IsSatisfied(regoMetadata)
	}

	if len(r.categories) > 0 {
		catSpec := CategoryFilterSpecification{r.categories}
		isCategoryAllowed = catSpec.IsSatisfied(regoMetadata)
	}

	if len(r.scanRules) > 0 {
		refIDsSpec := RerefenceIDsFilterSpecification{r.scanRules}
		isScanRuleAllowed = refIDsSpec.IsSatisfied(regoMetadata)
	}

	return isSeverityAllowed && isCategoryAllowed && isScanRuleAllowed
}

type RegoDataFilter struct{}

func (r RegoDataFilter) Filter(rmap map[string]*policy.RegoData, input policy.EngineInput) map[string]*policy.RegoData {
	tempMap := make(map[string]*policy.RegoData)
	for resType := range *input.InputData {
		for k := range rmap {
			resFilterSpec := ResourceTypeFilterSpecification{resType}
			if resFilterSpec.IsSatisfied(&rmap[k].Metadata) {
				tempMap[k] = rmap[k]
			}
		}
	}
	return tempMap
}
