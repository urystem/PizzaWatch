-- Откат миграции orders и order_items

-- Сначала удаляем таблицу order_items, потому что она ссылается на orders
DROP TABLE IF EXISTS order_items;

-- Потом удаляем таблицу orders
DROP TABLE IF EXISTS orders;

-- Удаляем enum type order_status
DROP TYPE IF EXISTS order_status;

DROP TABLE IF EXISTS order_status_log;
