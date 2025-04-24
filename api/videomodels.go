package api

// todo 添加参数必填限制
type generateVideo struct {
	VideoSubject        string  `json:"video_subject"`
	VideoScript         string  `json:"video_script"`
	VideoTerms          string  `json:"video_terms"`
	VideoAspect         string  `json:"video_aspect"`
	VideoConCatMode     string  `json:"video_con_cat_mode"`
	VideoTransitionMode string  `json:"video_transition_mode"`
	VideoClipDuration   int     `json:"video_clip_duration"`
	VideoCount          int     `json:"video_count"`
	VideoSource         string  `json:"video_source"`
	VideoMaterals       string  `json:"video_materals"`
	VideoLanguage       string  `json:"video_language"`
	VoiceName           string  `json:"voice_name"`
	VoiceVolume         float32 `json:"voice_volume"`
	VoiceRate           float32 `json:"voice_rate"`
	BgmType             string  `json:"bgm_type"`
	BgmFile             string  `json:"bgm_file"`
	BgmVolume           float32 `json:"bgm_volume"`
	SubtitleEnabled     bool    `json:"subtitle_enabled"`
	SubtitlePosition    string  `json:"subtitle_position"`
	CustomPosition      float32 `json:"custom_position"`
	FontName            string  `json:"font_name"`
	TextForeColor       string  `json:"text_fore_color"`
	TextBackgroundColor bool    `json:"text_background_color"`
	FontSize            int     `json:"font_size"`
	StrokeColor         string  `json:"stroke_color"`
	StrokeWidth         float32 `json:"stroke_color"`
	NThreads            int     `json:"n_threads"`
	ParagraphNumber     int     `json:"paragraph_number"`
	FileName            string  `json:"file_name"`
}

type getVideosByUidReq struct {
	Page int `json:"page"`
	Num  int `json:"num"`
}

type getVideosByUidResp struct {
	Videos []Videos `json:"videos"`
	Count  int      `json:"count"`
}

type Videos struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Duration int    `json:"duration"`
}
