## To build the flux binary:

```bash
make flux
```

## To set password for non-root user:

Required for `main`, `cuda`, and `tex` images

Copy the .env.tmpl file to .env and set the password

```bash
cp .env.tmpl .env
```

Replace the value for the password in the .env file

```env
PASSWORD={{VALUE}}
```

## The machinelearning/devel:main image:

Use: Primary development environment

- Build essentials
- Python 3
    - Numfocus allied packages
    - Jax and allied packages
    - Pytorch and allied packages
    - Onnx
    - Quality of life packages
    - Jupyterlab
    - OpenCV
- Docker out of docker
- Rust
- Go
- Deno
- Node
- Elixir
- Zig
- Bazel
- upx

To build:

```bash
./flux build -t main
```
To run:

```bash
./flux up ${CONTAINER_NAME}
```

To remove:

```bash
./flux down ${CONTAINER_NAME}
```

To prune:

```bash
./flux down -p ${CONTAINER_NAME}
```

## The machinelearning/devel:cuda image:

Use: To develop cuda applications

- CUDA development essentials
- Build essentials
- Zig
- Bazel
- upx

To build:

```bash
./flux build -t cuda -f devel/cuda.Dockerfile
```
To run:

```bash
./flux up -t cuda -r none ${CONTAINER_NAME}
```

To remove:

```bash
./flux down ${CONTAINER_NAME}
```

To prune:

```bash
./flux down -p ${CONTAINER_NAME}
```

## The machinelearning/devel:func image:

Use: To develop container ecosystem applications

- Build essentials
- Go
- Containerd
- RunC
- CNI plugins
- upx

Configured to run tini as pid 1 and cacheable cinc setup

To build:

```bash
./flux build -t func -f devel/func.Dockerfile -u root
```
To run:

```bash
./flux up -t func -u root -r containerd ${CONTAINER_NAME}
```

To remove:

```bash
./flux down ${CONTAINER_NAME}
```

To prune:

```bash
./flux down -p ${CONTAINER_NAME}
docker volume prune
```

## The machinelearning/devel:tex image:

Use: To write tex documents and create manim visualizations

- Build essentials
- Python 3
    - Numfocus allied packages
    - Quality of life packages
    - Jupyterlab
    - OpenCV
- Texlive-full
- Manim


To build:

```bash
./flux build -t tex -f devel/tex.Dockerfile
```
To run:

```bash
./flux up -t tex -r none ${CONTAINER_NAME}
```

To remove:

```bash
./flux down ${CONTAINER_NAME}
```

To prune:

```bash
./flux down -p ${CONTAINER_NAME}
```

## To connect to an existing container:

```bash
./flux connect ${CONTAINER_NAME}
```

## Hooks

### SSH Key Generation

To push code to remote repositories, it is recommended to use ssh over https. To generate a key pair, the following command is typically run:

```bash
ssh-keygen -t ed25519 -C ${COMMENT}
```
Flux provides a thin utility wrapper to make this process easier:

```bash
./flux hooks ssh-keygen ${CONTAINER_NAME}
```

You can optionally specify a comment for the key.
```bash
./flux hooks ssh-keygen -c ${COMMENT} ${CONTAINER_NAME}
```

It is recommended you set up a passphrase different from user password.

If you accept the default prompt for the location, public key can be found at ~/.ssh/id_ed25519.pub, otherwise it would be at the location you specified.

Copy the public key and follow the instructions of your VCS provider (GitHub/GitLab etc.) to assign it as trusted.

### Neovim as IDE Setup

To setup neovim as an IDE, run the following command:

```bash
./flux hooks nvim-setup ${CONTAINER_NAME}
```

This installs nvchad and pulls a personal configuration that works in a polyglot environment.