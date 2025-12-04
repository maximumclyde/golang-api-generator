package utils

import (
	"reflect"

	"gorm.io/gorm"
)

/*
Attaches multiple query conditions to a gorm database reference.

For specific querying that doesn't rely only on equality, we can use an array
with the condition as the first parameter and the bindings as the rest of the parameters

@example

dbWithConditions := AttachQueryConditions(db, ["name = <> ? and created_at < ?", "Marc", "<some-date>"], model{Industry: "<some-industry>"})
*/
func AttachQueryConditions(db *gorm.DB, query ...any) *gorm.DB {
	if db == nil || len(query) == 0 {
		return db
	}

	conditions := db
	for _, q := range query {
		// if the type is a slice, we destructure the conditions for the update
		if reflect.TypeOf(q).Kind() == 23 {
			// we put the first argument as the main syntax for the query and destructure the bindings
			qC := q.([]any)
			if len(qC) == 0 {
				continue
			}

			if len(qC) > 1 {
				conditions = conditions.Where(qC[0], qC[1:]...)
			} else {
				conditions = conditions.Where(qC[0])
			}
		} else {
			// in case we are not dealing with an array, we can simply add the condition to the query
			conditions = conditions.Where(q)
		}
	}
	return conditions
}
