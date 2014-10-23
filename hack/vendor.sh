#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."

# Downloads dependencies into vendor/ directory
mkdir -p vendor
cd vendor

clone() {
	vcs=$1
	pkg=$2
	rev=$3

	pkg_url=https://$pkg
	target_dir=src/$pkg

	echo -n "$pkg @ $rev: "

	if [ -d $target_dir ]; then
		echo -n 'rm old, '
		rm -fr $target_dir
	fi

	echo -n 'clone, '
	case $vcs in
		git)
			git clone --quiet --no-checkout $pkg_url $target_dir
			( cd $target_dir && git reset --quiet --hard $rev )
			;;
		hg)
			hg clone --quiet --updaterev $rev $pkg_url $target_dir
			;;
	esac

	echo -n 'rm VCS, '
	( cd $target_dir && rm -rf .{git,hg} )

	echo done
}

clone git github.com/kr/pty 67e2db24c8

clone git github.com/gorilla/context 14f550f51a

clone git github.com/gorilla/mux 136d54f81f

clone git github.com/tchap/go-patricia v1.0.1

clone git github.com/digitalocean/godo v0.4.0

clone git github.com/rackspace/gophercloud a80ea95065ba07f59674f6189efc1de4d3218e32

clone git github.com/mitchellh/mapstructure 740c764bc6149d3f1806231418adb9f52c11bcbf

clone git github.com/racker/perigee 0c00cb0a026b71034ebc8205263c77dad3577db5

clone git github.com/google/go-querystring 30f7a39f4a218feb5325f3aebc60c32a572a8274

clone git github.com/tent/http-link-go ac974c61c2f990f4115b119354b5e0b47550e888

clone hg code.google.com/p/go.net 84a4013f96e0

clone hg code.google.com/p/gosqlite 74691fb6f837

clone git github.com/docker/libtrust d273ef2565ca

clone hg code.google.com/p/goauth2 afe77d958c70

clone git github.com/MSOpenTech/azure-sdk-for-go 43719165c8b42a47f6b9d13b206271842396bdaf

# get Go tip's archive/tar, for xattr support and improved performance
# TODO after Go 1.4 drops, bump our minimum supported version and drop this vendored dep
if [ "$1" = '--go' ]; then
	# Go takes forever and a half to clone, so we only redownload it when explicitly requested via the "--go" flag to this script.
	clone hg code.google.com/p/go 1b17b3426e3c
	mv src/code.google.com/p/go/src/pkg/archive/tar tmp-tar
	rm -rf src/code.google.com/p/go
	mkdir -p src/code.google.com/p/go/src/pkg/archive
	mv tmp-tar src/code.google.com/p/go/src/pkg/archive/tar
fi

clone git github.com/docker/libcontainer 8d1d0ba38a7348c5cfdc05aea3be34d75aadc8de
# see src/github.com/docker/libcontainer/update-vendor.sh which is the "source of truth" for libcontainer deps (just like this file)
rm -rf src/github.com/docker/libcontainer/vendor
eval "$(grep '^clone ' src/github.com/docker/libcontainer/update-vendor.sh | grep -v 'github.com/codegangsta/cli')"
# we exclude "github.com/codegangsta/cli" here because it's only needed for "nsinit", which Docker doesn't include
