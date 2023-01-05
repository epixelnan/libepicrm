package epicrm_apiparts

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPasswordHash(pass string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass),
                                           bcrypt.DefaultCost)
	
	if err != nil {
		return nil, fmt.Errorf("NewBcryptPasswordHash(): bcrypt: %w", err)
	}

	return hash, nil
}
