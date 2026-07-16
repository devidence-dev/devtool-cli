# Graph Report - .  (2026-07-15)

## Corpus Check
- Corpus is ~1,977 words - fits in a single context window. You may not need a graph.

## Summary
- 61 nodes · 96 edges · 8 communities (6 shown, 2 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 9 edges (avg confidence: 0.81)
- Token cost: 41,001 input · 0 output

## Community Hubs (Navigation)
- Vagrant Cleanup Command
- README Build & Install Docs
- Graphify Workflow
- Vagrant Cleanup Menu
- VSCode Kill Command
- UI Printer Utilities
- CLI Root Entrypoint
- Go Module Path

## God Nodes (most connected - your core abstractions)
1. `runVagrantCleanup()` - 14 edges
2. `runCmd()` - 10 edges
3. `devtool-cli (Go developer CLI)` - 9 edges
4. `devtool vagrant cleanup` - 8 edges
5. `runVSCodeKill()` - 7 edges
6. `graphify (knowledge graph tool)` - 7 edges
7. `make install` - 5 edges
8. `showCurrentState()` - 4 edges
9. `destroyAllVMs()` - 4 edges
10. `removeAllBoxes()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `runVagrantCleanup()` --calls--> `Warning()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `runVagrantCleanup()` --calls--> `Info()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `runVagrantCleanup()` --calls--> `Success()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `showCurrentState()` --calls--> `Info()`  [INFERRED]
  cmd/vagrant.go → internal/ui/printer.go
- `runVSCodeKill()` --calls--> `Info()`  [INFERRED]
  cmd/vscode.go → internal/ui/printer.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **graphify CLI workflow commands** — claude_graphify_query, claude_graphify_path, claude_graphify_explain, claude_graphify_update [EXTRACTED 1.00]
- **devtool vagrant cleanup interactive menu options** — readme_devtool_vagrant_cleanup, readme_menu_destroy_vms, readme_menu_prune_boxes, readme_menu_remove_all_boxes, readme_menu_clean_orphan_volumes, readme_menu_clean_inactive_domains, readme_menu_full_cleanup [EXTRACTED 1.00]
- **Devcontainer development Make targets** — readme_make_build, readme_make_run, readme_make_lint, readme_make_clean [EXTRACTED 1.00]

## Communities (8 total, 2 thin omitted)

### Community 0 - "Vagrant Cleanup Command"
Cohesion: 0.40
Nodes (12): cleanLibvirtDomains(), cleanLibvirtDomainsForce(), cleanLibvirtNetworks(), cleanLibvirtVolumes(), commandExists(), destroyAllVMs(), Command, removeAllBoxes() (+4 more)

### Community 1 - "README Build & Install Docs"
Cohesion: 0.20
Nodes (12): devtool-cli (Go developer CLI), devtool-gomodcache named volume, devtool vscode kill, Docker (build dependency), golangci-lint, make build, make clean, make install (+4 more)

### Community 2 - "Graphify Workflow"
Cohesion: 0.36
Nodes (8): graphify-out/graph.json, graphify-out/GRAPH_REPORT.md, graphify (knowledge graph tool), graphify explain <concept>, graphify path <A> <B>, graphify query <question>, graphify update . (AST-only refresh), graphify-out/wiki/index.md

### Community 3 - "Vagrant Cleanup Menu"
Cohesion: 0.32
Nodes (8): devtool vagrant cleanup, Menu option: Clean inactive libvirt domains, Menu option: Clean orphan libvirt volumes, Menu option: Destroy all Vagrant VMs, Menu option: Full cleanup with progress bar, Menu option: Remove unused Vagrant boxes (vagrant box prune), Menu option: Remove all Vagrant boxes, virsh (libvirt CLI dependency)

### Community 4 - "VSCode Kill Command"
Cohesion: 0.52
Nodes (5): findVSCodePIDs(), Command, printProcessTable(), processExists(), runVSCodeKill()

### Community 5 - "UI Printer Utilities"
Cohesion: 0.33
Nodes (3): Info(), Success(), Warning()

## Knowledge Gaps
- **14 isolated node(s):** `github.com/pxblo-x1/devtool-cli`, `graphify update . (AST-only refresh)`, `graphify-out/graph.json`, `graphify-out/wiki/index.md`, `make build` (+9 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `devtool-cli (Go developer CLI)` connect `README Build & Install Docs` to `Graphify Workflow`, `Vagrant Cleanup Menu`?**
  _High betweenness centrality (0.160) - this node is a cross-community bridge._
- **Why does `runVagrantCleanup()` connect `Vagrant Cleanup Command` to `UI Printer Utilities`?**
  _High betweenness centrality (0.093) - this node is a cross-community bridge._
- **Why does `devtool vagrant cleanup` connect `Vagrant Cleanup Menu` to `README Build & Install Docs`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `runVagrantCleanup()` (e.g. with `Info()` and `Success()`) actually correct?**
  _`runVagrantCleanup()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/pxblo-x1/devtool-cli`, `graphify update . (AST-only refresh)`, `graphify-out/graph.json` to the rest of the system?**
  _14 weakly-connected nodes found - possible documentation gaps or missing edges._