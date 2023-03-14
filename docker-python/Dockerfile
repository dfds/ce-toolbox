# ========================================
# CREATE UPDATED BASE IMAGE
# ========================================

FROM python:3.9 AS base

RUN apt-get update \
    && apt-get dist-upgrade -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# # ========================================
# # GENERAL PREREQUISITES
# # ========================================

FROM base
RUN apt-get install -y curl unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN pip3 install poetry


# ========================================
# AWS CLI v2
# ========================================

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm awscliv2.zip


# ========================================
# GO-AWS-SSO
# ========================================

ENV GAS_VERSION=1.3.0
ENV GAS_OS_ARCH=amd64
RUN curl -Lo gas.tar.gz https://github.com/theurichde/go-aws-sso/releases/download/v${GAS_VERSION}/go-aws-sso_${GAS_VERSION}_linux_${GAS_OS_ARCH}.tar.gz \
    && tar zxvf gas.tar.gz && rm gas.tar.gz \
    && mv go-aws-sso /usr/local/bin

# ========================================
# END
# ========================================

CMD [ "bash" ]
