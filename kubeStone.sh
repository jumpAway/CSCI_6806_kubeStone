#!/bin/bash

if [ "$#" != 1 ]; then
    echo "Usage: ./kubeStone.sh install|start|stop"
    exit -1
fi

install_service() {
    echo "Installing service..."
    echo "Checking Golang environment..."
    if which go >/dev/null 2>&1; then
        go build -o kubeStone >/dev/null 2>&1
        if [ ! -f "kubeStone" ]; then
          echo "Failed to generate binary file"
          exit -1
        fi
        go get golang.org/x/crypto/ssh
        go get github.com/go-sql-driver/mysql
    else
        echo "Please setup Golang First"
        exit -1
    fi
    yum install -y httpd >/dev/null 2>&1
    cp index.html /var/www/html/
    echo "Installation successful"
    echo "Note: you should install mysql manually at this version"
}

start_service() {
    echo "Starting back-end service..."
    ./kubeStone &
    if [[ $(netstat -tlunp | grep 8888 | wc -l) == 0 ]]; then
        echo "Failed to start back-end service"
        exit -1
    fi

    systemctl enable httpd
    systemctl start httpd
    if [[ $? != 0 || $(netstat -tlunp | grep 80 | wc -) == 0 ]]; then
        echo "Failed to start front-end service"
        exit -1
    fi
    echo "Back-end is running at port 8888, Front-end is running at 80. You can access from http://127.0.0.1:80"
}

stop_service() {
    echo "Stopping service..."
    systemctl stop httpd
    fuser -k 8888/tcp
}

case $1 in
    install)
        install_service
        ;;
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    *)
        echo "Invalid command. Usage: ./kubeStone.sh install|start|stop"
        exit 1
        ;;
esac