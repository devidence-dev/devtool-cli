# 🛠️ devtool-cli

> A Go-based developer CLI for managing VS Code Server orphan processes and cleaning up Vagrant/libvirt environments — with an interactive terminal UI.

---

## 📋 Requirements

- 🐳 **Docker** — to compile from the host without installing Go
- 🔐 **sudo** — to install into `/usr/local/bin`

---

## 🚀 Installation

From the project directory on the **host** (no need to enter the devcontainer):

```bash
make install
```

This compiles the binary using an ephemeral Go container and installs it at `/usr/local/bin/devtool`.

> 💡 The first run pulls the Go image and downloads dependencies. Subsequent builds are fast thanks to the `devtool-gomodcache` named volume.

### 🗑️ Uninstall

```bash
make uninstall
```

---

## 📖 Usage

### 💻 `devtool vscode kill`

Finds and terminates orphan VS Code Remote SSH processes for the current user.

```bash
devtool vscode kill
```

**Execution flow:**

1. 🔍 Spinner while searching for `.vscode-server` processes
2. 📊 Table listing found processes (PID, %CPU, %MEM, command)
3. ❓ Interactive confirmation before taking action
4. ⚡ Sends `SIGTERM` — escalates to `SIGKILL` if any process resists

---

### 📦 `devtool vagrant cleanup`

Displays the current state of VMs, boxes, and libvirt resources, then offers an interactive cleanup menu.

```bash
devtool vagrant cleanup
```

**Menu options** (navigate with ↑↓ arrows):

| # | Action |
|---|--------|
| 1️⃣ | Destroy all Vagrant VMs |
| 2️⃣ | Remove unused Vagrant boxes (`vagrant box prune`) |
| 3️⃣ | Remove **all** Vagrant boxes |
| 4️⃣ | Clean orphan libvirt volumes |
| 5️⃣ | Clean inactive libvirt domains |
| 6️⃣ | 🧹 Full cleanup with progress bar |
| 7️⃣ | Exit |

> ⚠️ Options 4 and 5 require `virsh` to be installed on the system.

---

## 🧑‍💻 Development

Available targets **inside the devcontainer**:

```bash
make build                   # Compile ./devtool locally
make run ARGS="vscode kill"  # Run without compiling
make lint                    # Run golangci-lint
make clean                   # Remove local binary
```
