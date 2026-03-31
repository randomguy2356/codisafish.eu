const login_screen = document.getElementById("login_screen")

const login_form = document.getElementById("login_form")

// const login_btn = document.getElementById("login")

const username_inpt = document.getElementById("username")
const password_inpt = document.getElementById("password")

const error_msg = document.getElementById("error_msg")
const msg = document.getElementById("msg")

login_form.addEventListener("submit", log_in)


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

  if (!response.ok){
    error_response = (body && typeof body.error === "string" && body.error.trim()) || "no error response"
    error_msg.textContent = `error ${response.status}: ${error_response}`
    return
  }
  
  let body = null;
  try {
    body = await response.json()
  }catch {
    console.error("failed to parse json")
    login_screen.style.display = "none"
    return
  }
  
  if (body.error){
    msg.textContent = `${body.error.trim()}`
    return
  }
  login_screen.style.display = "none"
}