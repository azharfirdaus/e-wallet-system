CREATE TYPE public.currency AS ENUM (
    'USD',
    'IDR',
    'SGD',
    'JPY',
    'AUD',
    'CNY'
);

CREATE TYPE public.wallet_status AS ENUM (
    'ACTIVATE',
    'SUSPENDED'
);
