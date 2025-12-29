package main

type Foo struct {
	Bar []int
	Baz map[string]bool
}

func main() {
	foo := Foo{
		Bar: []int{69, 420, 1337},
		Baz: map[string]bool{
			"urmom": false,
		},
	}

	if err := GobDump("gob.bin", foo); err != nil {
		panic(err)
	}
	if err := GobDumpAtomic("gobatomic.bin", foo); err != nil {
		panic(err)
	}
}
