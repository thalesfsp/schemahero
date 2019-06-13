package generate

import (
	"testing"

	"github.com/schemahero/schemahero/pkg/database/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sanitizeName(t *testing.T) {
	sanitizeNameTests := map[string]string{
		"two_words": "two-words",
	}

	for unsanitized, expectedSanitized := range sanitizeNameTests {
		t.Run(unsanitized, func(t *testing.T) {
			sanitized := sanitizeName(unsanitized)
			assert.Equal(t, expectedSanitized, sanitized)
		})
	}
}

func Test_writeTableFile(t *testing.T) {
	tests := []struct {
		name         string
		driver       string
		dbName       string
		tableName    string
		primaryKey   []string
		foreignKeys  []*types.ForeignKey
		columns      []*types.Column
		expectedYAML string
	}{
		{
			name:        "pg 1 col",
			driver:      "postgres",
			dbName:      "db",
			tableName:   "simple",
			primaryKey:  []string{"one"},
			foreignKeys: []*types.ForeignKey{},
			columns: []*types.Column{
				&types.Column{
					Name:     "id",
					DataType: "integer",
				},
			},
			expectedYAML: `apiVersion: schemas.schemahero.io/v1alpha2
kind: Table
metadata:
  name: simple
spec:
  database: db
  name: simple
  schema:
    postgres:
      primaryKey:
      - one
      columns:
      - name: id
        type: integer
`,
		},
		{
			name:       "pg foreign key",
			driver:     "postgres",
			dbName:     "db",
			tableName:  "withfk",
			primaryKey: []string{"pk"},
			foreignKeys: []*types.ForeignKey{
				&types.ForeignKey{
					ChildColumns:  []string{"cc"},
					ParentTable:   "p",
					ParentColumns: []string{"pc"},
					Name:          "fk_pc_cc",
				},
			},
			columns: []*types.Column{
				&types.Column{
					Name:     "pk",
					DataType: "integer",
				},
				&types.Column{
					Name:     "cc",
					DataType: "integer",
				},
			},
			expectedYAML: `apiVersion: schemas.schemahero.io/v1alpha2
kind: Table
metadata:
  name: withfk
spec:
  database: db
  name: withfk
  schema:
    postgres:
      primaryKey:
      - pk
      foreignKeys:
      - columns:
        - cc
        references:
          table: p
          columns:
          - pc
        name: fk_pc_cc
      columns:
      - name: pk
        type: integer
      - name: cc
        type: integer
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			y, err := generateTableYAML(test.driver, test.dbName, test.tableName, test.primaryKey, test.foreignKeys, test.columns)
			req.NoError(err)
			assert.Equal(t, test.expectedYAML, y)
		})
	}
}