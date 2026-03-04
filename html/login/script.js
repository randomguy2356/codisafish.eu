const login_btn = document.getElementById("login")

const username_inpt = document.getElementById("username")
const password_inpt = document.getElementById("password")

const error_msg = document.getElementById("error_msg")
const msg = document.getElementById("msg")

login_btn.addEventListener("click", log_in)

check_user()

async function check_user() {
  const response = await fetch("/api/auth/userinfo", {
    method: "GET",
  })
  
  let body = null
  try {
    body = await response.json()
  } catch {
    console.log("bad response")
  }
  
  if (body.exists) {
    window.location.href = "/"
  }
}

async function log_in(event){
  event.preventDefault()
  
  error_msg.textContent = ""
  msg.textContent = ""
  
  const response = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
  		"Content-Type": "application/json"
    },
    body: JSON.stringify({username: username_inpt.value, password: password_inpt.value})
  })
  
  let body = null;
  try {
    body = await response.json()
  }catch {
    if (!response.ok){
      error_msg.textContent = `error ${response.status}`
      return
    }
    window.location.href = "/"
  }
  
  if (!response.ok){
    error_response = (body && typeof body.error === "string" && body.error.trim()) || "no error response"
    error_msg.textContent = `error ${response.status}: ${error_response}`
    return
  }
  
  if (body.error){
    msg.textContent = `${body.error.trim()}`
    return
  }

  window.location.href = "/"
}