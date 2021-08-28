วิธีการเชื่อมต่อ oracle database
###### for ubuntu ติดตั้ง Oracle Instant Client ######

ดาวน์โหลดไฟล์ Oracle Instant Client ที่ ลิ้งก์
https://oracle.github.io/odpi/doc/installation.html#linux

สร้างโฟเดอร์ oracle
sudo mkdir -p /opt/oracle

แตกไฟล์ แล้ว copy instantclient_21_3 to /opt/oracle
sudo cp -r instantclient_21_3 /opt/oracle

cd /opt/oracle

sudo apt-get install libaio1

sudo sh -c "echo /opt/oracle/instantclient_21_3 > /etc/ld.so.conf.d/oracle-instantclient.conf"
sudo ldconfig

export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_3:$LD_LIBRARY_PATH

========== for lib godror ============
ต้องติดตั้ง Oracle Instant Client ด้านบนมาก่อนแล้วเพิ่มคำสั่งนี้เข้าไป
ติดตั้ง lib godror for golang
go get github.com/godror/godror

========== เซ็ตเครื่อง for lib go-oci8 ============
ต้องติดตั้ง Oracle Instant Client ด้านบนมาก่อนแล้วเพิ่มคำสั่งนี้เข้าไป
สร้างไฟล์ oci8.pc และเพิ่มคำสั่งนี้เข้าไป
===============================================
prefixdir=/opt/oracle/instantclient_21_3
exec_prefix=${prefixdir}
libdir=${prefixdir}
includedir=${prefixdir}/sdk/include

glib_genmarshal=glib-genmarshal
gobject_query=gobject-query
glib_mkenums=glib-mkenums

Name: oci8
Description: oci8 library
Libs: -L${libdir} -lclntsh
Cflags: -I${includedir}
Version: 12.2
==============================================
ต่อไป Run cmd
sudo cp oci8.pc /usr/lib/pkgconfig
export LD_LIBRARY_PATH=/opt/oracle/instantclient_21_3:$LD_LIBRARY_PATH
export PKG_CONFIG_PATH=/opt/oracle/instantclient_21_3:PKG_CONFIG_PATH

ติดตั้ง lib go-oci8 for golang
go get github.com/mattn/go-oci8