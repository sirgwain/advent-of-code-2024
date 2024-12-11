package advent

import "testing"

func TestDay11_blink(t *testing.T) {
	type args struct {
		stone int
		times int
	}
	tests := []struct {
		name   string
		args   args
		want   int
	}{
		{ "0 -> 1", args{stone: 0, times: 1}, 1},
		{ "1 -> 2024", args{stone: 1, times: 1}, 1},
		{ "1 -> 2024 -> 20 24", args{stone: 1, times: 2}, 2},
		{ "1 -> 2024 -> 20 24 -> 2 0 2 4", args{stone: 1, times: 3}, 4},
		{ "1 -> 2024 -> 20 24 -> 2 0 2 4 -> 4048 1 4048 8096", args{stone: 1, times: 4}, 4},
		{ "1 -> 2024 -> 20 24 -> 2 0 2 4 -> 4048 1 4048 8096 -> 40 48 2024 40 48 80 96", args{stone: 1, times: 5}, 7},
		{ "10 -> 1 0", args{stone: 10, times: 1}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Day11{
				stoneCache: make(map[day11CacheEntry]int),
			}
			if got := d.blink(tt.args.stone, tt.args.times); got != tt.want {
				t.Errorf("Day11.blink() = %v, want %v", got, tt.want)
			}
		})
	}
}
