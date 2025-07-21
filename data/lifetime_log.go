package main

import (
	"fmt"
	"time"
	"github.com/shirou/gopsutil/v4/process"
)

func main() {
	fmt.Println("hello")
	processIDs, err := process.Processes()
	if err != nil {
		fmt.Println("Error getting processes:", err)
		return
	}
	fmt.Println("Processes:", len(processIDs))

	for pid := range processIDs { //proc is the process ID
		proc, _ := process.NewProcess(int32(pid))
		if err != nil {
			fmt.Println("Error creating process:", err)
			continue
		}
		now := time.Now()
		println("Process ID:", pid, "Timestamp:", now.Unix())
		
		name, err := proc.Name()
    	if err != nil {
        	fmt.Println("Error getting process name:", err)
        	name = "unknown"
    	} else{
    		fmt.Printf("Process ID: %d, Name: %s, Timestamp: %d\n", pid, name, now.Unix())
		}

		createTime, err := proc.CreateTime()
		if err != nil {
			fmt.Println("Error getting create time:", err)
		} else {
			lifetime := now.UnixMilli() - createTime
			fmt.Printf("  Lifetime (ms): %d\n", lifetime)
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			fmt.Println("Error getting CPU percent:", err)
		} else {
			fmt.Printf("  CPU Usage (%%): %.2f\n", cpuPercent)
		}

		cpuTimes, err := proc.Times()
		if err != nil {
			fmt.Println("Error getting CPU times:", err)
		} else {
			fmt.Printf("  CPU Times: %+v\n", cpuTimes)
		}

		ioCounters, err := proc.IOCounters()
		if err != nil {
			fmt.Println("Error getting IO counters:", err)
		} else {
			fmt.Printf("  IO Counters: %+v\n", ioCounters)
		}

		status, err := proc.Status()
		if err != nil {
			fmt.Println("Error getting status:", err)
		} else {
			fmt.Printf("  Status: %s\n", status)
		}
		
	}

}