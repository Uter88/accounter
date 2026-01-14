package renderers

import (
	"accounter/internal/domain/task"
	"accounter/pkg/utils"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"text/template"
)

var (
	ErrNotImplemented    = errors.New("not implemented")
	ErrUnsupportedFormat = errors.New("format is not supported")
)

// Tasks export data model
type exportData struct {
	PricePerHour float32
	TotalPrice   float32
	Rows         []exportRow
}

// Tasks export data row
type exportRow struct {
	Index       int     `json:"-"`
	Description string  `json:"description"`
	TaskID      string  `json:"task_id"`
	Date        string  `json:"date"`
	WorkBegin   string  `json:"work_begin"`
	WorkEnd     string  `json:"work_end"`
	Price       float32 `json:"price"`
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

	d.TotalPrice = utils.ToFixed(d.TotalPrice, 2)

	if l := float32(len(d.Rows)); l > 0 {
		d.PricePerHour = utils.ToFixed(d.PricePerHour/l, 2)
	}

	return
}

// NewTaskRenderer creates new taskRenderer
func NewTaskRenderer() *taskRenderer {
	renderer := &taskRenderer{
		renderers: make(map[utils.FileFormat]renderer),
	}

	return renderer.registerDefaultFormats()
}

// Renderer of Task
type taskRenderer struct {
	renderers map[utils.FileFormat]renderer
}

// data renderer
type renderer = func(data exportData) (*bytes.Buffer, error)

// registerDefaultFormats register supported formats for rendering Task
func (tr *taskRenderer) registerDefaultFormats() *taskRenderer {
	tr.RegisterFormat(utils.FileFormatCSV, tr.toCSV)
	tr.RegisterFormat(utils.FileFormatDocX, tr.toDOCX)
	tr.RegisterFormat(utils.FileFormatXLSX, tr.toXLSX)
	tr.RegisterFormat(utils.FileFormatPDF, tr.toPDF)
	tr.RegisterFormat(utils.FileFormatJSON, tr.toJSON)
	tr.RegisterFormat(utils.FileFormatHTML, tr.toHTML)
	tr.RegisterFormat(utils.FileFormatXML, tr.toXML)

	return tr
}

// RegisterFormat register new render format
func (tr *taskRenderer) RegisterFormat(format utils.FileFormat, renderer renderer) {
	tr.renderers[format] = renderer
}

// Render Tasks to specified format
func (tr *taskRenderer) Render(format utils.FileFormat, tasks task.Tasks) (*bytes.Buffer, error) {
	data := tasksToExportData(tasks)

	if renderer, ok := tr.renderers[format]; ok {
		return renderer(data)
	}

	return nil, fmt.Errorf("%w: want: %s, have: %s", ErrUnsupportedFormat, utils.MapKeys(tr.renderers), format)

}

// toCSV render data to CSV
func (r *taskRenderer) toCSV(data exportData) (*bytes.Buffer, error) {
	return nil, ErrNotImplemented
}

// toPDF render data to PDF
func (r *taskRenderer) toPDF(data exportData) (*bytes.Buffer, error) {
	return nil, ErrNotImplemented
}

// toDOCX render data to DOCX
func (r *taskRenderer) toDOCX(data exportData) (*bytes.Buffer, error) {
	return nil, ErrNotImplemented
}

// toXML render data to XML
func (r *taskRenderer) toXML(data exportData) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := xml.NewEncoder(buf).Encode(data)

	return buf, err
}

// toJSON render data to JSON
func (r *taskRenderer) toJSON(data exportData) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(data)

	return buf, err
}

// toXLSX render data to XLSX
func (r *taskRenderer) toXLSX(data exportData) (*bytes.Buffer, error) {
	return nil, ErrNotImplemented
}

// toHTML render data to HTML
func (r *taskRenderer) toHTML(data exportData) (*bytes.Buffer, error) {
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
