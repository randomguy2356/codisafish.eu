package auth

import "net/http"

func (handler *LogoutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sidcookie, err := request.Cookie("sid")

	if err != nil {
		if err == http.ErrNoCookie {
			sidcookie = &http.Cookie{Value: ""}
		} else {
			http.Error(writer, "cookie error", http.StatusBadRequest)
			return
		}
	}

	InvalidateSID(sidcookie.Value, writer)
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}
