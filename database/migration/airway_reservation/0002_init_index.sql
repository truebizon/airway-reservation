-- +migrate Up
-- CREATE UNIQUE INDEX idx_feature_id ON airway_reservation.features USING btree (id);
-- CREATE INDEX idx_features_type ON airway_reservation.features USING btree (type);

-- +migrate Down

