#!/usr/bin/make -f

release:
	go build -ldflags="-w -s" -o blutils

deb:
	mkdir -p debian/blutils
	mkdir -p debian/blutils/usr/bin
	cp blutils debian/blutils/usr/bin/blutils
	fakeroot dh_gencontrol
	fakeroot dh_builddeb --destdir=.

clean:
# 	Clean debian package
	rm -rf ./blutils*.deb
	rm -rf ./debian/blutils
	rm -f ./debian/blutils.substvars
	rm -f ./debian/files

#	Clean general files
	rm -f ./blutils

deb_old:
	dpkg-buildpackage -b