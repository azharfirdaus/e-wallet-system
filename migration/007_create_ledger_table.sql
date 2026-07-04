CREATE TYPE public.ledger_operation AS ENUM (
    'TOPUP',
    'PAYMENT',
    'SEND',
    'RECEIVE'
);

CREATE TABLE public.ledger (
    id BIGSERIAL PRIMARY KEY,
    wallet_currency_id BIGINT NOT NULL,
    debit BIGINT NOT NULL DEFAULT 0,
    credit BIGINT NOT NULL DEFAULT 0,
    reference public.ledger_reference NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ledger_wallet_currency_id
        FOREIGN KEY (wallet_currency_id)
        REFERENCES public.wallets_currency (id),
    CONSTRAINT ledger_debit_non_negative
        CHECK (debit >= 0),
    CONSTRAINT ledger_credit_non_negative
        CHECK (credit >= 0)
);

CREATE INDEX idx_ledger_wallet_currency_id
ON public.ledger (wallet_currency_id);

CREATE INDEX idx_ledger_created_at
ON public.ledger (created_at);
