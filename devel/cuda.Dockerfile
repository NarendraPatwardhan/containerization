# -------------------------------
# machinelearning-one/devel:cuda
# -------------------------------
# Set the base image
# -------------------
FROM nvidia/cuda:12.1.0-cudnn8-devel-ubuntu22.04

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

# Install neovim
# ---------------
RUN add-apt-repository ppa:neovim-ppa/unstable && \
    apt update && \
    $APT_INSTALL \
    neovim


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