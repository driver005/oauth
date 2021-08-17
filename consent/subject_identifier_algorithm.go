package consent

import "github.com/driver005/oauth/models"

type SubjectIdentifierAlgorithm interface {
	// Obfuscate derives a pairwise subject identifier from the given string.
	Obfuscate(subject string, client *models.Client) (string, error)
}
