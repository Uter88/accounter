package models

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type SeriesType string

const (
	SeriesBar  SeriesType = "bar"
	SeriesLine SeriesType = "line"
	SeriesPie  SeriesType = "pie"
)

var defaultTextStyle = TextStyle{
	LineHeight: 18,
	FontWeight: 400,
	FontFamily: "RobotoRegular",
	FontSize:   11,
	Color:      "#545454",
}

var defaultLineStyle = LineStyle{
	Width: 2,
	Color: "#EFEFEF",
}

func NewChartOptions() *ChartOptions {
	return &ChartOptions{
		Title: Title{
			TextStyle: defaultTextStyle,
		},
		Tooltip: Tooltip{
			Show:      true,
			Trigger:   "axis",
			Position:  "top",
			ClassName: "tooltip-custom",
			TextStyle: defaultTextStyle,
		},
	}
}

type ChartOptions struct {
	Title     Title     `json:"title"`
	Tooltip   Tooltip   `json:"tooltip"`
	Legend    Legend    `json:"legend"`
	XAxis     XAxis     `json:"xAxis"`
	YAxis     YAxis     `json:"yAxis"`
	Series    []Series  `json:"series"`
	TextStyle TextStyle `json:"textStyle"`
}

type Title struct {
	Show      bool      `json:"show,omitempty"`
	Text      string    `json:"text,omitempty"`
	Top       string    `json:"top,omitempty"`
	Left      string    `json:"left,omitempty"`
	TextStyle TextStyle `json:"textStyle"`
}

type Tooltip struct {
	Show      bool      `json:"show,omitempty"`
	Trigger   string    `json:"trigger,omitempty"`
	Position  string    `json:"position,omitempty"`
	ClassName string    `json:"className,omitempty"`
	TextStyle TextStyle `json:"textStyle"`
}

type TextStyle struct {
	LineHeight int    `json:"lineHeight,omitempty"`
	FontWeight int    `json:"fontWeight,omitempty"`
	FontFamily string `json:"fontFamily,omitempty"`
	FontSize   int    `json:"fontSize,omitempty"`
	Color      string `json:"color,omitempty"`
}

type Legend struct {
	Data []string `json:"data,omitempty"`
}

type XAxis struct {
	Type        string    `json:"type,omitempty"`
	Name        string    `json:"name,omitempty"`
	Min         float64   `json:"min,omitempty"`
	Max         float64   `json:"max,omitempty"`
	Offset      float64   `json:"offset,omitempty"`
	SplitNumber int       `json:"splitNumer,omitempty"`
	MinInterval int       `json:"minInterval,omitempty"`
	MaxInterval int       `json:"maxInterval,omitempty"`
	Data        Array     `json:"data,omitempty"`
	AxisLabel   TextStyle `json:"axisLabel"`
	AxisLine    AxisLine  `json:"axisLine"`
}

func NewXAxis(tp string) XAxis {
	return XAxis{
		Type: tp,
		AxisLine: AxisLine{
			Show:      true,
			LineStyle: defaultLineStyle,
		},
		AxisLabel: defaultTextStyle,
	}
}

type YAxis struct {
	Type      string    `json:"type,omitempty"`
	AxisLabel TextStyle `json:"axisLabel"`
	AxisLine  AxisLine  `json:"axisLine"`
}

func NewYAxis(tp string) YAxis {
	return YAxis{
		Type: tp,
		AxisLine: AxisLine{
			Show:      true,
			LineStyle: defaultLineStyle,
		},
		AxisLabel: defaultTextStyle,
	}
}

type AxisLine struct {
	Show      bool      `json:"show,omitempty"`
	LineStyle LineStyle `json:"lineStyle"`
}

type LineStyle struct {
	Color string `json:"color,omitempty"`
	Width int    `json:"width,omitempty"`
}

type Series struct {
	Name  string     `json:"name,omitempty"`
	Type  SeriesType `json:"type,omitempty"`
	Data  Array      `json:"data,omitempty"`
	Label Label      `json:"label"`
}

type Label struct {
	Show      bool     `json:"show,omitempty"`
	FontSize  int      `json:"fontSize,omitempty"`
	Formatter app.Func `json:"formatter,omitempty"`
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
