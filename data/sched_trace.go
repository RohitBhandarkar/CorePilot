// File: main.go
// Description: Capture Linux sched_switch events via bpftrace and write directly to Parquet using parquet-go.

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// SchedSwitchEvent models one context-switch event for Parquet storage.
type SchedSwitchEvent struct {
    TimestampNs     int64  `parquet:"name=ts_ns, type=INT64"`
    CPU             int32  `parquet:"name=cpu, type=INT32"`
    PrevPID         int32  `parquet:"name=prev_pid, type=INT32"`
    PrevComm        string `parquet:"name=prev_comm, type=BYTE_ARRAY, logicaltype=STRING, encoding=PLAIN_DICTIONARY"`
    PrevPrio        int32  `parquet:"name=prev_prio, type=INT32"`
    NextPID         int32  `parquet:"name=next_pid, type=INT32"`
    NextComm        string `parquet:"name=next_comm, type=BYTE_ARRAY, logicaltype=STRING, encoding=PLAIN_DICTIONARY"`
    NextPrio        int32  `parquet:"name=next_prio, type=INT32"`
    WallUnixMicros  int64  `parquet:"name=wall_us, type=INT64, convertedtype=TIMESTAMP_MICROS"`
}

func main() {
	// 1. Prepare output Parquet file
	outFile := fmt.Sprintf("sched_trace_%s.parquet", time.Now().Format("2006-01-02_15-04-05"))
	fw, err := local.NewLocalFileWriter(outFile)
	check(err, "opening Parquet file")
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(SchedSwitchEvent), 4)
	check(err, "creating Parquet writer")
	pw.RowGroupSize = 128 * 1024 * 1024 // 128 MiB
	pw.CompressionType = parquet.CompressionCodec_SNAPPY
	defer func() {
		if err := pw.WriteStop(); err != nil {
			log.Fatalf("Parquet write stop error: %v", err)
		}
	}()

	// 2. Build bpftrace probe script
	probe := `
tracepoint:sched:sched_switch
{
printf("%llu,%d,%d,%s,%d,%d,%s,%d\n",
nsecs,
cpu,
args->prev_pid,
args->prev_comm,
args->prev_prio,
args->next_pid,
args->next_comm,
args->next_prio
);
}`

	// 3. Launch bpftrace
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "bpftrace", "-e", probe)
	stdout, err := cmd.StdoutPipe()
	check(err, "getting bpftrace stdout")
	cmd.Stderr = os.Stderr

	check(cmd.Start(), "starting bpftrace")
	log.Printf("bpftrace started (PID %d), writing to %s", cmd.Process.Pid, outFile)

	// 4. Scan and parse each line, write to Parquet
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		event, err := parseLine(line)
		if err != nil {
			log.Printf("warning: skipping line (%v): %s", err, line)
			continue
		}
		if err := pw.Write(event); err != nil {
			log.Fatalf("Parquet write error: %v", err)
		}
	}

	check(scanner.Err(), "reading bpftrace output")
	check(cmd.Wait(), "waiting for bpftrace to finish")

	log.Printf("Finished. Parquet file: %s", outFile)
}

// parseLine parses a CSV line from bpftrace into a SchedSwitchEvent.
func parseLine(s string) (*SchedSwitchEvent, error) {
	// Expect 8 comma-separated fields
	parts := strings.SplitN(s, ",", 8)
	if len(parts) != 8 {
		return nil, fmt.Errorf("expected 8 fields, got %d", len(parts))
	}

	// Helper to parse int64
	parseInt := func(i int) (int64, error) {
		return strconv.ParseInt(parts[i], 10, 64)
	}

	tsNs, err := parseInt(0)
	if err != nil {
		return nil, fmt.Errorf("invalid ts_ns: %v", err)
	}
	cpuVal, err := parseInt(1)
	if err != nil {
		return nil, fmt.Errorf("invalid cpu: %v", err)
	}
	prevPID, err := parseInt(2)
	if err != nil {
		return nil, fmt.Errorf("invalid prev_pid: %v", err)
	}
	prevComm := parts[3]
	prevPrio64, err := parseInt(4)
	if err != nil {
		return nil, fmt.Errorf("invalid prev_prio: %v", err)
	}
	nextPID64, err := parseInt(5)
	if err != nil {
		return nil, fmt.Errorf("invalid next_pid: %v", err)
	}
	nextComm := parts[6]
	nextPrio64, err := parseInt(7)
	if err != nil {
		return nil, fmt.Errorf("invalid next_prio: %v", err)
	}

	return &SchedSwitchEvent{
		TimestampNs:    tsNs,
		CPU:            int32(cpuVal),
		PrevPID:        int32(prevPID),
		PrevComm:       prevComm,
		PrevPrio:       int32(prevPrio64),
		NextPID:        int32(nextPID64),
		NextComm:       nextComm,
		NextPrio:       int32(nextPrio64),
		WallUnixMicros: time.Now().UnixMicro(),
	}, nil
}

// check logs and exits on error.
func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
