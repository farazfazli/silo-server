"use strict";
// TODO: Token refreshing
// TODO: Switch to sessionStorage for more security

/**
    There is no React, no Angular, no Bootstrap, no jQuery,
    nothing. Just standard ES5 syntax, I didn't want to use
    any libraries for this project, and kept simplicity with
    bare code as a key component. It's a single page app,
    written without using any libraries.
**/

// Check if token exists
if (localStorage.getItem("silo-authorization")) {
    sendToken(localStorage.getItem("silo-authorization"));
} else {
    document.getElementById("login").style.display = "inherit";
}

// Get DOM elements
var logoutLink = document.getElementById("logout");
var switchLink = document.getElementById("switch");
var submitButton = document.getElementById("submit");
var uploadButton = document.getElementById("apk");

// Listen for "Logout"
addClickEvent(logoutLink, function() {
logout(localStorage.getItem("silo-authorization"));
});

// Listen for "Login" -> "Register"
addClickEvent(switchLink, function() {
switchPage(switchLink, submitButton);
});

// Listens for change to file uploader
function onUploadChange() {
    var filename = document.getElementById('apk').files[0].name;
    if (filename.substr(filename.length - 3).toUpperCase() === "APK") {
        document.getElementById('upload').value = "⇪ Upload " + document.getElementById('apk').files[0].name;
    } else {
        alert("You can't upload that, it's not an APK.");
    }
}

// Submits upload form if changed
function submitUploadFormWithData() {
    // Because I didn't want to add Font Awesome for 1 icon
    // http://www.alanwood.net/unicode/arrows.html
    if (document.getElementById("upload").value !== "Choose APK") {
        sendFile();
    } else {
        document.getElementById('apk').click()
    }
}

function getPackageName() {
    return document.getElementById('package').value;
}

function getPackageHash() {
    return document.getElementById('hash').value;
}

// Sends a file up to the server
function sendFile() {
    if (getPackageName() === "" || getPackageHash() === "" || getPackageName().indexOf('.') === -1) {
        setRed();
        alert("Please enter a valid package name and package hash.");
        return;
    }
    var request = new XMLHttpRequest();
    request.addEventListener('progress', updateProgress);
    request.addEventListener('load', onSubmitLoad());
    request.onreadystatechange = function() {
        if (request.readyState == XMLHttpRequest.DONE) {
            if (request.status == 200) {
                document.getElementById("package").value = null;
                document.getElementById("hash").value = null;
                document.getElementById("apk").value = null;
                document.getElementById("upload").value = "Choose APK";
                alert("Successfully uploaded"); // TODO: change
            }
        }
    };
    var formData = new FormData();
    request.open("POST", "/upload");
    request.setRequestHeader("Authorization", localStorage.getItem("silo-authorization"));
    formData.append("name", getPackageName());
    formData.append("hash", getPackageHash());
    formData.append("apk", document.getElementById('apk').files[0]);
    request.send(formData);
}

// Submit Login form with data
function submitLoginFormWithData() {
    submitForm(getEmail(), getPassword(), submitButton.value == "Register");
}

// From StackOverflow, increased browser support for onclick
function addClickEvent(el, fn) {
    return ((el.attachEvent) ? el.attachEvent('onclick', fn) : el.addEventListener("click", fn, false));
}

// Tracking progress isn't needed for now, maybe later
function updateProgress(progress) {
    // if (progress.loaded > 0) {
    //     console.log((progress.loaded / progress.total) * 100 + "%")
    // } else {
    //     console.log("Error updating progress");
    // }
}

function onSubmitLoad() {
    // Called when request is sent
}

function getEmail() {
    return document.getElementById('email').value;
}

function getPassword() {
    return document.getElementById('password').value;
}

// Switches between "Login" and "Register" pages
function switchPage(switchLink, submitButton) {
    var login = document.getElementById("submit").value == "Login" ? true : false;
    switch (login) {
    case true:
        setRegisterAppearance(switchLink, submitButton);
        break;
    case false:
        setLoginAppearance(switchLink, submitButton);
        break;
    }
}

