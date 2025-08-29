-- Главная таблица заказов
CREATE TABLE IF NOT EXISTS orders (
  order_uid           text PRIMARY KEY,
  track_number        text,
  entry               text,
  locale              text,
  internal_signature  text,
  customer_id         text,
  delivery_service    text,
  shardkey            text,
  sm_id               int,
  date_created        timestamptz,
  oof_shard           text,
  raw_json            jsonb NOT NULL,    
  created_at          timestamptz DEFAULT now(),
  updated_at          timestamptz DEFAULT now()
);

-- Таблица доставки
CREATE TABLE IF NOT EXISTS deliveries (
  order_uid text PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
  name      text,
  phone     text,
  zip       text,
  city      text,
  address   text,
  region    text,
  email     text
);

-- Таблица оплаты
CREATE TABLE IF NOT EXISTS payments (
  order_uid     text PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction   text,
  request_id    text,
  currency      text,
  provider      text,
  amount        int,
  payment_dt    timestamptz,
  bank          text,
  delivery_cost int,
  goods_total   int,
  custom_fee    int
);

-- Таблица товаров (может быть несколько на один заказ)
CREATE TABLE IF NOT EXISTS items (
  id            bigserial PRIMARY KEY,
  order_uid     text REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id       bigint,
  track_number  text,
  price         int,
  rid           text,
  name          text,
  sale          int,
  size          text,
  total_price   int,
  nm_id         bigint,
  brand         text,
  status        int
);

-- Индексы для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_orders_updated_at ON orders(updated_at);
