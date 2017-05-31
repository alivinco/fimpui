version="0.0.6"
version_file=VERSION
working_dir=$(shell pwd)

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
	mv debian.deb fimpui_$(version)_arm.deb
	@echo "Done"

git-version:
	@echo "Getting latest version from git"
	@git describe --tags | sed 's/^v//' | sed 's/-\(alpha\|beta\|rc\)/~\1/' | sed 's/-/./' | sed 's/-g/+/' | sed 's/-/./' > "$(working_dir)/VERSION"
	@sed -i "s/^VERSION.*/VERSION = \"$(version)\"/" "$(version_file)"
	@echo "Set version to $(version)"

tar-arm: build-js build-go-arm package-deb-doc
	@echo "The application was packaged into tar archive "

make-arm : build-js build-go-arm package-deb-doc

.phony : clean