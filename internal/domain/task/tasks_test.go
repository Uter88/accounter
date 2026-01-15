package task

import (
	"accounter/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// TestTasksSuite creates new testTasksSuite
func TestTasksSuite(t *testing.T) {
	suite.Run(t, new(testTasksSuite))
}

// TestTasksSuite test suite for Tasks
type testTasksSuite struct {
	suite.Suite
	tasks         Tasks
	user1Task     Task
	user2Task     Task
	todayTask     Task
	yesterdayTask Task
}

// SetupTest test Tasks preparations
func (suite *testTasksSuite) SetupTest() {
	baseTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local)

	suite.user1Task = Task{
		ID:           1,
		UserID:       100,
		UserLabel:    "user_100",
		PricePerHour: 50.0,
		TaskID:       "TASK-001",
		Status:       TaskStatusInProgress,
		Description:  "Task 1",
		WorkBegin:    baseTime.Add(8 * time.Hour).Unix(),
		WorkEnd:      baseTime.Add(10 * time.Hour).Unix(),
		Date:         baseTime.Unix(),
	}

	suite.user2Task = Task{
		ID:           2,
		UserID:       200,
		UserLabel:    "user_200",
		PricePerHour: 60.0,
		TaskID:       "TASK-002",
		Status:       TaskStatusCompleted,
		Description:  "Task 2",
		WorkBegin:    baseTime.Add(9 * time.Hour).Unix(),
		WorkEnd:      baseTime.Add(12 * time.Hour).Unix(),
		Date:         baseTime.Unix(),
	}

	suite.todayTask = suite.user1Task

	suite.yesterdayTask = Task{
		ID:           3,
		UserID:       100,
		UserLabel:    "user_100",
		PricePerHour: 50.0,
		TaskID:       "TASK-003",
		Status:       TaskStatusInProgress,
		Description:  "Task 3",
		WorkBegin:    baseTime.Add(-24*time.Hour + 8*time.Hour).Unix(),
		WorkEnd:      baseTime.Add(-24*time.Hour + 11*time.Hour).Unix(),
		Date:         baseTime.Add(-24 * time.Hour).Unix(),
	}

	suite.tasks = Tasks{suite.user1Task, suite.user2Task, suite.yesterdayTask}
}

// TestTasksEmptyAndLen testing Empty and Len methods
func (suite *testTasksSuite) TestTasksEmptyAndLen() {
	suite.False(suite.tasks.Empty())
	suite.Equal(3, suite.tasks.Len())

	emptyTasks := Tasks{}
	suite.True(emptyTasks.Empty())
	suite.Equal(0, emptyTasks.Len())

	var nilTasks Tasks
	suite.True(nilTasks.Empty())
	suite.Equal(0, nilTasks.Len())
}

// TestTasksGetPrice testing price calculations
func (suite *testTasksSuite) TestTasksGetPrice() {
	expected := float32(430.0)
	actual := suite.tasks.GetPrice()

	suite.InEpsilon(expected, actual, 0.001)
	suite.Equal(expected, utils.ToFixed(expected, 2))
}

// TestTasksGetDuration testing common duration calculations
func (suite *testTasksSuite) TestTasksGetDuration() {
	expected := 8 * time.Hour
	actual := suite.tasks.GetDuration()

	suite.Equal(expected, actual)
}

// TestTasksGroupByUsers testing of grouping by user
func (suite *testTasksSuite) TestTasksGroupByUsers() {
	grouped := suite.tasks.GroupByUsers()

	suite.Equal(2, grouped.Len())

	keys := grouped.Keys()
	suite.Equal([]string{"user_100", "user_200"}, keys)

	user100Tasks, exists := grouped.Get("user_100")
	suite.True(exists)
	suite.Len(user100Tasks, 2)
	suite.Equal(int64(1), user100Tasks[0].ID)
	suite.Equal(int64(3), user100Tasks[1].ID)

	user200Tasks, exists := grouped.Get("user_200")
	suite.True(exists)
	suite.Len(user200Tasks, 1)
	suite.Equal(int64(2), user200Tasks[0].ID)
}

// TestTasksGroupByDates testing of grouping by date
func (suite *testTasksSuite) TestTasksGroupByDates() {
	grouped := suite.tasks.GroupByDates()

	suite.Equal(2, grouped.Len())

	keys := grouped.Keys()
	suite.Len(keys, 2)

	yesterday := time.Date(2024, 1, 14, 0, 0, 0, 0, time.Local).Unix()
	today := time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local).Unix()

	suite.Equal(yesterday, keys[0])
	suite.Equal(today, keys[1])

	yesterdayTasks, exists := grouped.Get(yesterday)
	suite.True(exists)
	suite.Len(yesterdayTasks, 1)
	suite.Equal(int64(3), yesterdayTasks[0].ID)

	todayTasks, exists := grouped.Get(today)
	suite.True(exists)
	suite.Len(todayTasks, 2)
	suite.Equal(int64(1), todayTasks[0].ID)
	suite.Equal(int64(2), todayTasks[1].ID)
}

// TestTasksEdgeCases testing of edge cases
func (suite *testTasksSuite) TestTasksEdgeCases() {
	emptyTasks := Tasks{}
	suite.Equal(float32(0), emptyTasks.GetPrice())
	suite.Equal(time.Duration(0), emptyTasks.GetDuration())

	userGroups := emptyTasks.GroupByUsers()
	suite.Equal(0, userGroups.Len())

	dateGroups := emptyTasks.GroupByDates()
	suite.Equal(0, dateGroups.Len())

	singleTask := Tasks{suite.user1Task}
	suite.Equal(float32(100), singleTask.GetPrice())

	zeroPriceTask := Task{
		PricePerHour: 0.0,
		WorkBegin:    suite.user1Task.WorkBegin,
		WorkEnd:      suite.user1Task.WorkEnd,
	}
	zeroPriceTasks := Tasks{zeroPriceTask, zeroPriceTask}
	suite.Equal(float32(0), zeroPriceTasks.GetPrice())
}
