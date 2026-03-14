import {connect} from "./common_scripts/sse_connect.js"

const year = document.getElementById("year")

const login_btn = document.getElementById("login")
const logout_btn = document.getElementById("logout")
const register_btn = document.getElementById("register")

const username = document.getElementById("username")

const notification_ul = document.getElementById("notifications")

load()

function load(){  
  year.textContent = new Date().getFullYear()
  
  check_account()
}

function showNotification(user) {
	const element = document.createElement("li")
  element.className = "notification"
	element.textContent = `you've been invited by user ${user}`
	
	const buttons = document.createElement("div")
	buttons.className = "buttons"
	
	const accept_btn = document.createElement("button")
	accept_btn.textContent = "✔"
	accept_btn.className = "accept_btn"
	const decline_btn = document.createElement("button")
	decline_btn.textContent = "X"
	decline_btn.className = "decline_btn"
	
	buttons.append(accept_btn)
	buttons.append(decline_btn)

	element.append(buttons)



  notification_ul.prepend(element)
  setTimeout(() => element.remove(), 5000)
}

connect(showNotification)

async function check_account(){
	const response = await fetch("/api/auth/userinfo", {
		method: "GET"
	})
	
	let body = null;
	try {
		body = await response.json()
	} catch {
		console.log(`error: bad response`)
	}
	
	if (!response.ok){
		error_response = (body && typeof body.error === 'string' && body.error.trim()) || `no error response`
		console.log(`eror ${response.status}: ${error_response}`)
		return
	}
	
	var exists = body.exists
	
	if (exists){
		login_btn.style.display = 'none'
		logout_btn.style.display = 'block'
		register_btn.style.display = 'none'

	  username.textContent = body.username
	}else{
		login_btn.style.display = 'block'
		logout_btn.style.display = 'none'
		register_btn.style.display = 'block'
    username.textContent = ""
	}
}