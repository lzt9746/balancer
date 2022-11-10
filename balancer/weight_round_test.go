package balancer

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewWeightRoundRobin(t *testing.T) {
	type args struct {
		hosts []string
	}
	tests := []struct {
		name string
		args args
		want Balancer
	}{
		{
			name: "test-1",
			args: args{
				[]string{
					"http://192.168.1.1:8000 1",
					"http://192.168.1.2:8000 2",
					"http://192.168.1.3:8000 3",
				},
			},
			want: &WeightRound{
				hosts: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				currentWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				initWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
			},
		},
		{
			name: "test-2",
			args: args{
				[]string{
					"http://192.168.1.1:8000  1",
					"http://192.168.1.2:8000  2",
					"http://192.168.1.3:8000  3",
				},
			},
			want: &WeightRound{
				hosts: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				currentWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				initWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
			},
		},
		{
			name: "test-3",
			args: args{
				[]string{
					"http://192.168.1.1:8000",
					"http://192.168.1.2:8000",
					"http://192.168.1.3:8000",
				},
			},
			want: &WeightRound{
				hosts: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 1,
					"192.168.1.3:8000": 1,
				},
				currentWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 1,
					"192.168.1.3:8000": 1,
				},
				initWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 1,
					"192.168.1.3:8000": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewWeightRoundRobin([]string{}, tt.args.hosts), "NewWeightRoundRobin(%v)", tt.args.hosts)
		})
	}
}

func TestWeightRound_Add(t *testing.T) {
	type fields struct {
		hosts         map[string]int
		currentWeight map[string]int
		initWeight    map[string]int
	}
	type args struct {
		host []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Balancer
	}{
		{
			name: "test-1",
			fields: fields{
				hosts: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				currentWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
				},
				initWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
					"192.168.1.4:8000": 4,
					"192.168.1.5:8000": 1,
					"192.168.1.6:8000": 6,
				},
			},
			args: args{
				host: []string{
					"192.168.1.4:8000",
					"192.168.1.5:8000",
					"192.168.1.6:8000",
					"192.168.1.7:8000",
				},
			},
			want: &WeightRound{
				hosts: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
					"192.168.1.4:8000": 4,
					"192.168.1.5:8000": 1,
					"192.168.1.6:8000": 6,
				},
				currentWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
					"192.168.1.4:8000": 4,
					"192.168.1.5:8000": 1,
					"192.168.1.6:8000": 6,
				},
				initWeight: map[string]int{
					"192.168.1.1:8000": 1,
					"192.168.1.2:8000": 2,
					"192.168.1.3:8000": 3,
					"192.168.1.4:8000": 4,
					"192.168.1.5:8000": 1,
					"192.168.1.6:8000": 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightRound{
				hosts:         tt.fields.hosts,
				currentWeight: tt.fields.currentWeight,
				initWeight:    tt.fields.initWeight,
			}
			for _, arg := range tt.args.host {
				r.Add(arg)
			}
			assert.Equalf(t, tt.want, r, "Add(%v)", tt.args)
		})
	}
}

func TestWeightRound_Balance(t *testing.T) {
	type args struct {
		in0 string
	}
	type want struct {
		host []string
		err  []error
	}
	tests := []struct {
		name string
		lb   Balancer
		args args
		want want
	}{
		{
			name: "test-1",
			lb:   NewWeightRoundRobin([]string{}, []string{"http://192.168.1.1:8000 1", "http://192.168.1.2:8000 2"}),
			args: args{
				in0: "",
			},
			want: want{
				[]string{
					"192.168.1.2:8000",
					"192.168.1.2:8000",
					"192.168.1.1:8000",
				},
				[]error{
					nil,
					nil,
					nil,
				},
			},
		},
		{
			name: "test-2",
			lb:   NewWeightRoundRobin([]string{}, nil),
			args: args{
				in0: "",
			},
			want: want{
				[]string{
					"",
				},
				[]error{
					NoHostError,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.lb
			for i := 0; i < 3; i++ {
				got, err := r.Balance(tt.args.in0)

				assert.Equalf(t, tt.want.host[i], got, "Balance(%v)", tt.args.in0)
				assert.Equalf(t, tt.want.err[i], err, "Balance(%v)", tt.args.in0)

			}
		})
	}
}

func TestWeightRound_Done(t *testing.T) {
	type fields struct {
		hosts         map[string]int
		currentWeight map[string]int
	}
	type args struct {
		host string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightRound{
				hosts:         tt.fields.hosts,
				currentWeight: tt.fields.currentWeight,
			}
			r.Done(tt.args.host)
		})
	}
}

func TestWeightRound_Remove(t *testing.T) {
	type fields struct {
		hosts         map[string]int
		currentWeight map[string]int
		RWMutex       sync.RWMutex
	}
	type args struct {
		host string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightRound{
				hosts:         tt.fields.hosts,
				currentWeight: tt.fields.currentWeight,
				RWMutex:       tt.fields.RWMutex,
			}
			r.Remove(tt.args.host)
		})
	}
}

func TestWeightRound_sumWeight(t *testing.T) {
	type fields struct {
		hosts         map[string]int
		currentWeight map[string]int
		RWMutex       sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightRound{
				hosts:         tt.fields.hosts,
				currentWeight: tt.fields.currentWeight,
				RWMutex:       tt.fields.RWMutex,
			}
			assert.Equalf(t, tt.want, r.sumWeight(), "sumWeight()")
		})
	}
}
