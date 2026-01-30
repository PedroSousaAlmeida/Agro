CREATE OR REPLACE VIEW client_stats AS
SELECT
    c.id,
    c.name,
    c.slug,
    c.max_users,
    COUNT(DISTINCT cu.user_id) as current_users,
    c.max_users - COUNT(DISTINCT cu.user_id) as available_slots,
    COUNT(DISTINCT m.id) as total_monitoramentos,
    COUNT(DISTINCT a.id) as total_areas,
    c.active,
    c.created_at
FROM clients c
LEFT JOIN client_users cu ON c.id = cu.client_id AND cu.active = true
LEFT JOIN monitoramentos m ON c.id = m.client_id
LEFT JOIN areas_monitoramento a ON c.id = a.client_id
GROUP BY c.id;
