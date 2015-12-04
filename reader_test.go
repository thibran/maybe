package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	var b bytes.Buffer
	b.WriteString("Hello")
	r := bufio.NewReader(&b)
	s, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	fmt.Println(s)
}

func TestByte(t *testing.T) {
	//[]string{"Hello", "World"}

	// row := []string{"Hello", "World"}
	// b := []byte(row)
	// r := []string(b)
	// fmt.Println(r)
}

func Read(r io.Reader, sep, lineEnd rune) {

}

// func FromReader(r io.Reader) {
// 	rd := bufio.NewReader(r)
// 	for {
// 		record, err := rd.ReadLine() //rd.ReadBytes('\n') //r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		arr := []string(record)
// 		fmt.Println(arr[0])
// 		fmt.Println(record)
// 	}
// }

func TestSCV(t *testing.T) {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"`
	r := csv.NewReader(strings.NewReader(in))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		arr := []string(record)
		fmt.Println(arr[0])
		fmt.Println(record)
	}
}

// func TestJSON(t *testing.T) {
// 	// type Message struct {
// 	// 	Name string
// 	// 	Body string
// 	// 	Arr  []string
// 	// }
// 	// m := &Message{"Rob", "loves go", []string{"a", "b"}}
// 	json := `["a","b"]`
//
// 	var b bytes.Buffer
// 	b.WriteString(json)
// 	enc := json.NewEncoder(&b)
//
// 	if err := enc.Encode(&v); err != nil {
// 		log.Println(err)
// 	}
// 	// b, err := json.Marshal([]string{"a", "b"})
// 	// if err != nil {
// 	// 	t.Fail()
// 	// }
// 	// fmt.Println(string(b))
// }
