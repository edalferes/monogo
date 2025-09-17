
SELECT r.id AS role_id, r.name AS role_name, p.id AS permission_id, p.name AS permission_name
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON p.id = rp.permission_id
ORDER BY r.id, p.id;
