package storage

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "strconv"
    "time"
    "simplesurance-group.de/counter/utils"
    "log"
)

type TimestampStorage struct {
	filePath string
	values []int64
	c_in chan bool
	c_out chan int
	c_log_timestamp chan int64
}

const (
    timeSpan = 60000
    logLimit = 10000
)

func NewTimestampStorage(filePath string) *TimestampStorage {
	timestampNow := utils.TimeNow()
	values := make([]int64, 0)
		
	c_in := make(chan bool, 100)
    c_out := make(chan int)
    c_log_timestamp := make(chan int64)

	store := TimestampStorage{filePath, values, c_in, c_out, c_log_timestamp}
	
	values = store.readTimestamps()
	values = filterTimestamps(timestampNow, values)
	store.values = values

    store.listenForRequests()
    
    store.listenAppendLogFile()
    
    store.autoRemoveOldLogFiles()
	
	return &store
}


func (store TimestampStorage) listenAppendLogFile(){
	lastLogTimestamps := readSingleFileTimestamps(store.filePath)
	lastLogCounter := len(lastLogTimestamps)
	
	lastLogFile, err := os.OpenFile(store.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)	
	
	if err != nil {
        log.Fatal(err)
    }
    
    go func() {
		for {
			timestamp := <- store.c_log_timestamp
			timestampStr := fmt.Sprintf("%v\n", timestamp)
			lastLogFile.WriteString(timestampStr)
			
			lastLogCounter++
			
			if lastLogCounter >= logLimit {
				lastLogCounter = 0
				lastLogFile.Close()
				
				// rename all filepaths  .log.0 > .log.1; .log > .log.0
				store.mapLogFiles(func(oldPath string, i int) {
					newPath := fmt.Sprintf("%s.%v", store.filePath, i+1)
					os.Rename(oldPath, newPath)
				})
				
				lastLogFile, err = os.OpenFile(store.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
}

func (store TimestampStorage) autoRemoveOldLogFiles(){
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			<- ticker.C
			store.mapLogFiles(func(logFile string, i int) {
				if i>= 0 {
					timestamps := readSingleFileTimestamps(logFile)
					if len(timestamps)==0 || (utils.TimeNow() - timeSpan) > timestamps[len(timestamps)-1] {
						os.Remove(logFile)
					}
				}
			})
		}
	}()
}
	
func (store *TimestampStorage) listenForRequests(){
	go func() {
        for {
            <- store.c_in
            
            timestamps := store.values
			timestampNow := utils.TimeNow()
			timestamps = append(timestamps, timestampNow)
			timestamps = filterTimestamps(timestampNow, timestamps)
			
			store.values = timestamps

			store.c_out <- len(timestamps)
			
			store.c_log_timestamp <- timestampNow
        }
    }()
}
	

func (store *TimestampStorage) CounterAddTimestampNow() chan int{
	store.c_in <- true
	return store.c_out
}

func filterTimestamps(timestampNow int64, timestamps []int64) []int64{
    minTimestamp := timestampNow - timeSpan
    
    for i, timestamp := range timestamps{
        if timestamp > minTimestamp{
			return timestamps[i:]
        }
    }
    
    return make([]int64, 0)
}

func (store TimestampStorage) readTimestamps() []int64 {
	res := make([]int64, 0)
    store.mapLogFiles(func(logFile string, i int) {
		res = append(res, readSingleFileTimestamps(logFile)...)
	})
	return res
}

func (store TimestampStorage) mapLogFiles(callback func (string, int)) {
	//defer reverses, therefore this will be last call
	defer callback(store.filePath, -1)
	for i:=0; i< 100; i++ {
		oldFilePath := fmt.Sprintf("%s.%v", store.filePath, i)
		if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		    break
		}
		defer callback(oldFilePath, i)
	}
}

func readSingleFileTimestamps(logFile string) []int64 {
	content, _ := ioutil.ReadFile(logFile)

    timestampStrings := strings.Split(string(content[:]), "\n")
    timestampStrings = timestampStrings[:len(timestampStrings)-1]
    
    timestamps := make([]int64, 0)
    
    for _, timestampString := range timestampStrings{
        timestamp, _ := strconv.ParseInt(timestampString, 10, 64)
        timestamps = append(timestamps, timestamp)
    }
    return timestamps
}
