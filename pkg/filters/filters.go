package filters

import (
	"github.com/accurics/terrascan/pkg/policy"
)

type RegoMetadataPreLoadFilter struct {
	scanRules   []string
	skipRules   []string
	categories  []string
	policyTypes []string
	severity    string
}

func NewRegoMetadataPreLoadFilter(scanRules, skipRules, categories, policyTypes []string, severity string) *RegoMetadataPreLoadFilter {
	return &RegoMetadataPreLoadFilter{
		scanRules:   scanRules,
		skipRules:   skipRules,
		categories:  categories,
		policyTypes: policyTypes,
		severity:    severity,
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
	isSeverityAllowed, isCategoryAllowed, isScanRuleAllowed, isPolicyTypeAllowed := true, true, true, true

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

	if len(r.policyTypes) > 0 {
		policyTypeSpec := PolicyTypeFilterSpecification{r.policyTypes}
		isPolicyTypeAllowed = policyTypeSpec.IsSatisfied(regoMetadata)
	}

	return isSeverityAllowed && isCategoryAllowed && isScanRuleAllowed && isPolicyTypeAllowed
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
