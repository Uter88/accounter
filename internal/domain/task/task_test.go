package task

import (
	"accounter/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// Task testing suite
type testTaskSuite struct {
	suite.Suite
	fixedTime  time.Time
	sampleTask Task
}

// TestTaskSuite run Task testing
func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(testTaskSuite))
}

// SetupTest testing preparations
func (suite *testTaskSuite) SetupTest() {
	suite.fixedTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.Local)
	suite.sampleTask = NewTask(suite.fixedTime)
	suite.sampleTask.ID = 1
}

// TestNewTask testing of Task creating
func (suite *testTaskSuite) TestNewTask() {
	now := time.Now()
	t := NewTask(now)

	suite.Empty(t.ID, "non-zero id")
	suite.Empty(t.UserID, "non-zero user id")
	suite.Empty(t.UserLabel, "non-empty user label")
	suite.Empty(t.TaskID, "non-empty task id")
	suite.Empty(t.Description, "non-empty desciption")
	suite.Empty(t.PricePerHour, "non-zero price per hour")
	suite.Equal(t.Status, TaskStatusCompleted, "incorrect status")
	suite.Equal(time.Unix(t.WorkBegin, 0).Hour(), 8, "incorrect work begin time")
	suite.Equal(time.Unix(t.WorkEnd, 0).Hour(), 18, "incorrect work end time")
	suite.Equal(t.GetDuration(), time.Hour*10, "incorrect work time duration")
	suite.Equal(t.Date, now.Unix(), "incorrect date")
}

// TestSetDate testing of set date
func (suite *testTaskSuite) TestSetDate() {
	t := suite.sampleTask

	newDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	t.SetDate(newDate)

	suite.Equal(t.Date, newDate, "incorrect set date")

	expectedWorkBegin := utils.SetDateTs(newDate, suite.sampleTask.WorkBegin)
	expectedWorkEnd := utils.SetDateTs(newDate, suite.sampleTask.WorkEnd)

	suite.Equal(expectedWorkBegin, t.WorkBegin, "incorrect work begin date")
	suite.Equal(expectedWorkEnd, t.WorkEnd, "incorrect work end date")
}

// TestIsValid testing of Task validation
func (suite *testTaskSuite) TestIsValid() {
	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name: "Valid",
			task: Task{
				UserID:      1,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   suite.fixedTime.Unix(),
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "Test",
				Status:      TaskStatusInProgress,
			},
			expected: true,
		},
		{
			name: "Invalid - missing UserID",
			task: Task{
				UserID:      0,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   suite.fixedTime.Unix(),
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "Test",
				Status:      TaskStatusInProgress,
			},
			expected: false,
		},
		{
			name: "Invalid - missing Date",
			task: Task{
				UserID:      100,
				Date:        0,
				WorkBegin:   suite.fixedTime.Unix(),
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "Test",
				Status:      TaskStatusInProgress,
			},
			expected: false,
		},
		{
			name: "Invalid - missing WorkBegin",
			task: Task{
				UserID:      100,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   0,
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "Test",
				Status:      TaskStatusInProgress,
			},
			expected: false,
		},
		{
			name: "Invalid - missing Description",
			task: Task{
				UserID:      100,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   suite.fixedTime.Unix(),
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "",
				Status:      TaskStatusInProgress,
			},
			expected: false,
		},
		{
			name: "Invalid - missing Status",
			task: Task{
				UserID:      100,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   suite.fixedTime.Unix(),
				WorkEnd:     suite.fixedTime.Add(1 * time.Hour).Unix(),
				Description: "Test",
				Status:      "",
			},
			expected: false,
		},
		{
			name: "Invalid - WorkEnd before WorkBegin",
			task: Task{
				UserID:      100,
				Date:        suite.fixedTime.Unix(),
				WorkBegin:   suite.fixedTime.Add(2 * time.Hour).Unix(),
				WorkEnd:     suite.fixedTime.Unix(),
				Description: "Test",
				Status:      TaskStatusInProgress,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.Equal(tt.expected, tt.task.IsValid())
		})
	}
}

// TestGetPrice test of price calculation of Task
func (suite *testTaskSuite) TestGetPrice() {
	tests := []struct {
		name     string
		task     Task
		expected float32
	}{
		{
			name: "2 hours at 50/hour",
			task: Task{
				PricePerHour: 50.0,
				WorkBegin:    suite.fixedTime.Unix(),
				WorkEnd:      suite.fixedTime.Add(2 * time.Hour).Unix(),
			},
			expected: 100.0,
		},
		{
			name: "1.5 hours at 30/hour",
			task: Task{
				PricePerHour: 30.0,
				WorkBegin:    suite.fixedTime.Unix(),
				WorkEnd:      suite.fixedTime.Add(90 * time.Minute).Unix(),
			},
			expected: 45.0,
		},
		{
			name: "Zero duration",
			task: Task{
				PricePerHour: 100.0,
				WorkBegin:    suite.fixedTime.Unix(),
				WorkEnd:      suite.fixedTime.Unix(),
			},
			expected: 0.0,
		},
		{
			name: "Zero price per hour",
			task: Task{
				PricePerHour: 0.0,
				WorkBegin:    suite.fixedTime.Unix(),
				WorkEnd:      suite.fixedTime.Add(2 * time.Hour).Unix(),
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.Equal(tt.expected, tt.task.GetPrice())
		})
	}
}

// TestGetDuration тестирует расчет продолжительности
func (suite *testTaskSuite) TestGetDuration() {
	task := Task{
		WorkBegin: suite.fixedTime.Unix(),
		WorkEnd:   suite.fixedTime.Add(3*time.Hour + 30*time.Minute).Unix(),
	}

	expected := 3*time.Hour + 30*time.Minute
	suite.Equal(expected, task.GetDuration())
}

// TestFormatMethods тестирует методы форматирования
func (suite *testTaskSuite) TestFormatMethods() {
	task := Task{
		PricePerHour: 75.5,
		WorkBegin:    suite.fixedTime.Unix(),
		WorkEnd:      suite.fixedTime.Add(2 * time.Hour).Unix(),
		Date:         suite.fixedTime.Unix(),
	}

	// Test FormatPrice
	suite.Equal("151.00", task.FormatPrice())

	// Test FormatDate
	suite.Equal("01.01.2026", task.FormatDate())

	// Test FormatWorkBegin
	suite.Equal("12:00", task.FormatWorkBegin())

	// Test FormatWorkEnd
	suite.Equal("14:00", task.FormatWorkEnd())

	// Test FormatDuration
	suite.Equal("02:00:00", task.FormatDuration(true))
	suite.Equal("02:00", task.FormatDuration(false))
}

// TestGetIDAndType тестирует методы GetID и GetType
func (suite *testTaskSuite) TestGetIDAndType() {
	task := suite.sampleTask

	suite.Equal(int64(1), task.GetID())
	suite.Equal("task", task.GetType())
}
