#!/bin/bash

cd static/fimpui
ng build --prod --deploy '/fimp/static/'
cd ../../
GOPATH=~/DevProjects/APPS/GOPATH GOOS=linux GOARCH=arm GOARM=6 go build -o fimpui_arm
tar cvzf fimpui.tar.gz fimpui_arm VERSION static/fimpui/dist
