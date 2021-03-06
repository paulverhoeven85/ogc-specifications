package ows

import (
	"encoding/xml"
	"testing"
)

func TestBoundingBoxBuildQueryString(t *testing.T) {
	var tests = []struct {
		boundingbox       BoundingBox
		boundingboxstring string
	}{
		// While 'not' correct this will we checked in the validation step
		0: {boundingbox: BoundingBox{}, boundingboxstring: `0.000000,0.000000,0.000000,0.000000`},
		1: {boundingbox: BoundingBox{LowerCorner: [2]float64{-180.0, -90.0}, UpperCorner: [2]float64{180.0, 90.0}}, boundingboxstring: `-180.000000,-90.000000,180.000000,90.000000`},
	}
	for k, test := range tests {
		str := test.boundingbox.BuildKVP()
		if str != test.boundingboxstring {
			t.Errorf("test: %d, expected: %v+,\n got: %v+", k, test.boundingboxstring, str)
		}
	}
}

func TestBuildBoundingBox(t *testing.T) {
	var tests = []struct {
		boundingbox string
		bbox        BoundingBox
		Exception   Exception
	}{
		0: {boundingbox: "0,0,100,100", bbox: BoundingBox{LowerCorner: [2]float64{0, 0}, UpperCorner: [2]float64{100, 100}}},
		1: {boundingbox: "0,0,-100,-100", bbox: BoundingBox{LowerCorner: [2]float64{0, 0}, UpperCorner: [2]float64{-100, -100}}}, // while this isn't correct, this will be 'addressed' in the validation step
		2: {boundingbox: "0,0,100", Exception: InvalidParameterValue(`0,0,100`, `boundingbox`)},
		3: {boundingbox: ",,,", Exception: InvalidParameterValue(`,,,`, `boundingbox`)},
		4: {boundingbox: ",,,100", Exception: InvalidParameterValue(`,,,100`, `boundingbox`)},
		5: {boundingbox: "number,,,100", Exception: InvalidParameterValue(`number,,,100`, `boundingbox`)},
	}

	for k, test := range tests {
		var bbox BoundingBox
		if err := bbox.ParseString(test.boundingbox); err != nil {
			if err != test.Exception {
				t.Errorf("test: %d, expected: %+v \ngot: %+v", k, test.Exception, err)
			}
		} else {
			if bbox != test.bbox {
				t.Errorf("test: %d, expected: %+v \ngot: %+v", k, test.bbox, bbox)
			}
		}
	}
}

func TestStripDuplicateAttr(t *testing.T) {
	var tests = []struct {
		attributes []xml.Attr
		expected   []xml.Attr
	}{
		0: {attributes: []xml.Attr{{Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}}, expected: []xml.Attr{{Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}}},
		1: {attributes: []xml.Attr{{Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}, {Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}, {Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}},
			expected: []xml.Attr{{Name: xml.Name{Local: "gml"}, Value: "http://www.opengis.net/gml/3.2"}}},
	}

	for k, test := range tests {
		stripped := StripDuplicateAttr(test.attributes)
		if len(test.expected) != len(stripped) {
			t.Errorf("test: %d, expected: %s,\n got: %s", k, test.expected, stripped)
		} else {
			c := false
			for _, exceptedAttr := range test.expected {
				for _, result := range stripped {
					if exceptedAttr == result {
						c = true
					}
				}
				if !c {
					t.Errorf("test: %d, expected: %s,\n got: %s", k, test.expected, stripped)
				}
				c = false
			}
		}
	}
}

func TestCRSParseString(t *testing.T) {
	var tests = []struct {
		input       string
		expectedCRS CRS
	}{
		0: {}, // Empty input == empty struct
		1: {input: `urn:ogc:def:crs:EPSG::4326`, expectedCRS: CRS{Code: 4326, Namespace: `EPSG`}},
		2: {input: `EPSG:4326`, expectedCRS: CRS{Code: 4326, Namespace: `EPSG`}},
	}

	for k, test := range tests {
		var crs CRS
		crs.parseString(test.input)

		if crs.Code != test.expectedCRS.Code || crs.Namespace != test.expectedCRS.Namespace {
			t.Errorf("test: %d, expected: %v,\n got: %v", k, test.expectedCRS, crs)
		}
	}
}
