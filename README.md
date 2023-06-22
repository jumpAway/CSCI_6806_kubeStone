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
1. Server Information
   Used to display the servers that have been managed by the platform and their information.
When you click on "Search Server" in the browser, the browser sends a request to the back-end service. The back-end service receives the request and routes it to the appropriate handler. The handler retrieves all the server information from the Server table in the database and converts it into json format to send back to the browser. At this stage the browser will print out the json content directly and no subsequent processing is done at this time.
<img width="667" alt="image" src="https://github.com/jumpAway/CSCI_6806_kubeStone/assets/134755433/7cd4ee6d-732e-4f4a-86d4-fd8d2ae2bab9">

2. Add Server
   This panel is used to test and add new servers to the database. The user enters the new server information in the form, after which he can click Test or Submit. The test is to verify the server's connectivity via ssh only, without doing any subsequent operations. Submit is for the backend to verify the server connectivity first, if it passes then add the server information to the database and return the result to the browser.
<img width="666" alt="image" src="https://github.com/jumpAway/CSCI_6806_kubeStone/assets/134755433/f7f48bbe-519d-46a9-87a3-03c13b2342a4">

3. Cluster Setup
   After the servers are joined, select the managed servers to create a Kubernetes cluster. The user needs to add the role (master/node) first and then select the server that is already in the database in the form. For master, you need to enter the Kubernetes-related initialization configuration, which allows you to customize the service subnet, pod subnet, and kubeproxy proxy mode (iptables/ipvs). After that, click Create and the backend service will automatically issue the deployment task to the target server.
<img width="790" alt="image" src="https://github.com/jumpAway/CSCI_6806_kubeStone/assets/134755433/73a1469b-1371-48ca-9c87-8d875558f2b0">


   
Note: No database tables are currently set up to store cluster information;
    not currently printing cluster deployment results in the browser;
    Cluster version currently only supports v1.26.5;
    the network mode currently only supports calico, and since the calico yaml is not configured to automatically modify it for now, the        pod subnet needs to be specified as the default 192.168.0.0/16. otherwise the cluster network is NotReady;
    Only single-master clusters are currently supported for deployment;







