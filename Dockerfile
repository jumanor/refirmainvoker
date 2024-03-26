# Usa la imagen base httpd:2.4
FROM httpd:2.4

COPY ./example/ /usr/local/apache2/htdocs

COPY ./public /opt/public
COPY main /opt
COPY config.properties.example /opt
RUN chmod +x /opt/main


COPY entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV CLIENT_ID=""
ENV CLIENT_SECRET=""

EXPOSE 80
EXPOSE 9091

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Actualiza el índice de paquetes y luego instala las dependencias necesarias
RUN apt-get update && \
    apt-get install -y wget && \
    apt-get install -y p7zip-full && \
    rm -rf /var/lib/apt/lists/*

# Descarga e instala Go 1.19
RUN wget https://golang.org/dl/go1.19.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz && \
    rm go1.19.linux-amd64.tar.gz

# Establece las variables de entorno para Go
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go




