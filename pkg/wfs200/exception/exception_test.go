package exception

import (
	"testing"

	"github.com/pdok/ogc-specifications/pkg/ows"
)

func TestWFSException(t *testing.T) {
	var tests = []struct {
		exception     ows.Exception
		exceptionText string
		exceptionCode string
		locatorCode   string
	}{
		0: {exception: WFSException{ExceptionCode: "", ExceptionText: "", LocatorCode: ""},
			exceptionText: "",
			exceptionCode: "",
			locatorCode:   "",
		},
		1: {exception: CannotLockAllFeatures(),
			exceptionCode: "CannotLockAllFeatures",
		},
		2: {exception: DuplicateStoredQueryIDValue(),
			exceptionCode: "DuplicateStoredQueryIDValue",
		},
		3: {exception: DuplicateStoredQueryParameterName(),
			exceptionCode: "DuplicateStoredQueryParameterName",
		},
		4: {exception: FeaturesNotLocked(),
			exceptionCode: "FeaturesNotLocked",
		},
		5: {exception: InvalidLockID(),
			exceptionCode: "InvalidLockID",
		},
		6: {exception: InvalidValue(),
			exceptionCode: "InvalidValue",
		},
		7: {exception: LockHasExpired(),
			exceptionCode: "LockHasExpired",
		},
		8: {exception: OperationParsingFailed("PARAMETER", "VALUE"),
			exceptionCode: "OperationParsingFailed",
			exceptionText: "Failed to parse the operation, found: PARAMETER",
			locatorCode:   "VALUE",
		},
		9: {exception: OperationProcessingFailed(),
			exceptionCode: "OperationProcessingFailed",
		},
		10: {exception: ResponseCacheExpired(),
			exceptionCode: "ResponseCacheExpired",
		},
	}

	for k, a := range tests {
		if a.exception.Error() != a.exceptionText {
			t.Errorf("test: %d, expected: %s\n got: %s", k, a.exceptionText, a.exception.Error())
		}
		if a.exception.Code() != a.exceptionCode {
			t.Errorf("test: %d, expected: %s\n got: %s", k, a.exceptionCode, a.exception.Code())
		}
		if a.exception.Locator() != a.locatorCode {
			t.Errorf("test: %d, expected: %s\n got: %s", k, a.locatorCode, a.exception.Locator())
		}
	}
}

func TestReport(t *testing.T) {
	var tests = []struct {
		exceptions []ows.Exception
		result     []byte
		err        error
	}{
		0: {exceptions: []ows.Exception{WFSException{ExceptionCode: "", ExceptionText: "", LocatorCode: ""}},
			result: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ExceptionReport xmlns:ows="http://www.opengis.net/ows/1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/ows/1.1 http://schemas.opengis.net/ows/1.1.0/owsExceptionReport.xsd" version="2.0.0" xml:lang="en">
 <Exception exceptionCode=""></Exception>
</ExceptionReport>`)},
		1: {exceptions: []ows.Exception{
			CannotLockAllFeatures(),
			DuplicateStoredQueryIDValue(),
		},
			result: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ExceptionReport xmlns:ows="http://www.opengis.net/ows/1.1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/ows/1.1 http://schemas.opengis.net/ows/1.1.0/owsExceptionReport.xsd" version="2.0.0" xml:lang="en">
 <Exception exceptionCode="CannotLockAllFeatures"></Exception>
 <Exception exceptionCode="DuplicateStoredQueryIDValue"></Exception>
</ExceptionReport>`)},
	}

	for k, a := range tests {
		report := WFSExceptionReport{}
		r := report.Report(a.exceptions)

		if string(r) != string(a.result) {
			t.Errorf("test: %d, expected: %s\n got: %s", k, a.result, r)
		}
	}
}
