# Intercampus Federated Learning (ICFL) - Link

A [Flower Framework](https://github.com/adap/flower) based federated learning coordinator.

## Installation

### Prerequisites
- **npm** (version 9 or higher) for the frontend
- **Python** 3.10+
- **Golang** 1.22 (to compile if desired)
- **MySQL**
- A server certificate to distribute to the participating nodes
- Server public and private keys (check `authentication/certificates` and `authentication/keys`)

### 1. Clone the Repository
```sh
git clone https://github.com/user/repo.git
cd repo
```

### 2. Frontend Installation
1. Install the required Node packages:
   ```sh
   npm install
   ```
2. Set the backend API URL in the `.env` file.
3. Run the frontend:
   ```sh
   npm run dev   # Development mode
   npm run build # Build for production
   ```

### 3. Backend Installation
1. Build the Go executable:
   ```sh
   go build -o backend link/cmd/server/main.go
   ```
2. Configure `config.yaml` within the `config` directory.
3. Place certificates and keys in their respective folders:
   - `authentication/certificates`
   - `authentication/keys`
   - Ensure `client_public_keys.csv` is present (it can be empty).
4. Run the binary:
   ```sh
   ./backend
   ```

### Default Credentials
- **User:** `admin`
- **Password:** `defaultpassword`
  
These credentials may be changed directly in the database.

The user may refer to the specifications in the following document: [link here]

All logs are saved to:
- `logs/`
- `uploads/experimentid/logs/`

---
## Experiments Usage

- It is **preferable** for experiments to have the same name as their Flower apps.
- Flower apps are uploaded as a zip with the following structure:
  ```
  experiment.zip
  └── experiment_name/
      ├── pyproject.toml
      └── experiment_name/
          ├── client_app.py
          ├── server_app.py
          ├── model.pt
          └── task.pt
  ```

### Flower App Configuration
The `pyproject.toml` file must contain the same `experiment_name`. The `root-certificates` and `address` must remain as follows:

```
[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"

[project]
name = "experiment_name"
version = "1.0.0"
description = "ICFL pyproject example"
license = "Apache-2.0"
dependencies = [
    "flwr>=1.15.0",
]

[tool.hatch.build.targets.wheel]
packages = ["."]

[tool.flwr.app]
publisher = "icfl"

[tool.flwr.app.components]
serverapp = "experiment_name.server_app:app"
clientapp = "experiment_name.client_app:app"

[tool.flwr.app.config]
num-server-rounds = 1
fraction-evaluate = 1

[tool.flwr.federations]
default = "experiment_name"

[tool.flwr.federations.experiment_name]
address = "127.0.0.1:9093" # Address of the Exec API
root-certificates = "../../../authentication/certificates/ca.crt"
```

The `client_app.py` is able to access the selected datasets on their respective nodes by including the following line in the `client_fn`:

```
dataset_path = context.node_config["dataset-path"]
```

For a sample flower deployment, look at 
[Flower-Authentication Example](https://github.com/adap/flower/tree/40f8a4a981967e67e3e3bac5f1ef7958a854ef45/examples/flower-authentication)

---
## Features
- Client node metadata registries
- Secure deployment of experiments
- Experiment management
- Support for federated YOLOv8 fine-tuning

---
## Known Bugs
- Node errors during training may cause discrepancies in experiment state.

