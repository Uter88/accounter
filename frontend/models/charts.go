package models

import "encoding/json"

type SeriesType string

const (
	SeriesBar  SeriesType = "bar"
	SeriesLine SeriesType = "line"
	SeriesPie  SeriesType = "pie"
)

func NewChartOptions() *ChartOptions {
	return &ChartOptions{}
}

type ChartOptions struct {
	Title   Title    `json:"title"`
	Tooltip Tooltip  `json:"tooltip"`
	Legend  Legend   `json:"legend"`
	XAxis   XAxis    `json:"xAxis"`
	YAxis   YAxis    `json:"yAxis"`
	Series  []Series `json:"series"`
}

type Title struct {
	Show bool   `json:"show,omitempty"`
	Text string `json:"text,omitempty"`
}

type Tooltip struct {
	Show bool `json:"show,omitempty"`
}

type Legend struct {
	Data []string `json:"data,omitempty"`
}

type XAxis struct {
	Type        string  `json:"type,omitempty"`
	Name        string  `json:"name,omitempty"`
	Min         float64 `json:"min,omitempty"`
	Max         float64 `json:"max,omitempty"`
	Offset      float64 `json:"offset,omitempty"`
	SplitNumber int     `json:"splitNumer,omitempty"`
	MinInterval int     `json:"minInterval,omitempty"`
	MaxInterval int     `json:"maxInterval,omitempty"`
	Data        Array   `json:"data,omitempty"`
}

type YAxis struct {
	Type string `json:"type,omitempty"`
}

type Series struct {
	Name string     `json:"name,omitempty"`
	Type SeriesType `json:"type,omitempty"`
	Data Array      `json:"data,omitempty"`
}

func (o *ChartOptions) WithTitle(text string) *ChartOptions {
	o.Title.Text = text
	o.Title.Show = true
	return o
}

func (o *ChartOptions) WithLegend(legend ...string) *ChartOptions {
	o.Legend.Data = append(o.Legend.Data, legend...)
	return o
}

func (o *ChartOptions) WithXAxis(xAxis XAxis) *ChartOptions {
	o.XAxis = xAxis
	return o
}

func (o *ChartOptions) WithYAxis(yAxis YAxis) *ChartOptions {
	o.YAxis = yAxis
	return o
}

func (o *ChartOptions) WithSeries(series ...Series) *ChartOptions {
	o.Series = append(o.Series, series...)
	return o
}

func (o *ChartOptions) Encode() Object {
	result := make(Object)

	data, _ := json.Marshal(o)
	json.Unmarshal(data, &result)

	return result
}

type Params struct {
	SeriesType string
}
