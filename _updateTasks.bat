rem This batch will stop WatchMe, update it's config file and start it again.

rem Update filepaths to the todo.txt and WatchMe config files
SET mytodotxt="p:\Phone\Documents\Notes\Lifeplan\_TODO\todo.txt"
rem SET mywatchme="WatchMeConfig.xml"


rem you probably do not need to modify anything below
TASKKILL /IM WatchMe.exe 
timeout /t 1
todo-watchme_bridge.exe -td %mytodotxt%
start watchme.exe
exit