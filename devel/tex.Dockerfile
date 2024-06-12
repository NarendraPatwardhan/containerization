# ------------------------------
# machinelearning-one/devel:tex
# ------------------------------
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
    PIP_INSTALL="python3 -m pip --no-cache-dir install --upgrade"

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

# Install python 3 and pip
# ---------------------------
ENV PATH=$PATH:~/.local/bin

RUN apt update && \
    $APT_INSTALL \
    python3 \
    python3-dev \
    python3-distutils \
    python3-venv \
    && \
    curl -o ~/get-pip.py https://bootstrap.pypa.io/get-pip.py && \
    python3 ~/get-pip.py

# Install numfocus and allied packages
# -------------------------------------
RUN $PIP_INSTALL \
    numpy \
    scipy \
    pandas \
    scikit-image \
    scikit-learn \
    matplotlib \
    seaborn

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

# Install texlive
# ----------------
RUN apt update && \
    $APT_INSTALL texlive-full

# Install Manim
# --------------
RUN apt update && \
    $APT_INSTALL \
    libcairo2-dev \
    libpango1.0-dev \
    ffmpeg && \
    $PIP_INSTALL \
    manim

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
