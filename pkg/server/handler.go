package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"kubeStone/m/v2/pkg/config"
	"kubeStone/m/v2/pkg/database"
	"kubeStone/m/v2/pkg/host"
	"log"
	"net/http"
)

func SearchSer(writer http.ResponseWriter, _ *http.Request) {
	var db *sql.DB
	db, err := database.InitDB(cfg)
	if err != nil {
		http.Error(writer, "Init Database ERROR", http.StatusInternalServerError)
		return
	}
	serRow, err := db.Query("SELECT * FROM server")
	if err != nil {
		log.Println(err)
		http.Error(writer, "Query servers from DB error", http.StatusInternalServerError)
		return
	}
	servers := make([]config.Server, 0)
	for serRow.Next() {
		var server config.Server
		var id int
		if err := serRow.Scan(&id, &server.Hostname, &server.IP, &server.Port, &server.Username, &server.Password); err != nil {
			http.Error(writer, "Show servers from DB error", http.StatusInternalServerError)
			return
		}
		servers = append(servers, server)
	}
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(servers); err != nil {
		http.Error(writer, "Response error", http.StatusInternalServerError)
		return
	}
}

func TestSer(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	var server config.Server
	if err := json.Unmarshal(body, &server); err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	err := host.ConnectSer(server)
	if err != nil {
		http.Error(writer, "Cannot access the server", http.StatusInternalServerError)
		return
	}
}

func AddSer(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	var server config.Server
	if err := json.Unmarshal(body, &server); err != nil {
		http.Error(writer, "Failed to parse Add server request body", http.StatusBadRequest)
		return
	}
	err := host.ConnectSer(server)
	if err != nil {
		http.Error(writer, "Cannot access the server", http.StatusInternalServerError)
		return
	}
	var db *sql.DB
	db, err = database.InitDB(cfg)
	if err != nil {
		http.Error(writer, "Init Database ERROR", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO server (name, ip, port, user, password) VALUES (?,?,?,?,?)", server.Hostname, server.IP, server.Port, server.Username, server.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("SELECT * FROM server WHERE ip = ?", server.IP)
	if err != nil {
		http.Error(writer, "Add Server not Success", http.StatusInternalServerError)
		return
	}
}
