package api

type AuthBlock struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

type AuthResponse struct {
	Data struct {
		AuthTelegramInitData AuthBlock `json:"authTelegramInitData"`
	} `json:"data"`
}

type MeBlock struct {
	FirstName        string  `json:"firstName"`
	Balance          float64 `json:"balance"`
	StarsBalance     float64 `json:"starsBalance"`
	NextFreeSpin     string  `json:"nextFreeSpin"`
	NextCaseFreeSpin string  `json:"nextCaseFreeSpin"`
}

type MeResponse struct {
	Data struct {
		Me MeBlock `json:"me"`
	} `json:"data"`
}

type PrizeBlock struct {
	Name string `json:"name"`
}

type SpinInner struct {
	Success     bool       `json:"success"`
	Prize       PrizeBlock `json:"prize"`
	StoryReward *int       `json:"storyReward"`
}

type SpinResponse struct {
	Data struct {
		StartRouletteSpin SpinInner `json:"startRouletteSpin"`
	} `json:"data"`
}

type CaseInner struct {
	Success bool       `json:"success"`
	Prize   PrizeBlock `json:"prize"`
}

type OpenCaseResponse struct {
	Data struct {
		OpenCase CaseInner `json:"openCase"`
	} `json:"data"`
}
