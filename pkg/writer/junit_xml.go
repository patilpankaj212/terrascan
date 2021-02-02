package writer

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/version"
)

const (
	junitXMLFormat supportedFormat = "junit-xml"
	testSuiteName                  = "TERRASCAN_POLICY_SUITE"
)

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	File      string   `xml:"file,attr"`
	Severity  string   `xml:"severity,attr"`
	Line      int      `xml:"line,attr"`
	Category  string   `xml:"category,attr"`
	// omit empty time because today we do not have this data
	Time        string            `xml:"time,attr,omitempty"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func newJunitTestSuite(summary results.ScanSummary) JUnitTestSuite {
	return JUnitTestSuite{Name: testSuiteName, Tests: summary.TotalPolicies, Time: fmt.Sprint(summary.TotalTime), Failures: summary.ViolatedPolicies, Properties: []JUnitProperty{
		{
			Name:  "Terrascan Version",
			Value: version.Get(),
		},
	}}
}

func init() {
	RegisterWriter(junitXMLFormat, JUnitXMLWriter)
}

// JUnitXMLWriter writes scan summary in junit xml format
func JUnitXMLWriter(data interface{}, writer io.Writer) error {
	output, ok := data.(policy.EngineOutput)
	if !ok {
		return fmt.Errorf("incorrect input for JunitXML writer, supportted type is policy.EngineOutput")
	}

	junitXMLOutput, err := convert(output)
	if err != nil {
		return err
	}

	return XMLWriter(junitXMLOutput, writer)
}

// convert is helper func to convert engine output to JUnitTestSuites
func convert(output policy.EngineOutput) (JUnitTestSuites, error) {
	testSuites := JUnitTestSuites{}
	suite := newJunitTestSuite(output.Summary)

	tests := violationsToTestCases(output.ViolationStore.Violations, false)
	if tests != nil {
		suite.TestCases = append(suite.TestCases, tests...)
	}

	skippedTests := violationsToTestCases(output.ViolationStore.SkippedViolations, true)
	if skippedTests != nil {
		suite.TestCases = append(suite.TestCases, skippedTests...)
	}

	testSuites.Suites = append(testSuites.Suites, suite)

	return testSuites, nil
}

// violationsToTestCases is helper func to convert scan violations to JunitTestCases
func violationsToTestCases(violations []*results.Violation, isSkipped bool) []JUnitTestCase {
	testCases := make([]JUnitTestCase, 0)
	for _, v := range violations {
		var testCase JUnitTestCase
		if isSkipped {
			testCase = JUnitTestCase{Failure: new(JUnitFailure), SkipMessage: new(JUnitSkipMessage)}
		} else {
			testCase = JUnitTestCase{Failure: new(JUnitFailure)}
		}
		testCase.Classname = v.RuleName
		testCase.Name = v.RuleID
		testCase.File = v.File
		testCase.Line = v.LineNumber
		testCase.Severity = v.Severity
		testCase.Category = v.Category
		testCase.Failure.Message = getViolationString(*v)
		testCase.Failure.Type = v.Severity
		if isSkipped {
			testCase.SkipMessage.Message = v.Comment
		}
		testCases = append(testCases, testCase)
	}
	return testCases
}

func getViolationString(v results.Violation) string {
	resourceName := v.ResourceName
	if resourceName == "" {
		resourceName = `""`
	}

	out := fmt.Sprintf("%s: %s, %s: %s, %s: %d, %s: %s, %s: %s, %s: %s, %s: %s, %s: %s, %s: %s",
		"Description", v.Description,
		"File", v.File,
		"Line", v.LineNumber,
		"Severity", v.Severity,
		"Rule Name", v.RuleName,
		"Rule ID", v.RuleID,
		"Resource Name", resourceName,
		"Resource Type", v.ResourceType,
		"Category", v.Category)
	return out
}
