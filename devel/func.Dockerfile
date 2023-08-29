# -------------------------------
# machinelearning-one/devel:func
# -------------------------------
# Set the base image
# -------------------
FROM ubuntu:22.04

# Define required arguments
# --------------------------
ARG AUTHOR
ARG EMAIL

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

# Install Go
# -----------
ARG GO_VERSION=1.20.5

ENV PATH=$PATH:/usr/local/go/bin

RUN curl -OL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Install Containerd
# -------------------
ENV CONTAINERD_VERSION="1.7.2"

RUN  wget https://github.com/containerd/containerd/releases/download/v${CONTAINERD_VERSION}/containerd-${CONTAINERD_VERSION}-linux-amd64.tar.gz -O /tmp/containerd.tar.gz && \
    tar -C /usr/local -xvf /tmp/containerd.tar.gz

# Install RunC
# -------------
ENV RUNC_VERSION="1.1.7"

RUN wget https://github.com/opencontainers/runc/releases/download/v${RUNC_VERSION}/runc.amd64 -O /tmp/runc.amd64 && \
    install -m 755 /tmp/runc.amd64 /usr/local/sbin/runc

# Install CNI plugins
# --------------------
ENV CNI_PLUGINS_VERSION="1.3.0"

RUN wget https://github.com/containernetworking/plugins/releases/download/v${CNI_PLUGINS_VERSION}/cni-plugins-linux-amd64-v${CNI_PLUGINS_VERSION}.tgz -O /tmp/cni-plugins.tgz && \
    mkdir -p /opt/cni/bin && \
    tar -xzvf /tmp/cni-plugins.tgz -C /opt/cni/bin    

# Install upx
# ------------
RUN apt update && \
    $APT_INSTALL upx

# Install neovim and tmux
# ------------------------
RUN apt update && \
    $APT_INSTALL \
    gpg-agent && \
    add-apt-repository ppa:neovim-ppa/unstable && \
    apt update && \
    $APT_INSTALL \
    neovim tmux


# Perform cleanup
# ---------------
RUN ldconfig && \
    apt clean && \
    apt autoremove && \
    rm -rf -- /var/lib/apt/lists/* /tmp/* ~/*

# Configure git
# --------------
RUN git config --global user.name "$AUTHOR" && \
    git config --global user.email "$EMAIL"

# Install Tini
# -------------
ENV TINI_VERSION v0.19.0

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini

# Copy the scripts/cinc.sh script to the container
# -------------------------------------------------
COPY scripts/cinc.sh /entrypoint.sh

# Change permissions
#--------------------
RUN chmod +x /tini /entrypoint.sh

# Set the entrypoint
# -------------------
ENTRYPOINT ["/tini", "--", "/entrypoint.sh"]
