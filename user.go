package tealtech

type BotUser struct {
	UserId    int    `json:"uid"`
	Login     string `json:"login"`
	State     string `json:"state"`
	ExtraData string `json:"extra_data"`
}

const (
	stateNewUser   = "state_new_user"
	stateUserSmth  = "state_smth"
	stateUserSmth2 = "state_smth"
	stateUserSmth3 = "state_smth"
	stateUserSmth4 = "state_smth"
)
