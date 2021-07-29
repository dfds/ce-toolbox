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
RUN pip3 install poetry


# ========================================
# END
# ========================================

CMD [ "bash" ]
