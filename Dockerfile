ARG project_dir=/go/src/golang-pdf-generator-lambda
ARG project_bin=${project_dir}/bin
ARG wkhtmltox_version=0.12.4
ARG wkhtmltox_dir=/opt/wkhtmltox

### Vendors
FROM golang:1-stretch AS vendor

ARG wkhtmltox_version
ARG wkhtmltox_dir

RUN apt-get update && \
    apt-get install -y \
    apt-transport-https \
    curl \
    xz-utils

RUN mkdir -p ${wkhtmltox_dir}

WORKDIR ${wkhtmltox_dir}

RUN curl -LsS https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/${wkhtmltox_version}/wkhtmltox-${wkhtmltox_version}_linux-generic-amd64.tar.xz \
    | tar xJ --strip-components=1

### Deploy
FROM node:8-alpine

ARG project_dir
ARG project_bin
ARG wkhtmltox_dir

ENV AWS_ACCESS_KEY_ID AKIAIKC#####SLT#####
ENV AWS_SECRET_ACCESS_KEY jpBU3kc9teP89d2fbhv7#####icGqEFPIL#####

RUN npm install -g serverless

RUN mkdir -p ${project_bin}

COPY --from=vendor ${wkhtmltox_dir}/bin/wkhtmltopdf ${project_bin}
COPY bin/ ${project_bin}/
COPY serverless.yml ${project_dir}/

RUN chmod -R 755 ${project_bin}/*

WORKDIR ${project_dir}

# RUN serverless deploy -v function --function pdf
RUN serverless deploy -v function --function inlinePdf

# RUN echo "Deployment on $(date '+%Y %b %d %H:%M')" && \
#     serverless deploy -v

# ENTRYPOINT [ "/bin/sh" ]
# CMD serverless deploy -v