package service

const (
	LINE  = "line"
	Kakao = "kakao"
)

var codeMap map[int]string

// ProfileType is struct of user profile info
type ProfileType struct {
	UserID      string
	DisplayName string
	PictureURL  string
}

func init() {
	codeMap = map[int]string{
		1: LINE,
		2: Kakao,
	}
}

func GetCode(service string) int {
	for k, v := range codeMap {
		if v == service {
			return k
		}
	}
	return 0
}

func GetName(code int) string {
	return codeMap[code]
}
