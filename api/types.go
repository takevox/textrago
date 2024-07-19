package api

/*
	TexTra 言語判定APIのレスポンス構造体
*/
type LangDetectResponse struct {
	ResultSet struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Request struct {
			Url  string `json:"url"`
			Text string `json:"text"`
		} `json:"request"`
		Result struct {
			LangDetect []struct {
				Lang string  `json:"lang"`
				Rate float32 `json:"rate"`
			} `json:"langdetect"`
		} `json:"result"`
	} `json:"resultset"`
}

/*
	TexTra 自動翻訳リクエストAPIのレスポンス構造体
*/
type TranslationResponse struct {
	ResultSet struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Request struct {
			Url         string `json:"url"`
			Text        string `json:"text"`
			Split       int    `json:"split"`
			History     string `json:"history"`
			Prompt      string `json:"prompt"`
			XML         string `json:"xml"`
			TermId      string `json:"term_id"`
			BilingualId string `json:"bilingual_id"`
			LogUse      int    `json:"log_use"`
			EditorUse   int    `json:"editor_use"`
			Data        string `json:"data"`
		} `json:"request"`
		Result struct {
			Text  string `json:"text"`
			Blank int    `json:"blank"`
		} `json:"result"`
	} `json:"resultset"`
}
