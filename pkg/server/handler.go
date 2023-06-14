package server

import (
	"database/sql"
	"encoding/json"
	"kubeStone/m/v2/pkg/config"
	"kubeStone/m/v2/pkg/database"
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
