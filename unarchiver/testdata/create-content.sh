#!/bin/bash

mkdir -p content/foo/bar
echo content > content/foo/bar/foobar1.txt
cp create-content.sh content/foo/bar

ln -s foobar1.txt content/foo/bar/foobar2.txt
ln content/foo/bar/foobar1.txt content/foo/bar/foobar3.txt

mkdir -p content/foo2/bar
ln -s ../../foo/bar/foobar1.txt content/foo2/bar/foobar2.txt
ln content/foo/bar/foobar1.txt content/foo2/bar/foobar3.txt

find content | sort | cpio -ov -H ustar -F content.tar
tar -tvf content.tar
gzip --keep content.tar
bzip2 --keep content.tar

pushd content
find foo foo2 | sort | cpio -ov -H ustar -F ../content2.tar
popd
tar -tvf content2.tar
gzip --keep content2.tar
bzip2 --keep content2.tar
