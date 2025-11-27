package components

import (
	"accounter/domain/task"
	"accounter/frontend/models"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TasksChart struct {
	app.Compo

	Tasks  task.Tasks
	Params *task.TaskParams

	lazy      bool
	chart     *Chart
	chartType models.SeriesType
}

func NewTasksChart(tasks task.Tasks, params *task.TaskParams) *TasksChart {
	return &TasksChart{
		lazy:      true,
		chartType: models.SeriesBar,
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

func (t *TasksChart) Type(tp models.SeriesType) *TasksChart {
	t.chartType = tp
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
					OnClick(func(ctx app.Context, e app.Event) {
						t.chartType = models.SeriesLine
						t.redraw(ctx)
					}),
				NewBtnIcon("bar_chart").
					OnClick(func(ctx app.Context, e app.Event) {
						t.chartType = models.SeriesBar
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
		xAxis   = models.XAxis{
			Type: "time",
			Min:  float64(t.Params.DateStart * 1000),
			Max:  float64(t.Params.DateEnd * 1000),
		}
	)

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
			s.Data = append(s.Data, models.Array{item.Key * 1000, int(item.Value.GetPrice())})
		}

		series = append(series, s)
	}

	return models.NewChartOptions().
		WithTitle("Tasks chart").
		WithLegend(legends...).
		WithXAxis(xAxis).
		WithYAxis(models.YAxis{Type: "value"}).
		WithSeries(series...)
}

func (t *TasksChart) makeBarSeries(tasks task.Tasks) *models.ChartOptions {
	var (
		xAxis   models.Array
		legends []string
		series  = models.Series{
			Type: models.SeriesBar,
		}
	)

	for _, item := range tasks.GroupByUsers().Items() {
		xAxis = append(xAxis, item.Key)
		legends = append(legends, item.Key)

		series.Data = append(series.Data, item.Value.GetPrice())
	}

	return models.NewChartOptions().
		WithTitle("Tasks chart").
		WithLegend(legends...).
		WithXAxis(models.XAxis{Type: "category", Data: xAxis}).
		WithYAxis(models.YAxis{Type: "value"}).
		WithSeries(series)
}
