package consent

import "github.com/driver005/oauth/models"

type SubjectIdentifierAlgorithmPublic struct{}

func NewSubjectIdentifierAlgorithmPublic() *SubjectIdentifierAlgorithmPublic {
	return &SubjectIdentifierAlgorithmPublic{}
}

func (g *SubjectIdentifierAlgorithmPublic) Obfuscate(subject string, client *models.Client) (string, error) {
	return subject, nil
}
