package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSplitUUID(t *testing.T) {
	type args struct {
		s          string
		sep        string
		uuidParser UUIDParser[uuid.UUID]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantOut []uuid.UUID
	}{
		{
			name: "split uuid success",
			args: args{
				s:          "B91D271D-5DF9-4920-A0CF-7A0DF0491932,78DF836E-45F2-4F62-BF84-5E29E69787D3",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			wantOut: []uuid.UUID{
				uuid.MustParse("B91D271D-5DF9-4920-A0CF-7A0DF0491932"),
				uuid.MustParse("78DF836E-45F2-4F62-BF84-5E29E69787D3"),
			},
			wantErr: false,
		},
		{
			name: "split uuid error",
			args: args{
				s:          "B91D271D-5DF9-4920-A0CF-7A0DF04919,78DF836E-45F2-4F62-BF84-5E29E69787D3",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			wantOut: nil,
			wantErr: true,
		},
		{
			name: "split uuid empty",
			args: args{
				s:          "",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			wantOut: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []uuid.UUID
			if err := SplitUUID(tt.args.s, tt.args.sep, tt.args.uuidParser, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitUUID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.wantOut) {
				t.Errorf("SplitUUID() not equal")
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{
			name: "split string success",
			args: args{
				s:   "1,2,3",
				sep: ",",
			},
			want:    []string{"1", "2", "3"},
			wantErr: false,
		},
		{
			name: "split string empty",
			args: args{
				s:   "",
				sep: ",",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			if err := Split(tt.args.s, tt.args.sep, &got); (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() not equal")
			}
		})
	}
}

func TestSplitDuration(t *testing.T) {
	type args struct {
		s          string
		sep        string
		uuidParser UUIDParser[uuid.UUID]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []time.Duration
	}{
		{
			name: "split duration success",
			args: args{
				s:          "1s,23m,32h3m",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			want: []time.Duration{
				must(time.ParseDuration("1s")),
				must(time.ParseDuration("23m")),
				must(time.ParseDuration("32h3m")),
			},
			wantErr: false,
		},
		{
			name: "split duration empty",
			args: args{
				s:          "",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "split duration error",
			args: args{
				s:          "qwert",
				sep:        ",",
				uuidParser: uuid.Parse,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []time.Duration
			if err := SplitDuration(tt.args.s, tt.args.sep, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitDuration() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitDuration() not equal")
			}
		})
	}
}

func TestSplitInt(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []int
	}{
		{
			name: "split int success",
			args: args{
				s:   "1,20,100,5000",
				sep: ",",
			},
			want:    []int{1, 20, 100, 5000},
			wantErr: false,
		},
		{
			name: "split int error",
			args: args{
				s:   "1,abc",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "split int empty",
			args: args{
				s:   "",
				sep: ",",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			if err := SplitInt(tt.args.s, tt.args.sep, 10, 64, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitInt() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitInt() not equal")
			}
		})
	}
}

func TestSplitFloat(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []float32
	}{
		{
			name: "split float success",
			args: args{
				s:   "3.14,45.77,300.1",
				sep: ",",
			},
			want:    []float32{3.14, 45.77, 300.1},
			wantErr: false,
		},
		{
			name: "split float error",
			args: args{
				s:   "3.14,b,300.1",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "split float empty",
			args: args{
				s:   "",
				sep: ",",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []float32
			if err := SplitFloat(tt.args.s, tt.args.sep, 32, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitFloat() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitFloat() not equal")
			}
		})
	}
}

func TestSplitUint(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []uint
	}{
		{
			name: "split uint success",
			args: args{
				s:   "0,1,2,3,4",
				sep: ",",
			},
			want:    []uint{0, 1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "split uint error",
			args: args{
				s:   "0,1,b,3,4",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "split uint empty",
			args: args{
				s:   "",
				sep: ",",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []uint
			if err := SplitUint(tt.args.s, tt.args.sep, 10, 64, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitUint() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitUint() not equal")
			}
		})
	}
}

func TestSplitTime(t *testing.T) {
	type args struct {
		s      string
		sep    string
		sepKV  string
		layout string
	}
	tests := []struct {
		name    string
		args    args
		want    []time.Time
		wantErr bool
	}{
		{
			name: "split time success",
			args: args{
				s:      "2025-01-02T15:04:05Z;2025-01-02T15:05:05Z",
				sep:    ";",
				sepKV:  "=",
				layout: time.RFC3339,
			},
			want: []time.Time{
				must(time.Parse(time.RFC3339, "2025-01-02T15:04:05Z")),
				must(time.Parse(time.RFC3339, "2025-01-02T15:05:05Z")),
			},
			wantErr: false,
		},
		{
			name: "split time error",
			args: args{
				s:      "2025-01-02T15:04:0;2025-01-02T15:05:05Z",
				sep:    ";",
				sepKV:  "=",
				layout: time.RFC3339,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "split time empty",
			args: args{
				s:      "",
				sep:    ";",
				sepKV:  "=",
				layout: time.RFC3339,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []time.Time
			if err := SplitTime(tt.args.s, tt.args.sep, tt.args.sepKV, tt.args.layout, &got); (err != nil) != tt.wantErr {
				t.Errorf("SplitTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitTime() not equal")
			}
		})
	}
}

func TestJoinInt(t *testing.T) {
	type args struct {
		values []int32
		sep    string
		base   int
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join int32 success",
			args: args{
				values: []int32{1, 2, 3, 4, 5},
				sep:    ",",
				base:   10,
			},
			wantResult: "1,2,3,4,5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinInt(tt.args.values, tt.args.sep, tt.args.base); gotResult != tt.wantResult {
				t.Errorf("JoinInt() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinFloat(t *testing.T) {
	type args struct {
		values  []float32
		sep     string
		fmt     byte
		prec    int
		bitSize int
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join float32 success",
			args: args{
				values:  []float32{10.1214, 34.34, 678.45345324},
				sep:     ",",
				fmt:     'f',
				prec:    2,
				bitSize: 64,
			},
			wantResult: "10.12,34.34,678.45",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinFloat(tt.args.values, tt.args.sep, tt.args.fmt, tt.args.prec, tt.args.bitSize); gotResult != tt.wantResult {
				t.Errorf("JoinFloat() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
