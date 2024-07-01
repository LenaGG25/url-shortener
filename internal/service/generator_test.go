package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLGenerator_GenerateShortURL(t *testing.T) {
	t.Parallel()

	type fields struct {
		maxShortKeySize int
	}
	type args struct {
		originalURL string
	}
	testCases := []struct {
		name     string
		fields   fields
		args     args
		expected string
	}{
		{
			name: "short key size 10",
			fields: fields{
				maxShortKeySize: 10,
			},
			args: args{
				originalURL: "https://www.google.ru/",
			},
			expected: "5JulPkzrEf",
		},
		{
			name: "empty url",
			fields: fields{
				maxShortKeySize: 10,
			},
			args: args{
				originalURL: "",
			},
			expected: "47DEQpj8HB",
		},
		{
			name: "short key size 6",
			fields: fields{
				maxShortKeySize: 6,
			},
			args: args{
				originalURL: "https://www.google.ru/",
			},
			expected: "5JulPk",
		},
		{
			name: "long original url",
			fields: fields{
				maxShortKeySize: 4,
			},
			args: args{
				originalURL: "https://www.ozon.ru/ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
			},
			expected: "xfWA",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			g := &URLGenerator{
				maxShortKeySize: test.fields.maxShortKeySize,
			}

			shortURL := g.GenerateShortURL(test.args.originalURL)

			assert.Equal(t, test.expected, shortURL)
		})
	}
}
