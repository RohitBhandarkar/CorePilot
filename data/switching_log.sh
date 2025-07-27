#!/bin/bash
btime=$(who -b)
echo "$btime" > ../artifacts/boot_time.txt

sudo bpftrace -e '
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
}' | tee ../artifacts/sched_trace.csv
 
