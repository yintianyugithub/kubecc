package script

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/collection"
	"reflect"
)

type test struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	C string `json:"c,omitempty"`
	D string `json:"d,omitempty"`
}

func (t *test) Compress() *test {
	v := reflect.ValueOf(t).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			if s := field.String(); s == "0" || s == "-" && field.CanSet() {
				field.SetString("")
			}
		}
	}
	return t
}

func (t *test) testGroup(tasks []*test) {
	groupRes := make(map[test]*collection.Set)

	for _, tk := range tasks {
		gKey := test{A: tk.A, B: tk.B, C: tk.C}
		if _, ok := groupRes[gKey]; !ok {
			groupRes[gKey] = collection.NewSet()
		}
		groupRes[gKey].Add(tk.D)
	}

	for gedKey, v := range groupRes {
		ds := v.KeysStr()
		fmt.Printf("%+v-%q", gedKey, ds)
	}
}

func main() {
	res := make([]*test, 0)
	s := `[{"a":"0","b":"b-1","c":"c-1","d":"-"},{"a":"0","b":"b-2","c":"c-2","d":"-"},{"a":"0","b":"b-3","c":"c-3","d":"-"},{"a":"0","b":"b-4","c":"c-4","d":"-"}]`
	_ = json.Unmarshal([]byte(s), &res)

	for i, v := range res {
		res[i] = v.Compress()
	}

	b, _ := json.Marshal(res)

	fmt.Println(string(b))

	t := &test{}
	t.testGroup([]*test{
		{A: "a-1", B: "b-1", C: "c-1", D: "d-1"},
		{A: "a-1", B: "b-1", C: "c-1", D: "d-2"},
		{A: "a-1", B: "b-1", C: "c-1", D: "d-3"},
		{A: "a-2", B: "b-2", C: "c-2", D: "d-4"},
		{A: "a-1", B: "b-1", C: "c-1", D: "d-5"},
	})
}
