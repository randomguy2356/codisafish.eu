let a_won_btn = document.getElementById("a_won")
let a_won_delta_label = document.getElementById("a_won_delta")
let draw_btn = document.getElementById("draw")
let draw_delta_label = document.getElementById("draw_delta")
let b_won_btn = document.getElementById("b_won")
let b_won_delta_label = document.getElementById("b_won_delta")

let a_elo_inpt = document.getElementById("a_elo")
let b_elo_inpt = document.getElementById("b_elo")

let error_msg = document.getElementById("error_msg")

a_won_btn.addEventListener("click", a_won)
draw_btn.addEventListener("click", draw)
b_won_btn.addEventListener("click", b_won)

async function game(K, a_elo, b_elo, score_a){
	const response = await fetch("/api/game", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({k: String(K), a_elo: String(a_elo), b_elo: String(b_elo), score_a: String(score_a)})
	})
	
	let body = null;
	try {
		body = await response.json()
	} catch{
		error_msg.textContent = `error: bad response`
	}

	if (!response.ok){
		error_response = (body && typeof body.error === 'string' && body.error.trim()) || `no error response`
		error_msg.textContent = `eror ${response.status}: ${error_response}`
		return
	}

	const { a_elo: a_elo_new, b_elo: b_elo_new } = body ?? {}
	if (a_elo_new === undefined || b_elo_new === undefined) {
		error_msg.textContent = "a or b undefined"
	}
	
	a_elo_inpt.value = a_elo_new
	b_elo_inpt.value = b_elo_new
	
	calc_deltas(K, a_elo_new, b_elo_new)
}

function calc_deltas(K, a_elo, b_elo){
	expected_score_a = 1 / (1 + (10 ** (b_elo - a_elo)))
	
	a_won_a_delta = k * (1 - expected_score_a)
	draw_a_delta = k * (0.5 - expected_score_a)
	b_won_a_delta = k * (0 - expected_score_a)

	a_delta_min = -(a_elo - 10)
	a_delta_max = b_elo - 10

	a_won_a_delta = Math.min(Math.max(a_won_a_delta, a_delta_min), a_delta_max)
	draw_a_delta = Math.min(Math.max(draw_a_delta, a_delta_min), a_delta_max)
	b_won_a_delta = Math.min(Math.max(b_won_a_delta, a_delta_min), a_delta_max)
	
	a_won_delta_label.textContent = ("a: " + a_won_a_delta + ", b: " + (-a_won_a_delta))
	draw_delta_label.textContent = ("a: " + draw_a_delta + ", b: " + (-draw_a_delta))
	b_won_delta_label.textContent = ("a: " + b_won_a_delta + ", b: " + (-b_won_a_delta))
	
}

function a_won(event){
	event.preventDefault()
	var a_elo = Number(a_elo_inpt.value)
	if (Number.isNaN(a_elo) || a_elo < 10){
		a_elo = 10
		a_elo_inpt.value = 10
	}
	var b_elo = Number(b_elo_inpt.value)
	if (Number.isNaN(b_elo) || b_elo < 10){
		b_elo = 10
		b_elo_inpt.value = 10
	}
	
	var score_a = 1
	var K = 32

	game(K, a_elo, b_elo, score_a)
}

function draw(event){
	event.preventDefault()
	var a_elo = Number(a_elo_inpt.value)
	if (Number.isNaN(a_elo) || a_elo < 10){
		a_elo = 10
		a_elo_inpt.value = 10
	}
	var b_elo = Number(b_elo_inpt.value)
	if (Number.isNaN(b_elo) || b_elo < 10){
		b_elo = 10
		b_elo_inpt.value = 10
	}
	
	var score_a = 0.5
	var K = 32

	game(K, a_elo, b_elo, score_a)

}

function b_won(event){
	event.preventDefault()
	var a_elo = Number(a_elo_inpt.value)
	if (Number.isNaN(a_elo) || a_elo < 10){
		a_elo = 10
		a_elo_inpt.value = 10
	}
	var b_elo = Number(b_elo_inpt.value)
	if (Number.isNaN(b_elo) || b_elo < 10){
		b_elo = 10
		b_elo_inpt.value = 10
	}
	
	var score_a = 0
	var K = 32

	game(K, a_elo, b_elo, score_a)
}