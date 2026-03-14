import {connect} from "../common_scripts/sse_connect.js"

const target_inpt = document.getElementById("target")
const ping_btn = document.getElementById("ping")

ping_btn.addEventListener("click", submit)

check_logged_in()

async function check_logged_in(){
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
		console.log(`eror ${response.status}: ${await response.text()}`)
		return
	}
	
  if (!body.exists) {
    window.location.href = "/"
  }
}

async function submit(event) {
  event.preventDefault()

  const target = encodeURIComponent(target_inpt.value)
  console.log(target)

  const response = await fetch(`/api/ping?target=${target}`, {
    method: "GET"
  })

  if(!response.ok){
    console.log("nay :(")
  } else {
    console.log("yay! :)")
  }
}

function showNotification(message) {
  const element = document.createElement("div")
  element.className = "notification"
  element.textContent = message
  document.body.appendChild(element)
  setTimeout(() => element.remove(), 10000)
}
connect(showNotification)