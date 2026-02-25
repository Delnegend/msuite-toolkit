FROM debian:trixie-slim

RUN apt-get update && apt-get install -y \
    curl git sudo unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1000 mstk && \
    useradd -m -u 1000 -g mstk mstk && \
    echo "mstk ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/mstk

USER mstk