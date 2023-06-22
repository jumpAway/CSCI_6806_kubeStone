# CSCI_6806_kubeStone
CSCI_6806_Guohao Li &amp; Hongyang Ruan &amp; Rohan Juneja
A platform that allows users to quickly deploy and operate Kubernetes graphi-
cally through a browser without access to a server and without using the com-
mand line. The platform also provides both online and offline installation modes
to meet the needs of different network environments.

Reading external configuration files and automatic database installation are not implemented at the mid-term stage yet.
Currently, only the database for servers is set up.
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

Please note that the directory containing "kubeStone" is a compiled binary file that can be run directly. 
However, to use this binary directly, you need to make sure that the system architecture is Arm64 and that mysql is running on the same server as kubeStone, and that the configuration parameters refer to config.json.
The simple front-end code is contained in the index.html file and is currently only for testing sending requests.



The front-end page shows the three parts of the functionality that have been completed so far.
1. Server Information：
   Used to display the servers that have been managed by the platform and their information.
When you click on "Search Server" in the browser, the browser sends a request to the back-end service. The back-end service receives the request and routes it to the appropriate handler. The handler retrieves all the server information from the Server table in the database and converts it into json format to send back to the browser. At this stage the browser will print out the json content directly and no subsequent processing is done at this time.
<img width="667" alt="image" src="https://github.com/jumpAway/CSCI_6806_kubeStone/assets/134755433/7cd4ee6d-732e-4f4a-86d4-fd8d2ae2bab9">

2. 






