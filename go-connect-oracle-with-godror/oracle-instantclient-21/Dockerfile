###### oraclelinux7/12.2.0.1 ######
FROM oraclelinux:7-slim

ADD oracle-instantclient*.rpm /tmp/

RUN  yum -y install /tmp/oracle-instantclient*.rpm && \
     rm -rf /var/cache/yum && \
     rm -f /tmp/oracle-instantclient*.rpm && \
     echo /usr/lib/oracle/12.2/client64/lib > /etc/ld.so.conf.d/oracle-instantclient12.2.conf && \
     ldconfig

ENV PATH=$PATH:/usr/lib/oracle/12.2/client64/bin

CMD ["sqlplus", "-v"]

###### oraclelinux7/21 ######
# FROM oraclelinux:7-slim

# RUN  yum -y install oracle-instantclient-release-el7 && \
#      yum -y install oracle-instantclient-basic oracle-instantclient-devel oracle-instantclient-sqlplus && \
#      rm -rf /var/cache/yum

# # Uncomment if the tools package is added
# # ENV PATH=$PATH:/usr/lib/oracle/21/client64/bin

# CMD ["sqlplus", "-v"]

###### oraclelinux8/21 ######
# FROM oraclelinux:8

# RUN  dnf -y install oracle-instantclient-release-el8 && \
#      dnf -y install oracle-instantclient-basic oracle-instantclient-devel oracle-instantclient-sqlplus && \
#      rm -rf /var/cache/dnf

# # Uncomment if the tools package is added
# # ENV PATH=$PATH:/usr/lib/oracle/21/client64/bin

# CMD ["sqlplus", "-v"]

