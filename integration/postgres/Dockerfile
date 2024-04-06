FROM postgres:12.4

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install wget -y \
    && apt-get install postgresql-12-postgis-3 -y \
    && apt-get install postgis -y