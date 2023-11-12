package s3errors

import (
	"encoding/xml"
	"net/http"

	"github.com/lvjp/s3impl/pkg/s3consts"
)

type ErrorWriter interface {
	Write(*S3Error, http.ResponseWriter)
}

type apierror struct {
	Code      string
	Message   string
	RequestID string
}

func (a *apierror) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "Error"
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	type Entry struct {
		Name  string
		Value string
	}

	entries := []Entry{
		{"Code", a.Code},
		{"Message", a.Message},
		{"RequestId", a.RequestID},
		{"HostId", a.RequestID},
	}

	for _, entry := range entries {
		if err := enc.EncodeElement(entry.Value, xml.StartElement{Name: xml.Name{Local: entry.Name}}); err != nil {
			return err
		}
	}

	return enc.EncodeToken(start.End())
}

type APIWriter struct{}

func (*APIWriter) Write(err *S3Error, w http.ResponseWriter) error {
	header := w.Header()
	header.Set("x-amz-request-id", err.RequestID)
	header.Set("x-amz-id-2", err.RequestID)
	header.Set("Content-Type", s3consts.MimetypeApplicationXML)
	w.WriteHeader(err.HTTPStatusCode)

	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}

	payload := &apierror{
		Code:      err.Code,
		Message:   err.Message,
		RequestID: err.RequestID,
	}

	return xml.NewEncoder(w).Encode(payload)
}
