# SchedulAI (CorePilot)

**An ML-Driven Adaptive CPU Scheduler for Linux**

---

## Project Overview

SchedulAI (CorePilot) is an experimental machine learning-based CPU scheduler designed to learn and optimize CPU scheduling patterns based on your unique usage style (e.g., gaming, development, streaming). The project aims to prototype a userland scheduler and, ultimately, integrate it as a Linux kernel module.

---

## Features

- **Data Collection:** Logs process lifetimes, CPU usage, wait times, and I/O statistics at regular intervals.
- **Data Cleaning:** Normalizes collected data for time of day, background noise, and system state.
- **ML Modeling:** Uses reinforcement learning or multi-armed bandit algorithms to optimize scheduling decisions.
- **Evaluation:** Compares the ML scheduler’s performance against standard Linux schedulers (CFS, BFS, etc.) with clear statistics and visualizations.
- **MLOps:** Supports weekly retraining and monitoring to adapt to changing usage patterns.

---

## Project Structure

```
CorePilot/
├── data/                # Collected logs and datasets
├── scripts/             # Data collection scripts
├── src/                 # ML models and scheduler code
├── notebooks/           # Jupyter notebooks for analysis/visualization
├── evaluation/          # Benchmarking and comparison scripts
├── docs/                # Documentation and project plan
├── tests/               # Unit and integration tests
├── requirements.txt     # Python dependencies
├── README.md            # Project overview and instructions
└── .gitignore           # Ignore data, logs, etc.
```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/CorePilot.git
cd CorePilot
```

### 2. Install Dependencies

```bash
pip install -r requirements.txt
```

### 3. Collect Data

Run the data collection script to start logging system metrics:

```bash
python scripts/collect_data.py
```

Sample logs will be saved in the `data/` directory.

### 4. Train and Evaluate

- Use scripts in `src/` to train the ML scheduler.
- Use scripts in `evaluation/` to benchmark against standard Linux schedulers.

---

## Checkpoints & Progress

See `docs/project_plan.md` (or `convo.txt`) for a detailed project diary, checkpoints, and progress log.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for suggestions and improvements.

---

## License

[MIT License](LICENSE)

---

## Acknowledgments

- Inspired by Linux kernel scheduling research and modern ML techniques.
- Built with [psutil](https://github.com/giampaolo/psutil), [scikit-learn](https://scikit-learn.org/), and