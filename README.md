# devtool-cli

CLI de herramientas de desarrollo. Migración a Go de scripts de shell para gestión de procesos VS Code Server y limpieza de entornos Vagrant/libvirt.

## Requisitos

- Docker (para compilar desde el host)
- `sudo` (para instalar en `/usr/local/bin`)

## Instalación

Desde el directorio del proyecto en el **host** (no es necesario entrar al devcontainer):

```bash
make install
```

Esto compila el binario usando un contenedor Go efímero y lo instala en `/usr/local/bin/devtool`.

> La primera ejecución descarga la imagen de Go y las dependencias. Las siguientes son rápidas gracias al volumen de caché `devtool-gomodcache`.

## Desinstalación

```bash
make uninstall
```

## Uso

### `devtool vscode kill`

Busca y termina procesos huérfanos de VS Code Remote SSH para el usuario actual.

```bash
devtool vscode kill
```

Flujo:
1. Muestra un spinner mientras busca procesos de `.vscode-server`
2. Lista los procesos encontrados en una tabla (PID, %CPU, %MEM, comando)
3. Pide confirmación antes de actuar
4. Envía `SIGTERM`; si algún proceso resiste, envía `SIGKILL`

### `devtool vagrant cleanup`

Muestra el estado actual de VMs, boxes y recursos de libvirt, y ofrece un menú interactivo de limpieza.

```bash
devtool vagrant cleanup
```

Opciones del menú (navegación con flechas ↑↓):

| Opción | Acción |
|--------|--------|
| 1 | Destruir todas las VMs de Vagrant |
| 2 | Eliminar boxes no utilizados (`vagrant box prune`) |
| 3 | Eliminar **todos** los boxes |
| 4 | Limpiar volúmenes huérfanos de libvirt |
| 5 | Limpiar dominios inactivos de libvirt |
| 6 | Limpieza completa con barra de progreso |
| 7 | Salir |

> Las opciones 4 y 5 requieren `virsh` instalado en el sistema.

## Desarrollo

Targets disponibles dentro del devcontainer:

```bash
make build          # Compila ./devtool localmente
make run ARGS="vscode kill"  # Ejecuta sin compilar
make lint           # Corre golangci-lint
make clean          # Elimina el binario local
```
