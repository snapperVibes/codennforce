CREATE TABLE IF NOT EXISTS realestateportal(
    id  int GENERATED ALWAYS AS IDENTITY,
    creationts   timestamp DEFAULT now(),
    url          text NOT NULL,
    address      bytea[],
    municipality bytea[],
    owner        bytea[],
    changemail   bytea[]
);
