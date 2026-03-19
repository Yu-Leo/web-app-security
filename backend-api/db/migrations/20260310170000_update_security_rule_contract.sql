-- +goose Up
UPDATE security_profiles
SET base_action = 'block'
WHERE base_action NOT IN ('allow', 'block');

UPDATE security_rules
SET action = 'block'
WHERE action NOT IN ('allow', 'block');

UPDATE security_rules
SET rule_type = CASE
    WHEN ml_model_id IS NOT NULL OR ml_threshold IS NOT NULL THEN 'ml'
    ELSE 'deterministic'
END;

ALTER TABLE security_profiles
    DROP CONSTRAINT IF EXISTS security_profiles_base_action_chk,
    ADD CONSTRAINT security_profiles_base_action_chk
        CHECK (base_action IN ('allow', 'block'));

ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_action_chk,
    ADD CONSTRAINT security_rules_action_chk
        CHECK (action IN ('allow', 'block')),
    DROP CONSTRAINT IF EXISTS security_rules_rule_type_chk,
    ADD CONSTRAINT security_rules_rule_type_chk
        CHECK (rule_type IN ('deterministic', 'ml')),
    DROP CONSTRAINT IF EXISTS security_rules_ml_pair_chk,
    ADD CONSTRAINT security_rules_ml_pair_chk
        CHECK (
            (rule_type = 'ml' AND ml_model_id IS NOT NULL AND ml_threshold IS NOT NULL) OR
            (rule_type = 'deterministic' AND ml_model_id IS NULL AND ml_threshold IS NULL)
        );

-- +goose Down
ALTER TABLE security_rules
    DROP CONSTRAINT IF EXISTS security_rules_ml_pair_chk,
    ADD CONSTRAINT security_rules_ml_pair_chk
        CHECK (
            (ml_model_id IS NULL AND ml_threshold IS NULL) OR
            (ml_model_id IS NOT NULL AND ml_threshold IS NOT NULL)
        ),
    DROP CONSTRAINT IF EXISTS security_rules_rule_type_chk,
    DROP CONSTRAINT IF EXISTS security_rules_action_chk;

ALTER TABLE security_profiles
    DROP CONSTRAINT IF EXISTS security_profiles_base_action_chk;

UPDATE security_rules
SET rule_type = CASE
    WHEN rule_type = 'deterministic' THEN 'base'
    ELSE rule_type
END;
