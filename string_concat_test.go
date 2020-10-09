package gobench 

import (
    "fmt"
    "testing"
)

var result string

func sprintf(a, b string) {
    result = fmt.Sprintf("%s %s", a, b) 
}

func plus(a, b string) {
    result = a + " " + b
}

func BenchmarkSprintf(b *testing.B) {
    for n := 0; n < b.N; n++ {
        sprintf("aaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb")
    }
}

func BenchmarkPlus(b *testing.B) {
    for n := 0; n < b.N; n++ {
        plus("aaaaaaaaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb")
    }
}
