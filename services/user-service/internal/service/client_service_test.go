package service

import (
	"common/pkg/auth"
	"context"
	"fmt"
	"testing"
	"time"

	"user-service/internal/dto"
	"user-service/internal/model"

	"github.com/stretchr/testify/require"
)

func newClientService(
	clientRepo *fakeClientRepo,
	identityRepo *fakeIdentityRepo,
	activationTokenRepo *fakeActivationTokenRepo,
	mailer *fakeMailer,
) *ClientService {
	return NewClientService(
		clientRepo,
		identityRepo,
		activationTokenRepo,
		mailer,
		testConfig(),
	)
}

func TestGetMobileVerificationSecret(t *testing.T) {
	t.Parallel()

	t.Run("returns secret for authenticated client", func(t *testing.T) {
		svc := newClientService(
			&fakeClientRepo{byID: &model.Client{ClientID: 2, MobileVerificationSecret: "JBSWY3DPEHPK3PXP"}},
			&fakeIdentityRepo{},
			&fakeActivationTokenRepo{},
			&fakeMailer{},
		)

		ctx := auth.SetAuthOnContext(context.Background(), &auth.AuthContext{ClientID: uintPtr(2)})
		secret, err := svc.GetMobileVerificationSecret(ctx)
		require.NoError(t, err)
		require.Equal(t, "JBSWY3DPEHPK3PXP", secret)
	})

	t.Run("forbidden when client auth missing", func(t *testing.T) {
		svc := newClientService(&fakeClientRepo{}, &fakeIdentityRepo{}, &fakeActivationTokenRepo{}, &fakeMailer{})

		secret, err := svc.GetMobileVerificationSecret(context.Background())
		require.Error(t, err)
		require.Empty(t, secret)
		require.Contains(t, err.Error(), "client access required")
	})

	t.Run("not found when secret is empty", func(t *testing.T) {
		svc := newClientService(
			&fakeClientRepo{byID: &model.Client{ClientID: 2}},
			&fakeIdentityRepo{},
			&fakeActivationTokenRepo{},
			&fakeMailer{},
		)

		ctx := auth.SetAuthOnContext(context.Background(), &auth.AuthContext{ClientID: uintPtr(2)})
		secret, err := svc.GetMobileVerificationSecret(ctx)
		require.Error(t, err)
		require.Empty(t, secret)
		require.Contains(t, err.Error(), "mobile verification secret not found")
	})
}

func TestClientRegister(t *testing.T) {
	t.Parallel()

	req := &dto.CreateClientRequest{
		FirstName:   "Jane",
		LastName:    "Client",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
		Gender:      "female",
		Email:       "client@example.com",
		Username:    "clientuser",
		PhoneNumber: "0603333333",
		Address:     "Client Address 1",
	}

	tests := []struct {
		name         string
		clientRepo   *fakeClientRepo
		identityRepo *fakeIdentityRepo
		tokenRepo    *fakeActivationTokenRepo
		mailer       *fakeMailer
		expectErr    bool
		errMsg       string
	}{
		{
			name:         "successful registration",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{},
		},
		{
			name:         "email already in use",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{emailExists: true},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{},
			expectErr:    true,
			errMsg:       "email already in use",
		},
		{
			name:         "username already in use",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{usernameExists: true},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{},
			expectErr:    true,
			errMsg:       "username already in use",
		},
		{
			name:         "identity create fails",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{createErr: fmt.Errorf("db error")},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{},
			expectErr:    true,
		},
		{
			name:         "client create fails",
			clientRepo:   &fakeClientRepo{createErr: fmt.Errorf("db error")},
			identityRepo: &fakeIdentityRepo{},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{},
			expectErr:    true,
		},
		{
			name:         "activation token create fails",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{},
			tokenRepo:    &fakeActivationTokenRepo{createErr: fmt.Errorf("db error")},
			mailer:       &fakeMailer{},
			expectErr:    true,
		},
		{
			name:         "email send fails",
			clientRepo:   &fakeClientRepo{},
			identityRepo: &fakeIdentityRepo{},
			tokenRepo:    &fakeActivationTokenRepo{},
			mailer:       &fakeMailer{sendErr: fmt.Errorf("smtp down")},
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newClientService(tt.clientRepo, tt.identityRepo, tt.tokenRepo, tt.mailer)

			client, err := svc.Register(context.Background(), req)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, client)
			require.Equal(t, "Jane", client.FirstName)
			require.Equal(t, auth.IdentityClient, client.Identity.Type)
			require.False(t, client.Identity.Active)
			require.Equal(t, uint(1), client.IdentityID)
			require.NotEmpty(t, client.MobileVerificationSecret)
		})
	}
}
