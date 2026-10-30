package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	CharacterName string
	AuthToken     string
}

func NewConfig() Config {
	return Config{
		CharacterName: "Buba",
		AuthToken:     "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6InRpbXVyMzMzOTlAdmsuY29tIiwicGFzc3dvcmRfY2hhbmdlZCI6bnVsbH0.qak98F18TUDBY8CLiK1ntzcD8SRGiGDhsgS916A1N1U",
	}
}

type MoveRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type CharacterInfo struct {
	Data struct {
		Name string `json:"name"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
	} `json:"data"`
}

type ArtifactsClient struct {
	config     Config
	httpClient *http.Client
	baseURL    string
}

func NewArtifactsClient(cfg Config) *ArtifactsClient {
	return &ArtifactsClient{
		config:     cfg,
		httpClient: &http.Client{},
		baseURL:    "https://api.artifactsmmo.com",
	}
}

func (c *ArtifactsClient) doRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AuthToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Status: %s | URL: %s", resp.Status, url)

	return respBody, nil
}

func (c *ArtifactsClient) GetCharacter() (*CharacterInfo, error) {
	url := fmt.Sprintf("%s/my/%s", c.baseURL, c.config.CharacterName)

	body, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var info CharacterInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &info, nil
}

func (c *ArtifactsClient) Move(moveReq MoveRequest) (string, error) {
	url := fmt.Sprintf("%s/my/%s/action/move", c.baseURL, c.config.CharacterName)

	jsonData, err := json.Marshal(moveReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	body, err := c.doRequest("POST", url, jsonData)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type PageData struct {
	Character *CharacterInfo
	Message   string
	Error     string
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Artifacts MMO - Управление персонажем</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            padding: 40px;
            max-width: 500px;
            width: 100%;
        }
        h1 {
            color: #333;
            margin-bottom: 8px;
            font-size: 24px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .info-card {
            background: #f7f8fc;
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 24px;
            border-left: 4px solid #667eea;
        }
        .info-card h3 {
            color: #667eea;
            font-size: 14px;
            margin-bottom: 12px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .info-row {
            display: flex;
            justify-content: space-between;
            padding: 8px 0;
            border-bottom: 1px solid #eaeef5;
        }
        .info-row:last-child { border-bottom: none; }
        .info-row span:first-child { color: #666; }
        .info-row span:last-child { font-weight: 600; color: #333; }
        form { display: flex; flex-direction: column; gap: 16px; }
        label {
            color: #555;
            font-size: 14px;
            font-weight: 500;
            margin-bottom: 4px;
        }
        .inputs {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 12px;
        }
        .input-group { display: flex; flex-direction: column; }
        input[type="number"] {
            padding: 12px 16px;
            border: 2px solid #e0e4ed;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.2s;
        }
        input[type="number"]:focus {
            outline: none;
            border-color: #667eea;
        }
        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 14px;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.1s, box-shadow 0.2s;
            margin-top: 8px;
        }
        button:hover {
            transform: translateY(-1px);
            box-shadow: 0 8px 16px rgba(102, 126, 234, 0.4);
        }
        button:active { transform: translateY(0); }
        .message {
            padding: 12px 16px;
            border-radius: 8px;
            margin-bottom: 20px;
            font-size: 14px;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border-left: 4px solid #28a745;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border-left: 4px solid #dc3545;
        }
        .response {
            margin-top: 20px;
            padding: 16px;
            background: #1e1e2e;
            color: #a6e3a1;
            border-radius: 8px;
            font-family: 'Consolas', monospace;
            font-size: 12px;
            white-space: pre-wrap;
            word-break: break-all;
            max-height: 300px;
            overflow-y: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎮 Artifacts MMO</h1>
        <p class="subtitle">Управление персонажем через веб-интерфейс</p>

        {{if .Message}}
        <div class="message success">{{.Message}}</div>
        {{end}}
        {{if .Error}}
        <div class="message error">{{.Error}}</div>
        {{end}}

        {{if .Character}}
        <div class="info-card">
            <h3>Текущее положение</h3>
            <div class="info-row">
                <span>Имя:</span>
                <span>{{.Character.Data.Name}}</span>
            </div>
            <div class="info-row">
                <span>X:</span>
                <span>{{.Character.Data.X}}</span>
            </div>
            <div class="info-row">
                <span>Y:</span>
                <span>{{.Character.Data.Y}}</span>
            </div>
        </div>
        {{end}}

        <form method="POST">
            <label>Куда идти:</label>
            <div class="inputs">
                <div class="input-group">
                    <label for="x">X</label>
                    <input type="number" id="x" name="x" value="4" required>
                </div>
                <div class="input-group">
                    <label for="y">Y</label>
                    <input type="number" id="y" name="y" value="-1" required>
                </div>
            </div>
            <button type="submit">🚶 Идти</button>
        </form>

        {{if .Message}}
        <div class="response">{{.Message}}</div>
        {{end}}
    </div>
</body>
</html>
`

func main() {
	cfg := NewConfig()
	client := NewArtifactsClient(cfg)

	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}

		charInfo, err := client.GetCharacter()
		if err != nil {
			data.Error = fmt.Sprintf("Не удалось получить данные персонажа: %v", err)
		} else {
			data.Character = charInfo
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				data.Error = fmt.Sprintf("Ошибка формы: %v", err)
			} else {
				x, _ := strconv.Atoi(r.FormValue("x"))
				y, _ := strconv.Atoi(r.FormValue("y"))

				action := MoveRequest{X: x, Y: y}
				response, err := client.Move(action)
				if err != nil {
					data.Error = fmt.Sprintf("Ошибка перемещения: %v", err)
				} else {
					data.Message = fmt.Sprintf("Успешно! X: %d, Y: %d\nОтвет: %s", x, y, response)

					if newInfo, err := client.GetCharacter(); err == nil {
						data.Character = newInfo
					}
				}
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		}
	})

	log.Println("Server started on :8080")
	log.Println("Open http://localhost:8080 in your browser")
	log.Fatal(http.ListenAndServe(":8080", nil))
}