## To build the flux binary:

```bash
make flux
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
- upx

To build:

```bash
./flux build -t main -p ${PASSWORD}
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
- upx

To build:

```bash
./flux build -t cuda -f devel/cuda.Dockerfile -p ${PASSWORD}
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
- cinc setup script

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
./flux build -t tex -f devel/tex.Dockerfile -p ${PASSWORD}
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