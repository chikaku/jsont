package jsont

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const caseDir = "testcase/test_parsing"

func iterTestcase(t *testing.T, prefix string, testFunc func(*assert.Assertions, []byte)) {
	ast := assert.New(t)
	dir, err := os.ReadDir(caseDir)
	ast.NoError(err)

	for _, entry := range dir {
		if !entry.IsDir() {
			filename := entry.Name()
			if strings.HasPrefix(filename, prefix) {
				f, err := os.Open(path.Join(caseDir, filename))
				ast.NoError(err)
				content, err := io.ReadAll(f)
				ast.NoError(err)

				t.Run(filename, func(t *testing.T) {
					testFunc(assert.New(t), content)
				})
			}
		}
	}
}

func TestCorrectCase(t *testing.T) {
	iterTestcase(t, "y_", func(ast *assert.Assertions, content []byte) {
		_, err := Decode(content)
		ast.NoError(err, string(content))
	})
}

func TestIncorrectCase(t *testing.T) {
	iterTestcase(t, "n_", func(ast *assert.Assertions, content []byte) {
		_, err := Decode(content)
		ast.NotNil(err, string(content))
	})
}

func TestReadElement(t *testing.T) {
	var err error
	ast := assert.New(t)

	_, _, err = ReadObject([]byte(`  {"a": "b"}...`))
	ast.NoError(err)
	_, _, err = ReadArray([]byte(`  [1, {}, ""],...`))
	ast.NoError(err)
	_, _, err = ReadString([]byte(`  "   ""...`))
	ast.NoError(err)
	_, _, err = ReadNumber([]byte(`  1.02e-10...`))
	ast.NoError(err)
	_, _, err = ReadNull([]byte(`  null\...`))
	ast.NoError(err)
	_, _, err = ReadBool([]byte(`  false\...`))
	ast.NoError(err)
}
