package main

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
)

func main() {
	// init database ด้วย bolt engine
	graph.InitQuadStore("bolt", "./db/bolt/cayleya.db", nil)

	// สร้าง connection ไปที่ database
	store, err := cayley.NewGraph("bolt", "./db/bolt/cayleya.db", nil)
	if err != nil {
		log.Fatalln(err)
	}

	// สร้าง quad ใหม่
	store.AddQuad(quad.Make(quad.IRI("alice"), quad.IRI("follows"), quad.IRI("bob"), nil))
	store.AddQuad(quad.Make(quad.IRI("bob"), quad.IRI("follows"), quad.IRI("fred"), nil))
	store.AddQuad(quad.Make(quad.IRI("bob"), quad.IRI("status"), quad.IRI("cool_person"), nil))
	store.AddQuad(quad.Make(quad.IRI("charlie"), quad.IRI("follows"), quad.IRI("bob"), nil))
	store.AddQuad(quad.Make(quad.IRI("charlie"), quad.IRI("follows"), quad.IRI("dani"), nil))
	store.AddQuad(quad.Make(quad.IRI("dani"), quad.IRI("follows"), quad.IRI("bob"), nil))
	store.AddQuad(quad.Make(quad.IRI("dani"), quad.IRI("follows"), quad.IRI("greg"), nil))
	store.AddQuad(quad.Make(quad.IRI("dani"), quad.IRI("status"), quad.IRI("cool_person"), nil))
	store.AddQuad(quad.Make(quad.IRI("emily"), quad.IRI("follows"), quad.IRI("fred"), nil))
	store.AddQuad(quad.Make(quad.IRI("fred"), quad.IRI("follows"), quad.IRI("greg"), nil))
	store.AddQuad(quad.Make(quad.IRI("greg"), quad.IRI("status"), quad.IRI("cool_person"), nil))
	store.AddQuad(quad.Make(quad.IRI("predicates"), quad.IRI("are"), quad.IRI("follows"), nil))
	store.AddQuad(quad.Make(quad.IRI("predicates"), quad.IRI("are"), quad.IRI("status"), nil))
	store.AddQuad(quad.Make(quad.IRI("emily"), quad.IRI("status"), quad.IRI("smart_person"), quad.IRI("smart_graph")))
	store.AddQuad(quad.Make(quad.IRI("greg"), quad.IRI("status"), quad.IRI("smart_person"), quad.IRI("smart_graph>")))

	// สร้าง path เพื่อใช้ดึงข้อมูล
	p := cayley.StartPath(store, quad.IRI("charlie")).Out(quad.IRI("follows"))

	// สร้าง optimized iterator
	it, _ := p.BuildIterator().Optimize()

	//เอา optimized iterator ไปชี้ที่ quad ใน graph
	it, _ = store.OptimizeIterator(it)

	// clear iterator
	defer it.Close()

	// ลูปดึงค่าออกมา
	for it.Next() {
		token := it.Result()                // ดึง token ออกมา (token เป็น reference)
		value := store.NameOf(token)        // ดึง value ที่ผูกกับ token นั้นอยู่
		nativeValue := quad.NativeOf(value) // แปลงเป็น go

		fmt.Println(nativeValue) // แสดงค่าออกมา
	}
	if err := it.Err(); err != nil {
		log.Fatalln(err)
	}
}
