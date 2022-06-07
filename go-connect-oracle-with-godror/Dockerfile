################## script v1 for library godror ########################
# FROM golang:1.17-alpine
# RUN apk update && apk add gcc git && apk add libc-dev
# RUN apk --no-cache add libaio libnsl libc6-compat curl && \
#     cd /tmp && \
#     curl -o instantclient-basiclite.zip https://download.oracle.com/otn_software/linux/instantclient/instantclient-basiclite-linuxx64.zip -SL && \
#     unzip instantclient-basiclite.zip && \
#     mv instantclient*/ /usr/lib/instantclient && \
#     rm instantclient-basiclite.zip && \
#     ln -s /usr/lib/instantclient/libclntsh.so.19.1 /usr/lib/libclntsh.so && \
#     ln -s /usr/lib/instantclient/libocci.so.19.1 /usr/lib/libocci.so && \
#     ln -s /usr/lib/instantclient/libociicus.so /usr/lib/libociicus.so && \
#     ln -s /usr/lib/instantclient/libnnz19.so /usr/lib/libnnz19.so && \
#     ln -s /usr/lib/libnsl.so.2 /usr/lib/libnsl.so.1 && \
#     ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 && \
#     ln -s /lib64/ld-linux-x86-64.so.2 /usr/lib/ld-linux-x86-64.so.2

# ENV ORACLE_BASE=/usr/lib/instantclient
# ENV LD_LIBRARY_PATH=/usr/lib/instantclient
# ENV TNS_ADMIN=/usr/lib/instantclient
# ENV ORACLE_HOME=/usr/lib/instantclient    

# WORKDIR /app
# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
# COPY *.go ./
# RUN go build -o main .
# EXPOSE 8000
# CMD [ "/app/main" ]
################## end script v1 for library godror ########################

#=============================================================
#--------------------- build stage ---------------------------
#=============================================================
FROM golang:1.16-stretch AS build_stage
ENV PACKAGE_PATH=api/test
RUN mkdir -p /go/src/
WORKDIR /go/src/$PACKAGE_PATH
COPY . /go/src/$PACKAGE_PATH/
RUN go mod download
RUN go build -o test
#=============================================================
#--------------------- final stage ---------------------------
#=============================================================
FROM oracle7/instantclient:12.2.0.1 AS final_stage
ENV PACKAGE_PATH=api/test
COPY --from=build_stage /go/src/$PACKAGE_PATH/test /go/src/$PACKAGE_PATH/
# COPY --from=build_stage /go/src/$PACKAGE_PATH/configs /go/src/$PACKAGE_PATH/configs
WORKDIR /go/src/$PACKAGE_PATH/
ENTRYPOINT ./test
EXPOSE 8000