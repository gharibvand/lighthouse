-- Seed default subscription plans
-- Run this after migrations are complete

INSERT INTO subscription_plans (name, display_name, price, currency, duration_days, max_profiles, max_quality, features) VALUES
('basic', 'Basic', 9.99, 'USD', 30, 1, '720p', '{"ad_supported": true, "download": false, "ultra_hd": false}'),
('standard', 'Standard', 15.99, 'USD', 30, 2, '1080p', '{"ad_supported": false, "download": true, "ultra_hd": false}'),
('premium', 'Premium', 19.99, 'USD', 30, 4, '4k', '{"ad_supported": false, "download": true, "ultra_hd": true}')
ON CONFLICT (name) DO NOTHING;
