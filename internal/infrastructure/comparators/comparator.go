package comparator

import (
	"accounter/internal/domain/event"
	"reflect"
	"slices"
	"strings"

	"github.com/r3labs/diff/v3"
)

// Comparator of objects differences
type comparator struct {
	tags   []string
	fields []string
}

// NewComparator creates new comporator
func NewComparator() *comparator {
	return &comparator{}
}

func (c *comparator) Fields(fields ...string) *comparator {
	c.fields = fields
	return c
}

func (c *comparator) Tags(tags ...string) *comparator {
	c.tags = tags
	return c
}

// Compare two objects
func (c *comparator) Compare(old, new any) (result []event.EventUpdateRecord, ok bool) {
	opts := []func(d *diff.Differ) error{}

	for _, tag := range c.tags {
		opts = append(opts, diff.TagName(tag))
	}

	opts = append(opts, diff.Filter(c.filter))
	changes, _ := diff.Diff(old, new, opts...)

	for _, ch := range changes {
		if len(ch.Path) > 0 {
			key := strings.Join(ch.Path, ".")

			result = append(result, event.EventUpdateRecord{
				Field:    key,
				OldValue: ch.From,
				NewValue: ch.To,
			})
			ok = true
		}
	}

	return
}

// filter tags for detect changes in field
func (c *comparator) filter(path []string, parent reflect.Type, field reflect.StructField) bool {
	if len(c.fields) > 0 {
		return slices.Contains(c.fields, "*") || slices.Contains(c.fields, field.Name)
	}

	for _, t := range c.tags {
		tag := field.Tag.Get(t)

		if tag == "-" {
			return false
		}

		if strings.Contains(t, "omitempty") {
			return false
		}
	}

	return true
}
