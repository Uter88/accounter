package task

import (
	"accounter/pkg/tools"
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

// Tasks export data model
type exportData struct {
	PricePerHour float32
	TotalPrice   float32
	Rows         []exportRow
}

// Tasks export data row
type exportRow struct {
	Index       int
	Description string
	TaskID      string
	Date        string
	WorkBegin   string
	WorkEnd     string
	Price       float32
}

// Export tasks to specified format
func exportTasks(tasks Tasks, format tools.FileFormat) (*bytes.Buffer, error) {
	data := tasksToExportData(tasks)

	switch format {
	case tools.FileFormatCSV:
		return data.convertToCSV()

	case tools.FileFormatDocX:
		return data.convertToDocX()

	case tools.FileFormatXLSX:
		return data.convertToXLSX()

	case tools.FileFormatPDF:
		return data.convertToPDF()

	case tools.FileFormatJSON:
		return data.convertToJSON()

	case tools.FileFormatHTML:
		return data.convertToHTML()

	default:
		return nil, errors.New("unexpected file format")
	}
}

// Convert Tasks to exportData
func tasksToExportData(tasks Tasks) (d exportData) {
	for i, t := range tasks {
		row := exportRow{
			Index:       i + 1,
			Description: t.Description,
			TaskID:      t.TaskID,
			Date:        t.FormatDate(),
			WorkBegin:   t.FormatWorkBegin(),
			WorkEnd:     t.FormatWorkEnd(),
			Price:       t.GetPrice(),
		}

		d.TotalPrice += row.Price
		d.PricePerHour += t.PricePerHour

		d.Rows = append(d.Rows, row)
	}

	d.TotalPrice = tools.ToFixed(d.TotalPrice, 2)

	if l := float32(len(d.Rows)); l > 0 {
		d.PricePerHour = tools.ToFixed(d.PricePerHour/l, 2)
	}

	return
}

func (data exportData) convertToCSV() (*bytes.Buffer, error) {
	return nil, nil
}

func (data exportData) convertToXLSX() (*bytes.Buffer, error) {
	return nil, nil
}

func (data exportData) convertToDocX() (*bytes.Buffer, error) {
	return nil, nil
}

func (data exportData) convertToPDF() (*bytes.Buffer, error) {
	return nil, nil
}

func (data exportData) convertToJSON() (*bytes.Buffer, error) {
	return nil, nil
}

// Convert exportData to HTML format
func (data exportData) convertToHTML() (*bytes.Buffer, error) {
	temp, err := template.ParseFiles("./templates/export_tasks.html")

	if err != nil {
		return nil, fmt.Errorf("error parse template file: %s", err.Error())
	}

	buf := bytes.NewBuffer(nil)

	if err := temp.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("error execute HTML template: %s", err.Error())
	}

	return buf, nil
}
