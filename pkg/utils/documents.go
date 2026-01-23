package utils

import "strings"

// File format
type FileFormat string

// Available file formats
const (
	FileFormatCSV  FileFormat = "csv"
	FileFormatXLSX FileFormat = "xlsx"
	FileFormatDocX FileFormat = "docx"
	FileFormatPDF  FileFormat = "pdf"
	FileFormatJSON FileFormat = "json"
	FileFormatHTML FileFormat = "html"
	FileFormatXML  FileFormat = "xml"
)

// Parse parse FileFormat from string
func (f FileFormat) Parse(v string) FileFormat {
	res := FileFormat(strings.TrimSpace(strings.ToLower(v)))

	switch res {
	case "xls", FileFormatXLSX:
		return FileFormatXLSX

	case "doc", FileFormatDocX:
		return FileFormatDocX

	case "js", FileFormatJSON:
		return FileFormatJSON

	case FileFormatCSV:
		return FileFormatCSV

	case FileFormatPDF:
		return FileFormatPDF

	case FileFormatXML:
		return FileFormatXML

	case "mhtml", FileFormatHTML:
		return FileFormatHTML
	}

	return f
}
