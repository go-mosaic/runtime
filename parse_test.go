package runtime

import (
	"net/url"
	"testing"
	"time"

	"golang.org/x/exp/constraints"
)

func TestParseBool(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[T constraints.Signed] struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "parse bool success",
			args: args{
				s: "true",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "parse bool error",
			args: args{
				s: "true2",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			err := ParseBool(tt.args.s, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	type args struct {
		s       string
		base    int
		bitSize int
	}
	type testCase[T constraints.Signed] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "test int",
			args: args{
				s:       "10",
				base:    10,
				bitSize: 64,
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int
			err := ParseInt(tt.args.s, tt.args.base, tt.args.bitSize, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt64(t *testing.T) {
	type args struct {
		s       string
		base    int
		bitSize int
	}
	type testCase[T constraints.Signed] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[int64]{
		{
			name: "test int64",
			args: args{
				s:       "10",
				base:    10,
				bitSize: 64,
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int64
			err := ParseInt[int64](tt.args.s, tt.args.base, tt.args.bitSize, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	type args struct {
		s       string
		bitSize int
	}
	type testCase[T constraints.Float] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[float32]{
		{
			name: "parse float32 success",
			args: args{
				s:       "10.1",
				bitSize: 64,
			},
			want:    10.1,
			wantErr: false,
		},
		{
			name: "parse float32 error",
			args: args{
				s:       "qw",
				bitSize: 64,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got float32
			err := ParseFloat[float32](tt.args.s, tt.args.bitSize, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseComplex(t *testing.T) {
	type args struct {
		s       string
		bitSize int
	}
	type testCase[T constraints.Complex] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[complex64]{
		{
			name: "parse complex64 success",
			args: args{
				s:       "10.1",
				bitSize: 64,
			},
			want:    10.1,
			wantErr: false,
		},
		{
			name: "parse complex64 error",
			args: args{
				s:       "qw",
				bitSize: 64,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got complex64
			err := ParseComplex[complex64](tt.args.s, tt.args.bitSize, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFloat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	type args struct {
		layout string
		s      string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "passe time success",
			args: args{
				layout: time.RFC3339,
				s:      "2025-01-02T15:04:05Z",
			},
			want:    must(time.Parse(time.RFC3339, "2025-01-02T15:04:05Z")),
			wantErr: false,
		},
		{
			name: "passe time error",
			args: args{
				layout: time.RFC3339,
				s:      "2025-01-02T15:04:05",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got time.Time
			if err := ParseTime(tt.args.layout, tt.args.s, &got); (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "passe duration success",
			args: args{
				s: "1h",
			},
			want:    must(time.ParseDuration("1h")),
			wantErr: false,
		},
		{
			name: "passe duration error",
			args: args{
				s: "a",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got time.Duration
			if err := ParseDuration(tt.args.s, &got); (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "parse URL success",
			args: args{
				s: "http://555f.ru",
			},
			want:    must(url.Parse("http://555f.ru")),
			wantErr: false,
		},
		{
			name: "parse URL error",
			args: args{
				s: string([]byte{0x7f}),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *url.URL
			if err := ParseURL(tt.args.s, &got); (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil && got.String() != tt.want.String() {
				t.Errorf("ParseURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
