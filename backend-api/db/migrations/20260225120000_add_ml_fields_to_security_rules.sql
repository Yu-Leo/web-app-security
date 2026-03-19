-- +goose Up
ALTER TABLE security_rules
    ADD COLUMN IF NOT EXISTS ml_model_id BIGINT NULL,
    ADD COLUMN IF NOT EXISTS ml_threshold SMALLINT NULL;

ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_ml_threshold_range_chk,
    ADD CONSTRAINT security_rules_ml_threshold_range_chk
        CHECK (ml_threshold IS NULL OR ml_threshold BETWEEN 0 AND 100);

ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_ml_pair_chk,
    ADD CONSTRAINT security_rules_ml_pair_chk
        CHECK (
            (ml_model_id IS NULL AND ml_threshold IS NULL) OR
            (ml_model_id IS NOT NULL AND ml_threshold IS NOT NULL)
        );

ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_ml_model_id_fkey,
    ADD CONSTRAINT security_rules_ml_model_id_fkey
        FOREIGN KEY (ml_model_id) REFERENCES ml_models(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_ml_model_id_fkey,
    DROP CONSTRAINT IF EXISTS security_rules_ml_pair_chk,
    DROP CONSTRAINT IF EXISTS security_rules_ml_threshold_range_chk;

ALTER TABLE security_rules
    DROP COLUMN IF EXISTS ml_threshold,
    DROP COLUMN IF EXISTS ml_model_id;
