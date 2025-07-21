package main

import (
	"log"
	"time"
	"github.com/shirou/gopsutil/v4/process"
	"encoding/json"
	"os"
)

type ProcessInfo struct {
	Pid        int32
	Pname      string
	Timestamp  time.Time
	CpuTimes   map[string]float64
	Lifetime   int64
	CpuUsage   float64
	IoCounters map[string]uint64
	Status     string
}

func main() {
	log.Println("Starting process info logger...")

	processes, err := process.Processes()
	if err != nil {
		log.Fatalf("Error getting processes: %v", err)
	}
	log.Printf("Found %d processes", len(processes))

	data := make([]ProcessInfo, 0, len(processes))

	for _, proc := range processes {
		now := time.Now()
		errors := false
		pid := proc.Pid
		name, err := proc.Name()
		if err != nil && !errors {
			log.Printf("PID %d: Error getting process name: %v", pid, err)
			name = "unknown"
			errors = true
		}

		createTime, err := proc.CreateTime()
		var lifetimeMs int64
		if err != nil {
			log.Printf("PID %d: Error getting create time: %v", pid, err)
			errors = true
		} else {
			lifetimeMs = now.UnixMilli() - createTime
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil{
			log.Printf("PID %d: Error getting CPU percent: %v", pid, err)
			errors = true
		}

		cpuTimesObj, err := proc.Times()
		if err != nil {
			log.Printf("PID %d: Error getting CPU times: %v", pid, err)
			errors = true
		}

		ioCountersObj, err := proc.IOCounters()
		if err != nil {
			log.Printf("PID %d: Error getting IO counters: %v", pid, err)
			errors = true
		}

		statusStr, err := proc.Status()
		if err != nil {
			log.Printf("PID %d: Error getting status: %v", pid, err)
			errors = true
		}

		log.Printf(
			"PID: %d | Name: %s | Timestamp: %d | Lifetime(ms): %d | CPU Usage(%%): %.2f | CPU Times: %+v | IO Counters: %+v | Status: %s",
			pid, name, now.Unix(), lifetimeMs, cpuPercent, cpuTimesObj, ioCountersObj, statusStr,
		)

		if !errors {
			processInfo := ProcessInfo{
				Pid:        pid,
				Pname:      name,
				Timestamp:  now,
				CpuTimes:   map[string]float64{"user": cpuTimesObj.User, "system": cpuTimesObj.System},
				Lifetime:   lifetimeMs,
				CpuUsage:   cpuPercent,
				IoCounters: map[string]uint64{"read_count": ioCountersObj.ReadCount, "write_count": ioCountersObj.WriteCount},
				Status:     statusStr[0],
			}
			data = append(data, processInfo)
			log.Printf("Process Info added")
		}
	}
	outputFile, err := os.Create("../artifacts/process_data.json")
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ") // For readable indentation
	err = encoder.Encode(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	log.Println("Process data successfully written to process_data.json")
}