package dto

type GoogleLoginDTO struct {
	IdToken string `json:"idToken" binding:"required"`
}

type GoogleLoginResponseDTO struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}
