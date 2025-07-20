# SchedulAI: ML-Driven CPU Scheduler Project Plan

---

## Project Summary

My goal is to build an ML module that learns optimal CPU scheduling patterns tailored to my unique machine usage style (such as gaming, development, or streaming). The project will start as a userland prototype and may eventually become a Linux kernel module.

---

## Key Steps & Checkpoints

1. **Data Collection:**  
   - I will write a script (Python or Bash) to log process lifetimes, CPU usage, wait times, and I/O stats from my system.
   - The script will log data at regular intervals and save it in a structured format (CSV, JSON, or SQLite).
   - **Checkpoint:** Paste the data logging script and sample output here for review.

2. **Data Cleaning:**  
   - I will normalize the collected data for time of day, background noise, and system state.

3. **Modeling:**  
   - I will use reinforcement learning or multi-armed bandit algorithms to optimize scheduling decisions.

4. **Deployment:**  
   - I will prototype the scheduler in userland, with the goal of integrating it as a Linux kernel scheduler module.

5. **MLOps:**  
   - I will monitor logs and retrain the model weekly to adapt to changing usage patterns.

6. **Evaluation:**  
   - I will compare the performance of my ML-driven scheduler against existing Linux schedulers (such as CFS and BFS) using key metrics.
   - I will visualize the results with charts and tables.
   - **Checkpoint:** Evaluation against current models is required before project completion.

---

## Progress Diary

**2025-07-19**  
- I initialized the project and defined the project structure and workflow.
- I scaffolded the README.md with an overview, structure, and instructions.
- I set up this project plan as a persistent diary to track progress and checkpoints.
- I decided that evaluation against current Linux schedulers (CFS, BFS, etc.) is a required checkpoint.

---

## Evaluation Plan

- I will collect and analyze performance metrics (CPU utilization, process wait time, throughput, etc.) for both my ML scheduler and standard Linux schedulers.
- I will present comparative statistics and visualizations (charts, tables) to clearly show improvements or trade-offs.
- I will summarize my findings in the project README and documentation.

---
