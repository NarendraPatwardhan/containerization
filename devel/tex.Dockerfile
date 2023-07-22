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