chmod 777 video-push
app_status=`ps aux | grep 'video-push' | grep -v grep | wc -l`
if [ $app_status -ne 0 ]
then
	echo -e "\t\t [\e[96mAlready running!\e[0m]"
	exit
fi
printf '%-20s %s' "Starting..."
nohup ./video-push   > error.log 2>&1 &
if [ $? -eq 0 ]
then
	echo -e "\t\t [\e[92mStart success\e[0m]"
	echo "PID is: $!"
else
	echo -e "\t\t [\e[91mStart failed!\e[0m]"
fi
