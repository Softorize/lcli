// Package model defines domain types for the LinkedIn API.
package model

// Profile represents a LinkedIn user profile with essential fields.
type Profile struct {
	ID             string `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Headline       string `json:"headline"`
	Vanity         string `json:"vanityName"`
	ProfilePicture string `json:"profilePicture"`
	Email          string `json:"email"`
}

// ProfileResponse maps the raw LinkedIn API response which uses
// localized field names for profile data.
type ProfileResponse struct {
	ID                 string              `json:"id"`
	LocalizedFirstName string              `json:"localizedFirstName"`
	LocalizedLastName  string              `json:"localizedLastName"`
	LocalizedHeadline  string              `json:"localizedHeadline"`
	VanityName         string              `json:"vanityName"`
	ProfilePicture     *profilePictureResp `json:"profilePicture"`
}

// profilePictureResp is the nested picture structure from the LinkedIn API.
type profilePictureResp struct {
	DisplayImage   string          `json:"displayImage"`
	DisplayImageV2 *displayImageV2 `json:"displayImage~"`
}

// displayImageV2 contains the resolved image elements.
type displayImageV2 struct {
	Elements []imageElement `json:"elements"`
}

// imageElement holds a single image variant from the API response.
type imageElement struct {
	Identifiers []imageIdentifier `json:"identifiers"`
}

// imageIdentifier contains the actual URL for an image variant.
type imageIdentifier struct {
	Identifier string `json:"identifier"`
}

// ToProfile converts a raw ProfileResponse into a clean Profile.
func (r *ProfileResponse) ToProfile() *Profile {
	p := &Profile{
		ID:        r.ID,
		FirstName: r.LocalizedFirstName,
		LastName:  r.LocalizedLastName,
		Headline:  r.LocalizedHeadline,
		Vanity:    r.VanityName,
	}

	if r.ProfilePicture != nil && r.ProfilePicture.DisplayImageV2 != nil {
		elems := r.ProfilePicture.DisplayImageV2.Elements
		if len(elems) > 0 && len(elems[0].Identifiers) > 0 {
			p.ProfilePicture = elems[0].Identifiers[0].Identifier
		}
	}

	return p
}
