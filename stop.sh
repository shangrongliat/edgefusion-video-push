printf '%-60s %s' "Stopping... video-push"
pid=`ps aux | grep 'video-push' | grep -v grep | awk '{print $2}'`
if [ -n "$pid" ]
then
  	#关闭进程
  	`kill -9 $pid`
  	#判断关闭是否成功
  	if [ $? -eq 0 ]
    	then
       		echo -e "\t\t [\e[92mStop success\e[0m]"
  	else
       		echo -e "\t\t [\e[91mStop failed!\e[0m]"
  	fi
  else
      echo -e "\t\t [\e[91mNot running\e[0m]"
  fi

