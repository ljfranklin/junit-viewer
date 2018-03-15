package junit

import (
	"encoding/xml"
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
	Time       float64
	Timestamp  time.Time
	Properties map[string]string
	TestCases  []TestCase
}

type TestCase struct {
	Name string
	Time float64
}

func Load(xmlPath string) ([]TestSuite, error) {
	xmlInput := struct {
		XMLName    xml.Name `xml:"testsuites"`
		TestSuites []struct {
			Tests      int     `xml:"tests,attr"`
			Failures   int     `xml:"failures,attr"`
			Time       float64 `xml:"time,attr"`
			Timestamp  string  `xml:"timestamp,attr"`
			Properties []struct {
				Property struct {
					Name  string `xml:"name,attr"`
					Value string `xml:"value,attr"`
				} `xml:"property"`
			} `xml:"properties"`
			TestCases []struct {
				XMLName xml.Name `xml:"testcase"`
				Name    string   `xml:"name,attr"`
				Time    float64  `xml:"time,attr"`
			} `xml:"testcase"`
		} `xml:"testsuite"`
	}{}

	input, err := os.Open(xmlPath)
	if err != nil {
		panic(err)
	}
	decoder := xml.NewDecoder(input)
	err = decoder.Decode(&xmlInput)
	if err != nil {
		panic(err)
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
			Tests:      inputSuite.Tests,
			Successes:  inputSuite.Tests - inputSuite.Failures,
			Failures:   inputSuite.Failures,
			Time:       inputSuite.Time,
			Timestamp:  timestamp,
			Properties: map[string]string{},
			TestCases:  []TestCase{},
		}
		for _, prop := range inputSuite.Properties {
			outputSuite.Properties[prop.Property.Name] = prop.Property.Value
		}
		for _, testCase := range inputSuite.TestCases {
			outputSuite.TestCases = append(outputSuite.TestCases, TestCase{
				Name: testCase.Name,
				Time: testCase.Time,
			})
		}
		output = append(output, outputSuite)
	}

	return output, nil
}
