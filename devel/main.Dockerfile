# -------------------------------
# machinelearning-one/devel:main
# -------------------------------
# Set the base image
# -------------------
FROM ubuntu:22.04

# Define required arguments
# --------------------------
ARG AUTHOR
ARG EMAIL
ARG USERNAME
ARG PASSWORD

# Define commands for headless installation
# ------------------------------------------
ENV DEBIAN_FRONTEND=noninteractive \
    APT_INSTALL="apt install -y --no-install-recommends" \
    PIP_INSTALL="python -m pip --no-cache-dir install --upgrade"

# Setup the locale
# -----------------
RUN apt update && \
    $APT_INSTALL locales && \
    locale-gen en_US.UTF-8 && \
    echo LANG=en_US.UTF-8 > /etc/default/locale

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# Install basic utilities
# ------------------------
RUN apt update && \
    $APT_INSTALL \
    sudo \
    build-essential \
    cmake \
    apt-utils \
    apt-transport-https \
    software-properties-common \
    ca-certificates \
    wget \
    curl \
    git \
    vim \
    libssl-dev \
    openssh-client \
    unzip \
    unrar

# Install python 3.10 and pip
# ---------------------------
ENV PATH=$PATH:~/.local/bin

RUN apt update && \
    $APT_INSTALL \
    python3.10 \
    python3.10-dev \
    python3.10-distutils \
    && \
    curl -o ~/get-pip.py https://bootstrap.pypa.io/get-pip.py && \
    python3.10 ~/get-pip.py && \
    ln -s /usr/bin/python3.10 /usr/local/bin/python

# Install numfocus and allied packages
# -------------------------------------
RUN $PIP_INSTALL \
    numpy \
    scipy \
    pandas \
    scikit-image \
    scikit-learn \
    matplotlib \
    seaborn \
    einops \
    hydra-core \
    hydra-colorlog \
    hydra-optuna-sweeper

# Install jax and allied packages
# ------------------------------------
RUN $PIP_INSTALL \
    "jax[cuda12_pip]" -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html  && \
    $PIP_INSTALL \
    dm-haiku \
    optax \
    chex \
    rlax \
    jraph \
    distrax \
    mctx

# Install pytorch and allied packages
# ------------------------------------
RUN $PIP_INSTALL \
    torch torchvision torchaudio && \
    $PIP_INSTALL \
    lightning \
    torchmetrics

# Install onnx
# -------------
RUN apt update && \
    $APT_INSTALL \
    protobuf-compiler \
    libprotoc-dev && \
    $PIP_INSTALL \
    protobuf \
    onnx \
    onnxruntime-gpu

# Install quality of life packages
# ---------------------------------
RUN $PIP_INSTALL \
    Cython \
    typing \
    pre-commit \
    black[jupyter] \
    flake8 \
    isort \
    nbstripout \
    python-dotenv \
    tqdm \
    rich \
    pytest \
    sh \
    pudb \
    twine

# Install jupyterlab
# -------------------
RUN $PIP_INSTALL \
    jupyterlab \
    ipywidgets

# Install opencv
# ---------------
RUN apt update && \
    $APT_INSTALL \
    libopencv-dev python3-opencv

# Install docker-ce-cli
# ----------------------
RUN apt update && \
    $APT_INSTALL \
    gpg-agent && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu jammy stable" && \
    $APT_INSTALL \
    docker-ce-cli \
    docker-compose-plugin \
    docker-buildx-plugin

# Install docker sdk for python
# ------------------------------
RUN $PIP_INSTALL \
    docker

# Install Rust
# -------------
ARG RUSTUP_VERSION=1.26.0
ARG RUST_VERSION=1.69.0

ENV RUSTUP_HOME=/usr/local/rustup \
    CARGO_HOME=/usr/local/cargo \
    PATH=$PATH:/usr/local/cargo/bin

RUN RUST_ARCH=x86_64-unknown-linux-gnu && \
    url="https://static.rust-lang.org/rustup/archive/${RUSTUP_VERSION}/${RUST_ARCH}/rustup-init" && \
    curl -o rustup-init "$url" && \
    chmod +x rustup-init && \
    ./rustup-init -y --no-modify-path --profile minimal --default-toolchain $RUST_VERSION --default-host $RUST_ARCH && \
    rm rustup-init && \
    chmod -R a+w $RUSTUP_HOME $CARGO_HOME

# Install Go
# -----------
ARG GO_VERSION=1.20.3

ENV PATH=$PATH:/usr/local/go/bin

RUN curl -OL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Install Deno
# -------------
ARG DENO_VERSION=1.33.0

RUN curl -fsSL https://github.com/denoland/deno/releases/download/v${DENO_VERSION}/deno-x86_64-unknown-linux-gnu.zip \
    --output deno.zip \
    && unzip deno.zip -d /usr/local/bin/ \
    && chmod 755 /usr/local/bin/deno \
    && rm deno.zip

# Install NVM and Node, enable alternative package managers
# ----------------------------------------------------------
ARG NVM_VERSION=v0.39.3
ARG NODE_VERSION=18.16.0

ENV NVM_DIR=/usr/local/nvm \
    NODE_PATH=$NVM_DIR/v${NODE_VERSION}/lib/node_modules \
    PATH=$PATH:$NVM_DIR/versions/node/v${NODE_VERSION}/bin

RUN mkdir -p $NVM_DIR && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/${NVM_VERSION}/install.sh | /bin/bash && \
    . ${NVM_DIR}/nvm.sh && \
    . ${NVM_DIR}/bash_completion && \
    echo ". ${NVM_DIR}/nvm.sh" >> /etc/bash.bashrc && \
    echo ". ${NVM_DIR}/bash_completion" >> /etc/bash.bashrc && \
    nvm install $NODE_VERSION && \
    nvm alias default $NODE_VERSION && \
    nvm use default && \
    corepack enable

# Install Elixir
# ---------------
RUN apt update && \
    $APT_INSTALL \
    elixir

# Install Zig
# ------------
ARG ZIG_VERSION=0.10.1

ENV PATH=$PATH:/usr/local/zig

RUN curl -fsSL https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz \
    --output zig.tar.xz \
    && tar -xf zig.tar.xz && mv zig-linux-x86_64-${ZIG_VERSION} /usr/local/zig \
    && rm zig.tar.xz

# Install Bazel
# --------------
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor >bazel-archive-keyring.gpg && \
    mv bazel-archive-keyring.gpg /usr/share/keyrings && \
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/bazel-archive-keyring.gpg] https://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && \
    apt update && \
    $APT_INSTALL \
    bazel

# Install upx
# ------------
RUN apt update && \
    $APT_INSTALL upx

# Perform cleanup
# ---------------
RUN ldconfig && \
    apt clean && \
    apt autoremove && \
    rm -rf -- /var/lib/apt/lists/* /tmp/* ~/*

# Create and configure the user
# ------------------------------
RUN --mount=type=secret,id=password useradd -m -s /bin/bash -p $(openssl passwd -1 $(cat /run/secrets/password)) $USERNAME && \
    # Add the user to the sudo group
    # -------------------------------
    usermod -aG sudo $USERNAME && \
    # Configure git
    # --------------
    git config --global user.name "$AUTHOR" && \
    git config --global user.email "$EMAIL"

# Set the default user and working directory
# -------------------------------------------
USER $USERNAME
WORKDIR /home/$USERNAME