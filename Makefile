version="0.6.16"
version_file=VERSION
working_dir=$(shell pwd)
arch="armhf"

build-js:
	cd static/fimpui;ng build --prod --deploy '/fimp/static/'

build-go-arm:
	GOPATH=~/DevProjects/APPS/GOPATH GOOS=linux GOARCH=arm GOARM=6 go build -o fimpui

build-go:
	GOPATH=~/DevProjects/APPS/GOPATH go build -o fimpui

build-go-amd:
	GOPATH=~/DevProjects/APPS/GOPATH GOOS=linux GOARCH=amd64 go build -o fimpui

clean:
	-rm -R debian/opt/fimpui/static/fhcore/*
	-rm -R debian/opt/fimpui/static/fimpui/dist/*
	-rm debian/opt/fimpui/fimpui
	-rm fimpui

configure-arm:
	python ./scripts/config_env.py prod $(version) armhf

configure-amd64:
	python ./scripts/config_env.py prod $(version) amd64

configure-dev-js:
	python ./scripts/config_env.py dev $(version) armhf	

package-tar:
	tar cvzf fimpui.tar.gz fimpui VERSION static/fimpui/dist static/fhcore

package-deb-doc:
	@echo "Packaging application as debian package"
	chmod a+x debian/DEBIAN/*
	cp fimpui debian/opt/fimpui
	cp VERSION debian/opt/fimpui
	cp -R static/fimpui/dist debian/opt/fimpui/static/fimpui/
	cp -R static/fhcore debian/opt/fimpui/static/
	docker run --rm -v ${working_dir}:/build -w /build --name debuild debian dpkg-deb --build debian
	@echo "Done"

tar-arm: build-js build-go-arm package-deb-doc
	@echo "The application was packaged into tar archive "

deb-arm : clean configure-arm build-js build-go-arm package-deb-doc
	mv debian.deb fimpui_$(version)_armhf.deb

deb-amd : configure-amd64 build-js build-go-amd package-deb-doc
	mv debian.deb fimpui_$(version)_amd64.deb

set-dev : configure-dev-js build-go

build-mac : build-js build-go

.phony : clean