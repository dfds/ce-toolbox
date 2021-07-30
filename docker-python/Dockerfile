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
RUN apt-get install -y curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
RUN pip3 install poetry


# ========================================
# AWS IAM AUTHENTICATOR
# ========================================

ENV AWSIAMAUTH_VERSION=1.15.10/2020-02-22

RUN curl -L https://amazon-eks.s3-us-west-2.amazonaws.com/${AWSIAMAUTH_VERSION}/bin/linux/amd64/aws-iam-authenticator -o aws-iam-authenticator \
    && chmod +x aws-iam-authenticator \
    && mv aws-iam-authenticator /usr/local/bin/


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
