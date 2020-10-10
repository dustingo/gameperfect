#!/bin/bash
status=`systemctl is-active iptables`
if [[ $status == $"active" ]];then
  echo -e "\033[1;36m iptables actived!\033[0m"
else
  echo -e "\033[1;31m iptables not actived\033[0m"
fi

enable=`systemctl is-enabled iptables`
if [[ $enable == $"enabled" ]];then
  echo -e "\033[1;36m iptables is enabled when startup!\033[0m"
else
  echo -e "\033[1;36m iptables is disabled when startup!\033[0m"
  echo "enable iptabes ...."
  systemctl enable iptables
fi
/usr/sbin/iptables -vnL|grep DROP >/dev/null
if [ $? != 0 ];then
  echo -e "\033[1;31m iptables has no drop rule!!\033[0m"
fi
exit 0