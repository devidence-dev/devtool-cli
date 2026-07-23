# Graph Report - devtool-cli  (2026-07-15)

## Corpus Check
- 11 files · ~2,025 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 55 nodes · 87 edges · 11 communities (7 shown, 4 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 9 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `809ff5b9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Vagrant Cleanup Command
- README Build & Install Docs
- Graphify Workflow
- Vagrant Cleanup Menu
- VSCode Kill Command
- UI Printer Utilities
- CLI Root Entrypoint
- Go Module Path
- runCmd
- make install
- make lint

## God Nodes (most connected - your core abstractions)
1. `runVagrantCleanup()` - 14 edges
2. `runCmd()` - 10 edges
3. `devtool-cli (Go developer CLI)` - 9 edges
4. `devtool vagrant cleanup` - 8 edges
5. `runVSCodeKill()` - 7 edges
6. `make install` - 5 edges
7. `showCurrentState()` - 4 edges
8. `destroyAllVMs()` - 4 edges
9. `removeAllBoxes()` - 4 edges
10. `cleanLibvirtVolumes()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `graphify` --conceptually_related_to--> `devtool-cli (Go developer CLI)`  [INFERRED]
  CLAUDE.md → README.md
- `runVagrantCleanup()` --calls--> `Warning()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `runVagrantCleanup()` --calls--> `Info()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `runVagrantCleanup()` --calls--> `Success()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `showCurrentState()` --calls--> `Info()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **devtool vagrant cleanup interactive menu options** — readme_devtool_vagrant_cleanup, readme_menu_destroy_vms, readme_menu_prune_boxes, readme_menu_remove_all_boxes, readme_menu_clean_orphan_volumes, readme_menu_clean_inactive_domains, readme_menu_full_cleanup [EXTRACTED 1.00]
- **Devcontainer development Make targets** — readme_make_build, readme_make_run, readme_make_lint, readme_make_clean [EXTRACTED 1.00]

## Communities (11 total, 4 thin omitted)

### Community 0 - "Vagrant Cleanup Command"
Cohesion: 0.48
Nodes (5): cleanLibvirtNetworks(), commandExists(), runVagrantCleanup(), showCurrentState(), withSpinner()

### Community 1 - "README Build & Install Docs"
Cohesion: 0.50
Nodes (5): devtool-cli (Go developer CLI), devtool vscode kill, make build, make clean, make run ARGS="vscode kill"

### Community 3 - "Vagrant Cleanup Menu"
Cohesion: 0.32
Nodes (8): devtool vagrant cleanup, Menu option: Clean inactive libvirt domains, Menu option: Clean orphan libvirt volumes, Menu option: Destroy all Vagrant VMs, Menu option: Full cleanup with progress bar, Menu option: Remove unused Vagrant boxes (vagrant box prune), Menu option: Remove all Vagrant boxes, virsh (libvirt CLI dependency)

### Community 4 - "VSCode Kill Command"
Cohesion: 0.52
Nodes (5): findVSCodePIDs(), Command, printProcessTable(), processExists(), runVSCodeKill()

### Community 5 - "UI Printer Utilities"
Cohesion: 0.33
Nodes (3): Info(), Success(), Warning()

### Community 8 - "runCmd"
Cohesion: 0.52
Nodes (7): cleanLibvirtDomains(), cleanLibvirtDomainsForce(), cleanLibvirtVolumes(), destroyAllVMs(), Command, removeAllBoxes(), runCmd()

### Community 9 - "make install"
Cohesion: 0.40
Nodes (5): devtool-gomodcache named volume, Docker (build dependency), make install, make uninstall, sudo (install dependency)

## Knowledge Gaps
- **11 isolated node(s):** `github.com/pxblo-x1/devtool-cli`, `make build`, `make clean`, `Docker (build dependency)`, `sudo (install dependency)` (+6 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `devtool-cli (Go developer CLI)` connect `README Build & Install Docs` to `make lint`, `make install`, `Graphify Workflow`, `Vagrant Cleanup Menu`?**
  _High betweenness centrality (0.118) - this node is a cross-community bridge._
- **Why does `runVagrantCleanup()` connect `Vagrant Cleanup Command` to `runCmd`, `UI Printer Utilities`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Why does `runVSCodeKill()` connect `VSCode Kill Command` to `UI Printer Utilities`?**
  _High betweenness centrality (0.086) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `runVagrantCleanup()` (e.g. with `Info()` and `Success()`) actually correct?**
  _`runVagrantCleanup()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/pxblo-x1/devtool-cli`, `make build`, `make clean` to the rest of the system?**
  _11 weakly-connected nodes found - possible documentation gaps or missing edges._