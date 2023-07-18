//xorEncrypt is a simple XOR encryption function
function xorEncrypt(str) {
    let key = 'this is a secret key.'
    let result = '';
    for (let i = 0; i < str.length; i++) {
        const charCode = str.charCodeAt(i) ^ key.charCodeAt(i % key.length);
        result += String.fromCharCode(charCode);
    }
    return stringToHex(result);
}
//stringToHex convert strings to hexadecimal representation
function stringToHex(str) {
    let hex = '';
    for (let i = 0; i < str.length; i++) {
        const charCode = str.charCodeAt(i).toString(16);
        hex += ('00' + charCode).slice(-2);
    }
    return hex;
}

document.getElementById("serverForm").addEventListener("submit", function(event) {event.preventDefault();});
function InitSer(){
    return {
        hostname: document.getElementById("hostname").value,
        ipaddress: document.getElementById("ipaddress").value,
        port: document.getElementById("port").value,
        username: document.getElementById("username").value,
        password: xorEncrypt(document.getElementById("password").value)
    }
}
function testConnectivity() {
    const formData = InitSer();
    console.log(formData.password);
    fetch("http://kubestonebackend:8888/testServer", {
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
                console.log(formData.password);
            }
        })
        .catch(function(error) {
            console.log(error);
        });
}

function addServer() {
    const formData = InitSer();
    fetch("http://kubestonebackend:8888/addServer", {
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

