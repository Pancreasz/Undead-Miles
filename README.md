# 🚗 Undead Miles

**Undead Miles** is a cloud-native ride-sharing and carpooling platform designed to match drivers with passengers efficiently. It is built using a microservices architecture, running on Kubernetes, with a full CI/CD pipeline and observability stack.

![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Kubernetes](https://img.shields.io/badge/kubernetes-v1.24+-blue.svg)

## ✨ Features

* **Driver Mode:** Drivers can publish trips with origin, destination, price, and departure time.
* **Passenger Mode (Watcher):** Passengers can create "watchers" for specific routes.
* **Real-time Matching:** The system continuously scans for matches between trips and watchers.
* **Notifications:** Users are alerted when a ride match is found.
* **Observability:** Full monitoring of cluster health and application metrics via Grafana & Prometheus.

## 🛠 Tech Stack

### Frontend
* **Framework:** React (TypeScript) + Vite
* **Styling:** Tailwind CSS
* **Icons:** Lucide React

### Backend (Microservices)
* **Language:** Go (Golang) v1.24
* **Services:**
  * `marketplace-service`: Manages trip creation and listings.
  * `watcher-service`: Handles passenger requests and route matching.
  * `notification-service`: Manages user alerts.
* **Database:** PostgreSQL (with pgvector support).

### Infrastructure & DevOps
* **Containerization:** Docker
* **Orchestration:** Kubernetes (K8s)
* **CI/CD:** Jenkins (running in-cluster) + Kaniko (for building images without Docker daemon).
* **Ingress:** Nginx Ingress Controller.
* **Monitoring:** Prometheus (Metrics), Grafana (Visualization), Node Exporter.
* **Package Management:** Helm.

---

## 📂 Project Structure

```bash
UndeadMiles/
├── k8s/                     # Kubernetes manifests (Deployments, Services, Ingress)
├── marketplace-service/     # Go backend for Drivers/Trips
├── watcher-service/         # Go backend for Passengers/Matching
├── notification-service/    # Go backend for Alerts
├── undeadmiles-frontend/    # React Frontend application
├── jenkins_data/            # Local Jenkins persistence
└── README.md
```
