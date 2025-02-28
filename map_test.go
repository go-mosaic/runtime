package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

func TestSplitKeyValInt16(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V constraints.Integer] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[int16]{
		{
			name: "split key/value int16 empty",
			args: args{
				s: "",
			},
			wantResult: map[string]int16{},
			wantErr:    false,
		},
		{
			name: "split key/value int16 parse error",
			args: args{
				s: "a=dvssdvsd",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value int16 format error",
			args: args{
				s: "a",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value int16",
			args: args{
				s: "a=12312;b=13312",
			},
			wantResult: map[string]int16{
				"a": 12312,
				"b": 13312,
			},
			wantErr: false,
		},
		{
			name: "split key/value int16 failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},

			wantErr: true,
		},
		{
			name: "split key/value int16 failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]int16
			err := SplitKeyValInt[int16](tt.args.s, ";", "=", 10, 64, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValInt() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValInt8(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V constraints.Integer] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[int8]{
		{
			name: "test int8",
			args: args{
				s: "a=127;b=126",
			},
			wantResult: map[string]int8{
				"a": 127,
				"b": 126,
			},
			wantErr: false,
		},
		{
			name: "test int8 failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},

			wantErr: true,
		},
		{
			name: "test int8 failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]int8
			err := SplitKeyValInt[int8](tt.args.s, ";", "=", 10, 64, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValInt() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValInt(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V constraints.Integer] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[int]{
		{
			name: "test int",
			args: args{
				s: "a=8888888;b=555555",
			},
			wantResult: map[string]int{
				"a": 8888888,
				"b": 555555,
			},
			wantErr: false,
		},
		{
			name: "test int failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},

			wantErr: true,
		},
		{
			name: "test int failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]int
			err := SplitKeyValInt[int](tt.args.s, ";", "=", 10, 64, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValInt() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValUint(t *testing.T) {
	type args struct {
		s       string
		sep     string
		sepKV   string
		base    int
		bitSize int
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]uint
		wantErr    bool
	}{
		{
			name: "split key/value uint empty",
			args: args{
				s:     "",
				sep:   ";",
				sepKV: "=",
			},

			wantResult: map[string]uint{},
			wantErr:    false,
		},
		{
			name: "split key/value uint parse error",
			args: args{
				s:     "a=dvssdvsd",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value uint format error",
			args: args{
				s:     "a",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split uint",
			args: args{
				s:       "a=2;b=8",
				sep:     ";",
				sepKV:   "=",
				base:    10,
				bitSize: 64,
			},
			wantResult: map[string]uint{
				"a": 2,
				"b": 8,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]uint
			err := SplitKeyValUint[uint](tt.args.s, tt.args.sep, tt.args.sepKV, tt.args.base, tt.args.bitSize, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValUint() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValFloat32(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V constraints.Float] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[float32]{
		{
			name: "split key/value float32 empty",
			args: args{
				s: "",
			},
			wantResult: map[string]float32{},
			wantErr:    false,
		},
		{
			name: "split key/value float32 parse error",
			args: args{
				s: "a=qwerty",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value float32 format error",
			args: args{
				s: "a",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "test int",
			args: args{
				s: "a=888.888;b=555.555",
			},
			wantResult: map[string]float32{
				"a": 888.888,
				"b": 555.555,
			},
			wantErr: false,
		},
		{
			name: "test int failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},

			wantErr: true,
		},
		{
			name: "test int failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]float32
			err := SplitKeyValFloat[float32](tt.args.s, ";", "=", 64, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValFloat() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValUUID(t *testing.T) {
	type args struct {
		s     string
		sep   string
		sepKV string
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]uuid.UUID
		wantErr    bool
	}{
		{
			name: "split key/value uuid empty",
			args: args{
				s:     "",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: map[string]uuid.UUID{},
			wantErr:    false,
		},
		{
			name: "split key/value uuid parse error",
			args: args{
				s:     "a=8F853762-D357-47DD-8331-C8084FA4FA6",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value uuid format error",
			args: args{
				s:     "a",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value uuid success",
			args: args{
				s:     "a=8F853762-D357-47DD-8331-C8084FA4FA6C;b=E606C4F8-5A57-415D-83B9-B358C886557F",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: map[string]uuid.UUID{
				"a": must(uuid.Parse("8F853762-D357-47DD-8331-C8084FA4FA6C")),
				"b": must(uuid.Parse("E606C4F8-5A57-415D-83B9-B358C886557F")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]uuid.UUID
			err := SplitKeyValUUID(tt.args.s, tt.args.sep, tt.args.sepKV, uuid.Parse, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValUUID() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValTime(t *testing.T) {
	type args struct {
		s      string
		layout string
	}
	type testCase[V time.Time] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[time.Time]{
		{
			name: "split key/value time success",
			args: args{
				s:      "a=2023-03-03T13:30:00Z;b=2023-03-03T13:30:01Z",
				layout: time.RFC3339,
			},
			wantResult: map[string]time.Time{
				"a": time.Date(2023, time.March, 3, 13, 30, 0, 0, time.UTC),
				"b": time.Date(2023, time.March, 3, 13, 30, 1, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "split key/value time empty",
			args: args{
				s: "",
			},
			wantResult: map[string]time.Time{},
			wantErr:    false,
		},
		{
			name: "split key/value time failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},
			wantErr: true,
		},
		{
			name: "split key/value time failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]time.Time
			err := SplitKeyValTime[time.Time](tt.args.s, ";", "=", tt.args.layout, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValTime() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValDuration(t *testing.T) {
	type args struct {
		s string
	}
	type testCase[V time.Duration] struct {
		name       string
		args       args
		wantResult map[string]V
		wantErr    bool
	}
	tests := []testCase[time.Duration]{
		{
			name: "split key/value duration success",
			args: args{
				s: "a=1h1m1s;b=2h2m2s",
			},
			wantResult: map[string]time.Duration{
				"a": must(time.ParseDuration("1h1m1s")),
				"b": must(time.ParseDuration("2h2m2s")),
			},
			wantErr: false,
		},
		{
			name: "split key/value duration empty",
			args: args{
				s: "",
			},
			wantResult: map[string]time.Duration{},
			wantErr:    false,
		},
		{
			name: "split key/value duration failed kv sep",
			args: args{
				s: "a=12312;b-13312",
			},
			wantErr: true,
		},
		{
			name: "split key/value duration failed item sep",
			args: args{
				s: "a=12312,b=13312",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]time.Duration
			err := SplitKeyValDuration[time.Duration](tt.args.s, ";", "=", &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSplitKeyValDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("TestSplitKeyValDuration() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestSplitKeyValString(t *testing.T) {
	type args struct {
		s     string
		sep   string
		sepKV string
	}
	tests := []struct {
		name       string
		args       args
		wantResult map[string]string
		wantErr    bool
	}{
		{
			name: "split key/value string empty",
			args: args{
				s:     "",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: map[string]string{},
			wantErr:    false,
		},
		{
			name: "split key/value string format error",
			args: args{
				s:     "a",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "split key/value string success",
			args: args{
				s:     "a=test;b=bar",
				sep:   ";",
				sepKV: "=",
			},
			wantResult: map[string]string{
				"a": "test",
				"b": "bar",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotResult map[string]string
			err := SplitKeyValString[string](tt.args.s, tt.args.sep, tt.args.sepKV, &gotResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitKeyValString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("SplitKeyValString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValInt(t *testing.T) {
	type args struct {
		values map[string]int
		sep    string
		sepKV  string
		base   int
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join int",
			args: args{
				values: map[string]int{"a": 100, "b": 500, "c": 1000, "d": -100},
				sep:    ";",
				sepKV:  "=",
				base:   10,
			},
			wantResult: "a=100;b=500;c=1000;d=-100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValInt(tt.args.values, tt.args.sep, tt.args.sepKV, tt.args.base); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValInt() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValUint(t *testing.T) {
	type args struct {
		values map[string]uint
		sep    string
		sepKV  string
		base   int
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join uint",
			args: args{
				values: map[string]uint{"a": 1, "b": 200},
				sep:    ";",
				sepKV:  "=",
				base:   10,
			},
			wantResult: "a=1;b=200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValUint(tt.args.values, tt.args.sep, tt.args.sepKV, tt.args.base); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValUint() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValFloat(t *testing.T) {
	type args struct {
		values  map[string]float32
		sep     string
		sepKV   string
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
			name: "join float32",
			args: args{
				values:  map[string]float32{"a": 34.523, "b": 69.4},
				sep:     ";",
				sepKV:   "=",
				fmt:     'f',
				prec:    2,
				bitSize: 64,
			},
			wantResult: "a=34.52;b=69.40",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValFloat(tt.args.values, tt.args.sep, tt.args.sepKV, tt.args.fmt, tt.args.prec, tt.args.bitSize); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValFloat() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValString(t *testing.T) {
	type args struct {
		values map[string]string
		sep    string
		sepKV  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join string",
			args: args{
				values: map[string]string{"a": "test string", "b": "bar"},
				sep:    ";",
				sepKV:  "=",
			},
			wantResult: "a=test string;b=bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValString(tt.args.values, tt.args.sep, tt.args.sepKV); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValTime(t *testing.T) {
	type args struct {
		values map[string]time.Time
		sep    string
		sepKV  string
		layout string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join time",
			args: args{
				values: map[string]time.Time{
					"a": time.Date(1985, time.April, 2, 0, 0, 0, 0, time.UTC),
					"b": time.Date(1985, time.April, 2, 14, 59, 0, 0, time.UTC),
				},
				sep:    ";",
				sepKV:  "=",
				layout: time.RFC3339,
			},
			wantResult: "a=1985-04-02T00:00:00Z;b=1985-04-02T14:59:00Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValTime(tt.args.values, tt.args.sep, tt.args.sepKV, tt.args.layout); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValTime() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValDuration(t *testing.T) {
	type args struct {
		values map[string]time.Duration
		sep    string
		sepKV  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join duration",
			args: args{
				values: map[string]time.Duration{"a": time.Minute * 5, "b": time.Second * 3920},
				sep:    ";",
				sepKV:  "=",
			},
			wantResult: "a=5m0s;b=1h5m20s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValDuration(tt.args.values, tt.args.sep, tt.args.sepKV); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValDuration() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestJoinKeyValUUID(t *testing.T) {
	uuidTestA := uuid.Must(uuid.NewRandom())
	uuidTestB := uuid.Must(uuid.NewRandom())
	type args struct {
		values map[string]uuid.UUID
		sep    string
		sepKV  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "join uuid",
			args: args{
				values: map[string]uuid.UUID{"a": uuidTestA, "b": uuidTestB},
				sep:    ";",
				sepKV:  "=",
			},
			wantResult: "a=" + uuidTestA.String() + ";b=" + uuidTestB.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := JoinKeyValUUID(tt.args.values, tt.args.sep, tt.args.sepKV); gotResult != tt.wantResult {
				t.Errorf("JoinKeyValUUID() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
