package filters

import (
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

type PolicyTypeFilterSpecification struct {
	policyType string
}

func (p PolicyTypeFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return p.policyType == r.PolicyType
}

type ResourceTypeFilterSpecification struct {
	resourceType string
}

func (rs ResourceTypeFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return rs.resourceType == r.ResourceType
}

type RerefenceIDFilterSpecification struct {
	ReferenceID string
}

func (rs RerefenceIDFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return rs.ReferenceID == r.ReferenceID
}

type RerefenceIDsFilterSpecification struct {
	ReferenceIDs []string
}

func (rs RerefenceIDsFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	isSatisfied := false
	for _, refID := range rs.ReferenceIDs {
		rfIDSpec := RerefenceIDFilterSpecification{refID}
		if rfIDSpec.IsSatisfied(r) {
			isSatisfied = true
			break
		}
	}
	return isSatisfied
}

type CategoryFilterSpecification struct {
	categories []string
}

func (c CategoryFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return utils.CheckCategory(r.Category, c.categories)
}

type SeverityFilterSpecification struct {
	severity string
}

func (s SeverityFilterSpecification) IsSatisfied(r *policy.RegoMetadata) bool {
	return utils.CheckSeverity(r.Severity, s.severity)
}
