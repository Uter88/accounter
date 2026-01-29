package comparator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestComparator tests Comparator
func TestComparator(t *testing.T) {
	object1 := map[string]any{
		"name": "Alice",
		"age":  30,
	}
	object2 := map[string]any{
		"name": "Alice",
		"age":  30,
	}

	comparator := NewComparator()
	changes, ok := comparator.Compare(object1, object2)

	assert.False(t, ok)
	assert.Equal(t, len(changes), 0)

	object2["age"] = 31
	changes, ok = comparator.Compare(object1, object2)

	assert.True(t, ok)
	assert.Equal(t, len(changes), 1)
	assert.Equal(t, changes[0].Field, "age")
	assert.Equal(t, changes[0].OldValue, 30)
	assert.Equal(t, changes[0].NewValue, 31)
}

// TestComparatorWithFieldsAndTags tests Comparator with fields and tags
func TestComparatorWithFieldsAndTags(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	user1 := User{Name: "Bob", Email: "bob@example.com", Age: 25}
	user2 := User{Name: "Bob", Email: "bob@example.com", Age: 26}

	comparator := NewComparator()
	changes, ok := comparator.Compare(user1, user2)

	assert.True(t, ok)
	assert.Equal(t, len(changes), 1)
	assert.Equal(t, changes[0].Field, "Age")
	assert.Equal(t, changes[0].OldValue, 25)
	assert.Equal(t, changes[0].NewValue, 26)
}

// TestComparatorWithTags tests Comparator with tags
func TestComparatorWithNoChanges(t *testing.T) {
	type Product struct {
		ID    int     `db:"id"`
		Name  string  `db:"name"`
		Price float64 `db:"price"`
	}

	product1 := Product{ID: 1, Name: "Laptop", Price: 999.99}
	product2 := Product{ID: 1, Name: "Laptop", Price: 999.99}

	comparator := NewComparator().Tags("db")
	changes, ok := comparator.Compare(product1, product2)

	assert.False(t, ok)
	assert.Equal(t, len(changes), 0)
}

// TestComparatorWithMultipleChanges tests Comparator with multiple changes
func TestComparatorWithMultipleChanges(t *testing.T) {
	type Car struct {
		Make  string `json:"make"`
		Model string `json:"model"`
		Year  int    `json:"year"`
		Color string `json:"color"`
	}

	car1 := Car{Make: "Toyota", Model: "Camry", Year: 2020, Color: "Red"}
	car2 := Car{Make: "Toyota", Model: "Camry", Year: 2021, Color: "Blue"}

	comparator := NewComparator().Fields("Year", "Color")
	changes, ok := comparator.Compare(car1, car2)

	assert.True(t, ok)
	assert.Equal(t, len(changes), 2)

	expectedFields := map[string]struct {
		oldValue any
		newValue any
	}{
		"Year":  {oldValue: 2020, newValue: 2021},
		"Color": {oldValue: "Red", newValue: "Blue"},
	}

	for _, change := range changes {
		expected, exists := expectedFields[change.Field]
		assert.True(t, exists)
		assert.Equal(t, change.OldValue, expected.oldValue)
		assert.Equal(t, change.NewValue, expected.newValue)
	}
}
