#!/bin/bash

cd static/fimpui
ng build --prod --deploy '/fimp/static/'
cd ../../
GOOS=linux GOARCH=arm GOARM=6 go build -o fimpui_arm
tar cvzf fimpui.tar.gz fimpui_arm static/fimpui/dist
