#!/bin/bash
rpm -q ntp >/dev/null
if [ $? == 0 ];then
  echo "ntp has been installed"
else
  echo "package ntp not installed and  then install"
  yum install ntp >/dev/null
fi
echo "ntpd server:"
cat /etc/ntp.conf |grep -v "#"|grep -v ^$|grep server|awk '{print $2}'
echo "Status:"
systemctl is-active ntpd
echo "Enable:"
systemctl is-enabled ntpd
exit 0