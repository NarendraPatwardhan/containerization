# -------------------------------
# machinelearning-one/devel:cuda
# -------------------------------
# Set the base image
# -------------------
FROM nvidia/cuda:12.4.1-cudnn-devel-ubuntu22.04

# Define required arguments
# --------------------------
ARG AUTHOR
ARG EMAIL
ARG USERNAME
ARG PASSWORD

# Define commands for headless installation
# ------------------------------------------
ENV DEBIAN_FRONTEND=noninteractive \
    APT_INSTALL="apt install -y --no-install-recommends"

# Setup the locale
# -----------------
ENV LANG=C.UTF-8

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

# Install Zig
# ------------
ARG ZIG_VERSION=0.13.0

ENV PATH=$PATH:/usr/local/zig

RUN curl -fsSL https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz \
    --output zig.tar.xz \
    && tar -xf zig.tar.xz && mv zig-linux-x86_64-${ZIG_VERSION} /usr/local/zig \
    && rm zig.tar.xz

# Install Bazel

# --------------
# Install Bazel
# --------------
ARG BAZELISK_VERSION=1.21.0

RUN curl -fsSL -o /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/download/v${BAZELISK_VERSION}/bazelisk-linux-amd64 && chmod +x /usr/local/bin/bazel

# Install upx
# ------------
RUN apt update && \
    $APT_INSTALL upx

# Install neovim and tmux
# ------------------------
RUN add-apt-repository ppa:neovim-ppa/unstable && \
    apt update && \
    $APT_INSTALL \
    neovim tmux

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
