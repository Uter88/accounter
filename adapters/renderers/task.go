package renderers

import (
	"accounter/internal/domain/task"
	"accounter/pkg/tools"
	"bytes"
	"fmt"
	"text/template"
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

// Convert Tasks to exportData
func tasksToExportData(tasks task.Tasks) (d exportData) {
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

type taskRenderer struct{}

func NewTaskRenderer() taskRenderer {
	return taskRenderer{}
}

func (r taskRenderer) Render(format tools.FileFormat, tasks task.Tasks) (*bytes.Buffer, error) {
	data := tasksToExportData(tasks)

	switch format {
	case tools.FileFormatCSV:
		return r.ToCSV(data)

	case tools.FileFormatDocX:
		return r.ToDOCX(data)

	case tools.FileFormatXLSX:
		return r.ToXLSX(data)

	case tools.FileFormatPDF:
		return r.ToPDF(data)

	case tools.FileFormatJSON:
		return r.ToJSON(data)

	case tools.FileFormatHTML:
		return r.ToHTML(data)

	default:
		return nil, fmt.Errorf("file format %s is not supported", format)
	}
}

func (r taskRenderer) ToCSV(data exportData) (*bytes.Buffer, error) {
	return nil, nil
}

func (r taskRenderer) ToPDF(data exportData) (*bytes.Buffer, error) {
	return nil, nil
}

func (r taskRenderer) ToDOCX(data exportData) (*bytes.Buffer, error) {
	return nil, nil
}

func (r taskRenderer) ToJSON(data exportData) (*bytes.Buffer, error) {
	return nil, nil
}

func (r taskRenderer) ToXLSX(data exportData) (*bytes.Buffer, error) {
	return nil, nil
}

func (r taskRenderer) ToHTML(data exportData) (*bytes.Buffer, error) {
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
