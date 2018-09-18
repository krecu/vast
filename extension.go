package vast

import "encoding/xml"

// Extension represent arbitrary XML provided by the platform to extend the
// VAST response or by custom trackers.
type Extension struct {
	Type           string     `xml:"type,attr,omitempty"`
	Name           string     `xml:"name,attr,omitempty"`
	CustomTracking []Tracking `xml:"CustomTracking>Tracking,omitempty"`
	Data           []byte     `xml:",innerxml"`
	Attributes     map[string]string
}

// the extension type as a middleware in the encoding process.
type extension Extension

type extensionNoCT struct {
	Type string `xml:"type,attr,omitempty"`
	Name string `xml:"name,attr,omitempty"`
	Data []byte `xml:",innerxml"`
}

// MarshalXML implements xml.Marshaler interface.
func (e Extension) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	// create a temporary element from a wrapper Extension, copy what we need to
	// it and return it's encoding.
	var e2 interface{}
	// if we have custom trackers, we should ignore the data, if not, then we
	// should consider only the data.
	if len(e.CustomTracking) > 0 {
		e2 = extension{Type: e.Type, Name: e.Name, CustomTracking: e.CustomTracking}
	} else {
		e2 = extensionNoCT{Type: e.Type, Name: e.Name, Data: e.Data}
	}

	// custom attributes
	if len(e.Attributes) > 0 {
		for name, value := range e.Attributes {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  xml.Name{Space: "", Local: name},
				Value: value,
			})
		}
	}

	return enc.EncodeElement(e2, start)
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *Extension) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	// decode the extension into a temporary element from a wrapper Extension,
	// copy what we need over.
	var e2 extension
	if err := dec.DecodeElement(&e2, &start); err != nil {
		return err
	}
	// copy the type and the customTracking
	e.Type = e2.Type
	e.Name = e2.Name
	e.CustomTracking = e2.CustomTracking
	// copy the data only of customTraking is empty
	if len(e.CustomTracking) == 0 {
		e.Data = e2.Data
	}

	// if extension have attribute
	if len(start.Attr) > 0 {
		for name, value := range e.Attributes {
			if name != "name" && name != "type" {
				e.Attributes[name] = value
			}
		}
	}

	return nil
}
