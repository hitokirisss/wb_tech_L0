BEGIN;

-- 1) orders (включая полный raw_json)
INSERT INTO orders (
  order_uid, track_number, entry, locale, internal_signature,
  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, raw_json
) VALUES (
  'b563feb7b2b84b6test',
  'WBILMTESTTRACK',
  'WBIL',
  'en',
  '',
  'test',
  'meest',
  '9',
  99,
  '2021-11-26T06:22:19Z',
  '1',
  $${
     "order_uid": "b563feb7b2b84b6test",
     "track_number": "WBILMTESTTRACK",
     "entry": "WBIL",
     "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
     },
     "payment": {
        "transaction": "b563feb7b2b84b6test",
        "request_id": "",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
     },
     "items": [
        {
           "chrt_id": 9934930,
           "track_number": "WBILMTESTTRACK",
           "price": 453,
           "rid": "ab4219087a764ae0btest",
           "name": "Mascaras",
           "sale": 30,
           "size": "0",
           "total_price": 317,
           "nm_id": 2389212,
           "brand": "Vivienne Sabo",
           "status": 202
        }
     ],
     "locale": "en",
     "internal_signature": "",
     "customer_id": "test",
     "delivery_service": "meest",
     "shardkey": "9",
     "sm_id": 99,
     "date_created": "2021-11-26T06:22:19Z",
     "oof_shard": "1"
  }$$::jsonb
)
ON CONFLICT (order_uid) DO UPDATE SET
  track_number       = EXCLUDED.track_number,
  entry              = EXCLUDED.entry,
  locale             = EXCLUDED.locale,
  internal_signature = EXCLUDED.internal_signature,
  customer_id        = EXCLUDED.customer_id,
  delivery_service   = EXCLUDED.delivery_service,
  shardkey           = EXCLUDED.shardkey,
  sm_id              = EXCLUDED.sm_id,
  date_created       = EXCLUDED.date_created,
  oof_shard          = EXCLUDED.oof_shard,
  raw_json           = EXCLUDED.raw_json,
  updated_at         = now();

-- 2) deliveries
INSERT INTO deliveries (
  order_uid, name, phone, zip, city, address, region, email
) VALUES (
  'b563feb7b2b84b6test',
  'Test Testov',
  '+9720000000',
  '2639809',
  'Kiryat Mozkin',
  'Ploshad Mira 15',
  'Kraiot',
  'test@gmail.com'
)
ON CONFLICT (order_uid) DO UPDATE SET
  name    = EXCLUDED.name,
  phone   = EXCLUDED.phone,
  zip     = EXCLUDED.zip,
  city    = EXCLUDED.city,
  address = EXCLUDED.address,
  region  = EXCLUDED.region,
  email   = EXCLUDED.email;

-- 3) payments (payment_dt из epoch → to_timestamp(...))
INSERT INTO payments (
  order_uid, transaction, request_id, currency, provider,
  amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
) VALUES (
  'b563feb7b2b84b6test',
  'b563feb7b2b84b6test',
  '',
  'USD',
  'wbpay',
  1817,
  to_timestamp(1637907727),
  'alpha',
  1500,
  317,
  0
)
ON CONFLICT (order_uid) DO UPDATE SET
  transaction   = EXCLUDED.transaction,
  request_id    = EXCLUDED.request_id,
  currency      = EXCLUDED.currency,
  provider      = EXCLUDED.provider,
  amount        = EXCLUDED.amount,
  payment_dt    = EXCLUDED.payment_dt,
  bank          = EXCLUDED.bank,
  delivery_cost = EXCLUDED.delivery_cost,
  goods_total   = EXCLUDED.goods_total,
  custom_fee    = EXCLUDED.custom_fee;

-- 4) items (в примере один товар; если будет несколько — добавляешь ещё INSERT)
-- На практике перед вставкой списка удобно чистить старые позиции заказа:
DELETE FROM items WHERE order_uid = 'b563feb7b2b84b6test';

INSERT INTO items (
  order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
) VALUES (
  'b563feb7b2b84b6test',
  9934930,
  'WBILMTESTTRACK',
  453,
  'ab4219087a764ae0btest',
  'Mascaras',
  30,
  '0',
  317,
  2389212,
  'Vivienne Sabo',
  202
);

COMMIT;
