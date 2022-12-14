package model

type TokenDetails struct {
	AccessToken string
	AccessUUID  string
	AtExpires   int64
}

type JwtPayload struct {
	UserID     string
	Username   string
	Email      string
	AccessUUID string
}
