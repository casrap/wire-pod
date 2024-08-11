package vars

import (
	"encoding/json"
	"os"

	"github.com/kercre123/wire-pod/chipper/pkg/logger"
)

// a way to create a JSON configuration for wire-pod, rather than the use of env vars

var ApiConfigPath = "./apiConfig.json"

var APIConfig apiConfig

type apiConfig struct {
	Weather struct {
		Enable   bool   `json:"enable"`
		Provider string `json:"provider"`
		Key      string `json:"key"`
		Unit     string `json:"unit"`
	} `json:"weather"`
	Knowledge struct {
		Enable                 bool   `json:"enable"`
		Provider               string `json:"provider"`
		Key                    string `json:"key"`
		ID                     string `json:"id"`
		Model                  string `json:"model"`
		IntentGraph            bool   `json:"intentgraph"`
		RobotName              string `json:"robotName"`
		OpenAIPrompt           string `json:"openai_prompt"`
		OpenAIVoice            string `json:"openai_voice"`
		OpenAIVoiceWithEnglish bool   `json:"openai_voice_with_english"`
		SaveChat               bool   `json:"save_chat"`
		CommandsEnable         bool   `json:"commands_enable"`
		Endpoint               string `json:"endpoint"`
	} `json:"knowledge"`
	STT struct {
		Service  string `json:"provider"`
		Language string `json:"language"`
	} `json:"STT"`
	Server struct {
		// false for ip, true for escape pod
		EPConfig bool   `json:"epconfig"`
		Port     string `json:"port"`
	} `json:"server"`
	PastInitialSetup bool `json:"pastinitialsetup"`
}

func WriteConfigToDisk() {
	logger.Println("Configuration changed, writing to disk")
	writeBytes, _ := json.Marshal(APIConfig)
	os.WriteFile(ApiConfigPath, writeBytes, 0644)
}

func WriteSTT() {
	// was not part of the original code, so this is its own function
	// launched if stt not found in config
	APIConfig.STT.Service = os.Getenv("STT_SERVICE")
	if os.Getenv("STT_SERVICE") == "vosk" || os.Getenv("STT_SERVICE") == "whisper.cpp" {
		APIConfig.STT.Language = os.Getenv("STT_LANGUAGE")
	}
}

func ReadConfig() {
	if _, err := os.Stat(ApiConfigPath); err == nil {
		// read config
		configBytes, err := os.ReadFile(ApiConfigPath)
		if err != nil {
			APIConfig.Knowledge.Enable = false
			APIConfig.Weather.Enable = false
			logger.Println("Failed to read API config file")
			logger.Println(err)
			return
		}
		err = json.Unmarshal(configBytes, &APIConfig)
		if err != nil {
			APIConfig.Knowledge.Enable = false
			APIConfig.Weather.Enable = false
			logger.Println("Failed to unmarshal API config JSON")
			logger.Println(err)
			return
		}
		// stt service is the only thing controlled by shell
		if APIConfig.STT.Service != os.Getenv("STT_SERVICE") {
			WriteSTT()
		}

		if APIConfig.Knowledge.Model == "meta-llama/Llama-2-70b-chat-hf" {
			logger.Println("Setting Together model to Llama3")
			APIConfig.Knowledge.Model = "meta-llama/Llama-3-70b-chat-hf"
		}

		writeBytes, _ := json.Marshal(APIConfig)
		os.WriteFile(ApiConfigPath, writeBytes, 0644)
		logger.Println("API config successfully read")
	}
}
