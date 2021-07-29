# ========================================
# CREATE UPDATED BASE IMAGE
# ========================================

FROM python:3.9-slim AS base

RUN apt-get update \
    && apt-get dist-upgrade -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# # ========================================
# # GENERAL PREREQUISITES
# # ========================================

FROM base
RUN apt-get install -y curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN pip3 install poetry


# ========================================
# SAML2AWS
# ========================================

ENV SAML2AWS_VERSION=2.20.0

RUN curl -L "https://github.com/Versent/saml2aws/releases/download/v${SAML2AWS_VERSION}/saml2aws_${SAML2AWS_VERSION}_linux_amd64.tar.gz" -o saml2aws.tar.gz \
    && tar -zxvf saml2aws.tar.gz \
    && rm saml2aws.tar.gz \
    && mv saml2aws /usr/local/bin/


# ========================================
# END
# ========================================

CMD [ "bash" ]
