package bac

import (
	"fmt"
	"os"

	"github.com/kataras/golog"
	"github.com/skip2/go-qrcode"

	"github.com/jinil-ha/blind-msg/db/mysql"
	"github.com/jinil-ha/blind-msg/utils/config"
	"github.com/jinil-ha/blind-msg/utils/token"
)

var qrDir string
var qrURL string
var sendURL string

const bacLength = 16

func init() {
	qrDir = config.GetString("qrcode.dir")
	qrURL = config.GetString("qrcode.url")
	sendURL = config.GetString("qrcode.send_url")
}

// GetQRURL returns url of QR Image file
func GetQRURL(bac string) string {
	return fmt.Sprintf("%s/%s.png", qrURL, bac)
}

// GetSendURL returns url of send page
func GetSendURL(bac string) string {
	return fmt.Sprintf("%s?bac=%s", sendURL, bac)
}

// CreateQR create QR Image(PNG) if not exists.
func CreateQR(bac string) error {
	path := fmt.Sprintf("%s/%s.png", qrDir, bac)

	// check if path exists
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		// error! file may exists or not exists. need to check err
		return err
	}

	golog.Infof("create QR Code file: %s", path)
	url := GetSendURL(bac)
	err := qrcode.WriteFile(url, qrcode.Medium, 256, path)
	if err != nil {
		return err
	}

	err = os.Chmod(path, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetBAC return BAC(Blind Access Code) of user
func GetBAC(service string, userid string) (string, error) {
	var bac string
	err := mysql.GetBAC(service, userid, &bac)

	if err != nil {
		return "", err
	}

	if bac == "" {
		// create new BAC
		bac = token.Generate(bacLength)

		golog.Warnf("create new BAC: %s", bac)
		err = mysql.SetBAC(service, userid, bac)
		if err != nil {
			return "", err
		}
	}
	return bac, nil
}

// GetUserInfo return service and user ID
func GetUserInfo(bac string) (string, string) {
	var service string
	var userID string

	err := mysql.GetUserInfo(bac, &service, &userID)
	if err != nil {
		golog.Errorf("cannot get user info: bac(%s) %s", bac, err)
		return "", ""
	}

	return service, userID
}
