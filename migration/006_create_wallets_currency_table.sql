CREATE TABLE public.wallets_currency (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    currency public.currency NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wallets_currency_wallet_id
        FOREIGN KEY (wallet_id)
        REFERENCES public.wallets (id),
    CONSTRAINT wallets_currency_wallet_id_currency_unique
        UNIQUE (wallet_id, currency)
);

CREATE TRIGGER set_wallets_currency_updated_at
BEFORE UPDATE ON public.wallets_currency
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();
