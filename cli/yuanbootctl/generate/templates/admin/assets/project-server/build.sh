
#!/bin/sh
out_file="sendex-server"
VersionPath="sendex-server/version"
[ $# -lt 4 ] && {
	echo "Usage: $0 1.0.0-SNAPSHOT linux amd64 dev"
	exit 1
}
build() {
  local version="$1"
  local os="$2"
  local arch="$3"
  local profile="$4"
  local dir="build/$out_file-$os-$arch"
  out_file="${out_file}-${version}"
  go env -w GOSUMDB=off
  [ "$os" = "windows" ] && {
  		out_file="${out_file}.exe"
  	}
  rm -rf $dir
  mkdir -p $dir
  GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -ldflags="-s -w -X $VersionPath.env=$profile " -o "${dir}/${out_file}" .

}

main() {
  echo "mod download"
  go get -t .
  go mod download
  go get $out_file
  build $1 $2 $3 $4
}

main $1 $2 $3 $4
