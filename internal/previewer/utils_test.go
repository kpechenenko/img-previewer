package previewer

import (
	"testing"

	"github.com/stretchr/testify/assert" //nolint:depguard
)

func TestFilesHaveSameContent(t *testing.T) {
	tests := []struct {
		name      string
		filename1 string
		filename2 string
		sameFiles bool
		err       error
	}{
		{
			name:      "same files files",
			filename1: "./testdata/gopher_200_50.jpeg",
			filename2: "./testdata/gopher_200_50.jpeg",
			sameFiles: true,
			err:       nil,
		},
		{
			name:      "different files",
			filename1: "./testdata/gopher_200_50.jpeg",
			filename2: "./testdata/gopher_275_183.jpeg",
			sameFiles: false,
			err:       nil,
		},
		{
			name:      "file 1 does not exist",
			filename1: "./testdata/__.jpeg",
			filename2: "./testdata/gopher_275_183.jpeg",
			sameFiles: false,
			err:       assert.AnError,
		},
		{
			name:      "file 2 does not exist",
			filename1: "./testdata/gopher_200_50.jpeg",
			filename2: "./testdata/__.jpeg",
			sameFiles: false,
			err:       assert.AnError,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			equals, err := filesHaveSameContent(tc.filename1, tc.filename2)
			if tc.err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.sameFiles, equals)
			}
		})
	}
}
