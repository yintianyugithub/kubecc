package script

import (
	"github.com/zeromicro/go-zero/core/fx"
	"go.uber.org/ratelimit"
	"testing"
)

func TestFx(t *testing.T) {
	slices := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slicesRate := ratelimit.New(10, ratelimit.WithoutSlack)
	fx.From(func(source chan<- any) {
		for _, it := range slices {
			source <- it
		}
	}).Walk(func(item any, pipe chan<- any) {
		slicesRate.Take()
		if item.(int) == 1 || item.(int) == 2 {
			pipe <- item
			return
		}
		t.Log(item.(int))
	}).Split(2).ForEach(func(item any) {
		itemSplit := item.([]any)
		items := make([]int, 0, len(itemSplit))
		for _, it := range itemSplit {
			items = append(items, it.(int))
		}
		t.Log(items)
	})
}
