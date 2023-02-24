package validators

import (
	"aa-survey/internal/utils"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/core"
	"golang.org/x/crypto/ssh"
)

// Validate password
func ValidatePassword(ans interface{}) error {
	if len(ans.(string)) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// Validate multiselect
func ValidateMany(ans interface{}) error {
	if len(ans.([]core.OptionAnswer)) < 2 {
		return fmt.Errorf("select minimum two options")
	}
	return nil
}

// Validate path is file
func ValidateIsFile(ans interface{}) error {
	path := utils.ToAbsPath(ans)
	fi, err := os.Stat(path.(string))
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("%s is directory", path.(string))
	}
	return nil
}

// Validate path is dir
func ValidateIsDir(ans interface{}) error {
	path := utils.ToAbsPath(ans)
	d, err := os.Stat(path.(string))
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	if !d.IsDir() {
		return fmt.Errorf("%s is file", path.(string))
	}
	return nil
}

// Validate file is public key
func ValidatePubKey(ans interface{}) error {
	pub, err := os.ReadFile(utils.ToAbsPath(ans).(string))
	if err != nil {
		return fmt.Errorf("no public key found")
	}
	_, _, _, _, err = ssh.ParseAuthorizedKey(pub)
	if err != nil {
		return fmt.Errorf("%s - %s", ans.(string), err)
	}
	return nil
}

// Validate file is private key
func ValidatePrivKey(ans interface{}) error {
	priv, err := os.ReadFile(utils.ToAbsPath(ans).(string))
	if err != nil {
		return fmt.Errorf("no private key found")
	}
	_, err = ssh.ParsePrivateKey(priv)
	if err != nil {
		return fmt.Errorf("%s - %s", ans.(string), err)
	}
	return nil
}