function setRegisterAppearance(switchLink, submitButton) {
    switchLink.innerHTML = "Login to your account";
    switchLink.style = "color:purple;";
    submitButton.value = "Register";
}

function setLoginAppearance(switchLink, submitButton) {
    switchLink.innerHTML = "Create an account";
    switchLink.style = "color:blue;";
    submitButton.value = "Login";
}

function errorRegistering() {
    alert("Error registering, please try again.");
}

function errorLoggingIn() {
    alert("Error logging in, please check your info.");
}

function setRed() {
    document.getElementById('email').style = "border-color:red;";
    document.getElementById('password').style = "border-color:red;";
}

function displayDashboard(text) {
    document.getElementById('dashboard').style = "display:visible";
    document.getElementById('response').innerHTML = text;
}

function getAuthHeader(request) {
    var authHeader = request.getResponseHeader("Authorization");
    if (authHeader !== null) {
        console.log("Received header:", authHeader);
        localStorage.setItem("silo-authorization", authHeader);
        sendToken(authHeader);
    } else {
        console.log("No headers returned");
    }
}

// Deletes token from local storage, sends logout request
function logout(token) {
    localStorage.removeItem("silo-authorization");
    var request = new XMLHttpRequest();
    request.addEventListener('progress', updateProgress);
    request.addEventListener('load', onSubmitLoad());
    request.onreadystatechange = function() {
        if (request.readyState == XMLHttpRequest.DONE) {
            if (request.status == 200) {
                document.getElementById('login').style.display = "inherit";
                document.getElementById('dashboard').style.display = "none";
            }
        }
    };
    request.open("POST", "/logout");
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.setRequestHeader("Authorization", token);
    request.send();
}

// Sends token to server, called if token exists in browser, and after
// getting the authorization header.
function sendToken(token) {
    var request = new XMLHttpRequest();
    request.addEventListener('progress', updateProgress);
    request.addEventListener('load', onSubmitLoad());
    request.onreadystatechange = function() {
        if (request.readyState == XMLHttpRequest.DONE) {
            if (request.status == 200) {
                if (document.getElementById('login').style.display = "none") {
                    var response = request.responseText;
                    displayDashboard(response);
                }
            }
        }
    }
    request.open("POST", "/login");
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    request.setRequestHeader("Authorization", token);
    request.send();
}

// Submits the form, checking for errors and headers.
function submitForm(email, password, register) {
    if (email === "" || password === "") {
        setRed();
        alert("Please enter a longer email and password.");
        return;
    }
    var request = new XMLHttpRequest();
    request.addEventListener('progress', updateProgress);
    request.addEventListener('load', onSubmitLoad());
    request.onreadystatechange = function() {
        if (this.readyState == XMLHttpRequest.HEADERS_RECEIVED) {
            getAuthHeader(this);
        }
        if (request.readyState == XMLHttpRequest.DONE) {
            if (request.status !== 200) {
                setRed();
                if (register) {
                    errorRegistering();
                } else {
                    errorLoggingIn();
                }
            }
        }
    }
    request.open("POST", register ? "/register" : "/login");
    request.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    var authorizationHeader = localStorage.getItem("silo-authorization");
    request.setRequestHeader("Authorization", authorizationHeader);
    request.send("email=" + encodeURIComponent(email) + "&password=" + encodeURIComponent(password) + "&register=" + register);
}

// http://stackoverflow.com/a/7317311
// We aren't ever submitting the form,
// because this is a single page app
window.addEventListener("beforeunload", function(e) {
    if (document.getElementById("upload").value !== "Choose APK") {
        var confirmationMessage = "You have uploaded an APK but haven\'t published. Sure you want to leave?";
        (e || window.event).returnValue = confirmationMessage; //Gecko + IE
        return confirmationMessage; //Gecko + Webkit, Safari, Chrome etc.
    }
});