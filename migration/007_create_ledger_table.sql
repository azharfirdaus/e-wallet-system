CREATE TYPE public.ledger_operation AS ENUM (
    'TOPUP',
    'PAYMENT',
    'SEND',
    'RECEIVE'
);

CREATE TABLE public.ledger (
    id BIGSERIAL PRIMARY KEY,
    debit BIGINT NOT NULL DEFAULT 0,
    credit BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ledger_debit_non_negative
        CHECK (debit >= 0),
    CONSTRAINT ledger_credit_non_negative
        CHECK (credit >= 0)
);

CREATE TRIGGER set_ledger_updated_at
BEFORE UPDATE ON public.ledger
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();
