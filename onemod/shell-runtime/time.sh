
#!/bin/bash

#print time
for((i=0;i<100;i++))
do
    sleep 1
    echo $(date +"%Y-%m-%d %H:%M:%S")
done
