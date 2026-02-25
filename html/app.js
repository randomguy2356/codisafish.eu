let a_won_btn = document.getElementById("a_won")
let a_won_delta_label = document.getElementById("a_won_delta")
let draw_btn = document.getElementById("draw")
let draw_delta_label = document.getElementById("draw_delta")
let b_won_btn = document.getElementById("b_won")
let b_won_delta_label = document.getElementById("b_won_delta")

let expected_score_lable = document.getElementById("expected_score")

let a_elo_inpt = document.getElementById("a_elo")
let b_elo_inpt = document.getElementById("b_elo")

let list = document.getElementById("transactions")
var last_index = 0

let error_msg = document.getElementById("error_msg")

a_won_btn.addEventListener("click", a_won)
draw_btn.addEventListener("click", draw)
b_won_btn.addEventListener("click", b_won)

a_elo_inpt.addEventListener("input", calc_deltas)
b_elo_inpt.addEventListener("input", calc_deltas)

calc_deltas()
do_list()

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
	
	calc_deltas()
	do_list()
}

function calc_deltas(){
	a_elo_raw = Number.isFinite(Number(a_elo_inpt.value))? Number(a_elo_inpt.value) : 0
	b_elo_raw = Number.isFinite(Number(b_elo_inpt.value))? Number(b_elo_inpt.value) : 0
	a_elo = Math.max(a_elo_raw, 10)
	b_elo = Math.max(b_elo_raw , 10)
	
	k = 32

	expected_score_a = 1 / (1 + (10 ** ((b_elo - a_elo) / 400)))
	
	expected_score_percent = ((Math.round(expected_score_a * 10000)) / 100.0)
	expected_score_lable.textContent = "" + Number(expected_score_percent.toFixed(2)).toString() + " / " + Number((100 - expected_score_percent).toFixed(2)).toString()
	
	a_won_a_delta = k * (1 - expected_score_a)
	draw_a_delta = k * (0.5 - expected_score_a)
	b_won_a_delta = k * (0 - expected_score_a)

	a_delta_min = -(a_elo - 10)
	a_delta_max = b_elo - 10

	a_won_a_delta = Math.min(Math.max(a_won_a_delta, a_delta_min), a_delta_max)
	draw_a_delta = Math.min(Math.max(draw_a_delta, a_delta_min), a_delta_max)
	b_won_a_delta = Math.min(Math.max(b_won_a_delta, a_delta_min), a_delta_max)
	
	
	
	a_won_delta_label.textContent = ("A: " + -(a_elo_raw - (a_elo + a_won_a_delta)).toFixed(2) + "\n B: " + -(b_elo_raw - (b_elo - a_won_a_delta)).toFixed(2))
	draw_delta_label.textContent = ("A: " + -(a_elo_raw - (a_elo + draw_a_delta)).toFixed(2) + "\n B: " + -(b_elo_raw - (b_elo - draw_a_delta)).toFixed(2))
	b_won_delta_label.textContent = ("A: " + -(a_elo_raw - (a_elo + b_won_a_delta)).toFixed(2) + "\n B: " + -(b_elo_raw - (b_elo - b_won_a_delta)).toFixed(2))
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

async function do_list(){
	const response = await fetch("api/transactions", {
		method: "POST",
		headers: {
			"Content-Type": "application/json"
		},
		body: JSON.stringify({last_index: last_index, limit: 20})
	})

	let body = null
	try {
		body = await response.json()
	} catch{
		error_msg.textContent = `error: bad response`
	}

	if (!response.ok){
		error_response = (body && typeof body.error === 'string' && body.error.trim()) || 'no error response'
		error_msg.textContent = `error ${response.status}: ${error_response}`
		return
	}
	
	list.innerHTML = "";
	
	body.entries.forEach(element => {
		const li = document.createElement("li");
		li.textContent = element
		list.appendChild(li);
	});
}