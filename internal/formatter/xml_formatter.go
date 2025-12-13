package formatter

import (
	"encoding/xml"

	"github.com/kubex-ecosystem/logz/internal/module/kbx"
)

type XMLFormatter struct {
	Pretty bool
}

func NewXMLFormatter(pretty bool) Formatter {
	return &XMLFormatter{Pretty: pretty}
}

func (f *XMLFormatter) Name() string {
	return "xml"
}

func (f *XMLFormatter) Format(e kbx.Entry) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}

	var xmlOut []byte
	var err error
	eTmp := e.Clone().(kbx.LogzEntry).
		WithFields(nil).
		WithData(nil)
	if f.Pretty {
		xmlOut, err = xml.MarshalIndent(eTmp, eTmp.GetPrefix(), "  ")
		if err != nil {
			return nil, err
		}
	} else {
		xmlOut, err = xml.Marshal(e)
		if err != nil {
			return nil, err
		}
	}

	return xmlOut, nil

}
