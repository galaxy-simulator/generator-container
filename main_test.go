package main

import (
	"testing"
)

func Test_netNFW(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "small values",
			args: args{
				x: 10,
				y: 20,
				z: 30,
			},
			want: 1440.368257365208,
		},
		{
			name: "negative values",
			args: args{
				x: -30,
				y: -40,
				z: -530,
			},
			want: 1043.5804324231447,
		},
		{
			name: "big values",
			args: args{
				x: 10000,
				y: 200000,
				z: 5,
			},
			want: 0.015581605046317826,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := netNFW(tt.args.x, tt.args.y, tt.args.z); got != tt.want {
				t.Errorf("netNFW() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gen(t *testing.T) {
	type args struct {
		galaxyRange float64
	}
	tests := []struct {
		name string
		args args
		want point
	}{
		{
			name: "Generate a single star (range=1e4)",
			args: args{
				galaxyRange: 1e4,
			},
			want: point{
				x: 0,
				y: 0,
				z: 0,
			},
		},
		{
			name: "Generate a single star (range=1e5)",
			args: args{
				galaxyRange: 1e5,
			},
			want: point{
				x: 0,
				y: 0,
				z: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen(tt.args.galaxyRange)
			if got == (point{}) {
				t.Errorf("gen() = %v, want %v", got, tt.want)
			}
		})
	}
}
