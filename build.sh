SDK_NAME=imgPrint
VERSION=1.0.0

path_source=$(cd $(dirname $0); pwd);
project_name=$(basename $path_source)
# output_file_name_mac=${SDK_NAME}_mac;
# output_file_name_linux=${SDK_NAME}_linux;
# output_file_name_win=${SDK_NAME}_win;

path_release=$path_source/release
file_path_min_mac=$path_release/mac/$SDK_NAME
file_path_mac=$file_path_min_mac'_source'

file_path_min_linux=$path_release/linux/$SDK_NAME
file_path_linux=$file_path_min_linux'_source'

file_path_min_win=$path_release/win/$SDK_NAME'.exe'
file_path_win=$path_release/win/$SDK_NAME'_source.exe'

rm -rf $path_release/*

mkdir -p $path_release

go build -ldflags "-s -w -X main.version=${VERSION} -X main.buildtime=`date +%Y%m%d.%H%M%S`" -o $file_path_mac
upx -o $file_path_min_mac $file_path_mac

# env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X ${project_name}/command.version=${VERSION} -X ${project_name}/command.buildtime=`date +%Y%m%d.%H%M%S`" -o $file_path_linux
# upx -o $file_path_min_linux $file_path_linux

# env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X ${project_name}/command.version=${VERSION} -X ${project_name}/command.buildtime=`date +%Y%m%d.%H%M%S`" -o $file_path_win
# upx -o $file_path_min_win $file_path_win