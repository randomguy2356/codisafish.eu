const login_screen = document.getElementById("login_screen");

const new_game_btn = document.getElementById("new_game");
const user_list_btn = document.getElementById("user_list");
const game_history_btn = document.getElementById("game_history");

new_game_btn.addEventListener("click", new_game_prompt);

check_logged_in();

const login_display_style = "flex";

async function check_logged_in() {
	const response = await fetch("/api/auth/userinfo", { method: "GET" });

	if (!response.ok) {
		console.error("not ok response");
		return;
	}

	let body = null;
	try {
		body = await response.json();
	} catch {
		console.error("invalid json lol");
		return;
	}
	if (body.exists === false) {
		console.log("not logged in");
		login_screen.style.display = login_display_style;
		return;
	}
	console.log("logged in");
}

function new_game_prompt() {}
