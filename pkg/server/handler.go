package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"kubeStone/pkg/GPT"
	"kubeStone/pkg/config"
	"kubeStone/pkg/database"
	"kubeStone/pkg/host"
	"kubeStone/pkg/install"
	"log"
	"net/http"
	"os"
	"strings"
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

func CreateCluster(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	var cluster []config.ClusterInfo
	if err := json.Unmarshal(body, &cluster); err != nil {
		http.Error(writer, "Failed to parse Create cluster request body", http.StatusBadRequest)
		return
	}
	var db *sql.DB
	db, err := database.InitDB(cfg)
	if err != nil {
		http.Error(writer, "Init Database ERROR", http.StatusInternalServerError)
		return
	}
	var serMaster config.Server
	var id int
	err = db.QueryRow("SELECT * FROM server WHERE ip = ?", cluster[0].MasterIp).Scan(&id, &serMaster.Hostname, &serMaster.IP, &serMaster.Port, &serMaster.Username, &serMaster.Password)
	if err != nil {
		http.Error(writer, "Query servers from DB error", http.StatusInternalServerError)
		return
	}

	if err = install.SetMaster(serMaster, cluster[0]); err != nil {
		http.Error(writer, "Failed to set up environment of "+serMaster.IP+err.Error(), http.StatusInternalServerError)
		return
	}

	var seq int
	seq = 1
	for _, node := range cluster {
		if node.NodeIp != "" {
			var serNode config.Server
			err = db.QueryRow("SELECT * FROM server WHERE ip = ?", node.NodeIp).Scan(&id, &serNode.Hostname, &serNode.IP, &serNode.Port, &serNode.Username, &serNode.Password)
			if err != nil {
				http.Error(writer, "Query servers from DB error", http.StatusInternalServerError)
				return
			}
			if err := install.SetNode(serNode, serMaster, seq); err != nil {
				http.Error(writer, "Failed to setup node: "+err.Error(), http.StatusInternalServerError)
				return
			}
			seq++
		}
	}
	// add into database
}

func byGPT(writer http.ResponseWriter, request *http.Request) {
	gptApiKey, err := os.ReadFile("/root/API_KEY")
	if err != nil || len(gptApiKey) == 0 {
		http.Error(writer, "Set API key error", http.StatusInternalServerError)
		return
	}
	apiKey := strings.TrimSpace(string(gptApiKey))
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "receive error"+err.Error(), http.StatusInternalServerError)
		return
	}
	uuid := request.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(writer, "UUID not provided", http.StatusBadRequest)
		return
	}
	GPT.HistoryMutex.Lock()
	message, exist := GPT.HistoryMap[uuid]
	if !exist {
		message = []map[string]string{
			{
				"role":    "system",
				"content": "The namespace is default, give me the yaml directly.",
			},
		}
		GPT.HistoryMap[uuid] = message
	}
	GPT.HistoryMutex.Unlock()
	message = append(message, map[string]string{
		"role":    "user",
		"content": string(body),
	})
	userIdea := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    message,
		"temperature": 0.7,
	}

	userIdeaText, err := json.Marshal(userIdea)
	if err != nil {
		http.Error(writer, "Encode request error"+err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	request, err = http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(userIdeaText))
	if err != nil {
		http.Error(writer, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)
	gptResp, err := client.Do(request)
	if err != nil {
		http.Error(writer, "Error sending request to OpenAI API: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var content config.GPTResponse
	err = json.NewDecoder(gptResp.Body).Decode(&content)
	if err != nil {
		http.Error(writer, "Decode from GPT error"+err.Error(), http.StatusInternalServerError)
		return
	}

	answer := content.Choices[0].Message.Content
	_, _, err = GPT.ExecuteGPT(answer)
	if err != nil {
		http.Error(writer, "Execute GPT error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	message = append(message, map[string]string{
		"role":    "assistant",
		"content": answer,
	})
	_, err = writer.Write([]byte(answer))
	if err != nil {
		http.Error(writer, "write error"+err.Error(), http.StatusInternalServerError)
		return
	}
	GPT.HistoryMutex.Lock()
	GPT.HistoryMap[uuid] = message
	GPT.HistoryMutex.Unlock()
}
