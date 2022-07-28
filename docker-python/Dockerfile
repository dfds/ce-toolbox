# ========================================
#CREATE UPDATED BASE IMAGE
# ========================================

FROM python:3.9 AS base


RUN export BUILD_ARCHITECTURE_UNIX=$(uname -m); \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "x86_64" ]; then export BUILD_ARCHITECTURE=amd64; fi; \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "aarch64" ]; then export BUILD_ARCHITECTURE=arm64; fi;


RUN echo "build architecture:  $BUILD_ARCHITECTURE_UNIX"

RUN apt-get update \
    && apt-get dist-upgrade -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# # ========================================
# # GENERAL PREREQUISITES
# # ========================================

FROM base
RUN apt-get update
RUN apt-get install -y curl unzip cargo \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN export CRYPTOGRAPHY_DONT_BUILD_RUST=1
RUN pip3 install poetry


# ========================================
# AWS CLI v2
# ========================================

RUN export BUILD_ARCHITECTURE_UNIX=$(uname -m); \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "x86_64" ]; then export BUILD_ARCHITECTURE=amd64; fi; \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "aarch64" ]; then export BUILD_ARCHITECTURE=arm64; fi; \
    echo "attempting to get @ https://awscli.amazonaws.com/awscli-exe-linux-${BUILD_ARCHITECTURE_UNIX}.zip" \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-${BUILD_ARCHITECTURE_UNIX}.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm awscliv2.zip

# ========================================
# SAML2AWS
# ========================================

ENV SAML2AWS_VERSION=2.35.0

RUN export BUILD_ARCHITECTURE_UNIX=$(uname -m); \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "x86_64" ]; then export BUILD_ARCHITECTURE=amd64; fi; \
    if [ "$BUILD_ARCHITECTURE_UNIX" = "aarch64" ]; then export BUILD_ARCHITECTURE=arm64; fi; \
    echo "attempting to get @ https://github.com/Versent/saml2aws/releases/download/v${SAML2AWS_VERSION}/saml2aws_${SAML2AWS_VERSION}_linux_${BUILD_ARCHITECTURE}.tar.gz" \
    && curl -L "https://github.com/Versent/saml2aws/releases/download/v${SAML2AWS_VERSION}/saml2aws_${SAML2AWS_VERSION}_linux_${BUILD_ARCHITECTURE}.tar.gz" -o saml2aws.tar.gz \
    && tar -zxvf saml2aws.tar.gz \
    && rm saml2aws.tar.gz \
    && mv saml2aws /usr/local/bin/


# ========================================
# END
# ========================================

CMD [ "bash" ]