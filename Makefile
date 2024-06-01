#!/usr/bin/make -f

.PHONY: release release_arm64 clean

release:
	cp cmd/version cmd/version.bak

	echo "COMMIT=$(shell git rev-parse --short=7 HEAD)" >> cmd/version;
	echo "UNIX_TIMESTAMP=$(shell date +%s)" >> cmd/version;

    # Remove empty lines from the version file
	sed -i '/^$$/d' cmd/version

	go build -ldflags="-w -s" -o blutils

    # Reset versionfile
	mv cmd/version.bak cmd/version

deb:
	mkdir -p debian/blutils
	mkdir -p debian/blutils/usr/bin
	cp blutils debian/blutils/usr/bin/blutils
	fakeroot dh_gencontrol
	fakeroot dh_fixperms
	fakeroot dh_builddeb --destdir=.

clean:
    # Clean debian package
	rm -rf ./blutils*.deb
	rm -rf ./debian/blutils
	rm -f ./debian/blutils.substvars
	rm -f ./debian/files

    # Clean general files
	rm -rf ./build
	rm -f ./blutils

#* 	ARM64
release_arm64:
	GOARCH=arm64 go build -ldflags="-w -s" -o blutils

deb_arm64:
	mkdir -p debian/blutils
	mkdir -p debian/blutils/usr/bin
	cp blutils debian/blutils/usr/bin/blutils
	fakeroot dh_gencontrol -- -DArchitecture=arm64
	fakeroot dh_fixperms
	fakeroot dh_builddeb --destdir=. -a arm64

deb_old:
	dpkg-buildpackage -b
