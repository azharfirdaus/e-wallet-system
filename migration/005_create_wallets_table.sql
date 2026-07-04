CREATE TABLE public.wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status public.wallet_status NOT NULL,
    closing_balance BIGINT NOT NULL DEFAULT 0,
    closing_balance_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT wallets_closing_balance_non_negative
        CHECK (closing_balance >= 0),
    CONSTRAINT fk_wallets_user_id
        FOREIGN KEY (user_id)
        REFERENCES public.users (id)
);

CREATE TRIGGER set_wallets_updated_at
BEFORE UPDATE ON public.wallets
FOR EACH ROW
EXECUTE FUNCTION public.set_updated_at();
