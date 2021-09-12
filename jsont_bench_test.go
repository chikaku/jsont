package jsont

import (
	"encoding/json"
	"testing"
)

var demo = []byte(`
{
	"code": 200,
	"info": "ok",
	"data": {
		"token": "mDxJkBcVp9NqyVoM",
		"user": {
			"items": [{"a": 1}, {"b": 2}, {"c": 3}]
		}
	}
}
`)

func BenchmarkDecode(b *testing.B) {
	b.Run("std", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var res map[string]interface{}
			if err := json.Unmarshal(demo, &res); err != nil {
				panic(err)
			}
		}
	})

	b.Run("jsont", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := Decode(demo); err != nil {
				panic(err)
			}
		}
	})
}
