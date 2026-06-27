package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"time"
)

// Device represents a device in the database
type Device struct {
	ID               int
	SecretHash       []byte
	CreatedAt        time.Time
	LastActivatedAt  *time.Time
	TotalActivations int
	Blacklisted      bool
}

// CreateDevice creates a new device in the database and returns the created device with its ID
func (d *DBClient) CreateDevice(ctx context.Context, secret string) (*Device, error) {
	query := `
		INSERT INTO devices (secret_hash, created_at, total_activations, blacklisted)
		VALUES ($1, $2, $3, $4)
		RETURNING id, secret_hash, created_at, last_activated_at, total_activations, blacklisted
	`

	device := &Device{}
	now := time.Now().UTC()

	secretHash := sha256.Sum256([]byte(secret))
	secretHashBytes := secretHash[:]

	err := d.db.QueryRowContext(ctx, query, secretHashBytes, now, 0, false).Scan(
		&device.ID,
		&device.SecretHash,
		&device.CreatedAt,
		&device.LastActivatedAt,
		&device.TotalActivations,
		&device.Blacklisted,
	)

	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to create device")
		return nil, err
	}

	d.logger.Info().Int("deviceId", device.ID).Msg("Device created successfully")
	return device, nil
}

// GetDeviceByID retrieves a device from the database by its ID
func (d *DBClient) GetDeviceByID(ctx context.Context, id int) (*Device, error) {
	query := `
		SELECT id, secret_hash, created_at, last_activated_at, total_activations, blacklisted
		FROM devices
		WHERE id = $1
	`

	device := &Device{}

	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&device.ID,
		&device.SecretHash,
		&device.CreatedAt,
		&device.LastActivatedAt,
		&device.TotalActivations,
		&device.Blacklisted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Debug().Int("deviceId", id).Msg("Device not found")
			return nil, nil
		}
		d.logger.Error().Err(err).Int("deviceId", id).Msg("Failed to retrieve device")
		return nil, err
	}

	return device, nil
}

// BlacklistDevice marks a device as blacklisted
func (d *DBClient) BlacklistDevice(ctx context.Context, id int) error {
	query := `
		UPDATE devices
		SET blacklisted = true
		WHERE id = $1
	`

	result, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		d.logger.Error().Err(err).Int("deviceId", id).Msg("Failed to blacklist device")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		d.logger.Error().Err(err).Int("deviceId", id).Msg("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		d.logger.Debug().Int("deviceId", id).Msg("Device not found for blacklisting")
		return nil
	}

	d.logger.Info().Int("deviceId", id).Msg("Device blacklisted successfully")
	return nil
}

// ActivateDevice atomically updates the device's last activation time and increments the activation count
func (d *DBClient) ActivateDevice(ctx context.Context, id int) error {
	query := `
		UPDATE devices
		SET last_activated_at = $1, total_activations = total_activations + 1
		WHERE id = $2
	`

	result, err := d.db.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		d.logger.Error().Err(err).Int("deviceId", id).Msg("Failed to activate device")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		d.logger.Error().Err(err).Int("deviceId", id).Msg("Failed to get rows affected")
		return err
	}

	if rowsAffected == 0 {
		d.logger.Debug().Int("deviceId", id).Msg("Device not found for activation")
		return nil
	}

	d.logger.Info().Int("deviceId", id).Msg("Device activated successfully")
	return nil
}
