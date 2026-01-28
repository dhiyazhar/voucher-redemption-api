CREATE TABLE vouchers (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    quota BIGINT NOT NULL,
    valid_until TIMESTAMP NOT NULL
);

CREATE TABLE redemption_history (
    id BIGSERIAL PRIMARY KEY,
    voucher_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL, 
    redeemed_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_voucher FOREIGN KEY(voucher_id) REFERENCES vouchers(id)
);

CREATE UNIQUE INDEX idx_unique_claim ON redemption_history (voucher_id, user_id);

INSERT INTO vouchers (code, quota, valid_until) VALUES
    ('DISKON100', 10, NOW() + INTERVAL '7 days'),
    ('SOLDOUT', 0, NOW() + INTERVAL '5 days'),        
    ('EXPIRED1', 50, NOW() - INTERVAL '1 day'),
    ('LASTONE', 1, NOW() + INTERVAL '2 days'),        
    ('FLASH50', 100, NOW() + INTERVAL '1 day');

INSERT INTO redemption_history (voucher_id, user_id, redeemed_at)
VALUES ((SELECT id FROM vouchers WHERE code = 'LASTONE'), 123, NOW());