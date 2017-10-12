package storage

import (
    "testing"
    "../storage"
    "time"
)


func TestWriteRead(t *testing.T) {
    store1 := storage.NewTimestampStorage("test1.log")
    
    //verify counter increases
    for i:=1; i<10; i++ {
        counter := <- store1.CounterAddTimestampNow()
        if counter != i {
            t.Errorf("Counter is wrong %v, want %v", counter, i)
        }
    }
    
    //verify path is not ignored
    store2 := storage.NewTimestampStorage("test2.log")
    
    for i:=1; i<5; i++ {
        counter := <- store2.CounterAddTimestampNow()
        if counter != i {
            t.Errorf("Counter is wrong %v, want %v", counter, i)
        }
    }
    
    counter := <- store1.CounterAddTimestampNow()
    
    if counter != 10 {
		t.Errorf("Counter is wrong %v, want %v", counter, 10)
	}
	
	counter = <- store2.CounterAddTimestampNow()
    
    if counter != 5 {
		t.Errorf("Counter is wrong %v, want %v", counter, 5)
	}
    
    time.Sleep(5 * time.Second)
    
    counter = <- store1.CounterAddTimestampNow()
    
    time.Sleep(56 * time.Second)
    
	counter = <- store1.CounterAddTimestampNow()
    
    if counter != 2 {
		t.Errorf("Counter after 60s is wrong %v, want %v", counter, 2)
	}
}
