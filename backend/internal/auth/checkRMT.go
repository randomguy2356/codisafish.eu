package auth

//		import (
//			"database/sql"
//			"net/http"
//		)
//
//		func (handler *checkRMTHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//			if request.Method != http.MethodGet {
//				http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
//				return
//			}
//			rmt, err := request.Cookie("rmt")
//			if err != nil {
//				if err == http.ErrNoCookie {
//					writer.WriteHeader(http.StatusOK)
//					return
//				}
//				http.Error(writer, "cookies bad", http.StatusBadRequest)
//			}
//			username, error := ValidateRMT(rmt.Value, handler.DB, request.Context())
//			if error != nil {
//				http.Error(writer, "internal server error", http.StatusInternalServerError)
//				return
//			}
//			if username == nil {
//				writer.WriteHeader(http.StatusUnauthorized)
//				return
//			}
//			CreateSession(*username, writer)
//			writer.WriteHeader(http.StatusOK)
//		}

//		func CheckRMTCookie(db *sql.DB, writer http.ResponseWriter, request *http.Request) (string, error) {
//			rmt, err := request.Cookie("rmt")
//			if err != http.ErrNoCookie {
//				if err != nil {
//					return "", err
//				}
//				username, err := ValidateRMT(rmt.Value, db, request.Context())
//				if err != nil {
//					return "", err
//				}
//				if username != nil {
//					sid := CreateSession(*username, writer)
//					return sid, nil
//				}
//			}
//			return "", nil
//		}
