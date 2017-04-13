

**Build Angular 2** : 
cd static/fimpui
ng build 
 or 
ng build --prod --deploy '/fimp/static/'

cd ../../
**Run Angular dev app** :

ng serve 
open http://localhost:4200

**Rsync static files for deverlopment:**
 
GOOS=linux GOARCH=arm GOARM=6 go build -o fimpui_arm

rsync -a static/fimpui/dist fh@aleks.local:~/fimpui/static/fimpui/

scp fimpui_arm fh@aleks.local:~/fimpui/

**Package** :
tar cvzf fimpui.tar.gz fimpui_arm static/fimpui/dist
scp fimpui.tar.gz fh@aleks.local:~/fimpui/
**Unpackage** : 
tar -xvf fimpui.tar.gz
Update static/fimpui/dist/index.html
