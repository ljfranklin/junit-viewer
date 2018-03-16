package junit

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

const (
	TimeFormat = time.RFC3339
)

type TestSuite struct {
	Tests      int
	Successes  int
	Failures   int
	Errors     int
	Skips      int
	Time       float64
	Timestamp  time.Time
	Properties map[string]string
	TestCases  []TestCase
}

type TestCase struct {
	Name           string
	Time           float64
	Passed         bool
	Failed         bool
	Errored        bool
	Skipped        bool
	FailureMessage string
	FailureType    string
	FailureOutput  string
	ErrorMessage   string
	ErrorType      string
	ErrorOutput    string
	SkipMessage    string
	SystemOut      string
	SystemErr      string
}

type xmlTestSuite struct {
	XMLName    xml.Name `xml:"testsuite"`
	Time       float64  `xml:"time,attr"`
	Timestamp  string   `xml:"timestamp,attr"`
	Properties []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"properties>property"`
	TestCases []struct {
		Name    string  `xml:"name,attr"`
		Time    float64 `xml:"time,attr"`
		Failure struct {
			Message string `xml:"message,attr"`
			Type    string `xml:"type,attr"`
			Output  string `xml:",innerxml"`
		} `xml:"failure"`
		Error struct {
			Message string `xml:"message,attr"`
			Type    string `xml:"type,attr"`
			Output  string `xml:",innerxml"`
		} `xml:"error"`
		Skipped struct {
			Message string `xml:"message,attr"`
		} `xml:"skipped"`
		SystemOut string `xml:"system-out"`
		SystemErr string `xml:"system-err"`
	} `xml:"testcase"`
}
type xmlTestSuites struct {
	XMLName    xml.Name       `xml:"testsuites"`
	TestSuites []xmlTestSuite `xml:"testsuite"`
}

func Load(xmlPath string) ([]TestSuite, error) {
	xmlInput := xmlTestSuites{}

	err := unmarshalXML(xmlPath, &xmlInput)
	if err != nil {
		// some junit output is missing the root `testsuites` element
		singleTestSuiteInput := xmlTestSuite{}
		err = unmarshalXML(xmlPath, &singleTestSuiteInput)
		if err != nil {
			return nil, fmt.Errorf("failed to parse XML in '%s': %s", xmlPath, err)
		}
		xmlInput = xmlTestSuites{
			TestSuites: []xmlTestSuite{singleTestSuiteInput},
		}
	}

	output := []TestSuite{}
	for _, inputSuite := range xmlInput.TestSuites {
		timestamp := time.Time{}
		if len(inputSuite.Timestamp) > 0 {
			timestamp, err = time.Parse(TimeFormat, inputSuite.Timestamp)
			if err != nil {
				panic(err)
			}
		}
		outputSuite := TestSuite{
			Time:       inputSuite.Time,
			Timestamp:  timestamp,
			Properties: map[string]string{},
			TestCases:  []TestCase{},
		}
		for _, prop := range inputSuite.Properties {
			outputSuite.Properties[prop.Name] = prop.Value
		}
		for _, testCase := range inputSuite.TestCases {
			tc := TestCase{
				Name:           testCase.Name,
				Time:           testCase.Time,
				FailureMessage: testCase.Failure.Message,
				FailureType:    testCase.Failure.Type,
				FailureOutput:  testCase.Failure.Output,
				ErrorMessage:   testCase.Error.Message,
				ErrorType:      testCase.Error.Type,
				ErrorOutput:    testCase.Error.Output,
				SkipMessage:    testCase.Skipped.Message,
				SystemOut:      testCase.SystemOut,
				SystemErr:      testCase.SystemErr,
			}
			outputSuite.Tests++
			if len(tc.FailureOutput) > 0 {
				tc.Failed = true
				outputSuite.Failures++
			} else if len(tc.ErrorOutput) > 0 {
				tc.Errored = true
				outputSuite.Errors++
			} else if len(tc.SkipMessage) > 0 {
				tc.Skipped = true
				outputSuite.Skips++
			} else {
				tc.Passed = true
				outputSuite.Successes++
			}
			outputSuite.TestCases = append(outputSuite.TestCases, tc)
		}
		output = append(output, outputSuite)
	}

	return output, nil
}

func unmarshalXML(xmlPath string, obj interface{}) error {
	input, err := os.Open(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %s", xmlPath, err)
	}
	return xml.NewDecoder(input).Decode(&obj)
}
