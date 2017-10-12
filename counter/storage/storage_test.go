package storage

import (
    "testing"
    "reflect"
    "../storage"
)


func TestWriteRead(t *testing.T) {
	os.Remove("test.log")
    
    timestamps1 := []int64{1,2,3}
    utils.WriteTimestamps("test.log", timestamps1)
    timestamps2 := utils.ReadTimestamps("test.log")
    
    if !reflect.DeepEqual(timestamps1, timestamps2) {
        t.Errorf("ReadTimestamps data differs from WriteTimestamps %v, want %v", timestamps2, timestamps1)
    }
    
    //verify path is not ignored
    
    os.Remove("test2.log")
    
    timestamps11 := []int64{1,2,3}
    utils.WriteTimestamps("test1.log", timestamps11)
    timestamps21 := []int64{4,5,3}
    utils.WriteTimestamps("test2.log", timestamps21)

    timestamps12 := utils.ReadTimestamps("test1.log")
    timestamps22 := utils.ReadTimestamps("test.log")
    
    if !reflect.DeepEqual(timestamps11, timestamps12) {
        t.Errorf("ReadTimestamps data differs from WriteTimestamps %v, want %v", timestamps12, timestamps11)
    }
    
    if !reflect.DeepEqual(timestamps11, timestamps12) {
        t.Errorf("ReadTimestamps data differs from WriteTimestamps %v, want %v", timestamps12, timestamps11)
    }
}
