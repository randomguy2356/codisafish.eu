const username_inpt = document.getElementById("username")
const username_err = document.getElementById("username_err")
const email_inpt = document.getElementById("email")
const email_err = document.getElementById("email_err")
const password_inpt = document.getElementById("password")
const password_err = document.getElementById("password_err")

const register_btn = document.getElementById("register")

const error_msg = document.getElementById("error_msg")

register_btn.addEventListener("click", register)
username_inpt.addEventListener("change", check_username)
email_inpt.addEventListener("change", check_email)

async function check_username(event){
  event.preventDefault()
  
  if (username_inpt.value == "") {
    return
  }

  var params = new URLSearchParams({
    username: username_inpt.value
  })

  const response = await fetch(`/api/auth/exists?${params}`, {
    method: "GET"
  })
  if (!response.ok){
    error_msg.textContent = `error ${response.status} while checking username`
    return
  }

  let body = null
  try {
    body = await response.json()
  } catch {
    error_msg.textContent = "invalid return value"
  }
  
  if(body?.username){
    username_err.textContent = "username already exists"
  } else if (body?.username === false) {
    username_err.textContent = ""
  }
}

async function check_email(event){
  event.preventDefault()
  
  if (email_inpt.value == "") {
    return
  }

  var params = new URLSearchParams({
    email: email_inpt.value
  })

  const response = await fetch(`/api/auth/exists?${params}`, {
    method: "GET"
  })
  if (!response.ok){
    error_msg.textContent = `error ${response.status} while checking email`
    return
  }

  let body = null
  try {
    body = await response.json()
  } catch {
    error_msg.textContent = "invalid return value"
  }
  
  if(body?.email){
    email_err.textContent = "email already exists"
  } else if (body?.email === false) {
    email_err.textContent = ""
  } 
}

async function register(event){
  event.preventDefault()

  const response = await fetch("/api/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({username: username_inpt.value, email: email_inpt.value, password: password_inpt.value})
  })
  
  let message = await response.text()
  
  if (!response.ok) {
    error_msg.textContent = `error: ${response.status}`
    if (message != ""){
      error_msg.textContent = `${message}`
    }
    return
  }
  
  window.location.href = "/"
}