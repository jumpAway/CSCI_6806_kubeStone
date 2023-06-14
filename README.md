# CSCI_6806_kubeStone
CSCI_6806_Guohao Li &amp; Hongyang Ruan &amp; Rohan Juneja
A platform that allows users to quickly deploy and operate Kubernetes graphi-
cally through a browser without access to a server and without using the com-
mand line. The platform also provides both online and offline installation modes
to meet the needs of different network environments.

Reading external configuration files and automatic database installation are not implemented at the mid-term stage yet.
Currently, only the database for the server is set up.
Before running the service you need to configure the mysql database manually and write the database information to config.json.
After installing mysql, create the kubeStone database and server table.

    mysql>
        -> CREATE DATABASE kubeStone;
        -> CREATE TABLE server (
        -> id INT AUTO_INCREMENT PRIMARY KEY,
        -> name VARCHAR(255) NOT NULL,
        -> ip VARCHAR(15) NOT NULL UNIQUE,
        -> port VARCHAR(15) NOT NULL,
        -> user VARCHAR(255) NOT NULL,
        -> password VARCHAR(255) NOT NULL               //Temporarily stored in plaintext for testing purposes only.
        -> );

install golang package

        go get golang.org/x/crypto/ssh
        go get github.com/go-sql-driver/mysql

build & run

    cd kubeStone/
    go build .
    ./kubeStone

then can verify the service port

    netstat -tlunp | grep 8888

The simple front-end code is contained in the index.html file and is currently only for testing sending requests.








