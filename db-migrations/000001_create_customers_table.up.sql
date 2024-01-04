CREATE TABLE customers (
    id UUID NOT NULL PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    residence_address VARCHAR NOT NULL,
    birth_date DATE NOT NULL,
    kyc_status VARCHAR NOT NULL,

    updated_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);
