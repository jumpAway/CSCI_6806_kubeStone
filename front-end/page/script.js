

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

