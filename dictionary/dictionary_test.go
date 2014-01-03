package Dictionary

import "testing"
import "fmt"

var dict = NewDictionary()

func init(){
    dict.Set("HogeHoge","value")
    dict.Set("ahahaha","ihihi")
}
func TestA(t*testing.T){
    t.Log("Equal Test\n")
    if dict.Get("hogehoge") != "value" {
        t.Error("get(\"hogehoge\") != \"value\"\n")
    }
}

func TestB(t*testing.T){
    t.Log("Not Equal Test\n")
    if dict.Get("uhauha") != "" {
        t.Error("get(\"uhauha\") != \"\"")
    }
}

func TestC(t*testing.T){
    for pair := range(dict.Iter()){
        fmt.Printf("%s=%s\n",pair.Key,pair.Value)
    }
}

func TestD(t*testing.T){
    for pair := range(dict.SortIter()) {
        fmt.Printf("(sort)%s=%s\n",pair.Key,pair.Value)
    }
}
