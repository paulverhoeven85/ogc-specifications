package request

import (
	"net/url"

	"github.com/pdok/ogc-specifications/pkg/ows"
)

//
const (
	Service = `WFS`
	Version = `2.0.0`
)

// WFS 2.0.0 Tokens
const (
	SERVICE = `SERVICE`
	REQUEST = `REQUEST`
	VERSION = `VERSION`

	OUTPUTFORMAT = `OUTPUTFORMAT`
)

// BaseRequest based on Table 5 WFS2.0.0 spec
// Note: not usable for GetCapabilities request regarding deviation of Optional/Mandatory parameters SERVICE and VERSION
type BaseRequest struct {
	Service string           `xml:"service,attr" yaml:"service,omitempty"`
	Version string           `xml:"version,attr" yaml:"version"`
	Attr    ows.XMLAttribute `xml:",attr"`
}

// ParseKVP builds a BaseRequest Struct based on the given parameters
func (b *BaseRequest) ParseKVP(query url.Values) ows.Exceptions {
	if len(query[SERVICE]) > 0 {
		// Service is optional, because it's implicit for a GetFeature/DescribeFeatureType request
		b.Service = query[SERVICE][0]
	}
	if len(query[VERSION]) > 0 {
		b.Version = query[VERSION][0]
	} else {
		// Version is mandatory
		return ows.Exceptions{ows.MissingParameterValue(VERSION)}
	}
	return nil
}
