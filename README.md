# CSCI_6806_kubeStone
CSCI_6806_kubeStone is a platform that allows users to quickly deploy and operate Kubernetes graphically through a browser without access to a server and without using the command line.
The platform also provides both online and offline installation modes to meet the needs of different network environments.
## Table of content
1. [Installation](#installation)
2. [Usage](#usage)
3. [Contributing](#contributing)

## Installation
1.
Reading external configuration files and automatic database installation are not implemented at the mid-term stage yet.
Currently, only the database for servers is set up. Before running the service you need to configure the mysql database manually and write the database information to config.json.
After installing mysql, create the kubeStone database and server table.

    mysql>
        CREATE DATABASE kubeStone;
        CREATE TABLE server (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        ip VARCHAR(15) NOT NULL UNIQUE,
        port VARCHAR(15) NOT NULL,
        user VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL 
        );

        CREATE TABLE cluster (
        id INT AUTO_INCREMENT PRIMARY KEY,
        cluster_name VARCHAR(255) NOT NULL UNIQUE,
        version VARCHAR(15) NOT NULL,
        CNI VARCHAR(15) NOT NULL,
        ServiceSubnet VARCHAR(15) NOT NULL,
        PodSubnet VARCHAR(15) NOT NULL,
        ProxyMode VARCHAR(15) NOT NULL,
        master VARCHAR(15) NOT NULL,
        node VARCHAR(255) NOT NULL,
        context VARCHAR(255) NOT NULL UNIQUE
        );
        
        CREATE TABLE gptHistory (
        id INT AUTO_INCREMENT PRIMARY KEY,
        uuid VARCHAR(255) NOT NULL UNIQUE,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        cluster VARCHAR(255) NOT NULL,
        namespace VARCHAR(255) NOT NULL,
        model VARCHAR(255) NOT NULL,
        temperature VARCHAR(255) NOT NULL
        );
        
        CREATE TABLE gptMessage (
        id INT AUTO_INCREMENT PRIMARY KEY,
        history_id INT NOT NULL,
        role VARCHAR(255),
        content TEXT,
        FOREIGN KEY (history_id) REFERENCES gptHistory(id)
        );
2.
Enter the root directory

        cd /CSCI_6806_kubeStone
Execute installation script (Requires preconfigured golang environment https://go.dev/doc/install)

        ./kubeStone.sh install
## Usage
Execute running script

        ./kubeStone.sh start

Access web page

        http://127.0.0.1:80


        
Note: 
1. Please note that the mysql needs to run manually on the same server as kubeStone, and that the configuration parameters refer to config.json.
2. The simple front-end code is contained in the index.html file and is currently only for testing sending requests.
3. No database tables are currently set up to store cluster information.
4. Not currently printing cluster deployment results in the browser;
5. Cluster version currently only supports v1.26.5;
6. The network mode currently only supports calico, and since the calico yaml is not configured to automatically modify it for now, the pod subnet needs to be specified as the default 192.168.0.0/16. otherwise the cluster network will be NotReady;
7. Only single-master clusters are currently supported for deployment;









