package storage

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "strconv"
    "bufio"
    "time"
    "simplesurance-group.de/counter/utils"
)

type TimestampStorage struct {
	filePath string
	values []int64
	c_in chan bool
	c_out chan int
	c_update chan bool
}

const (
    timeSpan = 60000
)

func NewTimestampStorage(filePath string) *TimestampStorage {
	values := readTimestamps(filePath)
	
	c_in := make(chan bool, 100)
    c_out := make(chan int)
    c_update := make(chan bool)

	storage := TimestampStorage{filePath, values, c_in, c_out, c_update}

    go func() {
        for {
            <- c_in
            storage.reportTimespanCounter()
        }
    }()
    
    go func() {
		updated := true
		for {
			select {
				case <- c_update:
					updated = false
				case <- time.After(time.Second):
					if !updated {
						writeTimestamps(storage.filePath, storage.values)
						updated = true
					}
			}
		}
	}()
	
	return &storage
}

func (storage *TimestampStorage) CounterAddTimestampNow() chan int{
	storage.c_in <- true
	return storage.c_out
}

func (storage *TimestampStorage) reportTimespanCounter() {
	timestamps := storage.values
	timestampNow := utils.TimeNow()
    timestamps = append(timestamps, timestampNow)
    timestamps = filterTimestamps(timestampNow, timestamps)
    
    storage.values = timestamps

    storage.c_out <- len(timestamps)
    
    storage.c_update <- true
}

func filterTimestamps(timestampNow int64, timestamps []int64) []int64{
    resTimestamps := make([]int64, 0)
    minTimestamp := timestampNow - timeSpan
    
    for _, timestamp := range timestamps{
        if timestamp > minTimestamp{
            resTimestamps = append(resTimestamps, timestamp)
        }
    }
    
    return resTimestamps
}

func readTimestamps(filePath string) []int64 {
    content, _ := ioutil.ReadFile(filePath)

    timestampStrings := strings.Split(string(content[:]), "\n")
    timestampStrings = timestampStrings[:len(timestampStrings)-1]
    
    timestamps := make([]int64, 0)
    
    for _, timestampString := range timestampStrings{
        timestamp, _ := strconv.ParseInt(timestampString, 10, 64)
        timestamps = append(timestamps, timestamp)
    }
    return timestamps
}

func writeTimestamps(filePath string, timestamps []int64) error {	
    tempFilePath := filePath + ".temp"
    f, err := os.Create(tempFilePath)
    if err != nil {
        return err
    }
    
    defer os.Rename(tempFilePath, filePath)

    defer f.Close()

    w := bufio.NewWriter(f)
    defer w.Flush()
    for _, timestamp := range timestamps {
        timestampStr := fmt.Sprintf("%v\n", timestamp)
        _, err := w.WriteString(timestampStr)
        if err != nil {
            return err
        }
    }
    return nil
}
