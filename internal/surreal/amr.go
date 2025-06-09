package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"

	"github.com/supabase/auth/internal/utilities"
)

type AMRClaim struct {
	ID                   *models.RecordID `json:"id" db:"id"`
	SessionID            *models.RecordID `json:"session_id" db:"session_id"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at" db:"updated_at"`
	AuthenticationMethod *string          `json:"authentication_method" db:"authentication_method"`
}

func (AMRClaim) TableName() string {
	tableName := "mfa_amr_claims"
	return tableName
}

func (cl *AMRClaim) IsAAL2Claim() bool {
	return *cl.AuthenticationMethod == TOTPSignIn.String() || *cl.AuthenticationMethod == MFAPhone.String() || *cl.AuthenticationMethod == MFAWebAuthn.String()
}

func AddClaimToSession(tx *surrealdb.DB, sessionId *models.RecordID, authenticationMethod AuthenticationMethod) error {
	//	id := uuid.Must(uuid.NewV4())

	currentTime := time.Now()
	_, err := surrealdb.Insert[AMRClaim](tx, models.Table("mfa_amr_claims"), &AMRClaim{
		ID:                   nil,
		SessionID:            sessionId,
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		AuthenticationMethod: utilities.To(authenticationMethod.String()),
	})
	// return tx.RawQuery("INSERT INTO "+(&pop.Model{Value: AMRClaim{}}).TableName()+
	// 	`(id, session_id, created_at, updated_at, authentication_method) values (?, ?, ?, ?, ?)
	// 		ON CONFLICT ON CONSTRAINT mfa_amr_claims_session_id_authentication_method_pkey
	// 		DO UPDATE SET updated_at = ?;`, id, sessionId, currentTime, currentTime, authenticationMethod.String(), currentTime).Exec()
	return err
}

func (a *AMRClaim) GetAuthenticationMethod() string {
	if a.AuthenticationMethod == nil {
		return ""
	}
	return *(a.AuthenticationMethod)
}
