const year = document.getElementById("year")

const login_btn = document.getElementById("login")
const logout_btn = document.getElementById("logout")
const register_btn = document.getElementById("register")

const username = document.getElementById("username")

load()

function load(){  
  year.textContent = new Date().getFullYear()
  
  check_account()
}

async function check_account(){
	const response = await fetch("/api/auth/userinfo", {
		method: "GET"
	})
	
	let body = null;
	try {
		body = await response.json()
	} catch {
		error_msg.textContent = `error: bad response`
	}
	
	if (!response.ok){
		error_response = (body && typeof body.error === 'string' && body.error.trim()) || `no error response`
		error_msg.textContent = `eror ${response.status}: ${error_response}`
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