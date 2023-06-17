function searchServer() {
    fetch("http://192.168.102.132:8888/getServer")
        .then(response => response.json())
        .then(data => {
            // Create a table
            let table = document.createElement('table');
            table.className = 'table';

            // Create table header
            let thead = table.createTHead();
            let headerRow = thead.insertRow();
            ['Hostname', 'IP Address', 'Port', 'Username', 'Password'].forEach(header => {
                let th = document.createElement('th');
                th.innerText = header;
                headerRow.appendChild(th);
            });

            // Populate table with data
            let tbody = document.createElement('tbody');
            data.forEach(item => {
                let row = tbody.insertRow();
                let cell1 = row.insertCell();
                cell1.innerText = item.hostname;
                let cell2 = row.insertCell();
                cell2.innerText = item.ipaddress;
                let cell3 = row.insertCell();
                cell3.innerText = item.port;
                let cell4 = row.insertCell();
                cell4.innerText = item.username;
                let cell5 = row.insertCell();
                cell5.innerText = item.password;
            });

            // Append tbody to table
            table.appendChild(tbody);

            // Display the table in responseContainer
            const responseContainer = document.getElementById('responseContainer');
            responseContainer.appendChild(table);
            responseContainer.style.display = 'block';
        })
        .catch(error => {
            console.error("Error: " + error);
        });
}

// Execute the searchServer function as soon as the page loads
searchServer();

document.getElementById("serverForm").addEventListener("submit", function(event) {event.preventDefault();});
function InitSer(){
    return {
        hostname: document.getElementById("hostname").value,
        ipaddress: document.getElementById("ipaddress").value,
        port: document.getElementById("port").value,
        username: document.getElementById("username").value,
        password: document.getElementById("password").value
    }
}
function testConnectivity() {
    const formData = InitSer();
    fetch("http://192.168.102.132:8888/testServer", {
        method: "POST",
        body: JSON.stringify(formData),
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(function(response) {
            if (response.ok) {
                document.getElementById('result').textContent = 'ACCESS SERVER SUCCESS';
            } else {
                document.getElementById('result').textContent = 'ACCESS SERVER NOT SUCCESS';
            }
        })
        .catch(function(error) {
            console.log(error);
        });
}

function addServer() {
    const formData = InitSer();
    fetch("http://192.168.102.132:8888/addServer", {
        method: "POST",
        body: JSON.stringify(formData),
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(function(response) {
            if (response.ok) {
                document.getElementById('result').textContent = 'ADD SERVER SUCCESS';
            } else {
                response.text().then(errorMessage => {
                    document.getElementById('result').textContent = errorMessage;
                });
            }
        })
        .catch(function(error) {
            console.log(error);
        });
}

