#!/bin/bash

cd static/fimpui
ng build --prod --deploy '/fimp/static/'
cd ../../
GOPATH=~/DevProjects/APPS/GOPATH GOOS=linux GOARCH=arm GOARM=6 go build -o fimpui
cp fimpui debian/opt/fimpui
cp -R static/fimpui/dist debian/opt/fimpui/static/fimpui/dist
cp -R static/fhcore debian/opt/fimpui/static/fhcore
#tar cvzf fimpui.tar.gz fimpui VERSION static/fimpui/dist static/fhcore
#dpkg-deb --build debian