package components

import (
	"accounter/adapters/ui/wasm-goap/models"
	"accounter/internal/domain/task"
	"accounter/pkg/tools"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TasksChart struct {
	app.Compo

	Tasks  task.Tasks
	Params *task.TaskParams

	lazy      bool
	chart     *Chart
	chartType models.SeriesType
	valueType string
}

func NewTasksChart(tasks task.Tasks, params *task.TaskParams) *TasksChart {
	return &TasksChart{
		lazy:      true,
		chartType: models.SeriesBar,
		valueType: "money",
		Tasks:     tasks,
		Params:    params,
		chart: NewChart("tasks-chart").
			Width("100%").
			Height("100%").
			Lazy(true),
	}
}

func (t *TasksChart) Lazy(v bool) *TasksChart {
	t.lazy = v
	return t
}

func (t *TasksChart) ChartType(tp models.SeriesType) *TasksChart {
	t.chartType = tp
	return t
}

func (t *TasksChart) ValueType(tp string) *TasksChart {
	t.valueType = tp
	return t
}

func (t *TasksChart) OnMount(ctx app.Context) {
	ctx.Handle("setTasks", func(ctx app.Context, a app.Action) {
		if tasks, ok := a.Value.(task.Tasks); ok {
			t.Tasks = tasks
			t.redraw(ctx)
		}
	})

	t.redraw(ctx)
}

func (t *TasksChart) Render() app.UI {
	return app.Div().
		Class("card p-1 d-flex h-50 flex-column align-items-center mt-3 mx-3").
		Style("position", "relative").
		Body(
			app.Div().Class("d-flex flex-row").Body(
				NewBtnIcon("timeline").
					Tooltip("Line chart").
					OnClick(func(ctx app.Context, e app.Event) {
						t.chartType = models.SeriesLine
						t.redraw(ctx)
					}),
				NewBtnIcon("bar_chart").
					Tooltip("Bar chart").
					OnClick(func(ctx app.Context, e app.Event) {
						t.chartType = models.SeriesBar
						t.redraw(ctx)
					}),

				app.Div().Class("vr"),

				NewBtnIcon("money").
					Tooltip("Money mode").
					OnClick(func(ctx app.Context, e app.Event) {
						t.valueType = "money"
						t.redraw(ctx)
					}),

				NewBtnIcon("nest_clock_farsight_analog").
					Tooltip("Time mode").
					OnClick(func(ctx app.Context, e app.Event) {
						t.valueType = "time"
						t.redraw(ctx)
					}),
			),
			t.chart,
		)
}

func (t *TasksChart) redraw(ctx app.Context) {

	var options *models.ChartOptions

	switch t.chartType {
	case models.SeriesBar:
		options = t.makeBarSeries(t.Tasks)

	case models.SeriesLine:
		options = t.makeLineSeries(t.Tasks)
	default:
		return
	}

	t.chart.SetOptions(options).Draw(ctx)
}

func (t *TasksChart) makeLineSeries(tasks task.Tasks) *models.ChartOptions {
	var (
		legends []string
		series  []models.Series
		xAxis   = models.NewXAxis("time")
		yAxis   = models.NewYAxis("value")
	)

	xAxis.Min = float64(t.Params.DateStart * 1000)
	xAxis.Max = float64(t.Params.DateEnd * 1000)

	for _, item := range tasks.GroupByUsers().Items() {
		legends = append(legends, item.Key)

		s := models.Series{
			Type: models.SeriesLine,
			Name: item.Key,
			Label: models.Label{
				Show:     true,
				FontSize: 9,
			},
		}

		for _, item := range item.Value.GroupByDates().Items() {
			var value any

			switch t.valueType {
			case "money":
				value = int(item.Value.GetPrice())
			case "time":
				value = tools.ToFixed(item.Value.GetDuration().Hours(), 1)
			}

			s.Data = append(s.Data, models.Array{item.Key * 1000, value})
		}

		series = append(series, s)
	}

	return models.NewChartOptions().
		WithTitle("Tasks chart").
		WithLegend(legends...).
		WithXAxis(xAxis).
		WithYAxis(yAxis).
		WithSeries(series...)
}

func (t *TasksChart) makeBarSeries(tasks task.Tasks) *models.ChartOptions {
	var (
		legends []string
		series  = models.Series{
			Type: models.SeriesBar,
		}
		xAxis = models.NewXAxis("category")
		yAxis = models.NewYAxis("value")
	)

	for _, item := range tasks.GroupByUsers().Items() {
		xAxis.Data = append(xAxis.Data, item.Key)
		legends = append(legends, item.Key)

		var value any

		switch t.valueType {
		case "money":
			value = item.Value.GetPrice()
		case "time":
			value = tools.ToFixed(item.Value.GetDuration().Hours(), 1)
		}

		series.Data = append(series.Data, value)
	}

	return models.NewChartOptions().
		WithTitle("Tasks chart").
		WithLegend(legends...).
		WithXAxis(xAxis).
		WithYAxis(yAxis).
		WithSeries(series)
}
