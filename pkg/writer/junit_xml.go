package writer

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/version"
)

const (
	junitXMLFormat supportedFormat = "junit-xml"
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
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	File        string            `xml:"file,attr"`
	Severity    string            `xml:"severity,attr"`
	Line        int               `xml:"line,attr"`
	Category    string            `xml:"category,attr"`
	Time        string            `xml:"time,attr"`
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

func convert(output policy.EngineOutput) (JUnitTestSuites, error) {
	o := JUnitTestSuites{}
	suite := JUnitTestSuite{Name: "TERRASCAN_POLICY_SUITE", Tests: output.Summary.TotalPolicies, Time: fmt.Sprint(output.Summary.TotalTime), Failures: output.Summary.ViolatedPolicies, Properties: []JUnitProperty{
		{
			Name:  "Terrascan Version",
			Value: version.Get(),
		},
	}}
	for _, v := range output.ViolationStore.Violations {
		testCase := JUnitTestCase{Failure: new(JUnitFailure)}
		testCase.Classname = v.RuleID
		testCase.Name = v.RuleName
		testCase.File = v.File
		testCase.Line = v.LineNumber
		testCase.Severity = v.Severity
		testCase.Category = v.Category
		testCase.Failure.Message = v.Description
		suite.TestCases = append(suite.TestCases, testCase)
	}

	for _, v := range output.ViolationStore.SkippedViolations {
		testCase := JUnitTestCase{Failure: new(JUnitFailure), SkipMessage: new(JUnitSkipMessage)}
		testCase.Classname = v.RuleID
		testCase.Name = v.RuleName
		testCase.File = v.File
		testCase.Line = v.LineNumber
		testCase.Severity = v.Severity
		testCase.Category = v.Category
		testCase.Failure.Message = v.Description
		testCase.SkipMessage.Message = v.Comment
		suite.TestCases = append(suite.TestCases, testCase)
	}

	o.Suites = append(o.Suites, suite)

	return o, nil
}
